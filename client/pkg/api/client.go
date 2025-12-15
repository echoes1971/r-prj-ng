package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/echoes1971/r-prj-ng/client/pkg/models"
	"github.com/schollz/progressbar/v3"
)

// Client is an HTTP client for the ρBee API
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login authenticates with the API and returns a JWT token
func (c *Client) Login(username, password string) (string, error) {
	loginReq := models.LoginRequest{
		Login: username,
		Pwd:   password,
	}

	body, err := json.Marshal(loginReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/login", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var loginResp models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return loginResp.AccessToken, nil
}

// Get retrieves an object by ID
func (c *Client) Get(objectID string) (*models.DBObject, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/content/%s", c.BaseURL, objectID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Backend returns { data: {...}, metadata: {...} }
	var response struct {
		Data     models.DBObject        `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Copy classname from metadata to data if missing
	if response.Data.Classname == "" {
		if classname, ok := response.Metadata["classname"].(string); ok {
			response.Data.Classname = classname
		}
	}

	return &response.Data, nil
}

// Create creates a new object
func (c *Client) Create(obj *models.DBObject) (*models.DBObject, error) {
	body, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/objects", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Backend returns { data: {...}, metadata: {...} }
	var response struct {
		Data     models.DBObject        `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Copy classname from metadata to data if missing
	if response.Data.Classname == "" {
		if classname, ok := response.Metadata["classname"].(string); ok {
			response.Data.Classname = classname
		}
	}

	return &response.Data, nil
}

// Update updates an existing object
func (c *Client) Update(objectID string, obj *models.DBObject) error {
	// Marshal to JSON then unmarshal to map to remove classname field
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %w", err)
	}

	var objMap map[string]interface{}
	if err := json.Unmarshal(data, &objMap); err != nil {
		return fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	// Remove classname field (it belongs in metadata, not data)
	delete(objMap, "classname")

	body, err := json.Marshal(objMap)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/objects/%s", c.BaseURL, objectID), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Delete deletes an object
func (c *Client) Delete(objectID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/objects/%s", c.BaseURL, objectID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// UploadFile uploads a file to a folder
func (c *Client) UploadFile(filePath, folderID, name, description, permissions string, showProgress bool) (*models.DBObject, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file with optional progress bar
	if showProgress {
		bar := progressbar.DefaultBytes(
			fileInfo.Size(),
			"Uploading",
		)
		if _, err := io.Copy(io.MultiWriter(part, bar), file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		if _, err := io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
	}

	// Add form fields
	if name == "" {
		name = filepath.Base(filePath)
	}
	writer.WriteField("name", name)
	writer.WriteField("father_id", folderID)
	if description != "" {
		writer.WriteField("description", description)
	}
	if permissions == "" {
		permissions = "rw-r-----"
	}
	writer.WriteField("permissions", permissions)

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", c.BaseURL+"/objects", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Increase timeout for large files
	c.HTTPClient.Timeout = 5 * time.Minute

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Backend returns { data: {...}, metadata: {...} }
	var response struct {
		Data     models.DBObject        `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Copy classname from metadata to data if missing
	if response.Data.Classname == "" {
		if classname, ok := response.Metadata["classname"].(string); ok {
			response.Data.Classname = classname
		}
	}

	return &response.Data, nil
}

// UpdateFile updates the content and metadata of an existing DBFile
func (c *Client) UpdateFile(obj *models.DBObject, filePath string, showProgress bool) (*models.DBObject, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file with optional progress bar
	if showProgress {
		bar := progressbar.DefaultBytes(
			fileInfo.Size(),
			"Uploading",
		)
		if _, err := io.Copy(io.MultiWriter(part, bar), file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		if _, err := io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
	}

	// Add object fields to the form (if provided)
	if obj != nil {
		if obj.Name != "" {
			writer.WriteField("name", obj.Name)
		}
		if obj.Description != "" {
			writer.WriteField("description", obj.Description)
		}
		if obj.Permissions != "" {
			writer.WriteField("permissions", obj.Permissions)
		}
		if obj.FatherID != "" {
			writer.WriteField("father_id", obj.FatherID)
		}
		if obj.HTML != "" {
			writer.WriteField("html", obj.HTML)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request - use PUT to update existing file
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/objects/%s", c.BaseURL, obj.ID), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Increase timeout for large files
	c.HTTPClient.Timeout = 5 * time.Minute

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file update failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Backend returns { data: {...}, metadata: {...} }
	var response struct {
		Data     models.DBObject        `json:"data"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Copy classname from metadata to data if missing
	if response.Data.Classname == "" {
		if classname, ok := response.Metadata["classname"].(string); ok {
			response.Data.Classname = classname
		}
	}

	return &response.Data, nil
}

// DownloadFile downloads a file
func (c *Client) DownloadFile(fileID, outputPath string, showProgress bool) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/files/%s/download", c.BaseURL, fileID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download with optional progress bar
	if showProgress && resp.ContentLength > 0 {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"Downloading",
		)
		if _, err := io.Copy(outFile, io.TeeReader(resp.Body, bar)); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}
	} else {
		if _, err := io.Copy(outFile, resp.Body); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}
	}

	return nil
}

// Search searches for objects
func (c *Client) Search(classname, query string) ([]models.DBObject, error) {
	url := fmt.Sprintf("%s/objects/search?classname=%s", c.BaseURL, classname)
	if query != "" {
		url += fmt.Sprintf("&name=%s", query)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Backend returns { success: true, objects: [...] }
	var response struct {
		Success bool              `json:"success"`
		Objects []models.DBObject `json:"objects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Objects, nil
}

// SearchWithAllFields searches with searchJson to get all object fields including father_id
func (c *Client) SearchWithAllFields(classname, query string) ([]models.DBObject, error) {
	url := fmt.Sprintf("%s/objects/search?classname=%s&searchJson={}", c.BaseURL, classname)
	if query != "" {
		url += fmt.Sprintf("&name=%s", query)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		Success bool              `json:"success"`
		Objects []models.DBObject `json:"objects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Objects, nil
}

// GetChildren retrieves children of a folder
func (c *Client) GetChildren(folderID string) ([]models.DBObject, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nav/children/%s", c.BaseURL, folderID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get children failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Backend returns { children: [{data: {...}, metadata: {...}}], count: N }
	var response struct {
		Children []struct {
			Data     models.DBObject        `json:"data"`
			Metadata map[string]interface{} `json:"metadata"`
		} `json:"children"`
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract just the data objects
	var children []models.DBObject
	for _, child := range response.Children {
		// Copy classname from metadata to data if missing
		if child.Data.Classname == "" {
			if classname, ok := child.Metadata["classname"].(string); ok {
				child.Data.Classname = classname
			}
		}
		children = append(children, child.Data)
	}

	return children, nil
}

// GetAllChildren retrieves ALL children of a folder including index pages
// This uses multiple search queries since /nav/children filters out index pages
func (c *Client) GetAllChildren(folderID string) ([]models.DBObject, error) {
	// List of common object types to search for
	types := []string{"DBFile", "DBFolder", "DBPage", "DBCompany", "DBPerson", "DBNews", "DBNote"}

	allChildren := make([]models.DBObject, 0)
	seen := make(map[string]bool) // Avoid duplicates

	for _, objType := range types {
		// Search for objects of this type - use searchJson to get all fields including father_id
		req, err := http.NewRequest("GET",
			fmt.Sprintf("%s/objects/search?classname=%s&name=&searchJson={}", c.BaseURL, objType),
			nil)
		if err != nil {
			continue // Skip this type on error
		}
		req.Header.Set("Authorization", "Bearer "+c.Token)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		var searchResp struct {
			Success bool              `json:"success"`
			Objects []models.DBObject `json:"objects"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Filter by father_id and deduplicate
		for _, obj := range searchResp.Objects {
			if obj.FatherID == folderID && !seen[obj.ID] {
				allChildren = append(allChildren, obj)
				seen[obj.ID] = true
			}
		}
	}

	return allChildren, nil
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var ollamaURL string
var ollamaModel string

func OllamaInit(url, model string) error {
	ollamaURL = url
	ollamaModel = model

	go UpdateOllamaDefaultPageResponse("en")
	log.Printf("Ollama initialized with URL: %s and Model: %s\n", ollamaURL, ollamaModel)

	return nil
}

func OllamaHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Prompt string `json:"prompt"`
	}
	type Response struct {
		Response string `json:"response"`
		Error    string `json:"error,omitempty"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	respText, err := CallOllama(req.Prompt)
	if err != nil {
		res := Response{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	res := Response{Response: respText}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

var lastDefaultPageResponse string

func DefaultPageOllamaHandler(w http.ResponseWriter, r *http.Request) {
	// // Get Accept-Language header
	// acceptLang := r.Header.Get("Accept-Language")
	// tags, _, err := language.ParseAcceptLanguage(acceptLang)
	// if err != nil || len(tags) == 0 {
	// 	tags = []language.Tag{language.English}
	// }
	// langTag := tags[0].String()

	// Get lang query parameter
	langParam := r.URL.Query().Get("lang")
	if langParam == "" {
		langParam = "en"
	}
	log.Print("Requested language tag: ", langParam)

	// Just use English for now
	langParam = "en"
	go UpdateOllamaDefaultPageResponse(langParam)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": lastDefaultPageResponse,
	})
}

//	var modelList = []string{
//		"llama3.2:latest",
//		"codellama:7b",
//		"qwen3-vl:4b",
//		"mistral:7b", // NO
//		"gemma3:4b",
//	}
var prompts = []string{
	"Provide a short funny star trek based welcome message for the default page of a web application",
	"Provide a short funny star wars based welcome message for the default page of a web application",
	"Provide a short funny hitcher's guide to galaxy based welcome message for the default page of a web application",
	"Provide a short funny positive welcome message for the default page of a web application",
}

func UpdateOllamaDefaultPageResponse(languageTag string) {
	if ollamaURL == "" || ollamaModel == "" {
		log.Print("Ollama not configured, skipping default page response update.")
		return
	}

	// ollamaModel = modelList[0]
	prompt := prompts[2]
	switch languageTag {
	case "it", "it-IT":
		prompt += prompt + " in Italian."
	case "en", "en-US":
		prompt += prompt + " in English."
	case "fr", "fr-FR":
		prompt += prompt + " in French."
	case "de", "de-DE":
		prompt += prompt + " in German."
	}
	prompt += " Use only <h2>, <p>, <b>, <br> tags. No links <a>. Output html code only."

	log.Print("Prompt: ", prompt)

	respText, err := CallOllama(prompt)
	if err != nil {
		log.Printf("Error getting default page response from Ollama: %v\n", err)
		respText = "Welcome to our application!"
	}

	lastDefaultPageResponse = respText
}

// CallOllama sends a prompt to the Ollama API and returns the response
func CallOllama(prompt string) (string, error) {
	if ollamaURL == "" || ollamaModel == "" {
		return "", fmt.Errorf("Ollama not configured")
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": ollamaModel,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s", ollamaURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Ollama response: %s\n", string(body))

	// Parse the response
	// {"model":"llama3.2:latest","created_at":"2025-11-13T13:17:45.487900903Z","message":{"role":"assistant","content":"I'm an artificial intelligence model known as Llama. Llama stands for \"Large Language Model Meta AI.\""},"done":true,"done_reason":"stop","total_duration":1593273648,"load_duration":103413131,"prompt_eval_count":29,"prompt_eval_duration":132183696,"eval_count":23,"eval_duration":1338392262}

	var ollamaResponseSingle struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(body, &ollamaResponseSingle); err == nil && ollamaResponseSingle.Message.Content != "" {
		content := ollamaResponseSingle.Message.Content
		content = strings.ReplaceAll(content, "```html", "")
		content = strings.ReplaceAll(content, "```", "")
		content = strings.ReplaceAll(content, "[app name]", "our web site")
		return content, nil
	}

	return "", fmt.Errorf("invalid response from Ollama")

}

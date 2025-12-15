package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/echoes1971/r-prj-ng/client/pkg/api"
	"github.com/echoes1971/r-prj-ng/client/pkg/auth"
	"github.com/echoes1971/r-prj-ng/client/pkg/models"
	"github.com/spf13/cobra"
)

var (
	importFolder      string
	importPreserveIDs bool
	importForce       bool
)

var importCmd = &cobra.Command{
	Use:   "import <directory>",
	Short: "Import objects from exported directory",
	Long: `Import objects from an exported directory structure.
	
Examples:
  # Import to root folder
  rhobee import ./backup/ --folder 0

  # Import preserving original IDs
  rhobee import ./backup/ --folder 0 --preserve-ids
  
  # Force update if objects already exist (upsert mode)
  rhobee import ./backup/ --folder 0 --preserve-ids --force`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(&importFolder, "folder", "", "Target folder ID (required)")
	importCmd.Flags().BoolVar(&importPreserveIDs, "preserve-ids", false, "Preserve original object IDs")
	importCmd.Flags().BoolVar(&importForce, "force", false, "Update existing objects instead of failing")
	importCmd.MarkFlagRequired("folder")
}

func runImport(cmd *cobra.Command, args []string) error {
	importDir := args[0]

	// Check if directory exists
	if _, err := os.Stat(importDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", importDir)
	}

	// Read manifest
	manifestPath := filepath.Join(importDir, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest ExportManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Get token
	tokenManager, err := auth.NewTokenManager()
	if err != nil {
		return fmt.Errorf("failed to create token manager: %w", err)
	}

	instance, _ := cmd.Flags().GetString("instance")
	url, _, token, err := tokenManager.GetToken(instance)
	if err != nil {
		return fmt.Errorf("not logged in. Run 'rhobee login' first: %w", err)
	}

	// Create API client
	client := api.NewClient(url, token)

	fmt.Printf("Importing %d objects from %s...\n", manifest.TotalObjects, importDir)
	if importPreserveIDs {
		fmt.Println("  Preserving original IDs")
	}
	if importForce {
		fmt.Println("  Force mode: will update existing objects")
	}

	// Map old IDs to new IDs
	idMap := make(map[string]string)

	// Import objects (must maintain parent-child order)
	imported := make(map[string]bool)

	for len(imported) < len(manifest.Objects) {
		progress := false

		for _, objInfo := range manifest.Objects {
			if imported[objInfo.ID] {
				continue
			}

			// Read object
			objPath := filepath.Join(importDir, objInfo.Path)
			objData, err := os.ReadFile(objPath)
			if err != nil {
				return fmt.Errorf("failed to read object file %s: %w", objPath, err)
			}

			var obj models.DBObject
			if err := json.Unmarshal(objData, &obj); err != nil {
				return fmt.Errorf("failed to parse object: %w", err)
			}

			// Determine parent ID
			var parentID string
			if obj.ID == manifest.RootID {
				// This is the root object - import it under the target folder
				parentID = importFolder
			} else if newParentID, ok := idMap[obj.FatherID]; ok {
				// Parent has been imported, use its new ID
				parentID = newParentID
			} else {
				// Parent not yet imported, skip for now
				continue
			}

			// Update father_id
			obj.FatherID = parentID

			// Handle ID preservation
			originalID := obj.ID
			if !importPreserveIDs {
				obj.ID = "" // Let backend generate new ID
			}

			// Import based on type
			var newObj *models.DBObject
			var createErr error

			// Check if object already exists (only in force mode with preserve-ids)
			var objectExists bool
			if importForce && importPreserveIDs {
				existingObj, err := client.Get(originalID)
				objectExists = (err == nil && existingObj != nil)
			}

			if obj.Classname == "DBFile" && objInfo.FilePath != "" {
				// Upload file
				filePath := filepath.Join(importDir, objInfo.FilePath)

				if objectExists {
					// File exists - update it with new content and metadata
					fmt.Printf("  Updating file: %s\n", obj.Name)
					obj.ID = originalID // Ensure ID is set for update
					newObj, createErr = client.UpdateFile(&obj, filePath, false)
					if createErr != nil {
						fmt.Printf("    Warning: failed to update file: %v\n", createErr)
						continue
					}
				} else {
					// Upload new file
					fmt.Printf("  Uploading file: %s\n", obj.Name)
					newObj, createErr = client.UploadFile(filePath, parentID, obj.Name, obj.Description, obj.Permissions, false)

					if createErr != nil {
						fmt.Printf("    Warning: failed to import file: %v\n", createErr)
						continue
					}
				}
			} else {
				if objectExists {
					// Update existing object
					fmt.Printf("  Updating %s: %s\n", obj.Classname, obj.Name)
					err = client.Update(originalID, &obj)
					if err == nil {
						newObj = &obj
						newObj.ID = originalID
					} else {
						fmt.Printf("    Warning: update failed: %v\n", err)
						continue
					}
				} else {
					// Create new object
					fmt.Printf("  Creating %s: %s\n", obj.Classname, obj.Name)
					newObj, createErr = client.Create(&obj)
					if createErr != nil {
						fmt.Printf("    Warning: failed to create object: %v\n", createErr)
						continue
					}
				}
			}

			// Map old ID to new ID
			idMap[objInfo.ID] = newObj.ID
			imported[objInfo.ID] = true
			progress = true
		}

		if !progress {
			// No progress made, possibly broken references
			remaining := len(manifest.Objects) - len(imported)
			return fmt.Errorf("import stuck: %d objects remaining (possible broken references)", remaining)
		}
	}

	fmt.Printf("✓ Imported %d objects\n", len(imported))

	return nil
}

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/echoes1971/r-prj-ng/client/pkg/api"
	"github.com/echoes1971/r-prj-ng/client/pkg/auth"
	"github.com/spf13/cobra"
)

var (
	uploadFolder      string
	uploadFileID      string
	uploadName        string
	uploadDescription string
	uploadPermissions string
	uploadNoProgress  bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload <file> [file2] [file3] ...",
	Short: "Upload one or more files",
	Long: `Upload one or more files to ρBee, or update an existing DBFile with new content.
	
Examples:
  # Upload single file to folder
  rhobee upload photo.jpg --folder folder_id

  # Update existing DBFile with new content
  rhobee upload new_photo.jpg --file-id dbfile_id

  # Upload multiple files
  rhobee upload *.jpg --folder folder_id

  # Upload specific files
  rhobee upload file1.jpg file2.png file3.pdf --folder folder_id

  # Upload with custom description (applies to all files)
  rhobee upload *.pdf --folder folder_id --description "Documents"

  # Upload without progress bar
  rhobee upload large.zip --folder folder_id --no-progress`,
	Args: cobra.MinimumNArgs(1),
	RunE: runUpload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&uploadFolder, "folder", "", "Parent folder ID (for new uploads)")
	uploadCmd.Flags().StringVar(&uploadFileID, "file-id", "", "Existing DBFile ID (to update its content)")
	uploadCmd.Flags().StringVar(&uploadName, "name", "", "File name (only for single file upload)")
	uploadCmd.Flags().StringVar(&uploadDescription, "description", "", "File description (applies to all files)")
	uploadCmd.Flags().StringVar(&uploadPermissions, "permissions", "rw-r-----", "Permissions (default: rw-r-----)")
	uploadCmd.Flags().BoolVar(&uploadNoProgress, "no-progress", false, "Disable progress bar")

	// Make folder and file-id mutually exclusive - one must be provided
	uploadCmd.MarkFlagsMutuallyExclusive("folder", "file-id")
}

func runUpload(cmd *cobra.Command, args []string) error {
	filePaths := args

	// Validate flags
	if uploadFolder == "" && uploadFileID == "" {
		return fmt.Errorf("either --folder or --file-id must be specified")
	}

	if uploadFileID != "" && len(filePaths) > 1 {
		return fmt.Errorf("--file-id can only be used with a single file")
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

	showProgress := !uploadNoProgress

	// Handle file update mode
	if uploadFileID != "" {
		filePath := filePaths[0]
		fmt.Printf("Updating DBFile %s with new content...\n", uploadFileID)

		// Get existing object to preserve its metadata
		existingObj, err := client.Get(uploadFileID)
		if err != nil {
			return fmt.Errorf("failed to get existing file: %w", err)
		}

		// Update with new file content
		updated, err := client.UpdateFile(existingObj, filePath, showProgress)
		if err != nil {
			return fmt.Errorf("failed to update file: %w", err)
		}

		fmt.Printf("\n✓ File content updated successfully\n")
		fmt.Printf("  ID: %s\n", updated.ID)
		fmt.Printf("  Name: %s\n", updated.Name)
		fmt.Printf("  Type: %s\n", updated.Mime)

		return nil
	}

	// Handle normal upload mode
	successCount := 0
	failCount := 0

	fmt.Printf("Uploading %d file(s)...\n", len(filePaths))

	for _, filePath := range filePaths {
		// Use filename as name if not specified or uploading multiple files
		name := uploadName
		if name == "" || len(filePaths) > 1 {
			name = filepath.Base(filePath)
		}

		fmt.Printf("\n[%d/%d] Uploading: %s\n", successCount+failCount+1, len(filePaths), name)

		// Upload file
		uploaded, err := client.UploadFile(filePath, uploadFolder, name, uploadDescription, uploadPermissions, showProgress)
		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n", err)
			failCount++
			continue
		}

		fmt.Printf("  ✓ Uploaded successfully (ID: %s)\n", uploaded.ID)
		successCount++
	}

	fmt.Printf("\n")
	fmt.Printf("Summary: %d succeeded, %d failed\n", successCount, failCount)

	if failCount > 0 {
		return fmt.Errorf("%d file(s) failed to upload", failCount)
	}

	return nil
}

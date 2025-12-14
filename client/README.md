# ρBee CLI Client

Command-line interface for ρBee CMS. Automate content management, migrate objects between instances, and integrate ρBee with other tools.

## Features

- 🔐 **Authentication**: Login and token management
- 📄 **Object CRUD**: Create, read, update, delete objects
- 📁 **File Operations**: Upload and download files with progress
- 🔄 **Export/Import**: Migrate content between instances
- 🔍 **Search**: Find objects across your instance
- 🤖 **Automation**: Perfect for cronjobs and scripts

## Installation

```bash
cd client
go build -o rhobee
sudo mv rhobee /usr/local/bin/
```

Or install directly:
```bash
go install github.com/echoes1971/r-prj-ng/client@latest
```

## Quick Start

### Login
```bash
# Interactive login
rhobee login
> URL: https://mybee.com
> Username: admin
> Password: ****
✓ Logged in as admin
✓ Token saved to ~/.rhobee/config.yaml

# Non-interactive (for scripts)
rhobee login --url https://mybee.com --user admin --password secret
```

### Object Operations

**Get object**
```bash
# Get object details (JSON output)
rhobee get c123abc456

# Save to file
rhobee get c123abc456 --output page.json
```

**Create object**
```bash
# Create from command line
rhobee create --type DBPage \
  --name "My Page" \
  --description "Page description" \
  --father-id 0 \
  --permissions "rw-r-----"

# Create from JSON file
rhobee create --file page.json
```

**Update object**
```bash
# Update fields
rhobee update c123abc456 --name "New Name" --description "New description"

# Update from file
rhobee update c123abc456 --file updated.json
```

**Delete object**
```bash
rhobee delete c123abc456
rhobee delete c123abc456 --force  # skip confirmation
```

### File Operations

**Upload file**
```bash
# Upload single file
rhobee upload photo.jpg --folder f789def

# Upload with metadata
rhobee upload photo.jpg \
  --folder f789def \
  --name "My Photo" \
  --description "Sunset photo" \
  --permissions "rw-r-----"

# Upload multiple files
rhobee upload *.jpg --folder f789def
```

**Download file**
```bash
# Download file
rhobee download file123abc --output ./photo.jpg

# Auto-detect filename from server
rhobee download file123abc
```

### Export/Import

**Export folder**
```bash
# Export folder with all children and files
rhobee export f789def --output ./backup/

# Output:
# Exporting folder "Gallery"...
# ├─ 3 pages
# ├─ 12 files (45.2 MB)
# └─ 2 subfolders
# ✓ Exported to ./backup/f789def/
```

**Import folder**
```bash
# Import to another instance
rhobee login --url https://newinstance.com --user admin --password secret
rhobee import ./backup/f789def/ --folder 0  # import to root

# Preserve IDs (recommended for migration)
rhobee import ./backup/f789def/ --folder 0 --preserve-ids

# Output:
# Importing from ./backup/f789def/...
# ├─ Creating pages...
# ├─ Uploading files...
# └─ Preserving permissions...
# ✓ Import complete
```

### Search & List

**Search objects**
```bash
# Search by name/description
rhobee search "keyword"

# Filter by type
rhobee search "keyword" --type DBPage

# Search in specific folder
rhobee search "keyword" --folder f789def
```

**List folder children**
```bash
# List immediate children
rhobee list f789def

# Recursive listing
rhobee list f789def --recursive

# List all objects of a type (admin)
rhobee list --all --type DBPage
```

## Configuration

Config file location: `~/.rhobee/config.yaml`

```yaml
default_instance: prod

instances:
  prod:
    url: https://mybee.com
    token: eyJhbGc...
    user: admin
    
  dev:
    url: http://localhost:8080
    token: eyJhbGc...
    user: admin

preferences:
  output_format: json  # json, yaml, table
  confirm_delete: true
  upload_chunk_size: 5242880  # 5MB
```

### Multiple Instances

```bash
# Add instance
rhobee login --instance staging --url https://staging.mybee.com

# Use specific instance
rhobee get c123abc --instance staging

# Set default instance
rhobee config set-default staging
```

## Use Cases

### Automation Script
```bash
#!/bin/bash
# backup-content.sh

DATE=$(date +%Y-%m-%d)
BACKUP_DIR="./backups/${DATE}"

# Export main folders
rhobee export root_folder_id --output "${BACKUP_DIR}/root"
rhobee export gallery_id --output "${BACKUP_DIR}/gallery"

# Compress
tar -czf "backup_${DATE}.tar.gz" "${BACKUP_DIR}"
rm -rf "${BACKUP_DIR}"

echo "✓ Backup saved to backup_${DATE}.tar.gz"
```

### Content Migration
```bash
#!/bin/bash
# migrate-to-new-instance.sh

# Export from old instance
rhobee login --instance old --url https://old.mybee.com --user admin
rhobee export root_id --output ./migration/

# Import to new instance
rhobee login --instance new --url https://new.mybee.com --user admin
rhobee import ./migration/root_id/ --folder 0 --preserve-ids

echo "✓ Migration complete"
```

### Bulk Operations
```bash
#!/bin/bash
# bulk-upload-photos.sh

FOLDER_ID="gallery_folder_id"

for photo in /path/to/photos/*.jpg; do
    FILENAME=$(basename "$photo")
    rhobee upload "$photo" \
      --folder "$FOLDER_ID" \
      --name "$FILENAME" \
      --permissions "rw-r-----"
done

echo "✓ Uploaded $(ls /path/to/photos/*.jpg | wc -l) photos"
```

### CI/CD Integration
```bash
# .github/workflows/deploy-content.yml
name: Deploy Content

on:
  push:
    branches: [main]
    paths:
      - 'content/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install ρBee CLI
        run: |
          curl -L https://github.com/echoes1971/r-prj-ng/releases/latest/download/rhobee-linux-amd64 -o rhobee
          chmod +x rhobee
          sudo mv rhobee /usr/local/bin/
      
      - name: Deploy content
        env:
          RHOBEE_URL: ${{ secrets.RHOBEE_URL }}
          RHOBEE_USER: ${{ secrets.RHOBEE_USER }}
          RHOBEE_PASSWORD: ${{ secrets.RHOBEE_PASSWORD }}
        run: |
          rhobee login --url "$RHOBEE_URL" --user "$RHOBEE_USER" --password "$RHOBEE_PASSWORD"
          rhobee import ./content/ --folder root_folder_id --preserve-ids
```

## Export/Import Format

### Directory Structure
```
backup/
├── manifest.json          # Metadata, version, timestamps
├── objects/
│   ├── folder_root.json
│   ├── page_123.json
│   └── file_456.json
└── files/
    ├── file_456_photo.jpg
    └── file_789_document.pdf
```

### Manifest Format
```json
{
  "version": "1.0",
  "exported_at": "2025-12-14T10:30:00Z",
  "source_url": "https://old.rhobee.com",
  "root_object_id": "f789def",
  "total_objects": 15,
  "total_files": 12,
  "total_size_bytes": 47456789,
  "objects": [
    {
      "id": "f789def",
      "type": "DBFolder",
      "path": "objects/folder_root.json",
      "children": ["page_123", "file_456"]
    }
  ]
}
```

### Object JSON Format
```json
{
  "id": "c123abc456",
  "classname": "DBPage",
  "name": "My Page",
  "description": "Page description",
  "father_id": "0",
  "permissions": "rw-r-----",
  "language": "en",
  "html": "<h1>Content</h1>",
  "creator": 1,
  "group_id": "-2",
  "creation_date": "2025-12-14T10:00:00Z",
  "last_modify_date": "2025-12-14T11:00:00Z"
}
```

## ID Preservation

ρBee uses 16-character random IDs. The CLI can preserve these IDs during export/import, making it perfect for:

- **Instance migration**: Keep URLs and references intact
- **Disaster recovery**: Restore exact state
- **Multi-instance sync**: Same IDs across dev/staging/prod

```bash
# Export
rhobee export f789def --output ./backup/

# Import with same IDs
rhobee import ./backup/f789def/ --folder 0 --preserve-ids
```

**Note**: ID preservation requires that the target instance doesn't already have objects with those IDs.

## Security Considerations

### Token Storage
Tokens are stored in `~/.rhobee/config.yaml`:
- File permissions: `0600` (read/write for owner only)
- Tokens are short-lived JWT tokens
- Recommended for server/automation environments

### Best Practices
- ✅ Use dedicated service accounts for automation
- ✅ Set restrictive permissions on config file
- ✅ Use environment variables in CI/CD: `RHOBEE_TOKEN`
- ✅ Rotate passwords regularly
- ✅ Limit token lifetime in backend configuration

### Environment Variables
```bash
# Override config file
export RHOBEE_URL=https://mybee.com
export RHOBEE_TOKEN=eyJhbGc...

rhobee get c123abc  # uses env vars
```

## Development

### Build from source
```bash
cd client
go mod download
go build -o rhobee
```

### Run tests
```bash
go test ./...
```

### Cross-compile
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o rhobee-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o rhobee-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o rhobee-windows-amd64.exe
```

## Roadmap

### Phase 1 (MVP) ✅ COMPLETED
- [x] Project structure
- [x] Login command (interactive + non-interactive)
- [x] Token management (multi-instance support)
- [x] Get object (JSON output, file output)
- [x] Create object (from flags or JSON file)
- [x] Config file support (~/.rhobee/config.yaml)

### Phase 2 (File Support) ✅ COMPLETED
- [x] Upload file (with progress bar)
- [x] Download file (with progress bar)
- [x] Multiple file upload support
- [x] Custom permissions and metadata

### Phase 3 (Advanced) 🚧 PLANNED
- [x] Export folder (recursive with files)
- [x] Import folder (with --preserve-ids)
- [x] Search command (by name, type, folder)
- [x] List command (children, recursive)

### Phase 4 (Future) 📋 BACKLOG
- [x] Update/Delete commands
- [ ] Colored output (success/error/warning) // eye candy?
- [ ] Table formatting (ascii tables)
- [ ] Batch operations (JSON manifest)      // ???
- [ ] Clone objects between instances       // import/export is not enough?
- [ ] Watch mode (auto-sync on file change) // ???
- [ ] Shell completion (bash, zsh, fish)    // seems cool

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for development guidelines.

## License

Same as parent project (see [LICENSE](../LICENSE)).

## Support

- Documentation: [docs/](../docs/)
- Issues: https://github.com/echoes1971/r-prj-ng/issues
- Discussions: https://github.com/echoes1971/r-prj-ng/discussions

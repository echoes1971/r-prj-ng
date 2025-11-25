# Test Data Files

This directory contains test files for DBFile testing.

## Generating Test Files

Run the generator to create all test files:

```bash
cd be/dblayer
go run testdata/generate_test_files.go
```

## Generated Files

### Images (`testdata/images/`)
- `test_image.jpg` - 200x150 JPEG with red-yellow gradient
- `test_image.png` - 300x200 PNG with blue gradient and transparency
- `test_image.gif` - 150x100 GIF with green checkerboard pattern
- `small_image.jpg` - 50x50 JPEG (too small for thumbnail generation)

### Files (`testdata/files/`)
- `test_document.txt` - Plain text file for testing MIME detection
- `test_document.pdf` - Minimal valid PDF document

### Expected Thumbnails (`testdata/expected/`)
Generated automatically during tests for verification.

## Usage in Tests

```go
func TestDBFileUpload(t *testing.T) {
    testImagePath := "testdata/images/test_image.jpg"
    // ... test code
}
```

## Notes

- All files are generated programmatically - no external dependencies
- Images use simple gradients for predictable checksums
- PDF is minimal but valid according to PDF 1.4 spec
- Small image (50x50) tests the "no thumbnail needed" code path

// This program generates test files for DBFile testing.
// Run with: go run generate_test_files.go
package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func main() {
	// Generate JPEG image (200x150, red gradient)
	generateJPEG()

	// Generate PNG image (300x200, blue gradient with transparency)
	generatePNG()

	// Generate GIF image (150x100, green gradient)
	generateGIF()

	// Generate small image (50x50, shouldn't need thumbnail)
	generateSmallImage()

	// Generate text file
	generateTextFile()

	// Generate PDF file
	generatePDF()

	log.Println("All test files generated successfully!")
}

func generateJPEG() {
	width, height := 200, 150
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Red gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			img.Set(x, y, color.RGBA{R: r, G: g, B: 0, A: 255})
		}
	}

	f, err := os.Create("images/test_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 85}); err != nil {
		log.Fatal(err)
	}
	log.Println("Created test_image.jpg (200x150, red-yellow gradient)")
}

func generatePNG() {
	width, height := 300, 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Blue gradient with transparency
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b := uint8((x * 255) / width)
			a := uint8((y * 255) / height)
			img.Set(x, y, color.RGBA{R: 0, G: 50, B: b, A: a})
		}
	}

	f, err := os.Create("images/test_image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
	log.Println("Created test_image.png (300x200, blue gradient with transparency)")
}

func generateGIF() {
	width, height := 150, 100
	img := image.NewPaletted(image.Rect(0, 0, width, height), color.Palette{
		color.RGBA{0, 0, 0, 255},       // Black
		color.RGBA{0, 255, 0, 255},     // Green
		color.RGBA{0, 128, 0, 255},     // Dark green
		color.RGBA{128, 255, 128, 255}, // Light green
	})

	// Green pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := ((x / 10) + (y / 10)) % 4
			img.SetColorIndex(x, y, uint8(idx))
		}
	}

	f, err := os.Create("images/test_image.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := gif.Encode(f, img, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Created test_image.gif (150x100, green checkerboard)")
}

func generateSmallImage() {
	width, height := 50, 50
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Simple purple square
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 128, G: 0, B: 128, A: 255})
		}
	}

	f, err := os.Create("images/small_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 85}); err != nil {
		log.Fatal(err)
	}
	log.Println("Created small_image.jpg (50x50, purple square - no thumbnail needed)")
}

func generateTextFile() {
	content := `This is a test text file.
It contains multiple lines.
Used for testing file upload and MIME type detection.

SHA1 checksum should be calculated correctly.
MIME type should be: text/plain
`

	f, err := os.Create("files/test_document.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		log.Fatal(err)
	}
	log.Println("Created test_document.txt")
}

func generatePDF() {
	// Minimal valid PDF structure
	content := `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj

2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj

3 0 obj
<<
/Type /Page
/Parent 2 0 R
/Resources <<
/Font <<
/F1 <<
/Type /Font
/Subtype /Type1
/BaseFont /Helvetica
>>
>>
>>
/MediaBox [0 0 612 792]
/Contents 4 0 R
>>
endobj

4 0 obj
<<
/Length 44
>>
stream
BT
/F1 24 Tf
100 700 Td
(Test PDF Document) Tj
ET
endstream
endobj

xref
0 5
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000317 00000 n 
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
410
%%EOF
`

	f, err := os.Create("files/test_document.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		log.Fatal(err)
	}
	log.Println("Created test_document.pdf (minimal valid PDF with text)")
}

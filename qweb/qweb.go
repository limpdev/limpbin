package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

func main() {
	// Step 1: Read the Markdown file
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <markdown_file> <css_file>")
		return
	}
	markdownFilePath := os.Args[1]
	cssFilePath := os.Args[2]

	/// Step 2: Read the Markdown file
	markdownContent, err := os.ReadFile(markdownFilePath)
	if err != nil {
		fmt.Printf("Error reading the markdown file %s: %v\n", markdownFilePath, err)
		return
	}

	// Step 3: Convert Markdown to HTML
	html := blackfriday.Run(markdownContent)

	// Step 4: Read the CSS file
	cssContent, err := os.ReadFile(cssFilePath)
	if err != nil {
		fmt.Printf("Error reading the CSS file %s: %v\n", cssFilePath, err)
		return
	}
	// Step 5: Create a complete HTML document with embedded CSS
	htmlDocument := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Markdown Documentation</title>
<link href="prism.css" rel="stylesheet" />
<style>%s</style>
</head>
<body>
<script src="prism.js"></script>
%s
</body>
</html>
`, cssContent, html)

	// Step 6: Write the HTML to a file (optional)
	mdFileName := filepath.Base(markdownFilePath)
	outputFilePath := filepath.Join(mdFileName) + ".html"
	err = os.WriteFile(outputFilePath, []byte(htmlDocument), 0644)
	if err != nil {
		fmt.Printf("Error writing the output HTML file %s: %v\n", outputFilePath, err)
		return
	}

	fmt.Printf("HTML document generated successfully: %s\n", outputFilePath)
}

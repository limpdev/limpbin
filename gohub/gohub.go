package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

func main() {
	// Step 1: Read the Markdown file
	if len(os.Args) != 2 {
		fmt.Println("Usage: gohub <markdown_file>")
		return
	}
	markdownFilePath := os.Args[1]

	/// Step 2: Read the Markdown file
	markdownContent, err := os.ReadFile(markdownFilePath)
	if err != nil {
		fmt.Printf("Error reading the markdown file %s: %v\n", markdownFilePath, err)
		return
	}

	// Step 3: Convert Markdown to HTML
	html := blackfriday.Run(markdownContent)

	cssFile, err := os.ReadFile("C:\\Users\\drewg\\proj\\Github\\limpbin\\gohub\\GithubAPI.css")
	if err != nil {
		fmt.Printf("Error reading the CSS file: %v\n", err)
		return
	}

	// Step 4: Create a complete HTML document with embedded CSS and JS
	htmlDocument := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Markdown Documentation</title>
<style type="text/css">%s</style>
</head>
<body>
<article class="markdown-body">
%s
</article>
</body>
</html>
`, cssFile, html)

	// Step 5: Write the HTML to a file (optional)
	mdFileName := filepath.Base(markdownFilePath)
	outputFilePath := filepath.Join(mdFileName) + ".html"
	err = os.WriteFile(outputFilePath, []byte(htmlDocument), 0644)
	if err != nil {
		fmt.Printf("Error writing the output HTML file %s: %v\n", outputFilePath, err)
		return
	}

	fmt.Printf("HTML document generated successfully: %s\n", outputFilePath)
}

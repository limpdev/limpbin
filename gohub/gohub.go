package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gohub <mdFile>")
		return
	}

	mdFile := os.Args[1]

	// Read markdown file
	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		fmt.Printf("Could not read markdown file: %s\n", err)
		return
	}

	// Generate the filename
	mdFileName := filepath.Base(mdFile)
	mdNude := mdFileName[:len(mdFileName)-len(filepath.Ext(mdFileName))]
	outFile := mdNude + ".html"
	// Convert markdown to HTML
	htmlContent := blackfriday.Run(mdContent)

	// Read the template HTML file
	templatePath := "index.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Could not read template file %s: %s\n", templatePath, err)
		return
	}

	// Replace a placeholder in the template with the converted HTML content
	// Assuming index.html contains a {{CONTENT}} placeholder
	htmlPre := strings.Replace(string(templateContent), "{{TITLE}}", string(mdNude), 1)
	htmlDoc := strings.Replace(string(htmlPre), "{{CONTENT}}", string(htmlContent), 1)

	// Write the final HTML to the output file
	err = os.WriteFile(outFile, []byte(htmlDoc), 0644)
	if err != nil {
		fmt.Printf("Error writing the output HTML file %s: %v\n", outFile, err)
		return
	}

	fmt.Printf("HTML document generated successfully: %s\n", outFile)
}

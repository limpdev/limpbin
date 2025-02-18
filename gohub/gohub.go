package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gohub <mdFile>")
		return
	}

	mdFile := os.Args[1]
	cssGFM := "https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/GithubAPI.css"

	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		fmt.Printf("Could not read: %s", err)
		return
	}

	htmlContent := blackfriday.Run(mdContent)

	htmlDoc := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<title>Markdown Documentation</title>
	<link rel="stylesheet" type="text/css" href="%s">
	</head>
	<body>
	<article class="markdown-body">
	%s
	</article>
	</body>
	</html>`, cssGFM, htmlContent)

	mdFileName := filepath.Base(mdFile)
	mdNude := mdFileName[:len(mdFileName)-len(filepath.Ext(mdFileName))]
	outFile := mdNude + ".html"
	err = os.WriteFile(outFile, []byte(htmlDoc), 0644)
	if err != nil {
		fmt.Printf("Error writing the output HTML file %s: %v\n", outFile, err)
		return
	}

	fmt.Printf("HTML document generated successfully: %s\n", outFile)
}

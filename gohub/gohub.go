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

	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		fmt.Printf("Could not read: %s", err)
		return
	}

	htmlContent := blackfriday.Run(mdContent)
	mdFileName := filepath.Base(mdFile)
	mdNude := mdFileName[:len(mdFileName)-len(filepath.Ext(mdFileName))]

	htmlDoc := fmt.Sprintf(`
<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1" />
<html lang="en">
<head>
<meta charset="UTF-8">
<title>%s</title>
<link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/GithubAPI.css"/>
<link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/prismHL.css"/>
<link rel="stylesheet"type="text/css"href="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/clipb.css"/>
</head>
<body>
<script src="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/prismHL.js"></script>
<article class="markdown-body">
%s
</article>
<script src="https://cdn.jsdelivr.net/gh/limpdev/limpbin@main/css/clipb.js"></script>
</body>
</html>`, mdNude, htmlContent)

	outFile := mdNude + ".html"
	err = os.WriteFile(outFile, []byte(htmlDoc), 0644)
	if err != nil {
		fmt.Printf("Error writing the output HTML file %s: %v\n", outFile, err)
		return
	}

	fmt.Printf("HTML document generated successfully: %s\n", outFile)
}

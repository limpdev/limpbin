This file is a merged representation of the entire codebase, combined into a single document by Repomix.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded

## Additional Info

# Directory Structure
```
go.mod
golly.go
```

# Files

## File: go.mod
```
module web/golly

go 1.23.5

require (
	github.com/JohannesKaufmann/html-to-markdown v1.6.0
	github.com/gocolly/colly/v2 v2.1.0
)

require gopkg.in/yaml.v2 v2.4.0 // indirect

require (
	github.com/PuerkitoBio/goquery v1.9.2 // direct
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/antchfx/htmlquery v1.2.3 // indirect
	github.com/antchfx/xmlquery v1.2.4 // indirect
	github.com/antchfx/xpath v1.1.8 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/temoto/robotstxt v1.1.1 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
)
```

## File: golly.go
```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type PageContent struct {
	URL     string
	Title   string
	Content string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: golly <startURL>")
	}
	startURL := os.Args[1]

	// Create output directory
	outputDir := "gollydocs"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Compile the regular expression pattern for filtering URLs
	re, err := regexp.Compile(startURL)
	if err != nil {
		log.Fatalf("Failed to compile regex pattern: %v", err)
	}

	// Initialize HTML to Markdown converter
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	converter.Use(plugin.Table())
	converter.Use(plugin.Strikethrough(""))
	converter.Use(plugin.TaskListItems())

	// Initialize Colly collector
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
	)
	// Set limits for parallel scraping
	err = collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
	})
	if err != nil {
		log.Fatalf("Failed to set limit rule: %v", err)
	}

	pagesChan := make(chan PageContent, 100)

	// Handle the actual page content
	collector.OnResponse(func(r *colly.Response) {
		log.Printf("Processing page: %s", r.Request.URL.String())

		contentType := r.Headers.Get("Content-Type")
		log.Printf("Content-Type: %s", contentType)

		if strings.Contains(contentType, "text/html") {
			// Create a goquery document
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
			if err != nil {
				log.Printf("Error parsing HTML from %s: %v", r.Request.URL.String(), err)
				return
			}

			// Extract title
			title := doc.Find("title").First().Text()
			if title == "" {
				title = r.Request.URL.Path
			}
			log.Printf("Found title: %s", title)

			// Convert relative URLs to absolute
			doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
				if href, exists := s.Attr("href"); exists {
					absURL := r.Request.AbsoluteURL(href)
					if absURL != "" {
						s.SetAttr("href", absURL)
					}
				}
			})

			// Get the main content
			var content string
			mainContent := doc.Find("main, article, .content, .documentation, #content").First()
			if mainContent.Length() > 0 {
				html, err := mainContent.Html()
				if err != nil {
					log.Printf("Error getting HTML from %s: %v", r.Request.URL.String(), err)
					return
				}
				content = html
				log.Printf("Found main content section")
			} else {
				content, err = doc.Find("body").Html()
				if err != nil {
					log.Printf("Error getting HTML from %s: %v", r.Request.URL.String(), err)
					return
				}
				log.Printf("Using body content as fallback")
			}

			// Convert HTML to Markdown
			markdown, err := converter.ConvertString(content)
			if err != nil {
				log.Printf("Error converting HTML to Markdown for %s: %v", r.Request.URL.String(), err)
				return
			}

			// Add front matter
			frontMatter := fmt.Sprintf(`---
title: "%s"
source_url: "%s"
---

`, title, r.Request.URL.String())

			markdown = frontMatter + markdown

			pagesChan <- PageContent{
				URL:     r.Request.URL.String(),
				Title:   title,
				Content: markdown,
			}
			log.Printf("Successfully processed and queued page: %s", r.Request.URL.String())
		}
	})

	// Register link finder
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			absURL := e.Request.AbsoluteURL(link)
			if absURL != "" {
				log.Printf("Found link: %s", absURL)
				if re.MatchString(absURL) {
					log.Printf("Link matches pattern, queuing: %s", absURL)
					err := collector.Visit(absURL)
					if err != nil {
						log.Printf("Error visiting %s: %v", absURL, err)
					}
				}
			}
		}
	})

	// Start scraping
	log.Printf("Starting initial visit to %s", startURL)
	err = collector.Visit(startURL)
	if err != nil {
		log.Printf("Error visiting start URL %s: %v", startURL, err)
	}

	// Start a goroutine to close the channel when all scraping is done
	go func() {
		collector.Wait()
		log.Println("Scraping completed, closing channel")
		close(pagesChan)
	}()

	// Process and save pages
	for page := range pagesChan {
		// Determine file path
		fileName := sanitizeFilename(page.Title) + ".md"
		filePath := filepath.Join(outputDir, fileName)

		// Check for filename conflicts and make filename unique
		uniqueFilePath := getUniqueFilePath(filePath)

		// Write content to file
		err = os.WriteFile(uniqueFilePath, []byte(page.Content), 0644)
		if err != nil {
			log.Printf("Error writing file %s: %v", uniqueFilePath, err)
		} else {
			log.Printf("Saved page to %s", uniqueFilePath)
		}
	}

	log.Println("All pages have been saved.")
}

func sanitizeFilename(filename string) string {
	// Remove illegal characters
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]+`)
	sanitized := invalidChars.ReplaceAllString(filename, "_")
	// Trim spaces
	sanitized = strings.TrimSpace(sanitized)
	// If empty, use a default name
	if sanitized == "" {
		sanitized = "page"
	}
	return sanitized
}

func getUniqueFilePath(filePath string) string {
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filepath.Base(filePath), ext)
	dir := filepath.Dir(filePath)
	counter := 1
	uniquePath := filePath
	for {
		if _, err := os.Stat(uniquePath); os.IsNotExist(err) {
			break
		}
		uniqueName := fmt.Sprintf("%s_%d", name, counter)
		uniquePath = filepath.Join(dir, uniqueName+ext)
		counter++
	}
	return uniquePath
}
```

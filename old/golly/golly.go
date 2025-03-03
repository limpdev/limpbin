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

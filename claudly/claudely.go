package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

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

type Config struct {
	StartURL     string
	OutputDir    string
	Parallelism  int
	MaxDepth     int
	IncludeRegex string
	ExcludeRegex string
	Verbose      bool
	Timeout      time.Duration
}

func main() {
	// Parse command line arguments
	config := parseFlags()

	// Create output directory
	err := os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Compile the regular expressions
	includeRe, err := regexp.Compile(config.IncludeRegex)
	if err != nil {
		log.Fatalf("Failed to compile include regex pattern: %v", err)
	}

	var excludeRe *regexp.Regexp
	if config.ExcludeRegex != "" {
		excludeRe, err = regexp.Compile(config.ExcludeRegex)
		if err != nil {
			log.Fatalf("Failed to compile exclude regex pattern: %v", err)
		}
	}

	// Initialize HTML to Markdown converter
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	converter.Use(plugin.Table())
	converter.Use(plugin.Strikethrough(""))
	converter.Use(plugin.TaskListItems())

	// Create a map to track visited URLs
	visitedURLs := make(map[string]bool)
	var visitedMutex sync.Mutex

	// Initialize Colly collector
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.MaxDepth(config.MaxDepth),
	)

	// Set timeout
	collector.SetRequestTimeout(config.Timeout)

	// Set limits for parallel scraping
	err = collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.Parallelism,
	})
	if err != nil {
		log.Fatalf("Failed to set limit rule: %v", err)
	}

	pagesChan := make(chan PageContent, 100)

	// Create an index file to track all scraped pages
	indexFile, err := os.Create(filepath.Join(config.OutputDir, "index.md"))
	if err != nil {
		log.Printf("Warning: Failed to create index file: %v", err)
	} else {
		defer indexFile.Close()
		indexFile.WriteString("# Index of Scraped Pages\n\n")
	}

	var indexMutex sync.Mutex

	// Handle the actual page content
	collector.OnResponse(func(r *colly.Response) {
		if config.Verbose {
			log.Printf("Processing page: %s", r.Request.URL.String())
		}

		contentType := r.Headers.Get("Content-Type")
		if config.Verbose {
			log.Printf("Content-Type: %s", contentType)
		}

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
			if config.Verbose {
				log.Printf("Found title: %s", title)
			}

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

			// Try to find the main content using different common selectors
			mainContent := doc.Find("main, article, .content, .documentation, #content, [role=main], .main-content").First()
			if mainContent.Length() > 0 {
				html, err := mainContent.Html()
				if err != nil {
					log.Printf("Error getting HTML from %s: %v", r.Request.URL.String(), err)
					return
				}
				content = html
				if config.Verbose {
					log.Printf("Found main content section")
				}
			} else {
				// Fallback to body content
				content, err = doc.Find("body").Html()
				if err != nil {
					log.Printf("Error getting HTML from %s: %v", r.Request.URL.String(), err)
					return
				}
				if config.Verbose {
					log.Printf("Using body content as fallback")
				}
			}

			// Convert HTML to Markdown
			markdown, err := converter.ConvertString(content)
			if err != nil {
				log.Printf("Error converting HTML to Markdown for %s: %v", r.Request.URL.String(), err)
				return
			}

			// Clean up the markdown
			markdown = cleanMarkdown(markdown)

			// Add front matter
			frontMatter := fmt.Sprintf(`---
title: "%s"
source_url: "%s"
date_scraped: "%s"
---

`, title, r.Request.URL.String(), time.Now().Format(time.RFC3339))

			markdown = frontMatter + markdown

			pagesChan <- PageContent{
				URL:     r.Request.URL.String(),
				Title:   title,
				Content: markdown,
			}
			if config.Verbose {
				log.Printf("Successfully processed and queued page: %s", r.Request.URL.String())
			}
		} else if config.Verbose {
			log.Printf("Skipping non-HTML content at %s", r.Request.URL.String())
		}
	})

	// Register link finder
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "" {
			absURL := e.Request.AbsoluteURL(link)
			if absURL != "" {
				// Parse the URL to check if it's the same domain
				_, err := url.Parse(absURL)
				if err != nil {
					log.Printf("Error parsing URL %s: %v", absURL, err)
					return
				}

				_, err = url.Parse(config.StartURL)
				if err != nil {
					log.Printf("Error parsing base URL %s: %v", config.StartURL, err)
					return
				}

				// Skip if URL has already been visited
				visitedMutex.Lock()
				visited := visitedURLs[absURL]
				visitedMutex.Unlock()

				if visited {
					return
				}

				// Skip URLs if they don't match the include pattern
				if !includeRe.MatchString(absURL) {
					return
				}

				// Skip URLs if they match the exclude pattern
				if excludeRe != nil && excludeRe.MatchString(absURL) {
					return
				}

				// Mark as visited
				visitedMutex.Lock()
				visitedURLs[absURL] = true
				visitedMutex.Unlock()

				if config.Verbose {
					log.Printf("Queuing: %s", absURL)
				}

				err = collector.Visit(absURL)
				if err != nil {
					log.Printf("Error visiting %s: %v", absURL, err)
				}
			}
		}
	})

	// Start scraping
	log.Printf("Starting initial visit to %s", config.StartURL)
	err = collector.Visit(config.StartURL)
	if err != nil {
		log.Printf("Error visiting start URL %s: %v", config.StartURL, err)
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
		filePath := filepath.Join(config.OutputDir, fileName)

		// Check for filename conflicts and make filename unique
		uniqueFilePath := getUniqueFilePath(filePath)

		// Get relative path for the index
		relPath, err := filepath.Rel(config.OutputDir, uniqueFilePath)
		if err != nil {
			relPath = uniqueFilePath
		}

		// Write content to file
		err = os.WriteFile(uniqueFilePath, []byte(page.Content), 0644)
		if err != nil {
			log.Printf("Error writing file %s: %v", uniqueFilePath, err)
		} else {
			log.Printf("Saved page to %s", uniqueFilePath)

			// Update the index file
			if indexFile != nil {
				indexMutex.Lock()
				_, err = indexFile.WriteString(fmt.Sprintf("- [%s](%s) - [Source](%s)\n",
					page.Title,
					url.PathEscape(relPath),
					page.URL))
				if err != nil {
					log.Printf("Warning: Failed to update index: %v", err)
				}
				indexMutex.Unlock()
			}
		}
	}

	log.Println("All pages have been saved.")
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.StartURL, "url", "", "Starting URL to scrape (required)")
	flag.StringVar(&cfg.OutputDir, "output", "gollydocs", "Output directory for scraped content")
	flag.IntVar(&cfg.Parallelism, "parallel", 4, "Number of parallel scrapers")
	flag.IntVar(&cfg.MaxDepth, "depth", 3, "Maximum crawl depth (0 for unlimited)")
	flag.StringVar(&cfg.IncludeRegex, "include", "", "Regex pattern for URLs to include (defaults to start URL domain)")
	flag.StringVar(&cfg.ExcludeRegex, "exclude", "", "Regex pattern for URLs to exclude")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")
	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "Request timeout")

	flag.Parse()

	// If no URL is provided, check positional arguments or show usage
	if cfg.StartURL == "" && len(flag.Args()) > 0 {
		cfg.StartURL = flag.Args()[0]
	}

	if cfg.StartURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	// If no include pattern is specified, default to the domain of the start URL
	if cfg.IncludeRegex == "" {
		parsedURL, err := url.Parse(cfg.StartURL)
		if err != nil {
			log.Fatalf("Failed to parse start URL: %v", err)
		}
		cfg.IncludeRegex = fmt.Sprintf("^https?://%s", regexp.QuoteMeta(parsedURL.Hostname()))
		log.Printf("Using default include pattern: %s", cfg.IncludeRegex)
	}

	return cfg
}

func sanitizeFilename(filename string) string {
	// Remove illegal characters
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]+`)
	sanitized := invalidChars.ReplaceAllString(filename, "_")
	// Trim spaces
	sanitized = strings.TrimSpace(sanitized)
	// Limit filename length
	if len(sanitized) > 100 {
		sanitized = sanitized[:100]
	}
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

func cleanMarkdown(markdown string) string {
	// Remove excessive newlines (more than 2 in a row)
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")

	// Fix code blocks with missing language
	re = regexp.MustCompile("```\n")
	markdown = re.ReplaceAllString(markdown, "```text\n")

	return markdown
}

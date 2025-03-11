package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/ncruces/zenity"
)

func main() {
	url := promptForURL()
	if url == "" {
		return // User canceled or error occurred
	}

	// Generate a filename based on the URL
	filename := generateFilename(url)

	// Convert the URL and save to the downloads folder
	convertURL(url, filename)
}

func promptForURL() string {
	text, err := zenity.Entry("Enter URL",
		zenity.Title("Moka"),
		zenity.Width(400),
		zenity.Height(50),
		zenity.OKLabel("Rip"),
		zenity.CancelLabel("Cancel"),
		zenity.WindowIcon("icon.png"),
		zenity.Icon("icon.png"),
	)
	if err != nil {
		log.Fatal("Failed to show entry dialog:", err)
	}
	return text
}

//func isValidURL(text string) bool {
//	// Basic URL validation
//	return strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://")
//}

func generateFilename(url string) string {
	// Get user's download directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not find home directory:", err)
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")

	// Create a filename using the hostname from the URL and current timestamp
	hostname := ""
	if startIdx := strings.Index(url, "://"); startIdx != -1 {
		trimmedURL := url[startIdx+3:]
		if endIdx := strings.Index(trimmedURL, "/"); endIdx != -1 {
			hostname = trimmedURL[:endIdx]
		} else {
			hostname = trimmedURL
		}
	}

	// Clean the hostname to make it a valid filename
	hostname = strings.ReplaceAll(hostname, ".", "-")
	hostname = strings.ReplaceAll(hostname, ":", "-")

	// Add timestamp for uniqueness
	timestamp := time.Now().Format("20060102-150405")

	return filepath.Join(downloadsDir, hostname+"-"+timestamp+".md")
}

func convertURL(url string, filename string) {
	if url == "" {
		zenity.Error("No URL provided")
		return
	}

	// Show progress dialog
	progress, err := zenity.Progress(
		zenity.Title("Moka - Converting"),
		zenity.Width(300),
	)
	if err != nil {
		log.Fatal("Failed to show progress dialog:", err)
	}

	progress.Text("Downloading content...")
	progress.Value(10)

	// Fetch the URL content
	cmdCurl := exec.Command("curl", "--no-progress-meter", url)
	output, err := cmdCurl.CombinedOutput()
	if err != nil {
		progress.Close()
		zenity.Error("Error during the CURL process: " + err.Error())
		return
	}

	progress.Text("Converting to Markdown...")
	progress.Value(50)

	// Convert HTML to Markdown
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(
				commonmark.WithStrongDelimiter("**"),
				commonmark.WithBulletListMarker("-"),
				commonmark.WithCodeBlockFence("```"),
				commonmark.WithHeadingStyle(commonmark.HeadingStyleATX),
				commonmark.WithLinkEmptyContentBehavior(commonmark.LinkBehaviorSkip),
				commonmark.WithLinkEmptyHrefBehavior(commonmark.LinkBehaviorSkip),
			),
			table.NewTablePlugin(
				table.WithHeaderPromotion(true),
				table.WithSkipEmptyRows(true),
				table.WithSpanCellBehavior(table.SpanBehaviorEmpty),
			),
		),
	)

	markdown, err := conv.ConvertString(string(output))
	if err != nil {
		progress.Close()
		zenity.Error("Error converting HTML to Markdown: " + err.Error())
		return
	}

	progress.Text("Saving file...")
	progress.Value(90)

	// Ensure the Downloads directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		progress.Close()
		zenity.Error("Error creating directory: " + err.Error())
		return
	}

	// Write to file
	err = os.WriteFile(filename, []byte(markdown), 0644)
	if err != nil {
		progress.Close()
		zenity.Error("Error writing file: " + err.Error())
		return
	}

	progress.Value(100)
	progress.Close()

	// Show success message with file location
	zenity.Info("Conversion completed. File saved to:\n" + filename)
}

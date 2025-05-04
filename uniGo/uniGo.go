// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/unicode/runenames" // For character names
)

// CharacterInfo holds data about a single Unicode character
type CharacterInfo struct {
	Char       string `json:"char"`
	CodePoint  string `json:"codePoint"` // Hex representation like U+0041
	Name       string `json:"name"`
	Category   string `json:"category"`
	CategoryAb string `json:"categoryAb"` // Abbreviation like Lu
	// BlockName string `json:"blockName"` // Block info is harder to get reliably without external data
}

// APIResponse structures the JSON response for the characters endpoint
type APIResponse struct {
	Characters   []CharacterInfo `json:"characters"`
	TotalItems   int             `json:"totalItems"`
	CurrentPage  int             `json:"currentPage"`
	ItemsPerPage int             `json:"itemsPerPage"`
	TotalPages   int             `json:"totalPages"`
}

// MetadataResponse structures the JSON response for metadata
type MetadataResponse struct {
	Categories map[string]string `json:"categories"` // Map Abbreviation -> Full Name
	Blocks     []string          `json:"blocks"`
}

var (
	allCharacters []CharacterInfo
	categories    map[string]string // Map Abbreviation -> Full Name
	dataMutex     sync.RWMutex
)

// Map of Unicode categories with their full names
var categoryNames = map[string]string{
	"Lu": "Uppercase Letter",
	"Ll": "Lowercase Letter",
	"Lt": "Titlecase Letter",
	"Lm": "Modifier Letter",
	"Lo": "Other Letter",
	"Mn": "Nonspacing Mark",
	"Mc": "Spacing Mark",
	"Me": "Enclosing Mark",
	"Nd": "Decimal Number",
	"Nl": "Letter Number",
	"No": "Other Number",
	"Pc": "Connector Punctuation",
	"Pd": "Dash Punctuation",
	"Ps": "Open Punctuation",
	"Pe": "Close Punctuation",
	"Pi": "Initial Punctuation",
	"Pf": "Final Punctuation",
	"Po": "Other Punctuation",
	"Sm": "Math Symbol",
	"Sc": "Currency Symbol",
	"Sk": "Modifier Symbol",
	"So": "Other Symbol",
	"Zs": "Space Separator",
	"Zl": "Line Separator",
	"Zp": "Paragraph Separator",
	"Cc": "Control",
	"Cf": "Format",
	"Cs": "Surrogate",
	"Co": "Private Use",
	"Cn": "Unassigned",
}

// loadUnicodeData pre-populates the character list
func loadUnicodeData() {
	log.Println("Loading Unicode data...")
	dataMutex.Lock()
	defer dataMutex.Unlock()

	allCharacters = []CharacterInfo{}
	categories = make(map[string]string)
	addedCategories := make(map[string]bool)

	// Iterate through a relevant range (e.g., BMP 0x0000 to 0xFFFF)
	for r := rune(0); r <= 0xFFFF; r++ {
		if !unicode.IsPrint(r) || unicode.IsControl(r) || (unicode.IsSpace(r) && r != ' ') {
			// Skip non-printable, control chars (except space)
			// Add more exclusion logic if needed (e.g., surrogates)
			if r >= 0xD800 && r <= 0xDFFF { // Skip surrogate pairs
				continue
			}

			// Skip private use area for general browsing
			if r >= 0xE000 && r <= 0xF8FF {
				continue
			}

			// Skip combining marks unless you specifically want them displayed standalone
			// if unicode.In(r, unicode.Mn, unicode.Me, unicode.Mc) {
			// continue // Uncomment to skip combining marks
			// }

			if r == '\uFFFD' { // Skip replacement character
				continue
			}

			// Skip characters known to cause issues or be unrenderable in many contexts
			if r == 0xAD { // Soft hyphen
				continue
			}
			if r >= 0x2060 && r <= 0x206F { // General punctuation invisible operators
				continue
			}
			if r >= 0xFFF9 && r <= 0xFFFB { // Interlinear annotation anchors etc
				continue
			}
			continue // Default skip if not printable or otherwise undesirable
		}

		name := runenames.Name(r)
		if name == "" || strings.Contains(name, "<") { // Skip reserved/private use/control names
			continue
		}

		// Get character category
		catAb := getCategoryAbbreviation(r)
		catName := categoryNames[catAb]
		if catName == "" {
			catName = "Unknown Category"
		}

		info := CharacterInfo{
			Char:       string(r),
			CodePoint:  fmt.Sprintf("U+%04X", r),
			Name:       name,
			Category:   catName,
			CategoryAb: catAb,
		}
		allCharacters = append(allCharacters, info)

		// Collect unique categories
		if !addedCategories[catAb] {
			categories[catAb] = catName
			addedCategories[catAb] = true
		}
	}

	log.Printf("Loaded %d characters and %d categories.", len(allCharacters), len(categories))
}

// getCategoryAbbreviation returns the two-letter Unicode category for a rune
func getCategoryAbbreviation(r rune) string {
	// Check for specific category ranges
	switch {
	case unicode.IsLetter(r):
		if unicode.IsUpper(r) {
			return "Lu"
		} else if unicode.IsLower(r) {
			return "Ll"
		} else if unicode.IsTitle(r) {
			return "Lt"
		} else {
			// Further differentiate between Lm and Lo if needed
			return "Lo" // Default to "Other Letter"
		}
	case unicode.IsDigit(r):
		return "Nd"
	case unicode.IsPunct(r):
		// This is simplified - ideally you'd distinguish between different punctuation types
		return "Po"
	case unicode.IsSymbol(r):
		// This is simplified - ideally you'd distinguish between different symbol types
		if strings.Contains(runenames.Name(r), "CURRENCY") {
			return "Sc"
		}
		return "So"
	case unicode.IsSpace(r):
		return "Zs"
	case unicode.IsControl(r):
		return "Cc"
	default:
		return "Cn" // Unassigned as fallback
	}
}

// handleCharacters serves the character data based on query parameters
func handleCharacters(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock() // Use read lock for concurrent reads
	defer dataMutex.RUnlock()

	query := r.URL.Query()
	search := strings.ToLower(strings.TrimSpace(query.Get("search")))
	categoryFilter := query.Get("category") // Expecting Category Abbreviation (e.g., "Lu")

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100 // Default limit
	}

	// Filter characters
	filtered := make([]CharacterInfo, 0, len(allCharacters))
	for _, charInfo := range allCharacters {
		match := true

		// Search filter (checks name, codepoint)
		if search != "" {
			nameLower := strings.ToLower(charInfo.Name)
			codeLower := strings.ToLower(charInfo.CodePoint)
			// Basic substring search, could be improved (e.g., word boundary)
			if !strings.Contains(nameLower, search) && !strings.Contains(codeLower, search) && charInfo.Char != search {
				match = false
			}
		}

		// Category filter
		if match && categoryFilter != "" && charInfo.CategoryAb != categoryFilter {
			match = false
		}

		if match {
			filtered = append(filtered, charInfo)
		}
	}

	// Apply pagination
	totalItems := len(filtered)
	totalPages := (totalItems + limit - 1) / limit
	if page > totalPages && totalPages > 0 {
		page = totalPages // Adjust page if it's out of bounds
	}

	start := (page - 1) * limit
	end := start + limit
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedChars := []CharacterInfo{}
	if start < end { // Ensure indices are valid
		paginatedChars = filtered[start:end]
	}

	// Prepare response
	resp := APIResponse{
		Characters:   paginatedChars,
		TotalItems:   totalItems,
		CurrentPage:  page,
		ItemsPerPage: limit,
		TotalPages:   totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// handleMetadata serves category (and potentially block) lists
func handleMetadata(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	// Sort category names alphabetically for the dropdown
	sortedCategories := make(map[string]string)
	catKeys := make([]string, 0, len(categories))
	for k := range categories {
		catKeys = append(catKeys, k)
	}
	sort.Slice(catKeys, func(i, j int) bool {
		// Sort by full name for user-friendliness
		return categories[catKeys[i]] < categories[catKeys[j]]
	})
	for _, k := range catKeys {
		sortedCategories[k] = categories[k]
	}

	resp := MetadataResponse{
		Categories: sortedCategories,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding metadata JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// serveHTML serves the index.html file
func serveHTML(w http.ResponseWriter, r *http.Request) {
	// Basic security: Prevent path traversal
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "uniGo.html") // Assume index.html is in the same directory
}

func main() {
	// Load data once on startup
	loadUnicodeData()

	// --- HTTP Handlers ---
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/api/characters", handleCharacters)
	http.HandleFunc("/api/metadata", handleMetadata)

	// --- Start Server ---
	port := "6969"
	log.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

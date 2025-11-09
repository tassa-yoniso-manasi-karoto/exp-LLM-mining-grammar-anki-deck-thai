package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SentenceResult represents a single search result
type SentenceResult struct {
	Thai         string   `json:"thai"`
	English      string   `json:"english"`
	AudioFile    string   `json:"audio_file"`
	ImageFile    string   `json:"image_file"`
	SourceShow   string   `json:"source_show"`
	Episode      string   `json:"episode,omitempty"`
	Timestamp    string   `json:"timestamp,omitempty"`
	OriginalLang string   `json:"original_language"`
	WordCount    int      `json:"word_count"`
	MediaDir     string   `json:"media_dir"`
	Difficulty   int      `json:"difficulty"`
	OriginalRow  []string `json:"-"` // Store full TSV row for export
}

// SearchResults contains all search results and metadata
type SearchResults struct {
	Query       string            `json:"query"`
	Results     []SentenceResult  `json:"results"`
	TotalFound  int              `json:"total_found"`
	Returned    int              `json:"returned"`
}

// Command-line flags
var (
	query        = flag.String("query", "", "Search query (supports substring match)")
	regexMode    = flag.Bool("regex", false, "Enable regex mode for query")
	contains     arrayFlag
	exclude      arrayFlag
	nativeOnly   = flag.Bool("native-only", false, "Only return sentences from native Thai content")
	hasAudio     = flag.Bool("has-audio", false, "Only return sentences with audio files")
	wordCountMin = flag.Int("word-count-min", 0, "Minimum word count (0 = no limit)")
	wordCountMax = flag.Int("word-count-max", 0, "Maximum word count (0 = no limit)")
	limit        = flag.Int("limit", 100, "Maximum number of results to return")
	outputFormat = flag.String("output", "text", "Output format: text, json, csv, anki-csv (default: text)")
	uniqueOnly   = flag.Bool("unique-only", true, "Return only unique sentences")
	verbose        = flag.Bool("v", false, "Enable verbose logging")
	rootPath       = flag.String("root", "/home/voiduser/Videos/_______/", "Root path for TSV files")
	full           = flag.Bool("full", false, "Show both Thai and English text")
	showEnglish    = flag.Bool("english", false, "Show English instead of Thai in text output")
	showDifficulty = flag.Bool("show-difficulty", false, "Show difficulty scores in text output")
	showShow       = flag.Bool("show-show", false, "Show show names in text output")
	exportIndices  = flag.String("export-indices", "", "Export selected indices as TSV file (e.g., '0,3,5,12'). Must be used with query.")
	outputDir      = flag.String("output-dir", "", "Directory for exported TSV files (required when using --export-indices)")
)

// arrayFlag implements flag.Value for string arrays
type arrayFlag []string

func (a *arrayFlag) String() string {
	return strings.Join(*a, ",")
}

func (a *arrayFlag) Set(value string) error {
	*a = append(*a, value)
	return nil
}

// Global variables
var (
	reHTMLTag = regexp.MustCompile(`<[^>]*>`)
	reSeason  = regexp.MustCompile(`.*/(.*)/S[0-9]+`)
	reMovie   = regexp.MustCompile(`(.*?) [12][0-9].*`)
	rePicURL  = regexp.MustCompile(`.*""(.*)"".*`)
)

// Language mapping for shows
var showLangs = map[string]string{
	"Analog Squad":                                "th",
	"Bad Genius":                                  "th",
	"Bad Genius_ The Series":                      "th",
	"Bangkok Traffic (Love) Story":                "th",
	"Crazy Little Thing Called Love":              "th",
	"Emily in Paris":                              "en",
	"First Love":                                  "jp",
	"From Me to You_ Kimi ni Todoke":              "jp",
	"Girl from Nowhere":                           "th",
	"Hormones":                                    "th",
	"La La Land":                                  "en",
	"Let Me Eat Your Pancreas":                    "jp",
	"Love Destiny":                                "th",
	"Love Destiny The Movie":                      "th",
	"Maid":                                        "en",
	"Mon-rak Transistor":                          "th",
	"Money Heist":                                 "en",
	"My Girl (Remastered in 4K)":                  "th",
	"Once Upon a Star":                            "th",
	"Parasite":                                    "ko",
	"Sex Education":                               "en",
	"Spirited Away":                               "jp",
	"Squid Game":                                  "ko",
	"Stranger Things":                             "en",
	"The Billionaire":                             "th",
	"The Crown":                                   "en",
	"The End of the F___ing World":                "en",
	"The Journalist":                              "jp",
	"The Makanai_ Cooking for the Maiko House":    "jp",
	"Twenty Five Twenty One":                      "ko",
	"Uncle Boonmee Who Can Recall His Past Lives": "th",
	"VINLAND SAGA":                                "jp",
	"Wednesday":                                   "en",
	"6ixtynin9 The Series":                        "th",
	"Answer for Heaven":                           "th",
	"Bad Guys":                                    "th",
	"Oh My Ghost":                                 "th",
	"Remember You":                                "th",
	"Sing Again":                                  "th",
	"The Underclass":                              "th",
	"Tunnel":                                      "th",
	"Unlucky Ploy":                                "th",
	"4 Kings":                                     "th",
	"Blue Again":                                  "th",
	"Fast & Feel Love":                            "th",
	"Hunger":                                      "th",
	"Love You My Arrogance":                       "th",
	"The Lost Lotteries":                          "th",
	"The Love of Siam":                            "th",
	"Tootsies & The Fake":                         "th",
	"Bangkok Breaking":                            "th",
	"Bangkok Buddies":                             "th",
	"DELETE":                                      "th",
	"Happy Old Year":                              "th",
	"I Need Romance":                              "th",
	"Let's Fight Ghost":                           "th",
	"Miss Culinary":                               "th",
	"Chainsaw Man":                                "jp",
	"Castle in the Sky":                           "jp",
	"Howl's Moving Castle":                        "jp",
	"The Bride of Naga":                           "th",
	"To the Moon and Back":                        "th",
	"SPY x FAMILY":                                "jp",
	"Princess Mononoke":                           "jp",
}

func init() {
	flag.Var(&contains, "contains", "Term that must be contained (can be used multiple times)")
	flag.Var(&exclude, "exclude", "Term to exclude (can be used multiple times)")
}

func main() {
	flag.Parse()

	// Setup logging
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	// Validate input
	if *query == "" && len(contains) == 0 {
		log.Fatal().Msg("Either --query or --contains must be provided")
	}

	log.Debug().
		Str("query", *query).
		Strs("contains", contains).
		Bool("native-only", *nativeOnly).
		Int("limit", *limit).
		Msg("Starting search")

	// Search for sentences
	results := searchSentences()

	// Sort by difficulty (sentence length) in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Difficulty > results[j].Difficulty
	})

	// Apply limit
	if *limit > 0 && len(results) > *limit {
		results = results[:*limit]
	}

	// Handle export mode if requested
	if *exportIndices != "" {
		indices := parseIndices(*exportIndices)
		outputFile := exportToTSV(results, indices)
		// Print only the filepath and exit
		fmt.Println(outputFile)
		return
	}

	// Create search results object
	searchResults := SearchResults{
		Query:      *query,
		Results:    results,
		TotalFound: len(results),
		Returned:   len(results),
	}

	// Output results
	switch *outputFormat {
	case "text":
		outputText(results)
	case "json":
		outputJSON(searchResults)
	case "csv":
		outputCSV(results)
	case "anki-csv":
		outputAnkiCSV(results)
	default:
		log.Fatal().Str("format", *outputFormat).Msg("Unknown output format")
	}
}

func searchSentences() []SentenceResult {
	var results []SentenceResult
	seen := make(map[string]bool)

	// Compile regex if needed
	var queryRegex *regexp.Regexp
	if *regexMode && *query != "" {
		var err error
		queryRegex, err = regexp.Compile(*query)
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid regex pattern")
		}
	}

	err := filepath.Walk(*rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .media directories
		if info.IsDir() && strings.HasSuffix(info.Name(), ".media") {
			return filepath.SkipDir
		}

		filename := filepath.Base(path)

		// Skip non-TSV files and "3.tsv" files
		if strings.HasSuffix(path, "3.tsv") || !strings.HasSuffix(filename, ".tsv") {
			return nil
		}

		// Extract show name
		show := extractShowName(path, filename)

		// Check language and native-only filter
		lang, found := showLangs[show]
		if !found {
			log.Debug().Str("show", show).Msg("Show has no language attribution")
			// Assume Thai if not specified
			lang = "th"
		}

		// Apply native-only filter
		if *nativeOnly && lang != "th" {
			log.Debug().Str("show", show).Str("lang", lang).Msg("Skipping non-Thai show")
			return nil
		}

		// Skip non-Thai content "2.tsv" files if native-only
		if *nativeOnly && strings.HasSuffix(path, "2.tsv") {
			return nil
		}

		// Determine column indexes based on language
		idx := 4      // Thai text column
		idxTra := 5   // Translation column
		mediaDir := strings.TrimSuffix(path, ".tsv") + ".media/"

		if lang != "th" {
			if strings.HasSuffix(path, "2.tsv") {
				mediaDir = strings.TrimSuffix(path, "2.tsv") + ".media/"
				idx = 10
			} else {
				return nil
			}
		}

		// Process TSV file
		log.Debug().Str("file", path).Msg("Processing file")
		processedResults := processTSVFile(path, show, lang, idx, idxTra, mediaDir, queryRegex, seen)
		results = append(results, processedResults...)

		return nil
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Error walking directory")
	}

	return results
}

func extractShowName(path, filename string) string {
	if reSeason.MatchString(path) {
		return reSeason.FindStringSubmatch(path)[1]
	} else if reMovie.MatchString(filename) {
		return reMovie.FindStringSubmatch(filename)[1]
	}
	return strings.TrimSuffix(filename, ".tsv")
}

func processTSVFile(path, show, lang string, idx, idxTra int, mediaDir string, queryRegex *regexp.Regexp, seen map[string]bool) []SentenceResult {
	var results []SentenceResult

	file, err := os.Open(path)
	if err != nil {
		log.Warn().Err(err).Str("file", path).Msg("Failed to open file")
		return results
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		arr := strings.Split(line, "\t")

		// Check if we have enough columns
		if len(arr) <= idx || len(arr) <= idxTra {
			continue
		}

		thai := arr[idx]
		english := reHTMLTag.ReplaceAllString(arr[idxTra], "")

		// Skip empty translations or music
		if english == "" || strings.Contains(english, "♪") || strings.Contains(thai, "♪") {
			continue
		}

		// Check if sentence matches search criteria
		if !matchesCriteria(thai, english, queryRegex) {
			continue
		}

		// Check uniqueness if required
		if *uniqueOnly {
			if seen[thai] {
				continue
			}
			seen[thai] = true
		}

		// Count words
		wordCount := countThaiWords(thai)

		// Apply word count filter
		if *wordCountMin > 0 && wordCount < *wordCountMin {
			continue
		}
		if *wordCountMax > 0 && wordCount > *wordCountMax {
			continue
		}

		// Check audio file existence if required
		audioFile := arr[0]
		if *hasAudio && audioFile == "" {
			continue
		}

		// Extract image file
		imageFile := ""
		if len(arr) > 3 && arr[3] != "" {
			if matches := rePicURL.FindStringSubmatch(arr[3]); len(matches) > 1 {
				imageFile = matches[1]
			}
		}

		result := SentenceResult{
			Thai:         thai,
			English:      english,
			AudioFile:    audioFile,
			ImageFile:    imageFile,
			SourceShow:   show,
			OriginalLang: lang,
			WordCount:    wordCount,
			MediaDir:     mediaDir,
			Difficulty:   len(thai),
			OriginalRow:  arr, // Store full TSV row
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		log.Warn().Err(err).Str("file", path).Msg("Error reading file")
	}

	return results
}

func matchesCriteria(thai, english string, queryRegex *regexp.Regexp) bool {
	// Check main query
	if *query != "" {
		if *regexMode {
			if !queryRegex.MatchString(thai) {
				return false
			}
		} else {
			if !strings.Contains(thai, *query) {
				return false
			}
		}
	}

	// Check contains criteria
	for _, term := range contains {
		if !strings.Contains(thai, term) {
			return false
		}
	}

	// Check exclude criteria
	for _, term := range exclude {
		if strings.Contains(thai, term) || strings.Contains(english, term) {
			return false
		}
	}

	return true
}

func countThaiWords(text string) int {
	// Simple word count: count Thai character groups separated by spaces
	// This is a simplified approach - proper Thai word segmentation would require a library
	words := 0
	inWord := false

	for _, r := range text {
		if unicode.Is(unicode.Thai, r) {
			if !inWord {
				words++
				inWord = true
			}
		} else if unicode.IsSpace(r) {
			inWord = false
		}
	}

	return words
}

func outputJSON(results SearchResults) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		log.Fatal().Err(err).Msg("Failed to encode JSON")
	}
}

func outputCSV(results []SentenceResult) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header
	header := []string{"Thai", "English", "Audio", "Image", "Show", "Language", "WordCount"}
	if err := writer.Write(header); err != nil {
		log.Fatal().Err(err).Msg("Failed to write CSV header")
	}

	// Write data
	for _, r := range results {
		record := []string{
			r.Thai,
			r.English,
			r.AudioFile,
			r.ImageFile,
			r.SourceShow,
			r.OriginalLang,
			fmt.Sprintf("%d", r.WordCount),
		}
		if err := writer.Write(record); err != nil {
			log.Fatal().Err(err).Msg("Failed to write CSV record")
		}
	}
}

func outputAnkiCSV(results []SentenceResult) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Anki CSV format: Front,Thai,English,Audio,Image,Show
	header := []string{"Front", "Thai", "English", "Audio", "Image", "Show"}
	if err := writer.Write(header); err != nil {
		log.Fatal().Err(err).Msg("Failed to write Anki CSV header")
	}

	for _, r := range results {
		audioTag := ""
		if r.AudioFile != "" {
			// Check if audio file already has [sound: tag
			if strings.HasPrefix(r.AudioFile, "[sound:") {
				audioTag = r.AudioFile
			} else {
				audioTag = fmt.Sprintf("[sound:%s]", r.AudioFile)
			}
		}

		imageTag := ""
		if r.ImageFile != "" {
			imageTag = fmt.Sprintf(`<img src="%s">`, r.ImageFile)
		}

		record := []string{
			audioTag,    // Front (audio only)
			r.Thai,
			r.English,
			audioTag,    // Audio field
			imageTag,    // Image field
			r.SourceShow,
		}
		if err := writer.Write(record); err != nil {
			log.Fatal().Err(err).Msg("Failed to write Anki CSV record")
		}
	}
}

func outputText(results []SentenceResult) {
	// Output plain text format optimized for sub-agents
	// Default format: Index	Thai
	// With flags: [Difficulty]	[Show]	Index	Thai/English

	for i, r := range results {
		// Build the output line
		var line string

		// Add difficulty if requested
		if *showDifficulty {
			line = fmt.Sprintf("%d\t", r.Difficulty)
		}

		// Add show name if requested
		if *showShow {
			showName := truncateString(r.SourceShow, 17)
			line += fmt.Sprintf("%s\t", showName)
		}

		// Add index
		line += fmt.Sprintf("%d\t", i)

		// Add primary text (Thai by default, English if --english flag)
		if *showEnglish {
			english := truncateString(r.English, 129)
			line += english
		} else {
			line += r.Thai
		}

		fmt.Println(line)

		// If --full flag is set, show the other language on next line
		if *full {
			if *showEnglish {
				// If showing English, add Thai on next line
				fmt.Printf("\t\t\t%s\n", r.Thai)
			} else {
				// If showing Thai (default), add English on next line
				english := truncateString(r.English, 129)
				fmt.Printf("\t\t\t%s\n", english)
			}
		}
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

// parseIndices parses a comma-separated string of indices
func parseIndices(indicesStr string) []int {
	parts := strings.Split(indicesStr, ",")
	var indices []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if idx, err := strconv.Atoi(p); err == nil {
			indices = append(indices, idx)
		}
	}
	return indices
}

// exportToTSV exports selected results to a TSV file
func exportToTSV(results []SentenceResult, indices []int) string {
	// First validate all indices are in range
	for _, idx := range indices {
		if idx < 0 || idx >= len(results) {
			log.Fatal().Int("index", idx).Int("max", len(results)-1).Msg("Index out of range")
		}
	}

	// STRICT VALIDATION: Require --output-dir when exporting
	if *outputDir == "" {
		log.Fatal().Msg("--output-dir is required when using --export-indices. Create a working directory first (e.g., /tmp/thai_selector_TIMESTAMP_PID/)")
	}

	// STRICT VALIDATION: Reject lazy "/tmp" usage
	if *outputDir == "/tmp" {
		log.Fatal().Msg("Cannot use /tmp as output directory. Sub-agents must create their own working directory (e.g., /tmp/thai_selector_TIMESTAMP_PID/)")
	}

	// Verify output directory EXISTS (no auto-creation)
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		log.Fatal().Str("dir", *outputDir).Msg("Output directory does not exist. Create your working directory first.")
	}

	// Generate unique filename in specified output directory
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(*outputDir, fmt.Sprintf("tsv_export_%s.tsv", timestamp))

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create export file")
	}
	defer file.Close()

	// Write TSV
	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	// Write selected rows (no header, pure data)
	// Only export core columns (0-5), excluding context sentences (6+)
	for _, idx := range indices {
		if results[idx].OriginalRow != nil {
			row := results[idx].OriginalRow
			// Only take first 6 columns (audio, timestamp, episode, image, thai, english)
			if len(row) > 6 {
				row = row[:6]
			}
			writer.Write(row)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal().Err(err).Msg("Failed to write TSV file")
	}

	return filename
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	heroicons "github.com/indaco/templheroicons"
)

// Constants
const (
	Size16 heroicons.Size = "16"
	Size20 heroicons.Size = "20"
	Size24 heroicons.Size = "24"

	cacheDuration = 30 * 24 * time.Hour
	datasetURL    = "https://raw.githubusercontent.com/iconify/icon-sets/refs/heads/master/json/heroicons_cache.json"
	maxRetries    = 3
	retryDelay    = 5 * time.Second
	cacheFile     = "heroicons_cache.json"
	outputFile    = "heroicons_generated.go"
)

// Utility for consistent error logging
func logAndExit(err error, context string) {
	log.Fatalf("%s: %v", context, err)
}

// Converts a kebab-case string to PascalCase.
func toPascalCase(input string) string {
	var builder strings.Builder
	for _, part := range strings.Split(input, "-") {
		if len(part) > 0 {
			builder.WriteString(strings.ToUpper(part[:1]))
			builder.WriteString(part[1:])
		}
	}
	return builder.String()
}

// Verifies if the cache file is valid.
func isCacheValid(filepath string, maxAge time.Duration) bool {
	info, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < maxAge
}

// Reads the dataset from the cache.
func loadCache(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}

// Saves the dataset to the cache.
func saveCache(filepath string, data []byte) error {
	return os.WriteFile(filepath, data, 0644)
}

// Fetches the dataset with retry logic.
func fetchDatasetWithRetry(url string, maxRetries int, delay time.Duration) ([]byte, error) {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Fetching heroicons dataset (attempt %d/%d)...\n", attempt, maxRetries)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				resp.Body.Close()
			}
			lastErr = fmt.Errorf("attempt %d: %w", attempt, err)
		} else {
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}
		if attempt < maxRetries {
			log.Printf("Retrying in %s...\n", delay)
			time.Sleep(delay)
		}
	}
	return nil, lastErr
}

// Fetches and caches the dataset.
func fetchAndCacheDataset(url string, cachePath string, maxAge time.Duration) ([]byte, error) {
	if isCacheValid(cachePath, maxAge) {
		log.Println("Using cached dataset.")
		return loadCache(cachePath)
	}
	data, err := fetchDatasetWithRetry(url, maxRetries, retryDelay)
	if err != nil {
		return nil, err
	}

	err = saveCache(cachePath, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Parses icons from the JSON dataset.
func parseIcons(jsonData []byte) (map[string]heroicons.Icon, error) {
	var jsonDataStruct struct {
		Icons map[string]struct {
			Body string `json:"body"`
		} `json:"icons"`
	}

	if err := json.Unmarshal(jsonData, &jsonDataStruct); err != nil {
		return nil, err
	}

	icons := make(map[string]heroicons.Icon)
	for name, iconData := range jsonDataStruct.Icons {
		icon := heroicons.Icon{
			Name: name,
			Body: iconData.Body,
			Size: Size24,
			Type: "Outline",
		}

		if strings.Contains(name, "solid") {
			icon.Type = "Solid"
		}

		switch {
		case strings.Contains(name, "16"):
			icon.Size = Size16
			icon.Type = "Micro"
		case strings.Contains(name, "20"):
			icon.Size = Size20
			icon.Type = "Mini"
		}
		icons[name] = icon
	}

	return icons, nil
}

// Cleans and standardizes icon names.
func cleanIconName(name string) string {
	return strings.NewReplacer("-16", "", "-20", "", "-solid", "").Replace(name)
}

// Generates the Go struct name for an icon.
func generateStructName(icon heroicons.Icon) string {
	baseName := toPascalCase(cleanIconName(icon.Name))
	switch icon.Type {
	case "Micro":
		return baseName + "Micro"
	case "Mini":
		return baseName + "Mini"
	case "Solid":
		return baseName + "Solid"
	default:
		return baseName
	}
}

// Generates a Go file with icon definitions.
func generateGoFile(outputFilePath string, icons map[string]heroicons.Icon) error {
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var builder strings.Builder
	builder.WriteString("// Code generated by 'scripts/icons-maker.go'; DO NOT EDIT.\n")
	builder.WriteString("package templheroicons\n\nvar (\n")
	var structs []string
	for _, icon := range icons {
		structs = append(structs, fmt.Sprintf("\t%s = &Icon{Name: \"%s\", Body: `%s`, Size: \"%s\", Type: \"%s\"}\n",
			generateStructName(icon), icon.Name, icon.Body, icon.Size.String(), icon.Type))
	}
	sort.Strings(structs)
	for _, structDef := range structs {
		builder.WriteString(structDef)
	}
	builder.WriteString(")\n")
	_, err = outFile.WriteString(builder.String())
	return err
}

// ensureDir ensures that the specified directory exists. If it does not exist, it creates it.
func ensureDir(dir string) error {
	err := os.MkdirAll(dir, 0755) // Create the directory and all necessary parents.
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

func main() {
	cacheFilePath := path.Join("..", "data", cacheFile)
	outputFilePath := path.Join("..", outputFile)

	// Ensure the "data" directory exists.
	dataDir := path.Dir(cacheFilePath) // Get the directory from the path.
	if err := ensureDir(dataDir); err != nil {
		log.Fatalf("Error ensuring data directory exists: %v", err)
	}

	// Fetch and parse the JSON dataset.
	data, err := fetchAndCacheDataset(datasetURL, cacheFilePath, cacheDuration)
	if err != nil {
		logAndExit(err, "Fetching dataset")
	}

	icons, err := parseIcons(data)
	if err != nil {
		logAndExit(err, "Parsing icons")
	}

	// Generate Go file with icon definitions.
	if err := generateGoFile(outputFilePath, icons); err != nil {
		logAndExit(err, "Generating Go file")
	}

	log.Println("heroicons_generated.go successfully created.")
}

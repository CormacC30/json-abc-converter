package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Tune represents the structure of each tune in the JSON file
type Tune struct {
	TuneID    string `json:"tune_id"`
	SettingID string `json:"setting_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Meter     string `json:"meter"`
	Mode      string `json:"mode"`
	ABC       string `json:"abc"`
	Date      string `json:"date,omitempty"`
	Username  string `json:"username,omitempty"`
}

func main() {
	inputFile := flag.String("input", "", "Path to the input JSON file")
	outputDir := flag.String("output", ".", "Directory for output ABC files")
	singleFile := flag.Bool("single", false, "Output to a single file instead of multiple files")
	singleFilePath := flag.String("outfile", "all_tunes.abc", "Name of the single output file (used with -single)")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please provide an input file with the -input flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read the input file
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON data
	var tunes []Tune
	err = json.Unmarshal(data, &tunes)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d tunes in the input file\n", len(tunes))

	// Create output directory if it doesn't exist
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(*outputDir, 0755)
		if err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created output directory: %s\n", *outputDir)
	}

	if *singleFile {
		// Write all tunes to a single file
		outputToSingleFile(tunes, *outputDir, *singleFilePath)
	} else {
		// Write each tune to a separate file
		outputToMultipleFiles(tunes, *outputDir)
	}
}

func outputToSingleFile(tunes []Tune, outputDir, fileName string) {
	outputPath := filepath.Join(outputDir, fileName)
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	for _, tune := range tunes {
		// Write the standard ABC header fields
		fmt.Fprintf(f, "X:%s\n", tune.SettingID)
		fmt.Fprintf(f, "T:%s\n", tune.Name)
		fmt.Fprintf(f, "R:%s\n", tune.Type)
		fmt.Fprintf(f, "M:%s\n", tune.Meter)
		fmt.Fprintf(f, "K:%s\n", modeToBetter(tune.Mode))
		if tune.Username != "" {
			fmt.Fprintf(f, "Z:%s\n", tune.Username)
		}
		if tune.Date != "" {
			fmt.Fprintf(f, "H:Added %s\n", tune.Date)
		}

		// Write the ABC notation
		fmt.Fprintln(f, tune.ABC)
		fmt.Fprintln(f, "") // Add a blank line between tunes
	}

	fmt.Printf("Successfully wrote %d tunes to %s\n", len(tunes), outputPath)
}

func outputToMultipleFiles(tunes []Tune, outputDir string) {
	processed := 0
	
	for i, tune := range tunes {
		// Create a filename based on the tune's ID and name
		safeName := sanitizeFileName(tune.Name)
		fileName := fmt.Sprintf("%s_%s.abc", tune.SettingID, safeName)
		outputPath := filepath.Join(outputDir, fileName)

		// Ensure we create a new file for each tune
		f, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating output file %s: %v\n", outputPath, err)
			continue
		}

		// Write the standard ABC header fields
		fmt.Fprintf(f, "X:%s\n", tune.SettingID)
		fmt.Fprintf(f, "T:%s\n", tune.Name)
		fmt.Fprintf(f, "R:%s\n", tune.Type)
		fmt.Fprintf(f, "M:%s\n", tune.Meter)
		fmt.Fprintf(f, "K:%s\n", modeToBetter(tune.Mode))
		if tune.Username != "" {
			fmt.Fprintf(f, "Z:%s\n", tune.Username)
		}
		if tune.Date != "" {
			fmt.Fprintf(f, "H:Added %s\n", tune.Date)
		}

		// Write the ABC notation
		fmt.Fprintln(f, tune.ABC)

		// Make sure to close the file after writing
		f.Close()
		processed++

		if (i+1)%100 == 0 {
			fmt.Printf("Processed %d tunes...\n", i+1)
		}
	}

	fmt.Printf("Successfully wrote %d tunes to individual files in %s\n", processed, outputDir)
}

// sanitizeFileName removes characters that are not allowed in filenames
func sanitizeFileName(name string) string {
	// Replace problematic characters with underscores
	illegalChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name

	for _, char := range illegalChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Shorten length if needed
	if len(result) > 50 {
		result = result[:50]
	}

	return result
}

// modeToBetter converts mode formats like "Gmajor" to "G" for ABC notation
func modeToBetter(mode string) string {
	// Handle common mode formats
	mode = strings.ToLower(mode)
	if strings.HasSuffix(mode, "major") {
		return strings.TrimSuffix(mode, "major")
	} else if strings.HasSuffix(mode, "minor") {
		return strings.TrimSuffix(mode, "minor") + "m"
	} else if strings.Contains(mode, "mixolydian") {
		// For mixolydian modes, use K:D mix format
		return strings.Replace(mode, "mixolydian", " mix", 1)
	} else if strings.Contains(mode, "dorian") {
		return strings.Replace(mode, "dorian", " dor", 1)
	} else if strings.Contains(mode, "phrygian") {
		return strings.Replace(mode, "phrygian", " phr", 1)
	} else if strings.Contains(mode, "lydian") {
		return strings.Replace(mode, "lydian", " lyd", 1)
	} else if strings.Contains(mode, "locrian") {
		return strings.Replace(mode, "locrian", " loc", 1)
	}
	
	// Return as is if no known pattern
	return mode
}

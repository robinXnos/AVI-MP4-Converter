package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	jsonFile, err := os.Open("paths.json")
	if err != nil {
		fmt.Println("Error reading JsonFile", err)
		return
	}

	byteVal, err := io.ReadAll(jsonFile)

	var paths struct {
		InputDir      string
		HandbrakePath string
		OutputDir     string
	}
	json.Unmarshal(byteVal, &paths)

	fmt.Printf("paths: %v\n", paths)

	// Find AVI files (case-insensitive) in the input directory
	var files []string
	err = filepath.Walk(paths.InputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has an .avi or .AVI extension
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".avi") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error reading the input directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No AVI files found.")
		return
	}

	for _, inputFile := range files {
		outputFile := filepath.Join(paths.OutputDir, strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))+".mp4")
		fmt.Printf("Converting: %s -> %s\n", inputFile, outputFile)

		// HandBrakeCLI command
		cmd := exec.Command(
			paths.HandbrakePath,
			"-i", inputFile,
			"-o", outputFile,
			"-e", "x264",
			"-q", "16",
			"-B", "192",
		)

		// Execute the command and handle errors
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error converting %s: %v\n", inputFile, err)
			fmt.Printf("HandBrakeCLI output:\n%s\n", string(output))
		} else {
			fmt.Printf("Successfully converted: %s\n", outputFile)
		}
	}
}

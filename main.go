package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt the user for the path to HandBrakeCLI
	fmt.Print("Enter the path of HandBrakeCLI: ")
	handBrakeCLIPath, _ := reader.ReadString('\n')
	handBrakeCLIPath = strings.TrimSpace(handBrakeCLIPath)

	// Check if the HandBrakeCLI path exists
	if _, err := os.Stat(handBrakeCLIPath); os.IsNotExist(err) {
		fmt.Printf("Error: HandBrakeCLI not found at %s\n", handBrakeCLIPath)
		return
	}

	// Prompt the user for the input directory
	fmt.Print("Enter the input directory: ")
	inputDir, _ := reader.ReadString('\n')
	inputDir = strings.TrimSpace(inputDir)

	// Check if the input directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory not found at %s\n", inputDir)
		return
	}

	// Prompt the user for the output directory
	fmt.Print("Enter the output directory: ")
	outputDir, _ := reader.ReadString('\n')
	outputDir = strings.TrimSpace(outputDir)

	// Create the output directory if it does not exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating output folder: %v\n", err)
			return
		}
	}

	// Find AVI files (case-insensitive) in the input directory
	var files []string
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
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
		outputFile := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))+".mp4")
		fmt.Printf("Converting: %s -> %s\n", inputFile, outputFile)

		// HandBrakeCLI command
		cmd := exec.Command(
			handBrakeCLIPath,
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

package main

import (
	"encoding/json"
	"fmt"
	"os"

	urlmetadata "github.com/naruebaet/url-metadata-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url>")
		os.Exit(1)
	}

	url := os.Args[1]

	// Create a client with default settings
	client := urlmetadata.DefaultClient()

	// Fetch metadata
	metadata, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Output as JSON
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

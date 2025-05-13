# URL Metadata Go

A Go library for fetching metadata from URLs. Extract title, description, images, and more from web pages.

## Installation

```
go get github.com/naruebaet/url-metadata-go
```

## Features

- Extract page title and description
- Get Open Graph metadata
- Find favicon and image URLs
- Extract author, keywords, and language information
- Automatic resolution of relative URLs
- Context support for cancellation and timeouts
- Customizable HTTP client

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"log"

	urlmetadata "github.com/naruebaet/url-metadata-go"
)

func main() {
	// Create a client with default settings
	client := urlmetadata.DefaultClient()
	
	// Fetch metadata
	metadata, err := client.Get("https://example.com")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	// Use the metadata
	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("Description: %s\n", metadata.Description)
	fmt.Printf("Image: %s\n", metadata.ImageURL)
}
```

### With Custom Timeout

```go
package main

import (
	"time"
	
	urlmetadata "github.com/naruebaet/url-metadata-go"
)

func main() {
	// Create a client with 5-second timeout
	client := urlmetadata.NewClient(5 * time.Second)
	
	// Fetch metadata
	metadata, err := client.Get("https://example.com")
	// ...
}
```

### With Context

```go
package main

import (
	"context"
	"time"
	
	urlmetadata "github.com/naruebaet/url-metadata-go"
)

func main() {
	// Create a client
	client := urlmetadata.DefaultClient()
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// Fetch metadata with context
	metadata, err := client.GetWithContext(ctx, "https://example.com")
	// ...
}
```

## Metadata Fields

The `Metadata` struct contains the following fields:

- `URL`: The original URL
- `Title`: Page title (from og:title or title tag)
- `Description`: Page description (from og:description or meta description)
- `SiteName`: Site name (from og:site_name)
- `ImageURL`: Primary image URL (from og:image)
- `FaviconURL`: Favicon URL
- `Type`: Content type (from og:type)
- `Author`: Content author
- `Keywords`: Array of keywords
- `Language`: Content language

## License

MIT

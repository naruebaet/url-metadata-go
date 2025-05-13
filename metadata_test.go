package urlmetadata

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetadataExtraction(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Test Page</title>
				<meta name="description" content="This is a test page">
				<meta property="og:title" content="OG Test Title">
				<meta property="og:description" content="OG test description">
				<meta property="og:image" content="https://example.com/image.jpg">
				<meta property="og:site_name" content="Test Site">
				<meta property="og:type" content="website">
				<meta name="author" content="Test Author">
				<meta name="keywords" content="test, metadata, go">
				<meta name="language" content="en">
				<link rel="icon" href="/favicon.ico">
			</head>
			<body>
				<h1>Test Page Content</h1>
			</body>
			</html>
		`))
	}))
	defer server.Close()
	
	client := NewClient(5 * time.Second)
	metadata, err := client.Get(server.URL)
	
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}
	
	// Check the extracted metadata
	if metadata.Title != "OG Test Title" {
		t.Errorf("Expected title 'OG Test Title', got '%s'", metadata.Title)
	}
	
	if metadata.Description != "OG test description" {
		t.Errorf("Expected description 'OG test description', got '%s'", metadata.Description)
	}
	
	if metadata.ImageURL != "https://example.com/image.jpg" {
		t.Errorf("Expected image URL 'https://example.com/image.jpg', got '%s'", metadata.ImageURL)
	}
	
	if metadata.SiteName != "Test Site" {
		t.Errorf("Expected site name 'Test Site', got '%s'", metadata.SiteName)
	}
	
	if metadata.Type != "website" {
		t.Errorf("Expected type 'website', got '%s'", metadata.Type)
	}
	
	if metadata.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", metadata.Author)
	}
	
	if len(metadata.Keywords) != 3 || metadata.Keywords[0] != "test" {
		t.Errorf("Expected keywords '[test metadata go]', got '%v'", metadata.Keywords)
	}
	
	if metadata.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", metadata.Language)
	}
	
	// Favicon should be resolved to absolute URL
	expectedFavicon := server.URL + "/favicon.ico"
	if metadata.FaviconURL != expectedFavicon {
		t.Errorf("Expected favicon URL '%s', got '%s'", expectedFavicon, metadata.FaviconURL)
	}
}

func TestFallbackTitle(t *testing.T) {
	// Setup test server with HTML that has no og:title
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Test Page Title</title>
			</head>
			<body>
				<h1>Content</h1>
			</body>
			</html>
		`))
	}))
	defer server.Close()
	
	client := NewClient(5 * time.Second)
	metadata, err := client.Get(server.URL)
	
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}
	
	// Should fall back to the title tag
	if metadata.Title != "Test Page Title" {
		t.Errorf("Expected title 'Test Page Title', got '%s'", metadata.Title)
	}
}

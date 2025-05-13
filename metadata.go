package urlmetadata

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Metadata represents extracted information from a URL
type Metadata struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	SiteName    string   `json:"site_name"`
	ImageURL    string   `json:"image_url"`
	FaviconURL  string   `json:"favicon_url"`
	Type        string   `json:"type"`
	Author      string   `json:"author"`
	Keywords    []string `json:"keywords"`
	Language    string   `json:"language"`
}

// Client is responsible for fetching and parsing URL metadata
type Client struct {
	httpClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new metadata client with the given timeout
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// DefaultClient returns a client with reasonable defaults
func DefaultClient() *Client {
	return NewClient(10 * time.Second)
}

// Get fetches metadata from the given URL
func (c *Client) Get(url string) (*Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.GetWithContext(ctx, url)
}

// GetWithContext fetches metadata from the given URL with the provided context
func (c *Client) GetWithContext(ctx context.Context, url string) (*Metadata, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set a user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; UrlMetadataBot/1.0)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	metadata := &Metadata{
		URL: url,
	}

	// Extract metadata from HTML
	extractMetadata(doc, metadata)

	// Resolve relative URLs
	metadata.resolveRelativeURLs(url)

	return metadata, nil
}

// resolveRelativeURLs converts relative URLs to absolute URLs
func (m *Metadata) resolveRelativeURLs(baseURL string) {
	if m.ImageURL != "" && !strings.HasPrefix(m.ImageURL, "http") {
		m.ImageURL = resolveURL(baseURL, m.ImageURL)
	}

	if m.FaviconURL != "" && !strings.HasPrefix(m.FaviconURL, "http") {
		m.FaviconURL = resolveURL(baseURL, m.FaviconURL)
	}
}

// resolveURL converts a relative URL to an absolute URL
func resolveURL(baseURL, relURL string) string {
	// Simple implementation - can be improved
	if strings.HasPrefix(relURL, "/") {
		// Get the scheme and host
		parts := strings.SplitN(baseURL, "/", 4)
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2] + relURL
		}
	}
	return relURL
}

package urlmetadata

import (
	"strings"

	"golang.org/x/net/html"
)

// extractMetadata parses the HTML document and fills the metadata struct
func extractMetadata(n *html.Node, metadata *Metadata) {
	// First, try to extract Open Graph and other meta tags
	extractMetaTags(n, metadata)

	// If title is still empty, try to get it from the title tag
	if metadata.Title == "" {
		metadata.Title = extractTitle(n)
	}

	// If favicon is not found in meta tags, look for favicon link
	if metadata.FaviconURL == "" {
		metadata.FaviconURL = extractFavicon(n)
	}
}

// extractMetaTags extracts metadata from meta tags, including Open Graph tags
func extractMetaTags(n *html.Node, metadata *Metadata) {
	if n.Type == html.ElementNode {
		// Check for meta tags
		if n.Data == "meta" {
			var property, name, content string

			for _, attr := range n.Attr {
				switch attr.Key {
				case "property":
					property = attr.Val
				case "name":
					name = attr.Val
				case "content":
					content = attr.Val
				}
			}

			// Process Open Graph tags
			switch property {
			case "og:title":
				metadata.Title = content
			case "og:description":
				metadata.Description = content
			case "og:image":
				metadata.ImageURL = content
			case "og:site_name":
				metadata.SiteName = content
			case "og:type":
				metadata.Type = content
			}

			// Process standard meta tags
			switch name {
			case "description":
				if metadata.Description == "" {
					metadata.Description = content
				}
			case "author":
				metadata.Author = content
			case "keywords":
				if content != "" {
					keywords := strings.Split(content, ",")
					for i, k := range keywords {
						keywords[i] = strings.TrimSpace(k)
					}
					metadata.Keywords = keywords
				}
			case "language":
				metadata.Language = content
			}
		}

		// Check for link tags for favicon
		if n.Data == "link" {
			var rel, href string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "rel":
					rel = attr.Val
				case "href":
					href = attr.Val
				}
			}

			// Look for favicon in various formats
			if (rel == "icon" || rel == "shortcut icon") && metadata.FaviconURL == "" {
				metadata.FaviconURL = href
			}
		}
	}

	// Recursively process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractMetaTags(c, metadata)
	}
}

// extractTitle gets the title from the title tag
func extractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data
		}
	}

	// Recursively search for title
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title := extractTitle(c)
		if title != "" {
			return title
		}
	}

	return ""
}

// extractFavicon tries to find a favicon link
func extractFavicon(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "link" {
		var rel, href string
		for _, attr := range n.Attr {
			switch attr.Key {
			case "rel":
				rel = attr.Val
			case "href":
				href = attr.Val
			}
		}

		if rel == "icon" || rel == "shortcut icon" {
			return href
		}
	}

	// Recursively search for favicon
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		favicon := extractFavicon(c)
		if favicon != "" {
			return favicon
		}
	}

	return ""
}

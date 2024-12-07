package templheroicons

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/a-h/templ"
)

var (
	iconBodyCache = map[string]string{}
	cacheMutex    sync.Mutex
)

// Size represents the size of UI components (e.g., small, medium, large).
type Size string

// String returns the string representation of a Size.
func (s Size) String() string {
	return string(s)
}

// Icon represents a single icon with its attributes.
type Icon struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        Size   `json:"size"`
	Stroke      string `json:"stroke,omitempty"`
	StrokeWidth string `json:"strokeWidth,omitempty"`
	Fill        string `json:"fill,omitempty"`
	Attrs       templ.Attributes
	body        string // Cached Body
}

func (i *Icon) Render() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, makeSVGTag(i))
		return err
	})
	//return templ.Raw(i.String())
}

func makeSVGTag(icon *Icon) string {
	fill := "currentColor"
	if icon.Fill != "" {
		fill = icon.Fill // Explicitly set fill takes priority
	} else if icon.Type == "Outline" {
		fill = "none" // Fallback for Outline type
	}

	stroke := "currentColor"
	if icon.Stroke != "" {
		stroke = icon.Stroke
	}

	strokeWidth := "1.5"
	if icon.StrokeWidth != "" {
		strokeWidth = icon.StrokeWidth
	}

	// Fetch the body if it's not cached.
	if icon.body == "" {
		body, err := getIconBody(icon.Name)
		if err != nil {
			return fmt.Sprintf("<!-- Error: %s -->", err)
		}
		icon.body = body
	}

	var builder strings.Builder

	// Start the <svg> tag with common attributes.
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg"`)
	fmt.Fprintf(&builder, ` width="%[1]s" height="%[1]s" viewBox="0 0 %[2]s %[2]s"`, icon.Size.String(), getViewBox(icon.Type))

	// Add attributes based on the type.
	switch icon.Type {
	case "Outline":
		fmt.Fprintf(&builder, ` fill="%s" stroke-width="%s" stroke="%s"`, fill, strokeWidth, stroke)
	case "Solid", "Mini", "Micro":
		fmt.Fprintf(&builder, ` fill="%s"`, fill)
	default:
		fmt.Fprintf(&builder, ` fill="%s" stroke-width="%s" stroke="%s"`, fill, strokeWidth, stroke)
	}

	// Add user-defined attributes in deterministic order.
	addAttributesToSVG(&builder, icon.Attrs)

	// Close the opening SVG tag.
	builder.WriteString(">")

	// Append the SVG body and close the tag.
	builder.WriteString(icon.body)
	builder.WriteString(`</svg>`)

	return builder.String()
}

// getViewBox returns the appropriate viewBox size based on the icon type.
func getViewBox(iconType string) string {
	switch iconType {
	case "Mini":
		return "20"
	case "Micro":
		return "16"
	default:
		return "24" // Default for "Outline" and "Solid".
	}
}

// getIconBody retrieves the body of an icon by its name, with thread-safe caching.
func getIconBody(name string) (string, error) {

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the body is already cached.
	if body, found := iconBodyCache[name]; found {
		return body, nil
	}

	// Load and parse the JSON only once (controlled elsewhere).
	var parsedData struct {
		Icons map[string]struct {
			Body string `json:"body"`
		} `json:"icons"`
	}

	// Read and parse the JSON data.
	jsonFilename := "data/heroicons_cache.json"
	heroiconsData, _ := heroiconsJSONSource.Open(jsonFilename)
	defer heroiconsData.Close()

	data, _ := io.ReadAll(heroiconsData)

	if err := json.Unmarshal(data, &parsedData); err != nil {
		return "", fmt.Errorf("failed to parse heroicons JSON: %w", err)
	}

	// Populate the cache.
	for iconName, icon := range parsedData.Icons {
		iconBodyCache[iconName] = icon.Body
	}

	// Return the requested icon body.
	body, exists := iconBodyCache[name]
	if !exists {
		return "", fmt.Errorf("icon '%s' not found", name)
	}
	return body, nil
}

// Reserved attributes for SVG tags that should not be overwritten.
var reservedSVGAttributes = map[string]struct{}{
	"xmlns":        {},
	"viewBox":      {},
	"width":        {},
	"height":       {},
	"stroke-width": {},
	"stroke":       {},
	"fill":         {},
}

// sanitizeAttribute ensures that attribute keys and values are safe for inclusion in the SVG tag.
func sanitizeAttribute(key, value string) (string, string, bool) {
	// Define allowlist for event attributes
	allowedEventAttributes := map[string]struct{}{
		"onclick":  {},
		"onchange": {},
		"onhover":  {},
	}

	// Check for unsafe attributes
	if _, isEvent := allowedEventAttributes[key]; isEvent {
		// For event attributes, only allow simple JS functions (no <script> tags, eval, etc.)
		if strings.Contains(strings.ToLower(value), "<script>") || strings.Contains(strings.ToLower(value), "javascript:") {
			return "", "", false // Unsafe value
		}
	}

	// Escape any unsafe characters for all attributes
	escapedKey := html.EscapeString(key)
	escapedValue := html.EscapeString(value)

	return escapedKey, escapedValue, true // Safe attribute
}

// addAttributesToSVG adds templ.Attributes to the SVG tag, placing them at the end of the <svg> opening tag.
// Reserved attributes are skipped to avoid overwriting critical SVG settings.
// Attributes are sanitized to prevent XSS or injection attacks.
func addAttributesToSVG(builder *strings.Builder, attrs templ.Attributes) {
	if len(attrs) == 0 {
		return
	}

	// Extract keys and sort them for deterministic order
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Process attributes in sorted order
	for _, key := range keys {
		value, ok := attrs[key].(string) // Ensure value is a string
		if !ok {
			// Skip attributes with non-string values
			continue
		}

		// Skip reserved attributes
		if _, isReserved := reservedSVGAttributes[key]; isReserved {
			continue
		}

		// Sanitize the attribute
		sanitizedKey, sanitizedValue, ok := sanitizeAttribute(key, value)
		if !ok {
			// Skip attributes that are not safe
			continue
		}

		// Add the sanitized attribute to the SVG tag
		fmt.Fprintf(builder, ` %s="%v"`, sanitizedKey, sanitizedValue)
	}
}

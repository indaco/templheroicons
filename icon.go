package templheroicons

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
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
	mu          sync.Mutex
}

// _clone creates a new copy of the Icon.
func (i *Icon) clone() *Icon {
	return &Icon{
		Name:        i.Name,
		Type:        i.Type,
		Size:        i.Size,
		Stroke:      i.Stroke,
		StrokeWidth: i.StrokeWidth,
		Fill:        i.Fill,
		Attrs:       i.Attrs,
		body:        i.body,
	}
}

// SetSize sets the icon size in pixels.
func (i *Icon) SetSize(size int) *Icon {
	clone := i.clone()
	clone.Size = Size(strconv.Itoa(size))
	return clone
}

// SetStroke sets the stroke attribute for the SVG tag.
func (i *Icon) SetStroke(value string) *Icon {
	clone := i.clone()
	clone.Stroke = value
	return clone
}

// SetStrokeWidth sets the stroke-width attribute for the SVG tag.
func (i *Icon) SetStrokeWidth(value string) *Icon {
	clone := i.clone()
	clone.StrokeWidth = value
	return clone
}

// SetFill sets the fill attribute for the SVG tag.
func (i *Icon) SetFill(value string) *Icon {
	clone := i.clone()
	clone.Fill = value
	return clone
}

// SetAttrs sets the attributes for the SVG tag.
func (i *Icon) SetAttrs(attrs templ.Attributes) *Icon {
	clone := i.clone()
	clone.Attrs = attrs
	return clone
}

func (i *Icon) ensureDefaults() {
	// Set default fill based on the icon type
	if i.Fill == "" {
		switch i.Type {
		case "Outline":
			i.Fill = "none"
		default: // Applies to "Solid", "Mini", "Micro" and other types
			i.Fill = "currentColor"
		}
	}

	if i.Stroke == "" {
		i.Stroke = "currentColor"
	}

	if i.StrokeWidth == "" {
		i.StrokeWidth = "1.5"
	}
}

// String returns the SVG data of the Icon, including updated size and attributes.
func (i *Icon) String() string {
	// Ensure defaults are set.
	i.ensureDefaults()

	if i.body == "" {
		i.mu.Lock()
		defer i.mu.Unlock()

		if i.body == "" { // Double-check after locking
			body, err := getIconBody(i.Name)
			if err != nil {
				return fmt.Sprintf("<!-- Error: %s -->", err)
			}
			i.body = body
		}
	}

	var builder strings.Builder

	// Start the <svg> tag with common attributes.
	builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg"`)
	fmt.Fprintf(&builder, ` width="%[1]s" height="%[1]s" viewBox="0 0 %[2]s %[2]s"`, i.Size.String(), getViewBox(i.Type))

	// Add attributes based on the type.
	switch i.Type {
	case "Outline":
		fmt.Fprintf(&builder, ` fill="none" stroke-width="%s" stroke="%s"`, i.StrokeWidth, i.Stroke)
	case "Solid", "Mini", "Micro":
		fmt.Fprintf(&builder, ` fill="%s"`, i.Fill)
	default:
		fmt.Fprintf(&builder, ` fill="%s" stroke-width="%s" stroke="%s"`, i.Fill, i.StrokeWidth, i.Stroke)
	}

	// Add user-defined attributes in deterministic order.
	addAttributesToSVG(&builder, i.Attrs)

	// Close the opening SVG tag.
	builder.WriteString(">")

	// Append the SVG body and close the tag.
	builder.WriteString(i.body)
	builder.WriteString(`</svg>`)

	return builder.String()
}

func (i *Icon) Render() templ.Component {
	return templ.Raw(i.String())
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

// getIconBody retrieves the body of an icon by its name.
func getIconBody(name string) (string, error) {
	var loadError error

	// Load and parse the JSON only once
	iconDataOnce.Do(func() {
		iconData = make(map[string]string)

		var parsedData struct {
			Icons map[string]struct {
				Body string `json:"body"`
			} `json:"icons"`
		}

		jsonFilename := "data/heroicons_cache.json"
		heroiconsData, _ := heroiconsJSONSource.Open(jsonFilename)
		defer heroiconsData.Close()

		data, _ := io.ReadAll(heroiconsData)

		if err := json.Unmarshal(data, &parsedData); err != nil {
			loadError = fmt.Errorf("failed to parse heroicons JSON: %w", err)
			return
		}

		for iconName, icon := range parsedData.Icons {
			iconData[iconName] = icon.Body
		}
	})

	if loadError != nil {
		return "", loadError
	}

	// Lookup the icon body
	body, exists := iconData[name]
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

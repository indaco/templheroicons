package templheroicons

import (
	"fmt"
	"html"
	"sort"
	"strings"

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
	Type        string `json:"type"` // Type is either "Outline", "Solid", "Mini", or "Micro"
	Size        Size   `json:"size"`
	Body        string `json:"body"`
	Stroke      string
	StrokeWidth string
	Fill        string
	Attrs       templ.Attributes
}

// SetStroke sets the stroke attribute for the SVG tag.
func (i *Icon) SetStroke(value string) *Icon {
	i.Stroke = value
	return i
}

// SetStrokeWidth sets the stroke-width attribute for the SVG tag.
func (i *Icon) SetStrokeWidth(value string) *Icon {
	i.StrokeWidth = value
	return i
}

// SetFill sets the fill attribute for the SVG tag.
func (i *Icon) SetFill(value string) *Icon {
	i.Fill = value
	return i
}

// SetAttrs sets the attributes for the SVG tag.
func (i *Icon) SetAttrs(attrs templ.Attributes) *Icon {
	i.Attrs = attrs
	return i
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
	var builder strings.Builder

	// Ensure defaults are set.
	i.ensureDefaults()

	// Start the <svg> tag.
	if i.Type == "Outline" || i.Type == "Solid" || i.Type == "Mini" || i.Type == "Micro" {
		fmt.Fprintf(&builder, `<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 %[1]s %[1]s"`, i.Size.String())
	} else {
		builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg"`)
	}

	// Add attributes based on the type.
	switch i.Type {
	case "Outline":
		fmt.Fprintf(&builder, ` fill="none" stroke-width="%s" stroke="%s"`, i.StrokeWidth, i.Stroke)
	case "Solid", "Mini", "Micro":
		fmt.Fprintf(&builder, ` fill="%s"`, i.Fill)
	}

	// Add user-defined attributes in deterministic order.
	addAttributesToSVG(&builder, i.Attrs)

	// Close the opening SVG tag.
	builder.WriteString(">")

	// Append the SVG body and close the tag.
	builder.WriteString(i.Body)
	builder.WriteString(`</svg>`)

	return builder.String()
}

func (i *Icon) Render() templ.Component {
	return templ.Raw(i.String())
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

package templheroicons

import (
	"fmt"
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
	Name  string `json:"name"`
	Body  string `json:"body"`
	Size  Size   `json:"size"`
	Type  string `json:"type"` // Type is either "Outline", "Solid", "Mini", or "Micro"
	Attrs templ.Attributes
}

// SetAttrs sets the attributes for the SVG tag.
func (i *Icon) SetAttrs(attrs templ.Attributes) *Icon {
	i.Attrs = attrs
	return i
}

// String returns the SVG data of the Icon, including updated size and attributes.
func (i *Icon) String() string {
	var builder strings.Builder

	// Generate the appropriate SVG opening tag based on the Type.
	switch i.Type {
	case "Outline":
		fmt.Fprintf(&builder, `<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" fill="none" viewBox="0 0 %s %s" stroke-width="1.5" stroke="currentColor"`,
			i.Size.String(), i.Size.String(), i.Size.String(), i.Size.String())
	case "Solid", "Mini", "Micro":
		fmt.Fprintf(&builder, `<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 %s %s" fill="currentColor"`,
			i.Size.String(), i.Size.String(), i.Size.String(), i.Size.String())
	default:
		builder.WriteString(`<svg xmlns="http://www.w3.org/2000/svg"`) // Fallback
	}

	// Add attributes to the opening SVG tag.
	addAttributesToSVG(&builder, i.Attrs)

	// Close the opening SVG tag.
	builder.WriteString(">")

	// Append the SVG body and close the tag.
	builder.WriteString(i.Body)
	builder.WriteString(`</svg>`)

	return builder.String()
}

func (i *Icon) Render2() templ.Component {
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

// addAttributesToSVG adds templ.Attributes to the SVG tag, placing them at the end of the <svg> opening tag.
// Reserved attributes are skipped to avoid overwriting critical SVG settings.
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
		value := attrs[key]
		// Skip reserved attributes
		if _, isReserved := reservedSVGAttributes[key]; isReserved {
			continue
		}

		// Add the attribute to the SVG tag
		fmt.Fprintf(builder, ` %s="%v"`, key, value)
	}
}

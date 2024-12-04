package templheroicons

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestIcon_String(t *testing.T) {
	tests := []struct {
		name     string
		icon     *Icon
		expected string
	}{
		{
			name: "Outline icon with default attributes",
			icon: &Icon{
				Name: "academic-cap",
				Body: `<path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/>`,
				Size: "24",
				Type: "Outline",
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/></svg>`,
		},
		{
			name: "Solid icon with default attributes",
			icon: &Icon{
				Name: "academic-cap-solid",
				Body: `<path d="M12 20a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`,
				Size: "24",
				Type: "Solid",
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 20a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
		},
		{
			name: "Mini icon with attributes",
			icon: &Icon{
				Name: "academic-cap-mini",
				Body: `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`,
				Size: "20",
				Type: "Mini",
				Attrs: templ.Attributes{
					"focusable": "false",
				},
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="currentColor" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`,
		},
		{
			name: "Micro icon with attributes",
			icon: &Icon{
				Name: "academic-cap-micro",
				Body: `<path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`,
				Size: "16",
				Type: "Micro",
				Attrs: templ.Attributes{
					"aria-hidden": "true",
					"class":       "icon-micro",
				},
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" class="icon-micro"><path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
		},
		{
			name: "Fallback case",
			icon: &Icon{
				Name: "unknown-icon",
				Body: `<circle cx="12" cy="12" r="10"/>`,
				Size: "24",
				Type: "Unknown", // Unsupported type
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg"><circle cx="12" cy="12" r="10"/></svg>`,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run test in parallel.
			result := strings.TrimSpace(tt.icon.String())
			expected := strings.TrimSpace(tt.expected)

			t.Logf("Generated SVG: %s", result)
			if result != expected {
				t.Errorf("String() = %q, want %q", result, expected)
			}
		})
	}
}

func TestIcon_SetAttrs(t *testing.T) {
	t.Parallel() // Run test in parallel.

	icon := &Icon{
		Name: "test-icon",
		Body: `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`,
		Size: "24",
		Type: "Outline",
	}

	attrs := templ.Attributes{
		"aria-hidden": "true",
		"custom-attr": "custom-val",
		"focusable":   "false",
	}

	icon.SetAttrs(attrs)

	if len(icon.Attrs) != len(attrs) {
		t.Errorf("expected %d attributes, got %d", len(attrs), len(icon.Attrs))
	}

	for key, expectedValue := range attrs {
		if value, exists := icon.Attrs[key]; !exists || value != expectedValue {
			t.Errorf("expected attribute %s=%s, got %s", key, expectedValue, value)
		}
	}

	expectedSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true" custom-attr="custom-val" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`
	if svg := icon.String(); svg != expectedSVG {
		t.Errorf("String() = %s, want %s", svg, expectedSVG)
	}
}

func TestAddAttributesToSVG(t *testing.T) {
	tests := []struct {
		name     string
		attrs    templ.Attributes
		expected string
	}{
		{
			name: "Non-reserved attributes are added",
			attrs: templ.Attributes{
				"aria-hidden": "false",
				"focusable":   "false",
			},
			expected: ` aria-hidden="false" focusable="false"`,
		},
		{
			name: "Reserved attributes are skipped",
			attrs: templ.Attributes{
				"xmlns":        "http://www.w3.org/2000/svg",
				"viewBox":      "0 0 24 24",
				"width":        "24",
				"height":       "24",
				"stroke-width": "1.5",
				"stroke":       "currentColor",
				"fill":         "none",
			},
			expected: "",
		},
		{
			name: "Mixed attributes: reserved are skipped, non-reserved are added",
			attrs: templ.Attributes{
				"xmlns":        "http://www.w3.org/2000/svg",
				"viewBox":      "0 0 24 24",
				"aria-hidden":  "true",
				"focusable":    "false",
				"stroke-width": "1.5",
			},
			expected: ` aria-hidden="true" focusable="false"`,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run test in parallel.

			var builder strings.Builder
			addAttributesToSVG(&builder, tt.attrs)

			result := builder.String()
			t.Logf("Generated attributes: %s", result)
			if result != tt.expected {
				t.Errorf("addAttributesToSVG() = %q, want %q", result, tt.expected)
			}
		})
	}
}

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
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" stroke="currentColor"><path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/></svg>`,
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
			name: "Micro icon with stroke and fill attributes",
			icon: &Icon{
				Name:        "academic-cap-micro",
				Body:        `<path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`,
				Size:        "16",
				Type:        "Micro",
				Stroke:      "#000000",
				StrokeWidth: "2",
				Fill:        "#FFFFFF",
				Attrs: templ.Attributes{
					"aria-hidden": "true",
					"class":       "icon-micro",
				},
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="#FFFFFF" aria-hidden="true" class="icon-micro"><path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
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
		{
			name: "SetSize modifies size",
			icon: func() *Icon {
				icon := &Icon{
					Name: "resizable-icon",
					Body: `<circle cx="12" cy="12" r="10"/>`,
					Size: "24",
					Type: "Outline",
				}
				icon.SetSize(32)
				return icon
			}(),
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none" stroke-width="1.5" stroke="currentColor"><circle cx="12" cy="12" r="10"/></svg>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(tt.icon.String())
			expected := strings.TrimSpace(tt.expected)

			if result != expected {
				t.Errorf("String() = %q, want %q", result, expected)
			}
		})
	}
}

func TestIcon_SetSize(t *testing.T) {
	tests := []struct {
		name     string
		initial  Size
		newSize  int
		expected Size
	}{
		{
			name:     "Set size to 24",
			initial:  "16",
			newSize:  24,
			expected: "24",
		},
		{
			name:     "Set size to 32",
			initial:  "24",
			newSize:  32,
			expected: "32",
		},
		{
			name:     "Set size to 48",
			initial:  "20",
			newSize:  48,
			expected: "48",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := &Icon{
				Size: tt.initial,
			}

			icon.SetSize(tt.newSize)

			if icon.Size != tt.expected {
				t.Errorf("SetSize() = %q, want %q", icon.Size, tt.expected)
			}
		})
	}
}

func TestIcon_Setters(t *testing.T) {
	icon := &Icon{
		Name: "test-icon",
		Body: `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`,
		Size: "24",
		Type: "Outline",
	}

	// Test setters
	icon.SetStroke("#FF0000").
		SetStrokeWidth("2").
		SetFill("#0000FF").
		SetSize(32)

	// Validate the individual fields
	if icon.Stroke != "#FF0000" {
		t.Errorf("SetStroke failed: expected #FF0000, got %s", icon.Stroke)
	}
	if icon.StrokeWidth != "2" {
		t.Errorf("SetStrokeWidth failed: expected 2, got %s", icon.StrokeWidth)
	}
	if icon.Fill != "#0000FF" {
		t.Errorf("SetFill failed: expected #0000FF, got %s", icon.Fill)
	}
	if icon.Size.String() != "32" {
		t.Errorf("SetSize failed: expected 32, got %s", icon.Size.String())
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

	expectedSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" stroke="currentColor" aria-hidden="true" custom-attr="custom-val" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`
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
		{
			name: "Non-string values are skipped",
			attrs: templ.Attributes{
				"aria-hidden": "true",
				"data-count":  123, // Non-string value
				"data-bool":   true,
			},
			expected: ` aria-hidden="true"`,
		},
		{
			name: "Safe onclick event is allowed",
			attrs: templ.Attributes{
				"aria-hidden": "true",
				"onclick":     "handleClick()", // Safe event handler
			},
			expected: ` aria-hidden="true" onclick="handleClick()"`,
		},
		{
			name: "Unsafe onclick event is skipped",
			attrs: templ.Attributes{
				"aria-hidden": "true",
				"onclick":     "javascript:alert('XSS')", // Unsafe value
			},
			expected: ` aria-hidden="true"`, // Unsafe "onclick" is excluded
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel tests.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run test in parallel.

			var builder strings.Builder
			addAttributesToSVG(&builder, tt.attrs)

			result := builder.String()
			if result != tt.expected {
				t.Errorf("addAttributesToSVG() = %q, want %q", result, tt.expected)
			}
		})
	}
}

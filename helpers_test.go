package templheroicons

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestHelpers_addAttributesToSVG(t *testing.T) {
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

package templheroicons

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"sync"
	"testing"

	"github.com/a-h/templ"
)

// 1. Core Tests for Icon Methods
// These tests cover methods like `String`, `SetSize`, and `SetAttrs`.

func TestIcon_String(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Icon
		expected string
	}{
		{
			name: "Outline icon with default attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name: "academic-cap",
					Size: "24",
					Type: "Outline",
				}
				icon.body = `<path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" stroke="currentColor"><path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/></svg>`,
		},
		{
			name: "Solid icon with default attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name: "academic-cap-solid",
					Size: "24",
					Type: "Solid",
				}
				icon.body = `<path d="M12 20a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 20a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
		},
		{
			name: "Mini icon with attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name: "academic-cap-mini",
					Size: "20",
					Type: "Mini",
					Attrs: templ.Attributes{
						"focusable": "false",
					},
				}
				icon.body = `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="currentColor" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`,
		},
		{
			name: "Micro icon with stroke and fill attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name:        "academic-cap-micro",
					Size:        "16",
					Type:        "Micro",
					Stroke:      "#000000",
					StrokeWidth: "2",
					Fill:        "#FFFFFF",
					Attrs: templ.Attributes{
						"aria-hidden": "true",
						"class":       "icon-micro",
					},
				}
				icon.body = `<path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="#FFFFFF" aria-hidden="true" class="icon-micro"><path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
		},
		{
			name: "Fallback case",
			setup: func() *Icon {
				icon := &Icon{
					Name: "unknown-icon",
					Size: "24",
					Type: "Unknown",
				}
				icon.body = `<circle cx="12" cy="12" r="10"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke-width="1.5" stroke="currentColor"><circle cx="12" cy="12" r="10"/></svg>`,
		},
		{
			name: "SetSize modifies size",
			setup: func() *Icon {
				originalIcon := &Icon{
					Name: "resizable-icon",
					Size: "24",
					Type: "Outline",
				}
				originalIcon.body = `<circle cx="12" cy="12" r="10"/>`
				// Capture the returned icon after setting size
				return originalIcon.SetSize(32)
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke-width="1.5" stroke="currentColor"><circle cx="12" cy="12" r="10"/></svg>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := tt.setup()
			result := strings.TrimSpace(icon.String())
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
			originalIcon := &Icon{
				Size: tt.initial,
			}

			clonedIcon := originalIcon.SetSize(tt.newSize)

			if clonedIcon.Size != tt.expected {
				t.Errorf("SetSize() = %q, want %q", clonedIcon.Size, tt.expected)
			}

			// Ensure original icon is unchanged
			if originalIcon.Size != tt.initial {
				t.Errorf("Original icon size modified: got %q, want %q", originalIcon.Size, tt.initial)
			}
		})
	}
}

func TestIcon_Setters(t *testing.T) {
	originalIcon := &Icon{
		Name: "test-icon",
		Size: "24",
		Type: "Outline",
	}

	// Chain the setters and capture the returned icon
	finalIcon := originalIcon.SetStroke("#FF0000").
		SetStrokeWidth("2").
		SetFill("#0000FF").
		SetSize(32)

	// Validate the individual fields on the returned icon
	if finalIcon.Stroke != "#FF0000" {
		t.Errorf("SetStroke failed: expected #FF0000, got %s", finalIcon.Stroke)
	}
	if finalIcon.StrokeWidth != "2" {
		t.Errorf("SetStrokeWidth failed: expected 2, got %s", finalIcon.StrokeWidth)
	}
	if finalIcon.Fill != "#0000FF" {
		t.Errorf("SetFill failed: expected #0000FF, got %s", finalIcon.Fill)
	}
	if finalIcon.Size.String() != "32" {
		t.Errorf("SetSize failed: expected 32, got %s", finalIcon.Size.String())
	}

	// Ensure the original icon remains unchanged
	if originalIcon.Size.String() != "24" {
		t.Errorf("Original icon size modified: expected 24, got %s", originalIcon.Size.String())
	}
	if originalIcon.Stroke != "" {
		t.Errorf("Original icon stroke modified: expected empty, got %s", originalIcon.Stroke)
	}
	if originalIcon.StrokeWidth != "" {
		t.Errorf("Original icon stroke-width modified: expected empty, got %s", originalIcon.StrokeWidth)
	}
	if originalIcon.Fill != "" {
		t.Errorf("Original icon fill modified: expected empty, got %s", originalIcon.Fill)
	}
}

func TestIcon_SetAttrs(t *testing.T) {
	t.Parallel() // Run test in parallel.

	originalIcon := &Icon{
		Name: "test-icon",
		Size: "24",
		Type: "Outline",
	}
	originalIcon.body = `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`

	attrs := templ.Attributes{
		"aria-hidden": "true",
		"custom-attr": "custom-val",
		"focusable":   "false",
	}

	// Capture the returned icon after setting attributes
	finalIcon := originalIcon.SetAttrs(attrs)

	if len(finalIcon.Attrs) != len(attrs) {
		t.Errorf("expected %d attributes, got %d", len(attrs), len(finalIcon.Attrs))
	}

	for key, expectedValue := range attrs {
		if value, exists := finalIcon.Attrs[key]; !exists || value != expectedValue {
			t.Errorf("expected attribute %s=%s, got %s", key, expectedValue, value)
		}
	}

	expectedSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" stroke="currentColor" aria-hidden="true" custom-attr="custom-val" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`
	if svg := finalIcon.String(); svg != expectedSVG {
		t.Errorf("String() = %s, want %s", svg, expectedSVG)
	}

	// Ensure the original icon remains unchanged
	if len(originalIcon.Attrs) != 0 {
		t.Errorf("Original icon attributes modified: expected 0, got %d", len(originalIcon.Attrs))
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

// 2. Tests for JSON-Based Functionality
// These tests cover JSON parsing, caching, and error handling.

func TestGetIconBody_RealData(t *testing.T) {
	tests := []struct {
		name           string
		iconName       string
		expectedBody   string
		expectingError bool
	}{
		{
			name:           "Retrieve existing icon",
			iconName:       "academic-cap",
			expectedBody:   `<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347m-15.482 0a51 51 0 0 0-2.658-.813A60 60 0 0 1 12 3.493a60 60 0 0 1 10.399 5.84q-1.345.372-2.658.814m-15.482 0A51 51 0 0 1 12 13.489a50.7 50.7 0 0 1 7.74-3.342M6.75 15a.75.75 0 1 0 0-1.5a.75.75 0 0 0 0 1.5m0 0v-3.675A55 55 0 0 1 12 8.443m-7.007 11.55A5.98 5.98 0 0 0 6.75 15.75v-1.5"/>`,
			expectingError: false,
		},
		{
			name:           "Retrieve another existing icon",
			iconName:       "academic-cap-solid",
			expectedBody:   `<g fill="currentColor"><path d="M11.7 2.805a.75.75 0 0 1 .6 0A60.7 60.7 0 0 1 22.83 8.72a.75.75 0 0 1-.231 1.337a50 50 0 0 0-9.902 3.912l-.003.002l-.34.18a.75.75 0 0 1-.707 0A51 51 0 0 0 7.5 12.173v-.224a.36.36 0 0 1 .172-.311a55 55 0 0 1 4.653-2.52a.75.75 0 0 0-.65-1.352a56 56 0 0 0-4.78 2.589a1.86 1.86 0 0 0-.859 1.228a50 50 0 0 0-4.634-1.527a.75.75 0 0 1-.231-1.337A60.7 60.7 0 0 1 11.7 2.805"/><path d="M13.06 15.473a48.5 48.5 0 0 1 7.666-3.282q.202 2.122.255 4.284a.75.75 0 0 1-.46.711a48 48 0 0 0-8.105 4.342a.75.75 0 0 1-.832 0a48 48 0 0 0-8.104-4.342a.75.75 0 0 1-.461-.71q.053-2.163.255-4.286q1.382.456 2.726.99v1.27a1.5 1.5 0 0 0-.14 2.508c-.09.38-.222.753-.397 1.11q.678.32 1.346.66a6.7 6.7 0 0 0 .551-1.607a1.5 1.5 0 0 0 .14-2.67v-.645a49 49 0 0 1 3.44 1.667a2.25 2.25 0 0 0 2.12 0"/><path d="M4.462 19.462c.42-.419.753-.89 1-1.395q.68.321 1.347.662a6.7 6.7 0 0 1-1.286 1.794a.75.75 0 0 1-1.06-1.06"/></g>`,
			expectingError: false,
		},
		{
			name:           "Icon not found",
			iconName:       "non-existing-icon",
			expectedBody:   "",
			expectingError: true,
		},
		{
			name:           "Empty icon name",
			iconName:       "",
			expectedBody:   "",
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := getIconBody(tt.iconName)

			if tt.expectingError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if body != tt.expectedBody {
					t.Errorf("getIconBody() = %q, want %q", body, tt.expectedBody)
				}
			}
		})
	}
}

func TestGetIconBody_OnceWithRealData(t *testing.T) {
	// First call should initialize the data
	_, err := getIconBody("academic-cap")
	if err != nil {
		t.Fatalf("unexpected error during first call: %v", err)
	}

	// Ensure no error on subsequent calls for valid icons
	_, err = getIconBody("academic-cap-solid")
	if err != nil {
		t.Fatalf("unexpected error during subsequent call: %v", err)
	}
}

// 3. Tests for Mocked Data
// These tests cover cases where mocked FS and invalid JSON are used.

func TestIcon_String_FetchBody(t *testing.T) {
	// Reset `iconDataOnce` to ensure fresh parsing during the test
	resetTestState()

	// Mock the embedded JSON with valid data
	validJSON := `{
        "icons": {
            "academic-cap": { "body": "<path fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347m-15.482 0a51 51 0 0 0-2.658-.813A60 60 0 0 1 12 3.493a60 60 0 0 1 10.399 5.84q-1.345.372-2.658.814m-15.482 0A51 51 0 0 1 12 13.489a50.7 50.7 0 0 1 7.74-3.342M6.75 15a.75.75 0 1 0 0-1.5a.75.75 0 0 0 0 1.5m0 0v-3.675A55 55 0 0 1 12 8.443m-7.007 11.55A5.98 5.98 0 0 0 6.75 15.75v-1.5\"/>" }
        }
    }`
	heroiconsJSONSource = mockInvalidJSONFS(validJSON)
	defer func() {
		heroiconsJSONSource = heroiconsJSON // Restore original embedded JSON
	}()

	t.Run("Fetches and caches body", func(t *testing.T) {
		icon := &Icon{
			Name: "academic-cap",
			Size: "24",
			Type: "Outline",
		}

		// Call String() for the first time to trigger the body fetch
		result := icon.String()

		// Validate the resulting SVG
		expected := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" stroke="currentColor"><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347m-15.482 0a51 51 0 0 0-2.658-.813A60 60 0 0 1 12 3.493a60 60 0 0 1 10.399 5.84q-1.345.372-2.658.814m-15.482 0A51 51 0 0 1 12 13.489a50.7 50.7 0 0 1 7.74-3.342M6.75 15a.75.75 0 1 0 0-1.5a.75.75 0 0 0 0 1.5m0 0v-3.675A55 55 0 0 1 12 8.443m-7.007 11.55A5.98 5.98 0 0 0 6.75 15.75v-1.5"/></svg>`
		if result != expected {
			t.Errorf("String() = %q, want %q", result, expected)
		}

		// Validate that the body is cached
		if icon.body != `<path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347m-15.482 0a51 51 0 0 0-2.658-.813A60 60 0 0 1 12 3.493a60 60 0 0 1 10.399 5.84q-1.345.372-2.658.814m-15.482 0A51 51 0 0 1 12 13.489a50.7 50.7 0 0 1 7.74-3.342M6.75 15a.75.75 0 1 0 0-1.5a.75.75 0 0 0 0 1.5m0 0v-3.675A55 55 0 0 1 12 8.443m-7.007 11.55A5.98 5.98 0 0 0 6.75 15.75v-1.5"/>` {
			t.Errorf("Body was not cached correctly, got %q", icon.body)
		}
	})

	t.Run("Handles error when fetching body", func(t *testing.T) {
		icon := &Icon{
			Name: "non-existent-icon",
			Size: "24",
			Type: "Outline",
		}

		// Call String() for a non-existent icon
		result := icon.String()

		// Validate the error message in the SVG output
		if !strings.Contains(result, "Error: icon 'non-existent-icon' not found") {
			t.Errorf("Expected error message in output, got %q", result)
		}
	})
}

func TestGetIconBody_JSONParsing(t *testing.T) {
	tests := []struct {
		name           string
		mockJSON       string
		iconName       string
		expectedError  string
		expectedResult string
	}{
		{
			name:          "Invalid JSON format",
			mockJSON:      `{"icons": "invalid"}`, // Invalid JSON structure
			iconName:      "academic-cap",
			expectedError: "failed to parse heroicons JSON",
		},
		{
			name:          "Missing icons field",
			mockJSON:      `{"missingIcons": {}}`, // No `icons` key
			iconName:      "academic",
			expectedError: "icon 'academic' not found",
		},
		{
			name:           "Valid JSON",
			mockJSON:       `{"icons": {"academic-cap": {"body": "<path d='...'/>"}}}`,
			iconName:       "academic-cap",
			expectedError:  "",
			expectedResult: "<path d='...'/>",
		},
		{
			name:          "Icon not found",
			mockJSON:      `{"icons": {"academic-cap": {"body": "<path d='...'/>"}}}`,
			iconName:      "non-existent-icon",
			expectedError: "icon 'non-existent-icon' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetTestState()

			// Replace heroiconsJSONSource with a mocked FS
			heroiconsJSONSource = mockInvalidJSONFS(tt.mockJSON)
			defer func() {
				heroiconsJSONSource = heroiconsJSON // Restore original embedded FS
			}()

			result, err := getIconBody(tt.iconName)

			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error %q, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected result %q, got %q", tt.expectedResult, result)
				}
			}
		})
	}
}

// 4. Utility Functions for Testing
// These utilities mock data and manage state resets.

type mockFS struct {
	data map[string]string
}

func mockInvalidJSONFS(data string) fs.FS {
	return &mockFS{
		data: map[string]string{"data/heroicons_cache.json": data},
	}
}

func (m *mockFS) Open(name string) (fs.File, error) {
	content, exists := m.data[name]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", name)
	}
	return &mockFile{content: strings.NewReader(content)}, nil
}

type mockFile struct {
	content io.Reader
}

func (f *mockFile) Read(p []byte) (int, error) {
	return f.content.Read(p)
}

func (f *mockFile) Close() error {
	return nil
}

func (f *mockFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func resetTestState() {
	iconDataOnce = sync.Once{}
	iconData = nil
}

func TestMockFS(t *testing.T) {
	data := `{"icons": invalid}`
	mockFS := mockInvalidJSONFS(data)
	content, err := fs.ReadFile(mockFS, "data/heroicons_cache.json")
	if err != nil {
		t.Fatalf("Failed to read mock file: %v", err)
	}
	if string(content) != data {
		t.Fatalf("Expected mock content %q, got %q", data, string(content))
	}
}

package templheroicons

import (
	_ "embed"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/tidwall/gjson"
)

// Cache to store parsed icon body content for reuse
var (
	iconBodyCache = map[string]string{}
	cacheMutex    sync.Mutex
)

// Size represents the size of UI components.
type Size string

// String returns the string representation of a Size.
func (s Size) String() string {
	return string(s)
}

// Icon represents a single icon with its attributes.
type Icon struct {
	Name  string           `json:"name"` // Name of the icon (e.g., "moon")
	Type  string           `json:"type"` // Type of the icon (e.g., "Outline", "Solid")
	Size  Size             `json:"size"` // Size of the icon (e.g., "24", "48")
	Color string           // Optional color for the icon's fill
	Attrs templ.Attributes // Custom attributes to be added to the <svg> tag
	body  string           // Cached body of the icon's SVG path (immutable)
}

// Render generates the complete SVG tag for the icon.
func (i *Icon) Render() templ.Component {
	return templ.Raw(makeSVGTag(i))
}

// IconBuilder is a builder for configuring an Icon.
// It allows method chaining to update the icon's properties.
type IconBuilder struct {
	icon *Icon // Reference to the icon being configured
}

// Config returns an IconBuilder to allow chaining configuration methods on the icon.
func (icon *Icon) Config() *IconBuilder {
	return &IconBuilder{
		icon: icon.clone(), // Clone the icon to ensure immutability
	}
}

// ConfigureIcon creates a new builder from an existing icon.
func ConfigureIcon(icon *Icon) *IconBuilder {
	return &IconBuilder{
		icon: icon.clone(), // Clone the icon to ensure immutability
	}
}

// SetSize sets the size of the icon.
func (b *IconBuilder) SetSize(size int) *IconBuilder {
	b.icon.Size = Size(strconv.Itoa(size))
	return b
}

// SetColor sets the fill color of the icon.
func (b *IconBuilder) SetColor(value string) *IconBuilder {
	b.icon.Color = value
	return b
}

// SetAttrs sets custom attributes for the SVG tag (e.g., `aria-hidden`, `focusable`).
func (b *IconBuilder) SetAttrs(attrs templ.Attributes) *IconBuilder {
	b.icon.Attrs = attrs
	return b
}

// GetIcon returns the configured icon instance.
func (b *IconBuilder) GetIcon() *Icon {
	return b.icon
}

// Render generates the SVG for the configured icon.
func (b *IconBuilder) Render() templ.Component {
	return b.icon.Render()
}

// clone creates a deep copy of the Icon to prevent shared state.
func (i *Icon) clone() *Icon {
	// Deep copy the attributes to prevent shared references
	attrsCopy := make(templ.Attributes, len(i.Attrs))
	for k, v := range i.Attrs {
		attrsCopy[k] = v
	}
	return &Icon{
		Name:  i.Name,
		Type:  i.Type,
		Size:  i.Size,
		Color: i.Color,
		Attrs: attrsCopy, // Use the deep copy of the attributes
		body:  i.body,    // The body is shared since it's immutable
	}
}

// fetchBody ensures that the body of the icon is loaded from the cache or file.
func (i *Icon) fetchBody() error {
	if i.body != "" {
		return nil // Body is already cached
	}

	body, err := getIconBody(i.Name)
	if err != nil {
		return err
	}

	i.body = body
	return nil
}

// makeSVGTag generates the full SVG tag for the icon.
func makeSVGTag(icon *Icon) string {
	// Ensure the body is loaded before rendering
	if err := icon.fetchBody(); err != nil {
		return errorSVGComment(err)
	}

	// Determine the appropriate viewBox and type-based attributes
	viewBox := getViewBoxDimensions(icon.Type)
	typeAttributes := getTypeAttributes(icon.Type)

	var builder strings.Builder
	// Construct the opening <svg> tag with common attributes
	fmt.Fprintf(&builder, `<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 %[2]s %[2]s"%s`,
		icon.Size.String(),
		viewBox,
		typeAttributes,
	)

	// If a custom color is set, add it to the <svg> tag
	if icon.Color != "" {
		fmt.Fprintf(&builder, ` color="%s"`, icon.Color)
	}

	// Add user-defined attributes to the <svg> tag
	addAttributesToSVG(&builder, icon.Attrs)

	// Close the opening <svg> tag, add the body, and close the <svg> tag
	builder.WriteString(">")
	builder.WriteString(icon.body)
	builder.WriteString(`</svg>`)

	return builder.String()
}

// getIconBody retrieves the body of an icon by its name, with thread-safe caching.
var getIconBody = func(name string) (string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the body is already cached
	if body, found := iconBodyCache[name]; found {
		return body, nil
	}

	// Read and parse the JSON file containing icon data
	jsonFilename := "data/heroicons_cache.json"
	heroiconsData, _ := heroiconsJSONSource.Open(jsonFilename)
	defer heroiconsData.Close()

	data, _ := io.ReadAll(heroiconsData)

	// Check if the JSON data is valid
	if !gjson.ValidBytes(data) {
		return "", fmt.Errorf("failed to parse heroicons JSON")
	}

	// Extract the "icons" key from the JSON data
	iconsResult := gjson.GetBytes(data, "icons")

	// If the "icons" key exists, populate the cache
	if iconsResult.Exists() {
		iconsResult.ForEach(func(key, value gjson.Result) bool {
			iconBody := value.Get("body").String()
			iconBodyCache[key.String()] = iconBody
			return true
		})
	}

	// Return the requested icon body from the cache
	body, exists := iconBodyCache[name]
	if !exists {
		return "", fmt.Errorf("icon '%s' not found", name)
	}
	return body, nil
}

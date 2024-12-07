package templheroicons

import (
	"strconv"

	"github.com/a-h/templ"
)

type IconBuilder struct {
	icon Icon
}

// ConfigureIcon creates a new builder from an existing icon.
func ConfigureIcon(icon *Icon) *IconBuilder {
	return &IconBuilder{
		icon: *icon, // Copy the icon (not reference)
	}
}

// SetSize sets the size of the icon.
func (b *IconBuilder) SetSize(size int) *IconBuilder {
	b.icon.Size = Size(strconv.Itoa(size))
	return b
}

// SetStroke sets the stroke of the icon.
func (b *IconBuilder) SetStroke(value string) *IconBuilder {
	b.icon.Stroke = value
	return b
}

// SetStrokeWidth sets the stroke-width of the icon.
func (b *IconBuilder) SetStrokeWidth(value string) *IconBuilder {
	b.icon.StrokeWidth = value
	return b
}

// SetFill sets the fill of the icon.
func (b *IconBuilder) SetFill(value string) *IconBuilder {
	b.icon.Fill = value
	return b
}

// SetAttrs sets the attributes for the SVG tag.
func (b *IconBuilder) SetAttrs(attrs templ.Attributes) *IconBuilder {
	b.icon.Attrs = attrs
	return b
}

// Build returns the final configured Icon.
func (b *IconBuilder) Build() *Icon {
	return &b.icon
}

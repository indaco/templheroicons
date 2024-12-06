package templheroicons

import (
	"embed"
	"io/fs"
)

//go:embed data/heroicons_cache.json
var heroiconsJSON embed.FS

var heroiconsJSONSource fs.FS = heroiconsJSON // Default to the embedded FS

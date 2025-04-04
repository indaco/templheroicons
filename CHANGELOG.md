# Changelog

All notable changes to this project will be documented in this file.

The format adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html),
and is generated by [changelogen](https://github.com/unjs/changelogen) and managed with [Changie](https://github.com/miniscruff/changie).

## v0.6.0 - 2025-01-13

[compare changes](https://github.com/indaco/templheroicons/compare/v0.5.0...v0.6.0)

> [!IMPORTANT]
> Starting with [v0.3.819](https://github.com/a-h/templ/releases/tag/v0.3.819), `templ` requires Go 1.23. If you are using Go 1.22, you will need to remain on **v0.5.0**.

### 📦 Build

- Bump `templ` to `v0.3.819` ([2c56d55](https://github.com/indaco/templheroicons/commit/2c56d55))

### 🏡 Chore

- **demos:** Regenerate with templ v0.3.819 ([3a913d3](https://github.com/indaco/templheroicons/commit/3a913d3))

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

## v0.5.0 - 2024-12-14

[compare changes](https://github.com/indaco/templheroicons/compare/v0.4.0...v0.5.0)

### 💅 Refactors

- Enhance API with Config method for cleaner builder pattern ([91fa581](https://github.com/indaco/templheroicons/commit/91fa581))

### ❤️ Contributors

- Indaco <github@mircoveltri.me>

## v0.4.0 - 2024-12-13

[compare changes](https://github.com/indaco/templheroicons/compare/v0.3.0...v0.4.0)

### 🚀 Enhancements

- Simplify builder pattern by removing Build method and adding Render wrapper ([b9ccd20](https://github.com/indaco/templheroicons/commit/b9ccd20))

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

## v0.3.0

[compare changes](https://github.com/indaco/templheroicons/compare/v0.2.0...v0.3.0)

### 🚀 Enhancements

- Implement lazy loading for icon SVG bodies ([1819564](https://github.com/indaco/templheroicons/commit/1819564))
- **icons-maker.go:** Force dataset re-fetch when missing icons, empty json file ([cfa2ab1](https://github.com/indaco/templheroicons/commit/cfa2ab1))

### 🔥 Performance

- Use gjson for efficient icon body extraction and caching ([2edc5bb](https://github.com/indaco/templheroicons/commit/2edc5bb))

### 🩹 Fixes

- **go.mod:** Use the 1.N.P syntax for go toolchain versions ([240927a](https://github.com/indaco/templheroicons/commit/240927a))
- DRAFT - stroke, stroke width and fill on parent svg ([1350ba0](https://github.com/indaco/templheroicons/commit/1350ba0))

### 💅 Refactors

- Better icon body caching ([2b63410](https://github.com/indaco/templheroicons/commit/2b63410))
- Icon configuration with ConfigureIcon. README.md updated ([f48ec24](https://github.com/indaco/templheroicons/commit/f48ec24))
- Use color to set the fill color for the icons ([c3c2580](https://github.com/indaco/templheroicons/commit/c3c2580))

### 📦 Build

- Add scripts to devbox.json ([71d958d](https://github.com/indaco/templheroicons/commit/71d958d))
- Move icons-maker.go to cmd folder ([9bc66c4](https://github.com/indaco/templheroicons/commit/9bc66c4))
- Add demo target to the Taskfile and Makefile ([951defc](https://github.com/indaco/templheroicons/commit/951defc))

### 🏡 Chore

- Add basic demo ([cc943f0](https://github.com/indaco/templheroicons/commit/cc943f0))
- **demos:** Update ([a2fef24](https://github.com/indaco/templheroicons/commit/a2fef24))

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

## v0.2.0 - 2024-12-04

[compare changes](https://github.com/indaco/templheroicons/compare/v0.1.2...v0.2.0)

### 🚀 Enhancements

- Allow users to set fill, stroke, and stroke-width attrs for icons ([6235412](https://github.com/indaco/templheroicons/commit/6235412))
- Add SetSize method to allow users to set custom icon size ([4976300](https://github.com/indaco/templheroicons/commit/4976300))

### 🩹 Fixes

- Typo in datasetURL ([87a42b3](https://github.com/indaco/templheroicons/commit/87a42b3))

### 📦 Build

- Update devbox.json file ([070881d](https://github.com/indaco/templheroicons/commit/070881d))

### 🏡 Chore

- Reorder struct props ([e997702](https://github.com/indaco/templheroicons/commit/e997702))

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

## v0.1.2 - 2024-12-04

[compare changes](https://github.com/indaco/templheroicons/compare/v0.1.1...v0.1.2)

### 🩹 Fixes

- Add sanitization to `addAttributesToSVG` for enhanced security ([60d12bd](https://github.com/indaco/templheroicons/commit/60d12bd))

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

## v0.1.1 - 2024-12-04

[compare changes](https://github.com/indaco/templheroicons/compare/v0.1.0...v0.1.1)

### 🩹 Fixes

- Typo in `Render` func name ([0ad04e4](https://github.com/indaco/templheroicons/commit/0ad04e4))

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

## v0.1.0 - 2024-12-04

### 🏡 Chore

- Initial Release

### ❤️ Contributors

- Indaco ([@indaco](http://github.com/indaco))

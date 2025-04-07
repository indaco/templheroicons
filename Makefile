color_red     := $(shell printf "\e[31m")
color_green   := $(shell printf "\e[32m")
color_yellow  := $(shell printf "\e[33m")
color_blue    := $(shell printf "\e[34m")
color_magenta := $(shell printf "\e[35m")
color_cyan    := $(shell printf "\e[36m")

# Bold variants
color_bold_red     := $(shell printf "\e[1;31m")
color_bold_green   := $(shell printf "\e[1;32m")
color_bold_yellow  := $(shell printf "\e[1;33m")
color_bold_blue    := $(shell printf "\e[1;34m")
color_bold_magenta := $(shell printf "\e[1;35m")
color_bold_cyan    := $(shell printf "\e[1;36m")
color_reset        := $(shell printf "\e[0m")

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
.PHONY: help
help: ## Print this help message
	@echo ""
	@echo "Usage: make [action]"
	@echo ""
	@echo "Available Actions:"
	@echo ""
	@awk -v green="$(color_green)" -v reset="$(color_reset)" -F ':|##' \
		'/^[^\t].+?:.*?##/ {printf " %s* %-15s%s %s\n", green, $$1, reset, $$NF}' $(MAKEFILE_LIST) | sort
	@echo ""

# ==================================================================================== #
# PRIVATE TASKS
# ==================================================================================== #

templ: ## Run templ fmt and templ generate (private)
	@echo "$(color_cyan)Running templ fmt and templ generate in ./_demos/pages/$(color_reset)"
	@cd ./_demos/pages/ && templ fmt . && templ generate

# ==================================================================================== #
# PUBLIC TASKS
# ==================================================================================== #

test: ## Run go tests
	@go test -race -covermode=atomic .

test/coverage: ## Run go tests and use go tool cover
	@go test -coverprofile=coverage.txt .
	@go tool cover -html=coverage.txt

build: ## Generate the Go icon definitions based on parsed data/heroicons_cache.json file.
	@cd cmd && go run icons-maker.go

demo: templ ## Run the demo server
	@echo "$(color_cyan)Running the demo server in ./_demos/$(color_reset)"
	@cd ./_demos/ && go run main.go
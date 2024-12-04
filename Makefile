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
# PUBLIC TASKS
# ==================================================================================== #

test: ## Run go tests
	@go test -race -covermode=atomic .

test/coverage: ## Run go tests and use go tool cover
	@go test -coverprofile=coverage.out .
	@go tool cover -html=coverage.out

build:
	@cd _scripts && go run icons-maker.go
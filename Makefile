##:
# ============================================================================ #
# This is a multi-platform Makefile, trying to support both Unix-like
# and Windows operating systems.
#
# Every comment line that starts with two '#' is parsed by the 'help' target
# as part of the help message:
#  - use a single ':' to output an empty line;
#  - use '<target>:<description>' for a single-line description;
#  - use '<target>:<description>' with the following ':<description>'
#    for a multiline description.
#
# Note that to to split a PowerShell command line over multiple lines
# use a comment block with a backslash inside:
#   do things <#\
#   #> do other things
# ============================================================================ #

ifeq ($(OS),Windows_NT)
    SHELL = pwsh.exe
    LIST_SEP = ;
else
    SHELL = /bin/sh
    LIST_SEP = :
endif

# Explicitly say what the target is default; change it as necessary.
.DEFAULT_GOAL := help

## help: print this help message and exit
.PHONY: help
help:
ifeq ($(OS),Windows_NT)
	@ Write-Host "Usage:" -NoNewline
    # Hack: replace two '#' with the NULL character to force ConvertFrom-Csv
    # to print empty lines.
	@ (Get-Content $(MAKEFILE_LIST)) -match "^##" -replace "^##","$$([char]0x0)" <#\
 #> | ConvertFrom-Csv -Delimiter ":" -Header Target,Description <#\
 #> | Format-Table <#\
 #>     -AutoSize -HideTableHeaders <#\
 #>     -Property @{Expression=" "},Target,@{Expression=" "},Description
else
	@ echo 'Usage:'
	@ sed --quiet 's/^##//p' $(MAKEFILE_LIST) \
    | column --table --separator ':' \
    | sed --expression='s/^/ /' \
    && echo
endif

## all: run audit and tests
.PHONY: all
all: audit test ;

# ============================================================================ #
##:
##: Variables
##:
# These variables can be changed here directly by editing this file
# or by passing them into the `make` call,
# e.g., `make <variable_1>=<value_1> <variable_2>=<value_2> [...]`.
# ============================================================================ #

## BINARY_DIR: get the directory with binaries
BINARY_DIR = bin
__VARIABLES__ += BINARY_DIR

# Directory containing the Makefile.
PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

export GOBIN ?= $(PROJECT_ROOT)$(BINARY_DIR)
export PATH := $(GOBIN)$(LIST_SEP)$(PATH)

# Generate getters for all variables.
define make_get_variable
.PHONY: $1
$1:
	@ echo "$($1)"
endef
$(foreach var,$(__VARIABLES__), \
    $(eval \
        $(call make_get_variable,$(var)) \
    ) \
)

# ============================================================================ #
# Helpers
# ============================================================================ #

.PHONY: confirm
confirm:
ifeq ($(OS),Windows_NT)
	@ if ((Read-Host -Prompt "Are you sure? [y/N]") -cne "y") { throw "" }
else
	@ read -r -p 'Are you sure? [y/N] ' answer \
    && [ $${answer:-N} = 'y' ]
endif

.PHONY: git/no-dirty
git/no-dirty:
ifeq ($(OS),Windows_NT)
	@ if (![string]::IsNullOrEmpty("$(shell git status --porcelain)")) { throw "" }
else
	@ test -z "$(shell git status --porcelain)"
endif

.PHONY: create/binary_dir
create/binary_dir:
ifeq ($(OS),Windows_NT)
	@ [void](New-Item "$(BINARY_DIR)" -ItemType Directory -Force)
else
	@ mkdir -p "$(BINARY_DIR)"
endif

.PHONY: cgo/enable
cgo/enable:
	@ go env -w CGO_ENABLED=1

.PHONY: cgo/disable
cgo/disable:
	@ go env -w CGO_ENABLED=0

# ============================================================================ #
##:
##: Quality control
##:
# ============================================================================ #

## audit: run quality control checks
.PHONY: audit
audit: fmt/no-dirty mod/tidy-diff mod/verify govulncheck golangci-lint ;

## mod/tidy-diff: check missing and unused modules without modifying
##              : the `go.mod` and `go.sum` files
.PHONY: mod/tidy-diff
mod/tidy-diff:
	@ go mod tidy -diff

## mod/tidy: add missing and remove unused modules
.PHONY: mod/tidy
mod/tidy:
	@ go mod tidy -v

## mod/verify: verify that dependencies have expected content
.PHONY: mod/verify
mod/verify:
	@ go mod verify

.PHONY: install/govulncheck
install/govulncheck:
	@ cd tools && go get golang.org/x/vuln/cmd/govulncheck && go install golang.org/x/vuln/cmd/govulncheck

## govulncheck: report known vulnerabilities that affect Go code
.PHONY: govulncheck
govulncheck: install/govulncheck
	@ $(GOBIN)/govulncheck ./...

## fmt: gofmt (reformat) package sources
.PHONY: fmt
fmt:
	@ go fmt ./...

## fmt/no-dirty: gofmt (reformat) package sources and fail if there are some
##             : changes
.PHONY: fmt/no-dirty
fmt/no-dirty:
ifeq ($(OS),Windows_NT)
	@ if (![string]::IsNullOrEmpty("$(shell go fmt ./...)")) { throw "" }
else
	@ test -z "$(shell go fmt ./...)"
endif

## vet: report likely mistakes in packages
.PHONY: vet
vet:
	@ go vet ./...

## golangci-lint: a fast linters runner for Go
.PHONY: golangci-lint
golangci-lint:
	@ golangci-lint run ./...

## test: run all the tests
.PHONY: test
test: cgo/enable
	@ go test -v -race ./...

## test/cover: run all the tests and display coverage
.PHONY: test/cover
test/cover: create/binary_dir cgo/enable
	@ go test -v -race -coverpkg=./... -coverprofile='$(BINARY_DIR)/coverage.out' ./...
	@ go tool cover -html='$(BINARY_DIR)/coverage.out' -o '$(BINARY_DIR)/coverage.html'

# ============================================================================ #
##:
##: Build
##:
# ============================================================================ #

## mod/download: download modules to local cache
.PHONY: mod/download
mod/download:
	@ go mod download -x

## clean: remove files from the binary directory
.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@ if (Test-Path "$(BINARY_DIR)" -PathType Container) { <#\
 #>     Remove-Item "$(BINARY_DIR)\*" -Recurse -Force <#\
 #> }
else
	@ rm -rf $(BINARY_DIR)/*
endif

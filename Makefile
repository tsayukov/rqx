# This file is licensed under the terms of the MIT License (see LICENSE file)
# Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

# ============================================================================ #
# This is a multi-platform Makefile, trying to support both Unix-like
# and Windows operating systems.
#
# It is very much inspired by the Makefile made by Alex Edwards. Check that out:
#  - https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects
#  - https://gist.github.com/alexedwards/3b40775846535d0014ab1ff477e4a568
#
#                    HOW TO WRITE DOCUMENTATION FOR TARGETS
#
# Every comment line that starts with two '#' is parsed by the 'help' target
# as part of the help message:
#  - use '##:' to output an empty line;
#  - use '##<target>:<description>' for a single-line description;
#  - use '##<target>:<description>' with the following '##:<description>',
#    each on the next line, for a multiline description.
#
# Whitespaces between '##', '<target>', ':', and '<description>' do not matter.
#
# A standalone '##:<description>' with the surrounding '##:' at the top
# and bottom can be used to write a header, e.g., 'Variables', 'Build', etc.
#
#                          OPERATING SYSTEM DETECTION
#
# To detect the operating system:
#
# 1. Check whether the environment variable PATH contains the path separator:
#  - in Windows:
    override __semicolon__ := ;
#  - in a Unix-like operating system:
    override __colon__     := :
# Unlikely, but the PATH might contain only one element or a path with
# a semicolon in a Unix-like operating system. It can also be passed
# via the command line to prepend it with the specific user path.
#
# 2. Check the environment variable OS that holds the string 'Windows_NT'
# on the Windows NT family. However, the OS variable can be overwritten.
#
# 3. Call the `uname` command to get the name of the current operating system.
    override __sh_uname_or_unknown__ := sh -c 'uname 2>/dev/null || echo Unknown'
# Transform the Cygwin/MSYS verbose output of `uname` to 'Cygwin'/'MSYS'.
# See: https://stackoverflow.com/a/52062069/10537247
    override __shorten_os_name__ = \
        $(patsubst CYGWIN%,Cygwin,\
            $(patsubst MSYS%,MSYS,\
                $(patsubst MINGW%,MSYS,\
                    $1\
                )\
            )\
        )
#
# Nevertheless, it is unlikely that these pitfalls will occur in most cases.
# Otherwise, pass the OS variable into the `make` call to specify the current
# operating system:
#   make OS=<operating system name> [...]
#
    ifeq ($(origin OS),command line)
        override __OS__ := $(OS)
    else ifeq ($(OS),Windows_NT)
        # Distinguish between native Windows and Cygwin/MSYS.
        ifneq (,$(findstring $(__semicolon__),$(PATH))) # if semicolon is in PATH
            override __OS__ := Windows
        else
            override __OS__ := $(call __shorten_os_name__,$(shell $(__sh_uname_or_unknown__)))
        endif
    else
        override __OS__ := $(shell $(__sh_uname_or_unknown__))
    endif
#
# Use the __OS__ variable to match the detected operating system against
# Windows, Linux, Darwin, etc.
#
#                                 CONFIGURATION
#
# Explicitly say what the target is default; change it as necessary.
    .DEFAULT_GOAL := help
#
# Choosing the appropriate shell, path separator, and list separator.
    ifeq ($(__OS__),Unknown)
        $(error unknown operating system)
    else ifeq ($(__OS__),Windows)
        SHELL := pwsh.exe
        override __PATH_SEP__ := \\
        override __LIST_SEP__ := $(__semicolon__)
    else
        SHELL := /bin/sh
        override __PATH_SEP__ := /
        override __LIST_SEP__ := $(__colon__)
    endif
#
# The project root containing the Makefile.
    override __PROJECT_ROOT__ := $(subst /,$(__PATH_SEP__),$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
#
#                                     TIPS
#
# 1. To split a PowerShell command line over multiple lines use a comment block
# with a backslash inside:
#   do things <#\
#   #> do other things
#
# 2. Colorful output:
#   $(call __color_text__,$(__BLUE__),o ) && echo "Starting checks..."
#   $(call __color_text__,$(__GREEN__),v ) && echo "OK"
#   $(call __color_text__,$(__RED__),x ) && echo "Error"
#
    ifeq ($(__OS__),Windows)
        override __RED__   := Red
        override __GREEN__ := Green
        override __BLUE__  := Blue
        override __color_text__ = Write-Host "$2" -ForegroundColor $1 -NoNewline
    else
        override __RED__   := \\033[0;31m
        override __GREEN__ := \\033[0;32m
        override __BLUE__  := \\033[0;34m
        override __color_text__ = printf "%b%s%b" "$1" "$2" "\033[0m"
    endif
# ============================================================================ #

# The blank line below is necessary to get the same help message on different
# operating systems.
##:
## help: print this help message and exit
.PHONY: help
help:
	@ $(info )
	@ $(info :: Go Makefile)
	@ $(info :: OS: $(__OS__))
	@ $(info :: SHELL: $(SHELL))
	@ $(info )
ifeq ($(__OS__),Windows)
	@ Write-Host "Targets:" -NoNewline
    # Hack: replace two '#' with the NULL character to force ConvertFrom-Csv
    # to print empty lines.
	@ (Get-Content $(MAKEFILE_LIST)) -match "^##" -replace "^##","$$([char]0x0)" <#\
 #> | ConvertFrom-Csv -Delimiter ":" -Header Target,Description <#\
 #> | Format-Table <#\
 #>     -AutoSize -HideTableHeaders <#\
 #>     -Property @{Expression=" "},Target,@{Expression=" "},Description
else
	@ echo 'Targets:'
	@ sed --quiet 's/^##//p' $(MAKEFILE_LIST) \
	| sed --expression='s/[ \t]*:[ \t]*/:/' \
    | column --table --separator ':' \
    | sed --expression='s/^/ /' \
    && echo
endif

## all: run audit and tests
.PHONY: all
all: \
        audit \
        test \
        ;

# ============================================================================ #
##:
##:                                 Variables
##:
# These variables can be changed here directly by editing this file
# or by passing them into the `make` call:
#   make <variable_1>=<value_1> <variable_2>=<value_2> [...]
#
# To generate a target that prints the value of a variable, use the list below
# and append it with the variable name:
#   __variables__ += <variable name>
    override __variables__ :=
# ============================================================================ #

## BINARY_DIR: get the directory with binaries
BINARY_DIR := bin
__variables__ += BINARY_DIR

# The `go install` command installs binaries to GOBIN.
export GOBIN ?= $(__PROJECT_ROOT__)$(BINARY_DIR)

# Prepend PATH with GOBIN.
export PATH := $(GOBIN)$(__LIST_SEP__)$(PATH)

# Generate variable getters for all the variables in the last __variables__.
define make_variable_getter
.PHONY: $1
$1:
	@ echo "$($1)"
endef
$(foreach var,$(__variables__), \
    $(eval \
        $(call make_variable_getter,$(var)) \
    ) \
)

# ============================================================================ #
#                                    Helpers
# ============================================================================ #

.PHONY: confirm
confirm:
ifeq ($(__OS__),Windows)
	@ if ((Read-Host -Prompt "Are you sure? [y/N]") -cne "y") { throw "" }
else
	@ read -r -p 'Are you sure? [y/N] ' answer \
    && [ $${answer:-N} = 'y' ]
endif

.PHONY: git/no-dirty
git/no-dirty:
ifeq ($(__OS__),Windows)
	@ if (![string]::IsNullOrEmpty("$(shell git status --porcelain)")) { throw "" }
else
	@ test -z "$(shell git status --porcelain)"
endif

.PHONY: create/binary_dir
create/binary_dir:
ifeq ($(__OS__),Windows)
	@ [void](New-Item "$(BINARY_DIR)" -ItemType Directory -Force)
else
	@ mkdir -p "$(BINARY_DIR)"
endif

.PHONY: cgo/enable
cgo/enable:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go env -w CGO_ENABLED=1
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

.PHONY: cgo/disable
cgo/disable:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go env -w CGO_ENABLED=0
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

# ============================================================================ #
##:
##:                              Quality control
##:
# ============================================================================ #

## audit: run quality control checks
.PHONY: audit
audit: \
        fmt/no-dirty \
        mod/tidy-diff \
        mod/verify \
        tools/run \
        golangci-lint \
        ;

## mod/tidy-diff: check missing and unused modules without modifying
##              : the `go.mod` and `go.sum` files
.PHONY: mod/tidy-diff
mod/tidy-diff:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go mod tidy -diff
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## mod/tidy: add missing and remove unused modules
.PHONY: mod/tidy
mod/tidy:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go mod tidy -v
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## mod/verify: verify that dependencies have expected content
.PHONY: mod/verify
mod/verify:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go mod verify
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## fmt: gofmt (reformat) package sources
.PHONY: fmt
fmt:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go fmt ./...
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## fmt/no-dirty: gofmt (reformat) package sources and fail if there are some
##             : changes
.PHONY: fmt/no-dirty
fmt/no-dirty:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
ifeq ($(__OS__),Windows)
	@ if (![string]::IsNullOrEmpty("$(shell go fmt ./...)")) { throw "" }
else
	@ test -z "$(shell go fmt ./...)"
endif
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## vet: report likely mistakes in packages
.PHONY: vet
vet:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go vet ./...
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## golangci-lint: a fast linters runner for Go
.PHONY: golangci-lint
golangci-lint:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ golangci-lint run ./...
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## test: run all the tests
.PHONY: test
test: cgo/enable
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go test -v -race ./...
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

## test/cover: run all the tests and display coverage
.PHONY: test/cover
test/cover: \
        create/binary_dir \
        cgo/enable
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ go test -v -race -coverpkg=./... -coverprofile='$(BINARY_DIR)/coverage.out' ./...
	@ go tool cover -html='$(BINARY_DIR)/coverage.out' -o '$(BINARY_DIR)/coverage.html'
	@ go tool cover -func='$(BINARY_DIR)/coverage.out'
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

# ============================================================================ #
##:
##:                                   Tools
##:
# ============================================================================ #

.PHONY: tools/install
tools/install:
ifeq ($(__OS__),Windows)
	@ cd tools <#\
 #> && (Get-Content "tools.go") -match "^\s*_" -replace '\s|_|"|`',"" <#\
 #> | ForEach-Object { go install $$_ }
else
	@ cd tools \
    && for tool in $$(cat tools.go | grep -E '^\s*_' | grep -Po '(?<=(["`]))[\w./]+(?=\1)'); do \
        go install $${tool}; \
    done
endif

## tools/run: run all developer tools
.PHONY: tools/run
tools/run: \
        tools/install \
        govulncheck \
        ;

## govulncheck: report known vulnerabilities that affect Go code
.PHONY: govulncheck
govulncheck:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Running $@..."
	@ $(GOBIN)/govulncheck ./...
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Running $@ - done"

# ============================================================================ #
##:
##:                                   Build
##:
# ============================================================================ #

## mod/download: download modules to local cache
.PHONY: mod/download
mod/download:
	@ go mod download -x

## clean: remove files from the binary directory
.PHONY: clean
clean:
	@ $(call __color_text__,$(__BLUE__),o ) && echo "Cleaning $(BINARY_DIR)..."
ifeq ($(__OS__),Windows)
	@ if (Test-Path "$(BINARY_DIR)" -PathType Container) { <#\
 #>     Remove-Item "$(BINARY_DIR)\*" -Recurse -Force <#\
 #> }
else
	@ rm -rf $(BINARY_DIR)/*
endif
	@ $(call __color_text__,$(__GREEN__),v ) && echo "Cleaning $(BINARY_DIR) - done"

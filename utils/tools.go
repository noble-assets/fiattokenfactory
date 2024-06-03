//go:build tools

// This file uses the recommended method for tracking developer tools in a Go module.
//
// REF: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package utils

import (
	_ "mvdan.cc/gofumpt"
)
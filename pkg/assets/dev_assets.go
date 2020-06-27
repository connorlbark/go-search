// +build dev

package assets

//go:generate vfsgendev -source="github.com/porgull/go-search/internal/assets".Assets

import "net/http"

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("assets")

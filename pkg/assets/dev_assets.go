// +build dev

package assets

//go:generate vfsgendev -source="github.com/porgull/go-search/pkg/assets".Assets

import (
	"net/http"
	"os"
	"path"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir(path.Join(os.Getenv("HOME"), "go/src/github.com/porgull/go-search/assets"))

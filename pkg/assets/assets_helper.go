package assets

import "net/http"

// MustOpen returns an asset file, and panics
// on error
func MustOpen(name string) http.File {
	f, err := Assets.Open(name)
	if err != nil {
		panic(err)
	}
	return f
}

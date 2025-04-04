package petrock_example_feature_name

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed all:assets
var embeddedAssets embed.FS

// Assets contains the embedded filesystem for feature-specific static assets.
var Assets fs.FS

func init() {
	// Expose the 'assets' subdirectory from the embedded FS
	subFS, err := fs.Sub(embeddedAssets, "assets")
	if err != nil {
		// This should not happen if the 'assets' directory exists
		panic(fmt.Sprintf("failed to get sub-filesystem for feature assets: %v", err))
	}
	Assets = subFS
}

// GetAssetsFS returns the embedded filesystem containing the feature's assets.
// This can be used by the main application's HTTP server setup.
func GetAssetsFS() fs.FS {
	return Assets
}

// ServeAssets returns an http.Handler that serves static files from the embedded assets.
// It strips the provided prefix before looking up files in the embedded FS.
// Example Usage in main server setup:
// featureAssets := petrock_example_feature_name.ServeAssets("/assets/feature_name/")
// mux.Handle("/assets/feature_name/", featureAssets)
func ServeAssets(prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.FS(Assets)))
}

package petrock_example_feature_name

import (
	"embed"
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

// --- Helper (if needed) ---
// Placeholder for fmt if used elsewhere and init panics
import "fmt"
```

**10. Create `internal/skeleton/feature_template/assets/` directory**

This step just involves creating an empty directory. I cannot represent this directly in a *SEARCH/REPLACE* block, but you should ensure this directory exists:

```bash
mkdir -p internal/skeleton/feature_template/assets
```

**11. Modify `embeds.go`**

```go
embeds.go
<<<<<<< SEARCH
package petrock

import "embed"

// SkeletonFS holds the embedded filesystem for the project skeleton.
//go:embed all:internal/skeleton
var SkeletonFS embed.FS

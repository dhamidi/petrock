# Plan for posts/assets.go (Example Feature)

This file uses `go:embed` to bundle feature-specific static assets (like CSS, JavaScript, images) directly into the Go binary.

## Types

- None specific to this file.

## Variables

- `//go:embed assets`
- `Assets embed.FS`: This declares a variable `Assets` of type `embed.FS` and the `//go:embed assets` directive instructs the Go compiler to embed the contents of the `posts/assets/` directory into this variable.

## Functions

- `GetAssetsFS() fs.FS`: Returns the embedded filesystem (`Assets`) wrapped as an `fs.FS`. This allows the rest of the application (e.g., the HTTP server setup in `cmd/serve.go`) to serve files from this embedded filesystem, typically under a specific path prefix (e.g., `/assets/posts/`).

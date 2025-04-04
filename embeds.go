package petrock

import "embed"

// SkeletonFS holds the embedded filesystem for the main project skeleton.
//go:embed all:internal/skeleton
var SkeletonFS embed.FS

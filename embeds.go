package petrock

import "embed"

// SkeletonFS holds the embedded filesystem for the project skeleton.
//go:embed all:internal/skeleton
var SkeletonFS embed.FS

package skeletonfs

import "embed"

// SkeletonFS holds the embedded filesystem for the project skeleton.
// The path is relative from this file's directory (internal/skeletonfs)
// to the target directory (internal/skeleton).
//go:embed all:../skeleton
var SkeletonFS embed.FS

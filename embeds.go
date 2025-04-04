package petrock

import "embed"

// SkeletonFS holds the embedded filesystem for the main project skeleton.
//go:embed all:internal/skeleton
var SkeletonFS embed.FS

// FeatureTemplateFS holds the embedded filesystem for the feature template.
//go:embed all:internal/skeleton/feature_template
var FeatureTemplateFS embed.FS

// TODO: Consider combining these into a single embed of a parent directory
// if managing multiple distinct embedded roots becomes complex.
// For now, separate embeds are clear.

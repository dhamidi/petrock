package queries

import (
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure query results implement the marker interface
var _ core.QueryResult = (*Result)(nil)
var _ core.QueryResult = (*ListResult)(nil)

// Ensure queries implement the marker interface
var _ core.Query = (*GetQuery)(nil)
var _ core.Query = (*ListQuery)(nil)
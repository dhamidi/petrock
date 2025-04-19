package pages

import (
	"github.com/petrock/example_module_path/petrock_example_feature_name/handlers"
	"github.com/petrock/example_module_path/petrock_example_feature_name/queries"
)

// Result is a type alias for ItemResult from the queries package
type Result = queries.ItemResult

// ListResult is a type alias for ListQueryResult from the queries package
type ListResult = handlers.ListResult

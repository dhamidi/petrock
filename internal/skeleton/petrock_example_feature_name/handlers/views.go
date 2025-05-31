package handlers

import (
	g "maragu.dev/gomponents"
	
	"github.com/petrock/example_module_path/petrock_example_feature_name/queries"
	"github.com/petrock/example_module_path/petrock_example_feature_name/ui/pages"
)

// ItemView renders the HTML representation of a single item.
func ItemView(item queries.ItemResult) g.Node {
	return pages.ItemView(pages.Result(item))
}

// ItemsListView renders a list of items, potentially with pagination.
func ItemsListView(result queries.ListQueryResult) g.Node {
	// Convert from queries.ListQueryResult to pages.ListResult
	pageResult := pages.ListResult{
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalCount: result.TotalCount,
		Items:      make([]pages.Result, len(result.Items)),
	}
	
	// Convert each Result item
	for i, item := range result.Items {
		pageResult.Items[i] = pages.Result(item)
	}
	
	return pages.ItemsListView(pageResult)
}

// ItemForm renders an HTML <form> for creating or editing an item.
func ItemForm(form interface{}, item *queries.ItemResult, csrfToken string) g.Node {
	var pageItem *pages.Result
	if item != nil {
		converted := pages.Result(*item)
		pageItem = &converted
	}
	return pages.EditForm(form, pageItem, csrfToken)
}

// DeleteConfirmForm renders a form to confirm deletion of an item.
func DeleteConfirmForm(item *queries.ItemResult, csrfToken string) g.Node {
	var pageItem *pages.Result
	if item != nil {
		converted := pages.Result(*item)
		pageItem = &converted
	}
	return pages.DeleteForm(pageItem, csrfToken)
}

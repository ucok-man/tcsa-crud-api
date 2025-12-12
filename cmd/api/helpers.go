package main

import "strings"

type envelope map[string]any

func (app *application) SortColumn(value string) string {
	return strings.TrimPrefix(value, "-")
}

func (app *application) SortDirection(value string) string {
	if strings.HasPrefix(value, "-") {
		return "DESC"
	}

	return "ASC"
}

func (app *application) PageOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

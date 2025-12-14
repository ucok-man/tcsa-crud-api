package main

import (
	"strings"
)

type envelope map[string]any

func (app *application) SortColumn(value string) string {
	column := strings.TrimPrefix(value, "-")
	return column
}

func (app *application) SortDirection(value string) string {
	var direction string

	if strings.HasPrefix(value, "-") {
		direction = "DESC"
	} else {
		direction = "ASC"
	}
	return direction

}

func (app *application) PageOffset(page, pageSize int) int {
	offset := (page - 1) * pageSize
	return offset
}

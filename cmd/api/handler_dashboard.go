package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/tcsa/cmd/api/dto"
	"github.com/ucok-man/tcsa/internal/data"
	"github.com/ucok-man/tcsa/internal/utility"
)

func (app *application) summaryTransactionHandler(ctx echo.Context) error {
	var dto dto.TransactionSummaryDTO

	// Set Default Value
	dto.Pagination.Page = utility.SetPtrValue(1)
	dto.Pagination.PageSize = utility.SetPtrValue(10)
	dto.Sort.Value = utility.SetPtrValue("id")

	if err := ctx.Bind(&dto); err != nil {
		return app.ErrBadRequest(err.Error())
	}

	if err := ctx.Validate(&dto); err != nil {
		return app.ErrFailedValidation(err)
	}

	summary, metadata, err := app.models.Transactions.Summary(data.TransactionSummaryParam{
		Page:            *dto.Pagination.Page,
		PageSize:        *dto.Pagination.PageSize,
		PageOffset:      app.PageOffset(*dto.Pagination.Page, *dto.Pagination.PageSize),
		SortColumn:      app.SortColumn(*dto.Sort.Value),
		SortDirection:   app.SortDirection(*dto.Sort.Value),
		FilterDateRange: utility.DerefOrDefault(dto.Filter.DateRange, 0),
		FilterUserId:    utility.DerefOrDefault(dto.Filter.UserId, 0),
	})
	if err != nil {
		return app.ErrInternalServer(err, "failed get transaction summary", ctx.Request())
	}

	return ctx.JSON(http.StatusOK, envelope{
		"data":     summary,
		"metadata": metadata,
	})
}

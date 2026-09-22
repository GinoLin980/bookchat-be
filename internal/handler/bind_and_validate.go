package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func BindAndValidate[T any](c *echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	return req, c.Validate(req)
}

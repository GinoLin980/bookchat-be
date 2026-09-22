package handler

import "github.com/labstack/echo/v5"

func BindAndValidate[T any](c *echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return nil, echo.NewHTTPError(400, "invalid request body")
	}
	if err := c.Validate(req); err != nil {
		return nil, err // 400 + fields error
	}
	return req, nil
}

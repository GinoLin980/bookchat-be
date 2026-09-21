package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"hello": "world"})
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("can't start server")
	}
}

package handlers

import "github.com/labstack/echo"

// Handler interface defines generic operations for API handlers
type Handler interface {
	Get(c echo.Context) error

	List(c echo.Context) error

	Add(c echo.Context) error

	Delete(c echo.Context) error

	Update(c echo.Context) error
}

package handlers

import "github.com/labstack/echo"

type responseBody struct {
	Message string `json:"message"`
}

func Response(c echo.Context, code int, msg string) error {
	return c.JSON(code, responseBody{
		Message: msg,
	})
}

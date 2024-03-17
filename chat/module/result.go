package module

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

func FailedResponse(ctx echo.Context, status int, message string) error {
	return ctx.JSON(status, map[string]interface{}{
		"status": "fail",
		"data": map[string]interface{}{
			"message": message,
		},
	})
}

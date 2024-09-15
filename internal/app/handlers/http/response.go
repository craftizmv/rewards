package http

import (
	"github.com/labstack/echo/v4"
)

// StandardResponse represents the base structure of all API responses.
type StandardResponse[T any] struct {
	Message  string `json:"message"`
	Response T      `json:"response,omitempty"`
}

// SendResponseWithData sends a response with typed data.
func SendResponseWithData[T any](c echo.Context, responseCode int, message string, data T) error {
	response := &StandardResponse[T]{
		Message:  message,
		Response: data,
	}

	return c.JSON(responseCode, response)
}

// SendResponse sends a response without data.
func SendResponse(c echo.Context, responseCode int, message string) error {
	return SendResponseWithData(c, responseCode, message, struct{}{})
}

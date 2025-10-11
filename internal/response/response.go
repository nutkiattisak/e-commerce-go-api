package response

import "github.com/labstack/echo/v4"

type ResponseData struct {
	Data interface{} `json:"data"`
}

type ResponseSuccess struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ResponsePagination struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, ResponseSuccess{
		Message: message,
		Status:  "success",
		Data:    data,
	})
}

func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ResponseError{
		Message: message,
		Status:  "error",
	})
}

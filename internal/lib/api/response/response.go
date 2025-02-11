package response

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// Response base structure for all responses.
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK         = "OK"
	StatusError      = "Error"
	StatusBadRequest = "BadRequest"
	StatusNotFound   = "NotFound"
)

// OK creates a success response.
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// Error creates a generic error response with a custom message.
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// ValidationError creates a response for validation errors with detailed messages.
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.StructField()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.StructField()))
		}
	}

	return Response{
		Status: StatusBadRequest,
		Error:  strings.Join(errMsgs, ", "),
	}
}

// NotFound creates a response for not found errors.
func NotFound(msg string) Response {
	return Response{
		Status: StatusNotFound,
		Error:  msg,
	}
}

// InternalError creates a response for internal server errors.
func InternalError(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// WriteResponse writes the response as JSON with the appropriate HTTP status code.
func WriteResponse(c *gin.Context, statusCode int, resp Response) {
	c.JSON(statusCode, resp)
}

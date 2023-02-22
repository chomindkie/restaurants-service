package errs

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

var (
	BAD_PARAM      = Status{Code: "BAD_PARAM", Message: "Bad Parameters"}
	GENERAL_ERROR  = Status{Code: "GENERAL_ERROR", Message: "General error"}
	INTERNAL_ERROR = Status{Code: "INTERNAL_ERROR", Message: "Internal server error"}

	MESSAGE_FORMAT_ERROR = "httpStatusCode: %d, code: %s, message: %s"
)

type Error struct {
	HTTPStatusCode int          `json:"-"`
	Status         Status       `json:"status"`
	Data           *interface{} `json:"data,omitempty"`
}

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf(MESSAGE_FORMAT_ERROR, e.HTTPStatusCode, e.Status.Code, e.Status.Message)
}

func (s *Status) Error(httpStatusCode int) string {
	return fmt.Sprintf(MESSAGE_FORMAT_ERROR, httpStatusCode, s.Code, s.Message)
}

func (s *Status) ErrorWithMessage(httpStatusCode int, message string) string {
	return fmt.Sprintf(MESSAGE_FORMAT_ERROR, httpStatusCode, s.Code, message)
}

func New(httpStatusCode int, code, message string) error {
	return &Error{
		HTTPStatusCode: httpStatusCode,
		Status: Status{
			Code:    code,
			Message: message,
		},
	}
}

func NewStatus(httpStatusCode int, errorStatus Status) error {
	return &Error{
		HTTPStatusCode: httpStatusCode,
		Status: Status{
			Code:    errorStatus.Code,
			Message: errorStatus.Message,
		},
	}
}

func NewWithData(httpStatusCode int, code, message string, data interface{}) error {
	return &Error{
		HTTPStatusCode: httpStatusCode,
		Status: Status{
			Code:    code,
			Message: message,
		},
		Data: &data,
	}
}

func NewStatusWithData(httpStatusCode int, errorStatus Status, data interface{}) error {
	return &Error{
		HTTPStatusCode: httpStatusCode,
		Status: Status{
			Code:    errorStatus.Code,
			Message: errorStatus.Message,
		},
		Data: &data,
	}
}

func HTTPErrorHandler(err error, c echo.Context) {
	if _, ok := errors.Cause(err).(*Error); ok {
		JSON(c, err)
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	JSON(c, New(code, GENERAL_ERROR.Code, err.Error()))
}

func JSON(c echo.Context, err error) error {
	if errs, ok := errors.Cause(err).(*Error); ok {
		//acnlog.ErrorfWithID(c, "raw error => %+v", err)
		if errs.Data != nil {
			return c.JSON(errs.HTTPStatusCode, Error{
				Status: Status{
					Code:    errs.Status.Code,
					Message: errs.Status.Message,
				},
				Data: errs.Data,
			})
		} else {
			return c.JSON(errs.HTTPStatusCode, Error{
				Status: Status{
					Code:    errs.Status.Code,
					Message: errs.Status.Message,
				},
			})
		}
	}

	return c.JSON(http.StatusInternalServerError, Error{
		Status: Status{
			Code:    GENERAL_ERROR.Code,
			Message: err.Error(),
		},
	})
}

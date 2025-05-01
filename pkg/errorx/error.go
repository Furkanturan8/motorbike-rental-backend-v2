package errorx

import (
	"fmt"
	"net/http"
)

var (
	ErrValidation         = New(http.StatusUnprocessableEntity, "Doğrulama hatası")
	ErrUnauthorized       = New(http.StatusUnauthorized, "Yetkisiz erişim")
	ErrForbidden          = New(http.StatusForbidden, "Erişim reddedildi")
	ErrNotFound           = New(http.StatusNotFound, "Kaynak bulunamadı")
	ErrInternal           = New(http.StatusInternalServerError, "Sunucu hatası")
	ErrDuplicate          = New(http.StatusConflict, "Kaynak zaten mevcut")
	ErrInvalidRequest     = New(http.StatusBadRequest, "Geçersiz istek")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "Geçersiz kimlik bilgileri")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err ...error) *AppError {
	var inner error
	if len(err) > 0 {
		inner = err[0]
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     inner,
	}
}

func WrapErr(base *AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Err:     err,
	}
}

func WrapMsg(base *AppError, customMessage string) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: customMessage,
	}
}

func Wrap(base *AppError, err error, customMessage string) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: customMessage,
		Err:     err,
	}
}

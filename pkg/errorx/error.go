package errorx

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func WithDetails(err *Error, details string) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Details: details,
		Err:     err,
	}
}

func Wrap(err *Error, originalErr error) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Details: originalErr.Error(),
		Err:     originalErr,
	}
}

// Önceden tanımlanmış hatalar
var (
	ErrValidation = &Error{
		Code:    http.StatusUnprocessableEntity,
		Message: "Doğrulama hatası",
	}

	ErrUnauthorized = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Yetkisiz erişim",
	}

	ErrForbidden = &Error{
		Code:    http.StatusForbidden,
		Message: "Erişim reddedildi",
	}

	ErrNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Kaynak bulunamadı",
	}

	ErrInternal = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Sunucu hatası",
	}

	ErrDuplicate = &Error{
		Code:    http.StatusConflict,
		Message: "Kaynak zaten mevcut",
	}

	ErrInvalidRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz istek",
	}

	ErrDatabaseOperation = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Veritabanı işlemi başarısız",
	}

	ErrInvalidCredentials = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Geçersiz kimlik bilgileri",
	}

	ErrAccountInactive = &Error{
		Code:    http.StatusForbidden,
		Message: "Hesap aktif değil",
	}

	ErrPasswordHash = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Şifre hash'leme hatası",
	}

	ErrDuplicateEmail = &Error{
		Code:    http.StatusConflict,
		Message: "E-posta adresi zaten kullanımda",
	}

	ErrCacheNotInitialized = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Önbellek başlatılmadı",
	}

	ErrKeyNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Önbellekte anahtar bulunamadı",
	}

	ErrInvalidValue = &Error{
		Code:    http.StatusUnprocessableEntity,
		Message: "Geçersiz değer tipi",
	}

	ErrTokenGeneration = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Token oluşturma hatası",
	}

	ErrTokenValidation = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Token doğrulama hatası",
	}

	ErrSessionInvalid = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Oturum geçersiz",
	}

	ErrPasswordTooShort = &Error{
		Code:    http.StatusBadRequest,
		Message: "Şifre çok kısa",
	}

	ErrEmailInvalid = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz e-posta adresi",
	}

	ErrUserNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Kullanıcı bulunamadı",
	}

	ErrUserAlreadyExists = &Error{
		Code:    http.StatusConflict,
		Message: "Kullanıcı zaten mevcut",
	}

	ErrInvalidToken = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Geçersiz token",
	}

	ErrTokenExpired = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Token süresi dolmuş",
	}

	ErrTokenRevoked = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Token iptal edilmiş",
	}

	ErrSessionExpired = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Oturum süresi dolmuş",
	}

	ErrSessionBlocked = &Error{
		Code:    http.StatusForbidden,
		Message: "Oturum engellenmiş",
	}

	ErrPasswordMismatch = &Error{
		Code:    http.StatusBadRequest,
		Message: "Şifreler eşleşmiyor",
	}

	ErrInvalidPasswordFormat = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz şifre formatı",
	}

	ErrInvalidRole = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz rol",
	}

	ErrInvalidStatus = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz durum",
	}

	ErrInvalidFirstName = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz ad",
	}

	ErrInvalidLastName = &Error{
		Code:    http.StatusBadRequest,
		Message: "Geçersiz soyad",
	}
)

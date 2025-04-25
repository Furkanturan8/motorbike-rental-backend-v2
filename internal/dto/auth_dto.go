package dto

// Giriş isteği
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Token yanıtı
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Saniye cinsinden
}

// Token yenileme isteği
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Şifre sıfırlama isteği
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Şifre sıfırlama
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type RegisterResponse struct {
	Email string `json:"email"`
}

package data

import "time"

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	SecretKey    string `json:"secret_key"`
	Username     string `json:"username"`
}

type OTPResponse struct {
	OTPUrl    string `json:"auth_url"`
	SecretKey string `json:"secret_key"`
}

type RefreshTokenResponse struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

package data

type UserResponse struct {
	UID         string `json:"uid"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	IsVerified  bool   `json:"is_verified"`
	OTPVerified bool   `json:"otp_verified"`
	OTPEnabled  bool   `json:"otp_enabled"`
	OTPUrl      string `json:"otp_auth_url"`
	LastLogin   string `json:"last_login"`
}

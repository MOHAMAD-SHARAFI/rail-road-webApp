package models

type SignUpRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type SignInRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

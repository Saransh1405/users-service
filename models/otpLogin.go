package models

type SendOTPRequest struct {
	ClientName  string `json:"clientName"`
	CountryCode string `json:"countryCode" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
} // @name SendOTPRequest

type LoginWithOTPRequest struct {
	ClientName  string `json:"clientName" binding:"required"`
	CountryCode string `json:"countryCode" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
} // @name LoginWithOTPRequest

type LoginResponse struct {
	User             Users  `bson:"user" json:"user"`
	AccessToken      string `json:"accessToken,omitempty"`
	ExpiresIn        int64  `json:"expiresIn,omitempty"`
	RefreshExpiresIn int    `json:"refreshExpiresIn,omitempty"`
	RefreshToken     string `json:"refreshToken,omitempty"`
} // @name LoginResponse

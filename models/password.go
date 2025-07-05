package models

type UpdatePasswordRequest struct {
	ClientName  string `json:"clientName" binding:"required" example:"kecya"`
	Email       string `json:"email" binding:"required" example:"xyz@gmail.com"`
	Password    string `json:"password" binding:"required" example:"password123"`
	NewPassword string `json:"newPassword" binding:"required" example:"password123"`
} //@name UpdatePasswordRequest

type ForgotPasswordRequest struct {
	ClientName string `json:"clientName" binding:"required" example:"kecya"`
	Email      string `json:"email" binding:"required" example:"xyz@gmail.com"`
} //@name ForgotPasswordRequest

type ResetPasswordRequest struct {
	ClientName         string `json:"clientName" binding:"required" example:"kecya"`
	ResetToken         string `json:"resetToken" binding:"required" example:"1234567890"`
	NewPassword        string `json:"newPassword" binding:"required" example:"password123"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required" example:"password123"`
} //@name ResetPasswordRequest

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPostRequest struct {
	FirstName      string `json:"firstName" binding:"required" example:"John"`
	LastName       string `json:"lastName" binding:"required" example:"Doe"`
	Email          string `json:"email" binding:"required" example:"xyz@gmail.com"`
	Phone          string `json:"phone" binding:"required" example:"1234567890"`
	CountryCode    string `json:"countryCode" binding:"required" example:"+91"`
	UserProfileUrl string `json:"userProfileUrl" example:"www.xyz.com"`
	Password       string `json:"password" binding:"required" example:"password123"`
} //@name UserPostRequest

type UserPatchRequest struct {
	Id                  string `json:"id" binding:"required" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`
	ClientName          string `json:"client_name" binding:"required" example:"kecya"`
	FirstName           string `json:"firstName" example:"John"`
	LastName            string `json:"lastName" example:"Doe"`
	UserProfileUrl      string `json:"userProfileUrl" example:"www.xyz.com"`
	Status              string `json:"status" example:"active"`
	ReasonForSuspension string `json:"reasonForSuspension" example:"xyz"`
} //@name UserPatchRequest

type UserDeleteRequest struct {
	ID                string `json:"id" binding:"required" example:"5"`
	ReasonForDeletion string `json:"reasonForDeletion" bson:"reasonForDeletion" example:"not verified"`
	ClientName        string `json:"client_name" binding:"required" example:"kecya"`
} //@name UserDeleteRequest

type GetUserRequest struct {
	ClientName string `json:"client_name" binding:"required" form:"client_name" example:"kecya"`
} //@name GetUserRequest

type Users struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`
	SuspendedAt         int64              `bson:"suspendedAt"  json:"suspendedAt"`
	FirstName           string             `bson:"firstName" example:"John" json:"firstName"`
	LastName            string             `bson:"lastName" example:"Doe" json:"lastName"`
	Password            string             `bson:"password" example:"password123" json:"password"`
	Email               string             `bson:"email" example:"jhon@example.com" json:"email"`
	CountryCode         string             `bson:"countryCode" example:"+91" json:"countryCode"`
	Phone               string             `bson:"phone" example:"1234567890" json:"phone"`
	PhoneVerified       bool               `bson:"phoneVerified" json:"phoneVerified"`
	UserProfileUrl      string             `bson:"userProfileUrl" example:"www.xyz.com" json:"userProfileUrl"`
	TemporaryPassword   string             `bson:"temporaryPassword" example:"123456" json:"temporaryPassword"`
	Status              Status             `bson:"status"  example:"Active" json:"status"`
	ReasonForDeletion   string             `bson:"reasonForDeletion" example:"not verified" json:"reasonForDeletion"`
	ReasonForSuspension string             `bson:"reasonForSuspension" example:"not verified" json:"reasonForSuspension"`
	StatusLogs          []StatusLogs       `bson:"statusLogs" json:"statusLogs"`
	ClientName          string             `bson:"clientName" json:"clientName"`
	CreatedAt           int64              `bson:"createdAt" json:"createdAt"`
} //@name Users

type OTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp,omitempty"`
} //@name OTPRequest

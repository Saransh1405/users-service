package login

import (
	"errors"
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/google"
	"users-service/utils/helperfunctions"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// GoogleLogin handles Google OAuth login
func GoogleLogin(ctx *gin.Context, request *models.GoogleLoginRequest) (*models.LoginResponse, error) {
	log := logger.GetLoggerWithoutContext()

	// Initialize Google OAuth service
	googleService, err := google.NewGoogleOAuthService()
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to initialize Google OAuth service")
		return nil, errors.New("authentication service unavailable")
	}

	// Exchange authorization code for access token
	token, err := googleService.ExchangeCodeForToken(ctx, request.Code)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to exchange code for token")
		return nil, errors.New("invalid authorization code")
	}

	// Get user information from Google
	googleUserInfo, err := googleService.GetUserInfo(ctx, token)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to get user info from Google")
		return nil, errors.New("failed to retrieve user information")
	}

	// Get user collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	// Check if user exists by Google ID
	var existingUser models.Users
	err = userCol.FindOne(ctx, bson.M{"googleId": googleUserInfo.ID}).Decode(&existingUser)
	if err == nil {
		// User exists, update last login and return
		return handleExistingGoogleUser(ctx, &existingUser, log)
	}

	// Check if user exists by email
	err = userCol.FindOne(ctx, bson.M{"email": googleUserInfo.Email}).Decode(&existingUser)
	if err == nil {
		// User exists with email but no Google ID, link accounts
		return linkGoogleAccount(ctx, &existingUser, googleUserInfo, log)
	}

	loginResponse, err := createGoogleUser(ctx, googleUserInfo, request.ClientName, log)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to create Google user")
		return nil, errors.New("failed to create user account")
	}

	// go func() {
	// 	helperfunctions.SendNotifications(map[string]interface{}{
	// 		"type":    models.NotificationTypeEmail,
	// 		"userId":  googleUserInfo.ID,
	// 		"message": "User logged in successfully",
	// 		"data": map[string]interface{}{
	// 			"To":       loginResponse.User.Email,
	// 			"Message":  "User logged in successfully",
	// 			"Template": "login",
	// 			"Variables": map[string]interface{}{
	// 				"Name":  loginResponse.User.FirstName + " " + loginResponse.User.LastName,
	// 				"Email": loginResponse.User.Email,
	// 			},
	// 			"Subject": "Welcome to the TribeWithVibe",
	// 		},
	// 	}, string(models.NotificationTopic))
	// }()

	// Create new user
	return loginResponse, nil
}

// handleExistingGoogleUser handles login for existing Google users
func handleExistingGoogleUser(ctx *gin.Context, user *models.Users, log logger.Logger) (*models.LoginResponse, error) {
	// Check if user is active
	if user.Status != "Active" {
		return nil, errors.New("account is not active")
	}

	// Generate JWT token
	token, expireTime, err := helperfunctions.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Error("Failed to generate JWT token")
		return nil, errors.New("failed to generate authentication token")
	}

	response := models.LoginResponse{
		User:        *user,
		AccessToken: token,
		ExpiresIn:   int64(expireTime),
	}

	return &response, nil
}

// linkGoogleAccount links existing email account with Google account
func linkGoogleAccount(ctx *gin.Context, user *models.Users, googleUserInfo *models.GoogleUserInfo, log logger.Logger) (*models.LoginResponse, error) {
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	// Update user with Google information
	update := bson.M{
		"$set": bson.M{
			"googleId":       googleUserInfo.ID,
			"authProvider":   "google",
			"emailVerified":  googleUserInfo.VerifiedEmail,
			"firstName":      googleUserInfo.GivenName,
			"lastName":       googleUserInfo.FamilyName,
			"userProfileUrl": googleUserInfo.Picture,
		},
	}

	_, err := userCol.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to link Google account")
		return nil, errors.New("failed to link accounts")
	}

	// Update user object with new data
	user.GoogleID = googleUserInfo.ID
	user.AuthProvider = "google"
	user.EmailVerified = googleUserInfo.VerifiedEmail
	user.FirstName = googleUserInfo.GivenName
	user.LastName = googleUserInfo.FamilyName
	user.UserProfileUrl = googleUserInfo.Picture

	// Generate JWT token
	token, expireTime, err := helperfunctions.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Error("Failed to generate JWT token")
		return nil, errors.New("failed to generate authentication token")
	}

	response := models.LoginResponse{
		User:        *user,
		AccessToken: token,
		ExpiresIn:   int64(expireTime),
	}

	return &response, nil
}

// createGoogleUser creates a new user from Google OAuth
func createGoogleUser(ctx *gin.Context, googleUserInfo *models.GoogleUserInfo, clientName string, log logger.Logger) (*models.LoginResponse, error) {
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	// Create new user
	newUser := models.Users{
		ID:             primitive.NewObjectID(),
		FirstName:      googleUserInfo.GivenName,
		LastName:       googleUserInfo.FamilyName,
		Email:          googleUserInfo.Email,
		UserProfileUrl: googleUserInfo.Picture,
		GoogleID:       googleUserInfo.ID,
		AuthProvider:   "google",
		EmailVerified:  googleUserInfo.VerifiedEmail,
		PhoneVerified:  false,
		Status:         models.Active,
		ClientName:     clientName,
		CreatedAt:      time.Now().UnixMilli(),
		// Set default values for required fields
		CountryCode: "+91",       // Default country code
		Phone:       "",          // Empty phone for Google users
		Password:    "temporary", // No password for OAuth users
	}

	// Insert user into database
	_, err := userCol.InsertOne(ctx, newUser)
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to create new Google user")
		return nil, errors.New("failed to create user account")
	}

	// Generate JWT token
	token, expireTime, err := helperfunctions.GenerateJWT(newUser.ID, newUser.Email)
	if err != nil {
		log.Error("Failed to generate JWT token")
		return nil, errors.New("failed to generate authentication token")
	}

	response := models.LoginResponse{
		User:        newUser,
		AccessToken: token,
		ExpiresIn:   int64(expireTime),
	}

	return &response, nil
}

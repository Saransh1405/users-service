package login

import (
	"users-service/business/login"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils"
	"users-service/utils/google"
	"users-service/utils/helperfunctions"
	"users-service/utils/localization"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// GoogleLogin handles Google OAuth login
func GoogleLogin(ctx *gin.Context) {
	// Get language
	lang, _ := ctx.Get(constants.LanguageString)

	// Get logger
	log := logger.GetLogger(ctx)

	var request models.GoogleLoginRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}

	// Call business logic
	resp, err := login.GoogleLogin(ctx, &request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	// Send success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, resp)
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func GetGoogleAuthURL(ctx *gin.Context) {
	// Get language
	lang, _ := ctx.Get(constants.LanguageString)

	// Get logger
	log := logger.GetLogger(ctx)

	// Get state parameter (optional, for CSRF protection)
	state := ctx.Query("state")
	if state == "" {
		state = "default_state" // You might want to generate a random state
	}

	// Initialize Google OAuth service
	googleService, err := google.NewGoogleOAuthService()
	if err != nil {
		log.With(zap.Error(err)).Error("Failed to initialize Google OAuth service")
		msg := localization.GetMessage(lang, "authentication service unavailable", nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	// Get authorization URL
	authURL := googleService.GetAuthURL(state)

	response := map[string]string{
		"authUrl": authURL,
		"state":   state,
	}

	// Send success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, response)
}

package password

import (
	"users-service/business/password"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils"
	"users-service/utils/helperfunctions"
	"users-service/utils/localization"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UpdatePassword(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.UpdatePasswordRequest
	if validationErr := helperfunctions.ValidateRequestDataParam(ctx, &request); validationErr != nil {
		return
	}

	//call the business logic
	err := password.UpdatePassword(ctx, &request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "Password updated successfully")
}

func ForgotPassword(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.ForgotPasswordRequest
	if validationErr := helperfunctions.ValidateRequestDataParam(ctx, &request); validationErr != nil {
		return
	}

	//call the business logic
	err := password.ForgotPassword(ctx, &request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "Password sent successfully")
}

func ResetPassword(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.ResetPasswordRequest
	if validationErr := helperfunctions.ValidateRequestDataParam(ctx, &request); validationErr != nil {
		return
	}

	//call the business logic
	err := password.ResetPassword(ctx, &request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "Password reset successfully")
}

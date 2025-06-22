package signup

import (
	"users-service/business/signup"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils"
	"users-service/utils/helperfunctions"
	"users-service/utils/localization"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

func Post(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.UserPostRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}
	resp, err := signup.Post(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, resp)
}

func Patch(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.UserPatchRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}
	err := signup.Patch(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "user updated successfully")
}

func Delete(ctx *gin.Context) {
	//get the lang
	lang, _ := ctx.Get(constants.LanguageString)

	//get the logger
	log := logger.GetLogger(ctx)

	var request models.UserDeleteRequest
	if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil {
		return
	}
	err := signup.Delete(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		msg := localization.GetMessage(lang, err.Error(), nil)
		utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err)
		return
	}

	//sent the success message
	successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil)
	utils.SendStatusOK(ctx, constants.IsString, successMessage, "user deleted successfully")
}

// package modifier import ( "go-ecom-apis/business/modifier" "go-ecom-apis/constants" "go-ecom-apis/logger" "go-ecom-apis/models" "go-ecom-apis/utils" "go-ecom-apis/utils/helperfunctions" "go-ecom-apis/utils/localization" "github.com/gin-gonic/gin" "github.com/gin-gonic/gin/binding" "go.uber.org/zap" ) // Handler - Handle the post requests at /v1/product/modifiers // @Summary Create a new modifier // @Description Create a new modifier for products // @Tags Modifier // @Accept json // @Produce json // @Param Authorization header string true "Authorization header" // @Param modifier body models.Modifier true "Modifier" // @Success 201 {object} models.Modifier // @Failure 400 {object} models.APIResponse "The server could not understand the request that it was sent." // @Failure 409 {object} models.APIResponse "Modifier already exists" // @Failure 500 {object} models.ErrorResponse "The server encountered an unexplained problem which has prevented it from executing the given request" // @Router /product/modifiers [post]
// func Post(ctx *gin.Context)
// {
// 	lang, _ := ctx.Get(constants.LanguageString)
// 	 log := logger.GetLogger(ctx)
// 	  var request models.Modifier
// 	  if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil { return } resp, err := modifier.Post(ctx, &request)
// 	  if err != nil { log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError) msg := localization.GetMessage(lang, err.Error(), nil) utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err) return }

// 	  successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil) utils.SendStatusCreated(ctx, constants.IsString, successMessage, resp) }

// Handler - Handle the patch requests at /v1/product/modifiers // @Summary Patch a currency for a storefront // @Description Patch a modifier for a storefront // @Tags Modifier // @Accept json // @Produce json // @Param Authorization header string true "Authorization header" // @Param modifier body models.PatchModifier true "Modifier" // @Success 200 {object} models.APIResponse "modifier updated successfully" // @Failure 400 {object} models.APIResponse "The server could not understand the request that it was sent." // @Failure 401 {object} models.APIResponse "The request requires authentication to access the resource, and the client did not send any credentials, or the credentials that were sent are wrong." // @Failure 409 {object} models.APIResponse "Already exists" // @Failure 500 {object} models.ErrorResponse "The server encountered an unexplained problem which has prevented it from executing the given request" // @Router /product/modifiers [patch] func Patch(ctx *gin.Context) { lang, _ := ctx.Get(constants.LanguageString) log := logger.GetLogger(ctx) var request models.PatchModifier if validationErr := helperfunctions.ValidateRequestData(ctx, &request, binding.JSON); validationErr != nil { return } err := modifier.Patch(ctx, &request) if err != nil { log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError) msg := localization.GetMessage(lang, err.Error(), nil) utils.ErrorBasedOnResponse(ctx, msg, constants.IsString, err) return } successMessage := localization.GetMessage(lang, constants.SuccessMessage, nil) utils.SendStatusOK(ctx, constants.IsString, successMessage, "modifier updated successfully") } // Handler - Handle the get requests at /v1/product/modifiers // @Summary Get currencies for a storefront // @Description Get modifiers for a storefront // @Tags Modifier // @Accept json // @Produce json // @Param modifierGetRequest query models.ModifierGetRequest true "Modifier Get Request" // @Param Authorization header string true "Authorization header" // @Success 200 {object} []models.Modifier // @Failure 400 {object} models.APIResponse "The server could not understand the request that it was sent." // @Failure 401 {object} models.APIResponse "The request requires authentication to access the resource, and the client did not send any credentials, or the credentials that were sent are wrong." // @Failure 500 {object} models.ErrorResponse "The server

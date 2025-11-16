package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"users-service/api/accept"
	"users-service/api/campaign"
	"users-service/api/join"
	"users-service/api/login"
	"users-service/api/nearby"
	"users-service/api/password"
	"users-service/api/signup"
	"users-service/constants"
	"users-service/utils"
	"users-service/utils/localization"
	"users-service/utils/middleware"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// GetRouter is used to get the router configured with the middlewares and the routes.
func GetRouter(localizationMiddleware gin.HandlerFunc, loggerMiddleware gin.HandlerFunc, applicationConfig *viper.Viper) *gin.Engine {
	router := gin.New()

	router.Use(localizationMiddleware)
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware)

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

	middlewareFunc := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(time.Duration(12 * time.Hour).Seconds()),
	})

	router.Use(middlewareFunc)

	unAuthRoutes := router.Group("v1")
	{
		// ========================================
		// HTTP ENDPOINTS (OAuth & Public APIs)
		// ========================================
		// These endpoints are designed for frontend integration
		// and OAuth flows that require HTTP redirects

		// Handle the POST requests at /v1/login
		unAuthRoutes.POST(constants.Login, login.Post)

		// Google OAuth endpoints (HTTP-only for OAuth flow)
		unAuthRoutes.POST("/google/login", login.GoogleLogin)
		unAuthRoutes.GET("/google/auth-url", login.GetGoogleAuthURL)

		// Handle the POST requests at /v1/signup
		unAuthRoutes.POST(constants.Signup, signup.Post)

		// Handle the POST requests at /v1/resendOTP
		unAuthRoutes.GET(constants.ResendOTP, signup.ResendOTP)

		// Handle the POST requests at /v1/verifyOTP
		unAuthRoutes.GET(constants.VerifyOTP, signup.VerifyOTP)

		// Handle the POST requests at /v1/sendOTP
		unAuthRoutes.POST(constants.SendOTP, signup.PostSendOTP)

		// ========================================
		// INTERNAL SERVICE ENDPOINTS
		// ========================================
		// These endpoints are for authenticated users
		// and internal service communication

		// Handle the GET requests at /v1/statusNew
		unAuthRoutes.GET("/krakend.json", func(ctx *gin.Context) {
			lang := ctx.GetHeader("language")
			content, err := ioutil.ReadFile("utils/krakend/krakend.json")

			if err != nil {
				Msg := localization.GetMessage(lang, constants.InternalServerMessage, nil)
				utils.SendInternalServerError(ctx, Msg, "0", constants.IsJsonArray, nil)
				return
			}

			backendHost := applicationConfig.GetString(constants.ServerHost)

			krakendData := strings.ReplaceAll(string(content), "SERVER_HOST", backendHost)

			var result map[string]interface{}
			json.Unmarshal([]byte(krakendData), &result)

			// send success response
			ctx.JSON(http.StatusOK, result)
		})
	}

	v1Routes := router.Group("v1")
	{
		v1Routes.Use(middleware.KecyalokMiddleware())

		// Handle the GET requests at /v1/getMyDetails
		v1Routes.GET(constants.User, signup.GetMyDetails)

		// Handle the PATCH requests at /v1/users/me
		v1Routes.PATCH(constants.User, signup.UpdateUser)

		// Handle the DELETE requests at /v1/users/me
		v1Routes.DELETE(constants.User, signup.Delete)

		// Handle the PATCH requests at /v1/password
		v1Routes.PATCH(constants.Password, password.UpdatePassword)

		// Handle the POST requests at /v1/forgotPassword
		v1Routes.POST(constants.ForgotPassword, password.ForgotPassword)

		// Handle the POST requests at /v1/resetPassword
		v1Routes.POST(constants.ResetPassword, password.ResetPassword)

		// Handle the POST requests at /v1/campaign
		v1Routes.POST(constants.Campaign, campaign.CreateCampaign)

		// Handle the PATCH requests at /v1/campaign
		v1Routes.PATCH(constants.Campaign, campaign.UpdateCampaign)

		// Handle the GET requests at /v1/campaign
		v1Routes.GET(constants.Campaign, campaign.GetCampaign)

		// Handle the GET requests at /v1/campaign/nearby
		v1Routes.GET(constants.CampaignNearby, nearby.GetCampaign)

		// Handle the PATCH requests at /v1/campaign/join
		v1Routes.PATCH(constants.CampaignJoin, join.JoinCampaign)

		// Handle the PATCH requests at /v1/campaign/leave
		v1Routes.PATCH(constants.CampaignLeave, join.LeaveCampaign)

		// Handle the PATCH requests at /v1/campaign/accept
		v1Routes.PATCH(constants.CampaignAccept, accept.AcceptCampaign)
	}

	return router
}

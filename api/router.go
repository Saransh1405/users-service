package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"users-service/api/login"
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

	v1Routes := router.Group("v1")
	{
		v1Routes.Use(middleware.KecyalokMiddleware())

		// Handle the GET requests at /v1/getMyDetails
		v1Routes.GET("/getMyDetails", signup.GetMyDetails)

		// Handle the PATCH requests at /v1/users/me
		v1Routes.PATCH("/users/me", signup.UpdateUser)

		// Handle the DELETE requests at /v1/users/me
		v1Routes.DELETE("/users/me", signup.Delete)

	}

	unAuthRoutes := router.Group("v1")
	{
		// Handle the POST requests at /v1/login
		unAuthRoutes.POST(constants.Login, login.Post)

		// // Handle the POST requests at /v1/logout
		// unAuthRoutes.POST(constants.Logout, logout.Post)

		// // Handle the POST requests at /v1/loginWithOtp
		// unAuthRoutes.POST(constants.LoginWithOtp, otpLogin.Post)

		// // Handle the PATCH requests at /v1/loginWithOtp
		// unAuthRoutes.PATCH(constants.LoginWithOtp, otpLogin.Patch)

		// Handle the POST requests at /v1/signup
		unAuthRoutes.POST(constants.Signup, signup.Post)

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

	return router
}

package signup

import (
	"fmt"
	"io"
	"net/http"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/configs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ResendOTP(ctx *gin.Context, request *models.OTPRequest) error {
	//get the logger
	log := logger.GetLogger(ctx)

	// get application config
	applicationConfig, err1 := configs.Get(constants.ApplicationConfig)
	if err1 != nil {
		log.With(zap.Error(err1)).Error(constants.BindingFailedErrr)
	}

	//get the auth key
	authKey := applicationConfig.GetString(constants.OTPAuthKey)

	//get the phone
	phone := request.Phone

	url := fmt.Sprintf("https://control.msg91.com/api/v5/otp?authkey=%s&mobile=%s&retrytype=text", authKey, phone)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("ResendOTP response:", string(body))
	return nil

}

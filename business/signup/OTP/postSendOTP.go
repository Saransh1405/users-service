package signup

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/configs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PostSendOTP(ctx *gin.Context, request *models.OTPRequest) error {
	//get the logger
	log := logger.GetLogger(ctx)

	// get application config
	applicationConfig, err1 := configs.Get(constants.ApplicationConfig)
	if err1 != nil {
		log.With(zap.Error(err1)).Error(constants.BindingFailedErrr)
	}

	//get the auth key
	authKey := applicationConfig.GetString(constants.OTPAuthKey)

	//get the template id
	templateID := applicationConfig.GetString(constants.OTPTemplateID)

	//get the phone
	phone := request.Phone

	url := "https://control.msg91.com/api/v5/otp"
	payload := strings.NewReader(fmt.Sprintf(`{
			"authkey": "%s",
			"mobile": "%s",
			"template_id": "%s"
		}`, authKey, phone, templateID))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("SendOTP response:", string(body))
	return nil
}

package helperfunctions

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/utils"
	"users-service/utils/configs"
	"users-service/utils/localization"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GeneratePassword() string {
	rand.Seed(time.Now().Unix())
	lowerCharSet := constants.ABCDLower
	upperCharSet := constants.ABCDUpper
	specialCharSet := constants.SpecialCharSet2
	numberSet := constants.Number
	allCharSet := lowerCharSet + upperCharSet + specialCharSet + numberSet
	minSpecialChar := 2
	minNum := 2
	minUpperCase := 2
	passwordLength := 13

	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func ValidateRequestData(ctx *gin.Context, request interface{}, b binding.Binding) error {

	lang, _ := ctx.Get(constants.LanguageString)
	log := logger.GetLogger(ctx)

	err := ctx.ShouldBindWith(request, b)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
		var verr validator.ValidationErrors
		fields := []string{}
		if errors.As(err, &verr) {
			for _, f := range verr {
				fields = append(fields, f.Field())
			}
		}
		Badrequestmsg := localization.GetMessage(lang, constants.BadRequestMessage, map[string]interface{}{
			"Fields": strings.Join(fields, ", "),
		})
		utils.SendBadRequest(ctx, constants.BadRequestErr, Badrequestmsg, constants.IsString, err)
		return err
	}

	return nil
}

func ValidateRequestDataParam(ctx *gin.Context, request interface{}) error {

	lang, _ := ctx.Get(constants.LanguageString)
	log := logger.GetLogger(ctx)

	err := ctx.ShouldBind(request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
		var verr validator.ValidationErrors
		fields := []string{}
		if errors.As(err, &verr) {
			for _, f := range verr {
				fields = append(fields, f.Field())
			}
		}
		Badrequestmsg := localization.GetMessage(lang, constants.BadRequestMessage, map[string]interface{}{
			"Fields": strings.Join(fields, ", "),
		})
		utils.SendBadRequest(ctx, constants.BadRequestErr, Badrequestmsg, constants.IsString, err)
		return err
	}

	return nil
}

func GenerateID() string {
	rand.Seed(time.Now().UTC().UnixNano())
	lowerCharSet := constants.ABCDLower
	numberSet := constants.Number
	allCharSet := lowerCharSet + "-" + numberSet
	minNum := 7
	IDLength := 20

	var password strings.Builder

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	remainingLength := IDLength - minNum
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func ValidateEmail(ctx context.Context, email string) (bool, error) {
	//user col
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	if email == "" {
		return false, errors.New("email is empty")
	}

	filter := bson.M{"email": email}

	// Check if the email is already in use
	exists, err := userCol.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}

	return exists > 0, nil
}

func ValidatePhoneNumber(ctx context.Context, countryCode, phone string) (bool, error) {
	//user col
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	if countryCode == "" || phone == "" {
		return false, errors.New("country code or phone number is empty")
	}

	filter := bson.M{"phone": phone, "countryCode": countryCode}

	// Check if the phone number is already in use
	exists, err := userCol.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error checking phone number existence: %w", err)
	}

	return exists > 0, nil
}

// JWT Claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
} // @claims

// generateJWT creates a new JWT token for the user
func GenerateJWT(userID primitive.ObjectID, email string) (string, int64, error) {
	logger := logger.GetLoggerWithoutContext()

	// get application config
	applicationConfig, err1 := configs.Get(constants.ApplicationConfig)
	if err1 != nil {
		logger.With(zap.Error(err1)).Error(constants.BindingFailedErrr)
	}

	// Get JWT secret from environment variable
	jwtSecret := applicationConfig.GetString(constants.JwtSecret)

	// Set token expiration time (24 hours from now)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims
	claims := &Claims{
		UserID: userID.Hex(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "users-service",
			Subject:   userID.Hex(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.UnixMilli(), nil
}

func GenerateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

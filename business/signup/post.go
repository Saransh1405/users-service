package signup

import (
	"errors"
	"sync"
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/helperfunctions"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// takes user data and signs up the user in the db
func Post(ctx *gin.Context, request *models.UserPostRequest) (*models.Users, error) {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//use user id as clientName
	clientName := primitive.NewObjectID()

	//create a new status log
	statusLogs := []models.StatusLogs{{
		Status:         "ACTIVE",
		ActionByUserId: clientName.Hex(),
		Timestamp:      time.Now().UnixMilli(),
		Notes:          "user signed up successfully",
	},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	//create channels to send validation results
	emailChan := make(chan struct {
		exists bool
		err    error
	}, 1)
	phoneChan := make(chan struct {
		exists bool
		err    error
	}, 1)

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.With(zap.Error(err)).Error("Password hashing failed")
		return nil, err
	}

	go func() {
		defer wg.Done()

		//validate email
		exists, err := helperfunctions.ValidateEmail(ctx, request.Email)
		emailChan <- struct {
			exists bool
			err    error
		}{exists, err}
	}()

	go func() {
		defer wg.Done()

		if request.Phone == "" {
			phoneChan <- struct {
				exists bool
				err    error
			}{exists: false, err: nil}
		} else {
			//validate phone number
			exists, err := helperfunctions.ValidatePhoneNumber(ctx, request.CountryCode, request.Phone)
			phoneChan <- struct {
				exists bool
				err    error
			}{exists, err}
		}
	}()

	// Wait for both validations to complete
	wg.Wait()

	// Get results from both channels
	emailResult := <-emailChan
	phoneResult := <-phoneChan

	// Check for validation errors
	if emailResult.err != nil {
		log.With(zap.Error(errors.New(constants.ErrorInValidatingEmail))).Error(constants.ErrorInValidatingEmail)
		return nil, errors.New(constants.ErrorInValidatingEmail)
	}

	if phoneResult.err != nil {
		log.With(zap.Error(phoneResult.err)).Error("Error validating phone number")
		return nil, errors.New(constants.PhoneNumberValidationFailed)
	}

	// Check if email or phone already exists
	if emailResult.exists {
		log.Info("Email already exists")
		return nil, errors.New(constants.EmailAlreadyExists)
	}

	if phoneResult.exists {
		log.Error("Phone number already exists")
		return nil, errors.New(constants.PhoneNumberValidationFailed)
	}

	//create a new user
	user := models.Users{
		ID:             clientName,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Password:       string(passwordHash),
		Email:          request.Email,
		Phone:          request.Phone,
		CountryCode:    request.CountryCode,
		UserProfileUrl: request.UserProfileUrl,
		Status:         models.Active,
		StatusLogs:     statusLogs,
		ClientName:     clientName.Hex(),
		CreatedAt:      time.Now().UnixMilli(),
	}

	//insert the user into the collection
	_, err = userCol.InsertOne(ctx, user)
	if err != nil {
		log.With(zap.Error(errors.New(constants.ErrorInInertingData))).Error(constants.ErrorInInertingData)
		return nil, errors.New(constants.ErrorInInertingData)
	}

	return &user, err
}

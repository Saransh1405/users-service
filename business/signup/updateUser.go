package signup

import (
	"context"
	"errors"
	"sync"
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/helperfunctions"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func UpdateUser(ctx context.Context, request *models.UserPatchRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//get the client name from the request
	client := request.ClientName
	if client == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	//convert client name to object id
	clientName, err := primitive.ObjectIDFromHex(client)
	if err != nil {
		log.Error("Error converting string to object id")
	}

	qry := bson.M{}

	if request.FirstName != "" {
		qry["firstName"] = request.FirstName
	}

	if request.LastName != "" {
		qry["lastName"] = request.LastName
	}

	if request.UserProfileUrl != "" {
		qry["userProfileUrl"] = request.UserProfileUrl
	}

	if request.ReasonForSuspension != "" {
		qry["reasonForSuspension"] = request.ReasonForSuspension
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

	go func() {
		defer wg.Done()

		if request.Email == "" {
			emailChan <- struct {
				exists bool
				err    error
			}{exists: false, err: nil}
		} else {
			//validate email
			exists, err := helperfunctions.ValidateEmail(ctx, request.Email)
			emailChan <- struct {
				exists bool
				err    error
			}{exists, err}
		}
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
		return errors.New(constants.ErrorInValidatingEmail)
	}

	if phoneResult.err != nil {
		log.With(zap.Error(phoneResult.err)).Error("Error validating phone number")
		return errors.New(constants.PhoneNumberValidationFailed)
	}

	// Check if email or phone already exists
	if emailResult.exists {
		log.Info("Email already exists")
		return errors.New(constants.EmailAlreadyExists)
	}

	if phoneResult.exists {
		log.Error("Phone number already exists")
		return errors.New(constants.PhoneNumberValidationFailed)
	}

	var status string
	if request.Status != "" {
		status = request.Status
		qry["status"] = request.Status
	}

	//create a new status log
	statusLogs := []models.StatusLogs{{
		Status:         status,
		ActionByUserId: clientName.Hex(),
		Timestamp:      time.Now().UnixMilli(),
		Notes:          "user updated successfully",
	},
	}

	update := bson.M{
		"$set": qry,
		"$push": bson.M{
			"statusLogs": bson.M{
				"$each": statusLogs,
			},
		},
	}

	//update the user into the collection
	_, err = userCol.UpdateOne(ctx, bson.M{"_id": clientName}, update)
	if err != nil {
		log.With(zap.Error(errors.New(constants.ErrorInInertingData))).Error(constants.ErrorInInertingData)
		return errors.New(constants.ErrorInInertingData)
	}

	return nil
}

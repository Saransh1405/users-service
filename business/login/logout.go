package login

import (
	"context"
	"errors"
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func Logout(ctx context.Context, request *models.GetUserRequest) error {
	//get the logger
	log := logger.GetLogger(ctx)

	//get the user col
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//get the client name from the context
	clientName := request.ClientName
	if clientName == "" {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	//convert the client name to object id
	objectClientName, err := primitive.ObjectIDFromHex(clientName)
	if err != nil {
		log.With(zap.Error(errors.New(constants.ErrorInConvertingToObjectId))).Error(constants.ErrorInConvertingToObjectId)
		return errors.New(constants.ErrorInConvertingToObjectId)
	}

	update := bson.M{
		"$set": bson.M{
			"status": models.LoggedOut,
		},
		"$push": bson.M{
			"statusLogs": bson.M{
				"$each": []models.StatusLogs{
					{
						Status:         "LOGGED_OUT",
						ActionByUserId: clientName,
						Timestamp:      time.Now().UnixMilli(),
						Notes:          "user logged out successfully",
					},
				},
			},
		},
	}

	_, err = userCol.UpdateByID(ctx, objectClientName, update)
	if err != nil {
		log.With(zap.Error(errors.New(constants.ErrorInInertingData))).Error(constants.ErrorInInertingData)
		return errors.New(constants.ErrorInInertingData)
	}

	return nil
}

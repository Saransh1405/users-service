package signup

import (
	"context"
	"errors"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/helperfunctions"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GetMyDetails(ctx context.Context, request *models.GetUserRequest) (*models.Users, int64, error) {

	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the user col
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//get the user id from the context
	client := request.ClientName
	//validate the client name
	count, err := helperfunctions.ValidateUser(ctx, client)
	if err != nil {
		log.With(zap.Error(err)).Error("Error validating user")
	}
	if count {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return nil, 0, errors.New(constants.UserNotFoundMessage)
	}

	//convert the client name to object id
	objectClientName, err := primitive.ObjectIDFromHex(client)
	if err != nil {
		log.Error("Error converting string to object id", zap.Error(err))
		return nil, 0, errors.New(constants.ErrorInConvertingToObjectId)
	}

	//get the user details
	var user models.Users
	err = userCol.FindOne(ctx, bson.M{"_id": objectClientName}).Decode(&user)
	if err != nil {
		log.Error("Error getting user details", zap.Error(err))
		return nil, 0, errors.New(constants.ErrorInGettingData)
	}

	return &user, 1, nil
}

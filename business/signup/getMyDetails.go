package signup

import (
	"context"
	"fmt"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GetMyDetails(ctx context.Context, request *models.GetUserRequest) (*models.Users, int64, error) {

	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the user col
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	fmt.Printf("request: %v\n", request)
	//get the user id from the context
	client := request.ClientName
	if client == "" {
		return nil, 0, fmt.Errorf("client name is required")
	}

	//convert the client name to object id
	objectClientName, err := primitive.ObjectIDFromHex(client)
	if err != nil {
		log.Error("Error converting string to object id", zap.Error(err))
		return nil, 0, err
	}

	//get the user details
	var user models.Users
	err = userCol.FindOne(ctx, bson.M{"_id": objectClientName}).Decode(&user)
	if err != nil {
		log.Error("Error getting user details", zap.Error(err))
		return nil, 0, err
	}

	return &user, 1, nil
}

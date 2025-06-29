package login

import (
	"context"
	"fmt"
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
		return fmt.Errorf("client name is required")
	}

	//convert the client name to object id
	objectClientName, err := primitive.ObjectIDFromHex(clientName)
	if err != nil {
		log.Error("Error converting string to object id", zap.Error(err))
		return err
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
		log.Error("Error updating user details", zap.Error(err))
		return err
	}

	return nil
}

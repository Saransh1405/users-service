package signup

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

func Delete(ctx context.Context, request *models.UserDeleteRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

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

	qry := bson.M{
		"status":            models.Deleted,
		"reasonForDeletion": request.ReasonForDeletion,
	}

	//create a new status log
	statusLogs := []models.StatusLogs{{
		Status:         "DELETED",
		ActionByUserId: clientName,
		Timestamp:      time.Now().UnixMilli(),
		Notes:          "user deleted successfully",
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

	//delete the user into the collection
	_, err = userCol.UpdateByID(ctx, objectClientName, update)
	if err != nil {
		log.Error("error inserting user")
	}

	return nil
}

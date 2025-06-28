package signup

import (
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Delete(ctx *gin.Context, request *models.UserDeleteRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//conver user id to object id
	clientName, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Error("Error converting string to object id")
	}

	//delete the user into the collection
	_, err = userCol.DeleteOne(ctx, bson.M{"_id": clientName})
	if err != nil {
		log.Error("error inserting user")
	}

	return nil
}

package signup

import (
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Patch(ctx *gin.Context, request *models.UserPatchRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//conver user id to object id
	clientName, err := primitive.ObjectIDFromHex(request.Id)
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

	if request.ProfilePictureUrl != "" {
		qry["profilePictureUrl"] = request.ProfilePictureUrl
	}

	if request.ReasonForSuspension != "" {
		qry["reasonForSuspension"] = request.ReasonForSuspension
	}

	var status string
	if request.Status != "" {
		status = request.Status
		qry["status"] = request.Status
	} else {
		status = "ACTIVE"
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
		log.Error("error inserting user")
	}

	return nil
}

package password

import (
	"context"
	"errors"
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

func ForgotPassword(ctx context.Context, request *models.ForgotPasswordRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the user collection
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

	//get the user by email
	user := models.Users{}
	err = userCol.FindOne(ctx, bson.M{"_id": clientName, "email": request.Email}).Decode(&user)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	// 2. Generate reset token
	resetToken := helperfunctions.GenerateResetToken()
	tokenExpiry := time.Now().Add(15 * time.Minute).UnixMilli()

	// 3. Store reset token in database
	update := bson.M{
		"$set": bson.M{
			"resetToken":       resetToken,
			"resetTokenExpiry": tokenExpiry,
		},
	}

	result, err := userCol.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		return err
	}

	// 4. Send reset email
	go func() {
		if result.MatchedCount > 0 {
			helperfunctions.SendNotifications(map[string]interface{}{
				"type":    models.NotificationTypeEmail,
				"userId":  user.ID,
				"message": "Password reset request",
				"data": map[string]interface{}{
					"To":       user.Email,
					"Message":  "Password reset request",
					"Template": "password-reset",
					"Variables": map[string]interface{}{
						"Name":   user.FirstName + " " + user.LastName,
						"Email":  user.Email,
						"Expire": time.Now().Add(15 * time.Minute).UnixMilli(),
					},
					"Subject": "Password reset request",
				},
			}, string(models.NotificationTopic))
		}
	}()

	return nil
}

package password

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
	"golang.org/x/crypto/bcrypt"
)

func UpdatePassword(ctx context.Context, request *models.UpdatePasswordRequest) error {
	//get the logger
	log := logger.GetLoggerWithoutContext()

	//get the user collection
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	//get the client name from the request
	client := request.ClientName

	//validate the client name
	count, err := helperfunctions.ValidateUser(ctx, client)
	if err != nil {
		log.With(zap.Error(err)).Error("Error validating user")
	}
	if count {
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

	//check if the user is found
	if user.ID == primitive.NilObjectID {
		log.With(zap.Error(errors.New(constants.UserNotFoundMessage))).Error(constants.UserNotFoundMessage)
		return errors.New(constants.UserNotFoundMessage)
	}

	//compare the passwrod with old hashed password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		log.With(zap.Error(errors.New(constants.InvalidPasswordMessage))).Error(constants.InvalidPasswordMessage)
		return errors.New(constants.InvalidPasswordMessage)
	}

	//hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ErrorInConvertingPassword)
		return errors.New(constants.ErrorInConvertingPassword)
	}

	//update the password
	result, err := userCol.UpdateOne(ctx, bson.M{"_id": clientName}, bson.M{"$set": bson.M{"password": string(hashedPassword)}})
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ErrorInConvertingPassword)
		return errors.New(constants.ErrorInConvertingPassword)
	}

	go func() {
		if result.ModifiedCount > 0 {
			helperfunctions.SendNotifications(map[string]interface{}{
				"type":    models.NotificationTypeEmail,
				"userId":  user.ID,
				"message": "Password updated successfully",
				"data": map[string]interface{}{
					"To":       user.Email,
					"Message":  "Password updated successfully",
					"Template": "password-updated",
					"Variables": map[string]interface{}{
						"Name":  user.FirstName + " " + user.LastName,
						"Email": user.Email,
					},
					"Subject": "Password updated successfully",
				},
			}, string(models.NotificationTopic))
		}
	}()

	return nil
}

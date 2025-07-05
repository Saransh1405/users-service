package password

import (
	"context"
	"errors"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"

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
	_, err = userCol.UpdateOne(ctx, bson.M{"_id": clientName}, bson.M{"$set": bson.M{"password": string(hashedPassword)}})
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ErrorInConvertingPassword)
		return errors.New(constants.ErrorInConvertingPassword)
	}

	return nil
}

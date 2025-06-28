package login

import (
	"errors"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/logger"
	"users-service/models"

	"users-service/utils/helperfunctions"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *gin.Context, request *models.LoginRequest) (*models.LoginResponse, error) {
	//get the logger
	log := logger.GetLogger(ctx)

	//get the user col
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)

	// Find user by email
	var user models.Users
	err := userCol.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		log.Error("User not found")
		return nil, errors.New("invalid credentials") // Generic error for security
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)
	if err != nil {
		log.Error("Invalid password")
		return nil, errors.New("invalid credentials") // Same generic error
	}

	//generate the token
	token, expireTime, err := helperfunctions.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Error("Failed to generate JWT token")
		return nil, errors.New("failed to generate authentication token")
	}

	response := models.LoginResponse{
		User:        user,
		AccessToken: token,
		ExpiresIn:   int64(expireTime),
	}

	return &response, nil
}

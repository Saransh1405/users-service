package password

import (
	"context"
	"errors"
	"time"
	"users-service/constants"
	"users-service/library/mongoDb"
	"users-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(ctx context.Context, request *models.ResetPasswordRequest) error {
	// 1. Validate reset token
	userCol := mongoDb.GetCollection(constants.MongoUserCollection)
	var user models.Users
	err := userCol.FindOne(ctx, bson.M{
		"resetToken":       request.ResetToken,
		"resetTokenExpiry": bson.M{"$gt": time.Now().UnixMilli()},
	}).Decode(&user)
	if err != nil {
		return errors.New(constants.InvalidResetTokenMessage)
	}

	// 2. Validate new password
	if request.NewPassword != request.ConfirmNewPassword {
		return errors.New(constants.PasswordDoesNotMatchMessage)
	}

	// 3. Hash new password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(constants.ErrorInConvertingPassword)
	}

	// 4. Update password and clear reset token
	update := bson.M{
		"$set": bson.M{
			"password":         string(newPasswordHash),
			"resetToken":       "",
			"resetTokenExpiry": 0,
		},
		"$push": bson.M{
			"statusLogs": bson.M{
				"$each": []models.StatusLogs{{
					Status:         "PASSWORD_RESET",
					ActionByUserId: user.ID.Hex(),
					Timestamp:      time.Now().UnixMilli(),
					Notes:          "password reset via email",
				}},
			},
		},
	}

	_, err = userCol.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	return err
}

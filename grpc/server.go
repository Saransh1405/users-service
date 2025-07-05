package grpc

import (
	"context"
	"users-service/business/login"
	"users-service/business/password"
	"users-service/business/signup"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/proto"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
}

// Signup implements user signup via gRPC
func (s *UserServiceServer) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {
	log := logger.GetLoggerWithoutContext()

	log.Info("gRPC Signup called", zap.String("email", req.Email))

	// Convert proto request to your internal model
	userReq := models.UserPostRequest{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		CountryCode:    req.CountryCode,
		UserProfileUrl: req.UserProfileUrl,
		Password:       req.Password,
	}

	ginCtx := &gin.Context{}

	//call the signup business logic
	result, err := signup.Post(ginCtx, &userReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	// Convert business logic response to proto response
	return &proto.SignupResponse{
		Id:             result.ID.Hex(), // Convert ObjectID to string
		Email:          result.Email,
		Password:       result.Password,
		FirstName:      result.FirstName,
		LastName:       result.LastName,
		CountryCode:    result.CountryCode,
		UserProfileUrl: result.UserProfileUrl,
		ClientName:     result.ClientName,
		Status:         string(result.Status),
	}, nil
}

// Login implements user login via gRPC
func (s *UserServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	log := logger.GetLoggerWithoutContext()

	ginCtx := &gin.Context{}

	loginReq := models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	//call the login business logic
	result, err := login.Login(ginCtx, &loginReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to login")
	}

	return &proto.LoginResponse{
		User: &proto.SignupResponse{
			Email:          result.User.Email,
			Password:       result.User.Password,
			FirstName:      result.User.FirstName,
			LastName:       result.User.LastName,
			CountryCode:    result.User.CountryCode,
			UserProfileUrl: result.User.UserProfileUrl,
			ClientName:     result.User.ClientName,
			Status:         string(result.User.Status),
		},
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	}, nil
}

// Logout implements user logout via gRPC
func (s *UserServiceServer) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	log := logger.GetLoggerWithoutContext()

	request := models.GetUserRequest{
		ClientName: req.ClientName,
	}

	// User ID is already in context from the interceptor
	err := login.Logout(ctx, &request)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to logout")
	}

	return &proto.LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}

// getUserDetails implements user getUserDetails via gRPC
func (s *UserServiceServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	log := logger.GetLoggerWithoutContext()

	request := models.GetUserRequest{
		ClientName: req.ClientName,
	}

	// User ID is already in context from the interceptor
	result, count, err := signup.GetMyDetails(ctx, &request)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to get user details")
	}

	return &proto.GetUserResponse{
		User: &proto.SignupResponse{
			Email:          result.Email,
			Password:       result.Password,
			FirstName:      result.FirstName,
			LastName:       result.LastName,
			CountryCode:    result.CountryCode,
			UserProfileUrl: result.UserProfileUrl,
			ClientName:     result.ClientName,
			Status:         string(result.Status),
		},
		Count: count,
	}, nil
}

// UpdateUser implements user UpdateUser via gRPC
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	log := logger.GetLoggerWithoutContext()

	userPatchReq := models.UserPatchRequest{
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		UserProfileUrl:      req.UserProfileUrl,
		Status:              req.Status,
		ReasonForSuspension: req.ReasonForSuspension,
		ClientName:          req.ClientName,
	}

	//call the logout business logic
	err := signup.UpdateUser(ctx, &userPatchReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to logout")
	}

	return &proto.UpdateUserResponse{
		Message: "Successfully updated user",
	}, nil
}

// UpdateUser implements user UpdateUser via gRPC
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	log := logger.GetLoggerWithoutContext()

	userDeleteReq := models.UserDeleteRequest{
		ID:                req.Id,
		ReasonForDeletion: req.ReasonForDelete,
		ClientName:        req.ClientName,
	}

	//call the logout business logic
	err := signup.Delete(ctx, &userDeleteReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to logout")
	}

	return &proto.DeleteUserResponse{
		Message: "Successfully deleted user",
	}, nil
}

// ForgotPassword implements user ForgotPassword via gRPC
func (s *UserServiceServer) ForgotPassword(ctx context.Context, req *proto.ForgotPasswordRequest) (*proto.ForgotPasswordResponse, error) {
	log := logger.GetLoggerWithoutContext()

	userForgotPasswordReq := models.ForgotPasswordRequest{
		ClientName: req.ClientName,
		Email:      req.Email,
	}

	//call the logout business logic
	err := password.ForgotPassword(ctx, &userForgotPasswordReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to send password reset email")
	}

	return &proto.ForgotPasswordResponse{
		Message: "Successfully sent password reset email",
	}, nil
}

// ForgotPassword implements user ForgotPassword via gRPC
func (s *UserServiceServer) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.ResetPasswordResponse, error) {
	log := logger.GetLoggerWithoutContext()

	userResetPasswordReq := models.ResetPasswordRequest{
		ClientName:         req.ClientName,
		ResetToken:         req.ResetToken,
		ConfirmNewPassword: req.ConfirmNewPassword,
		NewPassword:        req.Password,
	}

	//call the logout business logic
	err := password.ResetPassword(ctx, &userResetPasswordReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to send password reset email")
	}

	return &proto.ResetPasswordResponse{
		Message: "Successfully reset password",
	}, nil
}

// UpdatePassword implements user UpdatePassword via gRPC
func (s *UserServiceServer) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordRequest) (*proto.UpdatePasswordResponse, error) {
	log := logger.GetLoggerWithoutContext()

	userUpdatePasswordReq := models.UpdatePasswordRequest{
		ClientName:  req.ClientName,
		Password:    req.Password,
		NewPassword: req.NewPassword,
	}

	//call the logout business logic
	err := password.UpdatePassword(ctx, &userUpdatePasswordReq)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
		return nil, status.Error(codes.Internal, "Failed to send password reset email")
	}

	return &proto.UpdatePasswordResponse{
		Message: "Successfully updated password",
	}, nil
}

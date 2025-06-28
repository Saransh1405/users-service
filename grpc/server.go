package grpc

import (
	"context"
	"users-service/business/login"
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
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		CountryCode:       req.CountryCode,
		ProfilePictureUrl: req.UserProfileUrl,
		Password:          req.Password,
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
		UserProfileUrl: result.ProfilePictureUrl,
		ClientName:     result.ClientName,
		Status:         string(result.Status),
	}, nil
}

// Login implements user login via gRPC
func (s *UserServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	log := logger.GetLoggerWithoutContext()

	log.Info("gRPC Login called", zap.String("email", req.Email))

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
			UserProfileUrl: result.User.ProfilePictureUrl,
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

	log.Info("gRPC Logout called")

	// TODO: Call your existing logout business logic here

	return &proto.LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}

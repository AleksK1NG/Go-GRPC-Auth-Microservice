package server

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/pkg/grpc_errors"
	"github.com/AleksK1NG/auth-microservice/pkg/utils"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Register new user
func (u *usersServer) Register(ctx context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

	user, err := u.registerReqToUserModel(r)
	if err != nil {
		u.logger.Errorf("registerReqToUserModel: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "registerReqToUserModel: %v", err)
	}

	if err := utils.ValidateStruct(ctx, user); err != nil {
		u.logger.Errorf("ValidateStruct: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "ValidateStruct: %v", err)
	}

	createdUser, err := u.userUC.Register(ctx, user)
	if err != nil {
		u.logger.Errorf("userUC.Register: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "Register: %v", err)
	}

	return &userService.RegisterResponse{User: u.userModelToProto(createdUser)}, nil
}

// Login user with email and password
func (u *usersServer) Login(ctx context.Context, r *userService.LoginRequest) (*userService.LoginResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

	email := r.GetEmail()
	if !utils.ValidateEmail(email) {
		u.logger.Errorf("ValidateEmail: %v", email)
		return nil, status.Errorf(codes.InvalidArgument, "ValidateEmail: %v", email)
	}

	user, err := u.userUC.Login(ctx, email, r.GetPassword())
	if err != nil {
		u.logger.Errorf("userUC.Login: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "Login: %v", err)
	}

	return &userService.LoginResponse{User: u.userModelToProto(user)}, err
}

// Find user by email address
func (u *usersServer) FindByEmail(ctx context.Context, r *userService.FindByEmailRequest) (*userService.FindByEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

	email := r.GetEmail()
	if !utils.ValidateEmail(email) {
		u.logger.Errorf("ValidateEmail: %v", email)
		return nil, status.Errorf(codes.InvalidArgument, "ValidateEmail: %v", email)
	}

	user, err := u.userUC.FindByEmail(ctx, email)
	if err != nil {
		u.logger.Errorf("userUC.FindByEmail: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.FindByEmail: %v", err)
	}

	return &userService.FindByEmailResponse{User: u.userModelToProto(user)}, err
}

// Find user by uuid
func (u *usersServer) FindByID(ctx context.Context, r *userService.FindByIDRequest) (*userService.FindByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

	userUUID, err := uuid.Parse(r.GetUuid())
	if err != nil {
		u.logger.Errorf("uuid.Parse: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "uuid.Parse: %v", err)
	}

	user, err := u.userUC.FindById(ctx, userUUID)
	if err != nil {
		u.logger.Errorf("userUC.FindById: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.FindById: %v", err)
	}

	return &userService.FindByIDResponse{User: u.userModelToProto(user)}, nil
}

func (u *usersServer) registerReqToUserModel(r *userService.RegisterRequest) (*models.User, error) {
	candidate := &models.User{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Role:      r.GetRole(),
		Avatar:    r.GetAvatar(),
		Password:  r.GetPassword(),
	}

	if err := candidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return candidate, nil
}

func (u *usersServer) userModelToProto(user *models.User) *userService.User {
	userProto := &userService.User{
		Uuid:      user.UserID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	return userProto
}

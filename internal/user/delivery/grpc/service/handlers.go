package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/pkg/grpc_errors"
	"github.com/AleksK1NG/auth-microservice/pkg/utils"
	userService "github.com/AleksK1NG/auth-microservice/proto"
)

// Register new user
func (u *usersService) Register(ctx context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
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
func (u *usersService) Login(ctx context.Context, r *userService.LoginRequest) (*userService.LoginResponse, error) {
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

	session, err := u.sessUC.CreateSession(ctx, &models.Session{
		UserID: user.UserID,
	}, u.cfg.Session.Expire)
	if err != nil {
		u.logger.Errorf("sessUC.CreateSession: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessUC.CreateSession: %v", err)
	}

	return &userService.LoginResponse{User: u.userModelToProto(user), SessionId: session}, err
}

// Find user by email address
func (u *usersService) FindByEmail(ctx context.Context, r *userService.FindByEmailRequest) (*userService.FindByEmailResponse, error) {
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
func (u *usersService) FindByID(ctx context.Context, r *userService.FindByIDRequest) (*userService.FindByIDResponse, error) {
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

// Get session id from, ctx metadata, find user by uuid and returns it
func (u *usersService) GetMe(ctx context.Context, r *userService.GetMeRequest) (*userService.GetMeResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

	sessID, err := u.getSessionIDFromCtx(ctx)
	if err != nil {
		u.logger.Errorf("getSessionIDFromCtx: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessUC.getSessionIDFromCtx: %v", err)
	}

	session, err := u.sessUC.GetSessionByID(ctx, sessID)
	if err != nil {
		u.logger.Errorf("sessUC.GetSessionByID: %v", err)
		if errors.Is(err, redis.Nil) {
			return nil, status.Errorf(codes.NotFound, "sessUC.GetSessionByID: %v", grpc_errors.ErrNotFound)
		}
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessUC.GetSessionByID: %v", err)
	}

	user, err := u.userUC.FindById(ctx, session.UserID)
	if err != nil {
		u.logger.Errorf("userUC.FindById: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "userUC.FindById: %v", err)
	}

	return &userService.GetMeResponse{User: u.userModelToProto(user)}, nil
}

// Logout user, delete current session
func (u *usersService) Logout(ctx context.Context, request *userService.LogoutRequest) (*userService.LogoutResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

	sessID, err := u.getSessionIDFromCtx(ctx)
	if err != nil {
		u.logger.Errorf("getSessionIDFromCtx: %v", err)
		return nil, err
	}

	if err := u.sessUC.DeleteByID(ctx, sessID); err != nil {
		u.logger.Errorf("sessUC.DeleteByID: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessUC.DeleteByID: %v", err)
	}

	return &userService.LogoutResponse{}, nil
}

func (u *usersService) registerReqToUserModel(r *userService.RegisterRequest) (*models.User, error) {
	avatar := r.GetAvatar()
	candidate := &models.User{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Role:      r.GetRole(),
		Avatar:    &avatar,
		Password:  r.GetPassword(),
	}

	if err := candidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return candidate, nil
}

func (u *usersService) userModelToProto(user *models.User) *userService.User {
	userProto := &userService.User{
		Uuid:      user.UserID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.GetAvatar(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	return userProto
}

func (u *usersService) getSessionIDFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext: %v", grpc_errors.ErrNoCtxMetaData)
	}

	sessionID := md.Get("session_id")
	if sessionID[0] == "" {
		return "", status.Errorf(codes.PermissionDenied, "md.Get sessionId: %v", grpc_errors.ErrInvalidSessionId)
	}

	return sessionID[0], nil
}

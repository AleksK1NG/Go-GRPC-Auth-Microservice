package server

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Register new user
func (u *usersServer) Register(ctx context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Register")
	defer span.Finish()

	u.logger.Infof("Get request %s\n", r.String())
	user, err := u.registerReqToUserModel(r)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "registerReqToUserModel: %#v", err)
	}

	createdUser, err := u.userUC.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userService.RegisterResponse{
		User: u.userModelToProto(createdUser),
	}, nil
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

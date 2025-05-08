package service

import (
	pb "common/proto/gen/user"
	"context"
	"errors"
	"user-service/model"
	"user-service/service/jwt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterResponse{
		UserId:  int64(user.ID),
		Message: "user registered successfully",
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user model.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to query user")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(int64(user.ID), user.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		UserId:  int64(user.ID),
		Token:   token,
		Message: "login successful",
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	var user model.User
	if err := s.db.First(&user, req.UserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to query user")
	}
	return &pb.UserInfoResponse{
		UserId:   int64(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Message:  "success",
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var user model.User
	if err := s.db.First(&user, req.UserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.UpdateUserResponse{Success: false, Message: "user not found"}, nil
		}
		return &pb.UpdateUserResponse{Success: false, Message: "failed to query user"}, nil
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return &pb.UpdateUserResponse{Success: false, Message: "failed to hash password"}, nil
		}
		user.Password = string(hashedPassword)
	}
	if err := s.db.Save(&user).Error; err != nil {
		return &pb.UpdateUserResponse{Success: false, Message: "failed to update user"}, nil
	}
	return &pb.UpdateUserResponse{Success: true, Message: "user updated successfully"}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.db.Delete(&model.User{}, req.UserId).Error; err != nil {
		return &pb.DeleteUserResponse{Success: false, Message: "failed to delete user"}, nil
	}
	return &pb.DeleteUserResponse{Success: true, Message: "user deleted successfully"}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var users []model.User
	var total int64
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to count users")
	}
	if err := s.db.Limit(pageSize).Offset(offset).Order("id desc").Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}
	var pbUsers []*pb.UserInfo
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.UserInfo{
			UserId:    int64(u.ID),
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &pb.ListUsersResponse{
		Users:   pbUsers,
		Total:   int32(total),
		Message: "success",
	}, nil
}

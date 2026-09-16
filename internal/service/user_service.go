package service

import (
	"errors"
	"ukm-hub/internal/dto"
	"ukm-hub/internal/entity"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/utils"

	"github.com/google/uuid"
)

type UserService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID string) (*dto.UserResponse, error)
	UpdateProfile(userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// 2. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Simpan ke database
	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return s.toResponse(&user), nil
}

func (s *userService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Verifikasi Password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *s.toResponse(user),
	}, nil
}

func (s *userService) GetProfile(userID string) (*dto.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.toResponse(user), nil
}

func (s *userService) UpdateProfile(userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toResponse(user), nil
}

func (s *userService) toResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

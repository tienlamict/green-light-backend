package usecase

import (
	"context"
	"errors"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/utils"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
)

type AuthUseCase struct {
	userRepo   domain.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtManager *utils.JWTManager) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
	User  *domain.User
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := uc.jwtManager.GenerateToken(user.UserID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
		User:  user,
	}, nil
}

func (uc *AuthUseCase) GetUserInfo(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	return user, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password string, role domain.UserRole) (*domain.User, error) {
	// Check if email already exists
	_, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		UserID:       utils.GenerateUUIDv7(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

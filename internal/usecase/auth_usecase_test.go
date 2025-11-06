package usecase

import (
	"context"
	"errors"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/utils"
	"testing"

	"gorm.io/gorm"
)

// Mock repository
type mockUserRepository struct {
	users map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("email already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	for _, user := range m.users {
		if user.UserID == userID {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, userID string) error {
	for email, user := range m.users {
		if user.UserID == userID {
			delete(m.users, email)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

func TestAuthUseCase_Login(t *testing.T) {
	// Setup
	mockRepo := newMockUserRepository()
	jwtManager := utils.NewJWTManager("test-secret", 24)
	authUC := NewAuthUseCase(mockRepo, jwtManager)

	// Create test user
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := &domain.User{
		UserID:       "test-user-id",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         domain.RoleAdmin,
	}
	mockRepo.users[testUser.Email] = testUser

	tests := []struct {
		name    string
		input   LoginInput
		wantErr bool
		errType error
	}{
		{
			name: "successful login",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			input: LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantErr: true,
			errType: ErrInvalidCredentials,
		},
		{
			name: "invalid password",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
			errType: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := authUC.Login(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("expected error %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if output == nil {
					t.Errorf("expected output, got nil")
				}
				if output != nil && output.Token == "" {
					t.Errorf("expected token, got empty string")
				}
			}
		})
	}
}

func TestAuthUseCase_Register(t *testing.T) {
	// Setup
	mockRepo := newMockUserRepository()
	jwtManager := utils.NewJWTManager("test-secret", 24)
	authUC := NewAuthUseCase(mockRepo, jwtManager)

	tests := []struct {
		name     string
		email    string
		password string
		role     domain.UserRole
		wantErr  bool
		errType  error
	}{
		{
			name:     "successful registration",
			email:    "newuser@example.com",
			password: "password123",
			role:     domain.RoleEditor,
			wantErr:  false,
		},
		{
			name:     "duplicate email",
			email:    "newuser@example.com",
			password: "password123",
			role:     domain.RoleEditor,
			wantErr:  true,
			errType:  ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authUC.Register(context.Background(), tt.email, tt.password, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("expected error %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("expected user, got nil")
				}
				if user != nil && user.Email != tt.email {
					t.Errorf("expected email %s, got %s", tt.email, user.Email)
				}
			}
		})
	}
}

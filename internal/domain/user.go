package domain

import (
	"context"
	"time"
)

// UserRole defines the role enum
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleEditor UserRole = "editor"
)

// User represents a user entity
type User struct {
	UserID       string    `gorm:"primaryKey;type:varchar(36)" json:"user_id"`
	Email        string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole  `gorm:"type:enum('admin','editor');default:'editor';not null" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, userID string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context, limit, offset int) ([]*User, int64, error)
}

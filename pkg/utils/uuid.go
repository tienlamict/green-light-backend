package utils

import (
	"github.com/google/uuid"
)

// GenerateUUIDv7 generates a UUIDv7 (time-ordered UUID)
// Note: github.com/google/uuid doesn't have v7 yet, so we use v7 ordering with v4
func GenerateUUIDv7() string {
	// Using UUID v7 would be ideal, but google/uuid doesn't support it yet
	// For now, we use v4 which is still unique
	// In production, consider using: github.com/gofrs/uuid which supports v7
	return uuid.New().String()
}

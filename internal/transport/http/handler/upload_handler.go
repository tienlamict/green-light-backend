package handler

import (
	"fmt"
	"green-light-backend/internal/transport/http/dto"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir   string
	maxFileSize int64
}

func NewUploadHandler(uploadDir string, maxFileSize int64) *UploadHandler {
	return &UploadHandler{
		uploadDir:   uploadDir,
		maxFileSize: maxFileSize,
	}
}

// Upload handles file upload (admin/editor only)
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("No file provided"))
		return
	}

	// Check file size
	if file.Size > h.maxFileSize {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", h.maxFileSize)))
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid file type. Allowed: jpg, jpeg, png, gif, webp"))
		return
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create upload directory"))
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(h.uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to save file"))
		return
	}

	// Return file URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)

	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{
		"url": fileURL,
	}, "File uploaded successfully"))
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
}

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func (s *Service) badRequest(c *gin.Context, message string) {
	s.respondError(c, 400, "BAD_REQUEST", message)
}

func (s *Service) notFound(c *gin.Context, message string) {
	s.respondError(c, 404, "NOT_FOUND", message)
}

func (s *Service) internalError(c *gin.Context, message string) {
	s.respondError(c, 500, "INTERNAL_ERROR", message)
}

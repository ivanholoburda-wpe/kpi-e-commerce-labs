package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lab2/internal/api/dto"
	"lab2/internal/api/middleware"
)

// RegistrationService defines the business operations the handler depends on.
type RegistrationService interface {
	Start(ctx context.Context, phone string, telegramChatID int64) error
	Verify(ctx context.Context, phone string, code string) (string, error)
	Setup(ctx context.Context, accountID uuid.UUID, password string) (string, error)
	ConfirmTOTP(ctx context.Context, accountID uuid.UUID, code string) error
}

type RegisterHandler struct {
	service RegistrationService
}

func NewRegisterHandler(service RegistrationService) *RegisterHandler {
	return &RegisterHandler{service: service}
}

func (h *RegisterHandler) Start(c *gin.Context) {
	var req dto.StartRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.Start(c.Request.Context(), req.Phone, req.TelegramChatID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "verification code sent"})
}

func (h *RegisterHandler) Verify(c *gin.Context) {
	var req dto.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.Verify(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func (h *RegisterHandler) Setup(c *gin.Context) {
	var req dto.SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	totpURI, err := h.service.Setup(c.Request.Context(), accountID, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SetupResponse{TOTPUri: totpURI})
}

func (h *RegisterHandler) ConfirmTOTP(c *gin.Context) {
	var req dto.ConfirmTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	if err := h.service.ConfirmTOTP(c.Request.Context(), accountID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "registration complete"})
}

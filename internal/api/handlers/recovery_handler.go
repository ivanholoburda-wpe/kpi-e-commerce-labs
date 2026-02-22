package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lab2/internal/api/dto"
	"lab2/internal/api/middleware"
)

// AccountRecoveryService defines the business operations the recovery handler depends on.
type AccountRecoveryService interface {
	Start(ctx context.Context, phone string) error
	VerifyOTP(ctx context.Context, phone, code string) (string, error)
	ConfirmPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, error)
	ResetPassword(ctx context.Context, accountID uuid.UUID, newPassword string) error
}

type RecoveryHandler struct {
	service AccountRecoveryService
}

func NewRecoveryHandler(service AccountRecoveryService) *RecoveryHandler {
	return &RecoveryHandler{service: service}
}

func (h *RecoveryHandler) Start(c *gin.Context) {
	var req dto.StartRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.Start(c.Request.Context(), req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "recovery code sent"})
}

func (h *RecoveryHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyRecoveryOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.VerifyOTP(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func (h *RecoveryHandler) ConfirmPIN(c *gin.Context) {
	var req dto.ConfirmRecoveryPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	token, err := h.service.ConfirmPIN(c.Request.Context(), accountID, req.Pin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func (h *RecoveryHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	if err := h.service.ResetPassword(c.Request.Context(), accountID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "password updated"})
}

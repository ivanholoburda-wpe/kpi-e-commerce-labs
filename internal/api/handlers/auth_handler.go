package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"lab2/internal/api/dto"
	"lab2/internal/api/middleware"
	"lab2/internal/models"
)

// AuthenticationService defines the business operations the auth handler depends on.
type AuthenticationService interface {
	Login(ctx context.Context, phone, password string) (string, error)
	Verify2FA(ctx context.Context, accountID uuid.UUID, code string) (token, nextStep string, err error)
	SetupPIN(ctx context.Context, accountID uuid.UUID, pin string) (sessionToken, accessToken string, err error)
	VerifyPIN(ctx context.Context, accountID uuid.UUID, pin string) (string, error)
	GetMe(ctx context.Context, accountID uuid.UUID) (models.Account, error)
}

type AuthHandler struct {
	service AuthenticationService
}

func NewAuthHandler(service AuthenticationService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req dto.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	token, nextStep, err := h.service.Verify2FA(c.Request.Context(), accountID, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Verify2FAResponse{
		Token:    token,
		NextStep: nextStep,
	})
}

func (h *AuthHandler) SetupPIN(c *gin.Context) {
	var req dto.SetupPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	sessionToken, accessToken, err := h.service.SetupPIN(c.Request.Context(), accountID, req.Pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthTokensResponse{
		SessionToken: sessionToken,
		AccessToken:  accessToken,
	})
}

func (h *AuthHandler) VerifyPIN(c *gin.Context) {
	var req dto.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	accessToken, err := h.service.VerifyPIN(c.Request.Context(), accountID, req.Pin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: accessToken})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	accountID := c.MustGet(middleware.AccountIDKey).(uuid.UUID)

	account, err := h.service.GetMe(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AccountResponse{
		ID:          account.ID.String(),
		Phone:       account.Phone,
		TOTPEnabled: account.TOTPEnabled,
		PinSet:      account.PinHash != nil,
		Status:      string(account.Status),
	})
}

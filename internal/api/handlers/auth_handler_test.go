package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"lab2/internal/api/dto"
	"lab2/internal/api/handlers"
	"lab2/internal/api/middleware"
	"lab2/internal/mocks"
	"lab2/internal/models"
)

func TestAuthHandler_Login(t *testing.T) {
	svc := new(mocks.AuthenticationService)
	h := handlers.NewAuthHandler(svc)

	svc.On("Login", mock.Anything, "+380991234567", "MyStr0ngP@ss").Return("2fa-token", nil)

	body, _ := json.Marshal(dto.LoginRequest{
		Phone:    "+380991234567",
		Password: "MyStr0ngP@ss",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "2fa-token", resp.Token)

	svc.AssertExpectations(t)
}

func TestAuthHandler_Verify2FA(t *testing.T) {
	svc := new(mocks.AuthenticationService)
	h := handlers.NewAuthHandler(svc)

	accountID := uuid.New()
	svc.On("Verify2FA", mock.Anything, accountID, "123456").Return("session-token", "verify_pin", nil)

	body, _ := json.Marshal(dto.Verify2FARequest{Code: "123456"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/verify-2fa", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.Verify2FA(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.Verify2FAResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "session-token", resp.Token)
	assert.Equal(t, "verify_pin", resp.NextStep)

	svc.AssertExpectations(t)
}

func TestAuthHandler_SetupPIN(t *testing.T) {
	svc := new(mocks.AuthenticationService)
	h := handlers.NewAuthHandler(svc)

	accountID := uuid.New()
	svc.On("SetupPIN", mock.Anything, accountID, "1234").Return("session-tok", "access-tok", nil)

	body, _ := json.Marshal(dto.SetupPINRequest{Pin: "1234"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/setup-pin", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.SetupPIN(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.AuthTokensResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "session-tok", resp.SessionToken)
	assert.Equal(t, "access-tok", resp.AccessToken)

	svc.AssertExpectations(t)
}

func TestAuthHandler_VerifyPIN(t *testing.T) {
	svc := new(mocks.AuthenticationService)
	h := handlers.NewAuthHandler(svc)

	accountID := uuid.New()
	svc.On("VerifyPIN", mock.Anything, accountID, "1234").Return("new-access-tok", nil)

	body, _ := json.Marshal(dto.VerifyPINRequest{Pin: "1234"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/verify-pin", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.VerifyPIN(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "new-access-tok", resp.Token)

	svc.AssertExpectations(t)
}

func TestAuthHandler_GetMe(t *testing.T) {
	svc := new(mocks.AuthenticationService)
	h := handlers.NewAuthHandler(svc)

	accountID := uuid.New()
	pin := "hashed"
	account := models.Account{
		ID:          accountID,
		Phone:       "+380991234567",
		TOTPEnabled: true,
		PinHash:     &pin,
		Status:      models.StatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	svc.On("GetMe", mock.Anything, accountID).Return(account, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	c.Set(middleware.AccountIDKey, accountID)

	h.GetMe(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.AccountResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, accountID.String(), resp.ID)
	assert.Equal(t, "+380991234567", resp.Phone)
	assert.True(t, resp.TOTPEnabled)
	assert.True(t, resp.PinSet)
	assert.Equal(t, "active", resp.Status)

	svc.AssertExpectations(t)
}

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"lab2/internal/api/dto"
	"lab2/internal/api/handlers"
	"lab2/internal/api/middleware"
	"lab2/internal/mocks"
)

func TestRecoveryHandler_Start(t *testing.T) {
	svc := new(mocks.AccountRecoveryService)
	h := handlers.NewRecoveryHandler(svc)

	svc.On("Start", mock.Anything, "+380991234567").Return(nil)

	body, _ := json.Marshal(dto.StartRecoveryRequest{Phone: "+380991234567"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/recovery/start", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Start(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.MessageResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "recovery code sent", resp.Message)

	svc.AssertExpectations(t)
}

func TestRecoveryHandler_VerifyOTP(t *testing.T) {
	svc := new(mocks.AccountRecoveryService)
	h := handlers.NewRecoveryHandler(svc)

	svc.On("VerifyOTP", mock.Anything, "+380991234567", "654321").Return("recovery-token", nil)

	body, _ := json.Marshal(dto.VerifyRecoveryOTPRequest{
		Phone: "+380991234567",
		Code:  "654321",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/recovery/verify-otp", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.VerifyOTP(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "recovery-token", resp.Token)

	svc.AssertExpectations(t)
}

func TestRecoveryHandler_ConfirmPIN(t *testing.T) {
	svc := new(mocks.AccountRecoveryService)
	h := handlers.NewRecoveryHandler(svc)

	accountID := uuid.New()
	svc.On("ConfirmPIN", mock.Anything, accountID, "1234").Return("confirmed-token", nil)

	body, _ := json.Marshal(dto.ConfirmRecoveryPINRequest{Pin: "1234"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/recovery/confirm-pin", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.ConfirmPIN(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "confirmed-token", resp.Token)

	svc.AssertExpectations(t)
}

func TestRecoveryHandler_ResetPassword(t *testing.T) {
	svc := new(mocks.AccountRecoveryService)
	h := handlers.NewRecoveryHandler(svc)

	accountID := uuid.New()
	svc.On("ResetPassword", mock.Anything, accountID, "NewStr0ngP@ss").Return(nil)

	body, _ := json.Marshal(dto.ResetPasswordRequest{NewPassword: "NewStr0ngP@ss"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/recovery/reset-password", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.ResetPassword(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.MessageResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "password updated", resp.Message)

	svc.AssertExpectations(t)
}

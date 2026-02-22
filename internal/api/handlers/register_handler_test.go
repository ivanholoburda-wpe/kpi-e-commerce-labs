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

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterHandler_Start(t *testing.T) {
	svc := new(mocks.RegistrationService)
	h := handlers.NewRegisterHandler(svc)

	svc.On("Start", mock.Anything, "+380991234567", int64(123)).Return(nil)

	body, _ := json.Marshal(dto.StartRegistrationRequest{
		Phone:          "+380991234567",
		TelegramChatID: 123,
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register/start", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Start(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.MessageResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "verification code sent", resp.Message)

	svc.AssertExpectations(t)
}

func TestRegisterHandler_Verify(t *testing.T) {
	svc := new(mocks.RegistrationService)
	h := handlers.NewRegisterHandler(svc)

	svc.On("Verify", mock.Anything, "+380991234567", "123456").Return("jwt-token", nil)

	body, _ := json.Marshal(dto.VerifyCodeRequest{
		Phone: "+380991234567",
		Code:  "123456",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register/verify", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Verify(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TokenResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "jwt-token", resp.Token)

	svc.AssertExpectations(t)
}

func TestRegisterHandler_Setup(t *testing.T) {
	svc := new(mocks.RegistrationService)
	h := handlers.NewRegisterHandler(svc)

	accountID := uuid.New()
	svc.On("Setup", mock.Anything, accountID, "MyStr0ngP@ss").Return("otpauth://totp/Lab2", nil)

	body, _ := json.Marshal(dto.SetupRequest{Password: "MyStr0ngP@ss"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register/setup", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.Setup(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SetupResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "otpauth://totp/Lab2", resp.TOTPUri)

	svc.AssertExpectations(t)
}

func TestRegisterHandler_ConfirmTOTP(t *testing.T) {
	svc := new(mocks.RegistrationService)
	h := handlers.NewRegisterHandler(svc)

	accountID := uuid.New()
	svc.On("ConfirmTOTP", mock.Anything, accountID, "482910").Return(nil)

	body, _ := json.Marshal(dto.ConfirmTOTPRequest{Code: "482910"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register/confirm-totp", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(middleware.AccountIDKey, accountID)

	h.ConfirmTOTP(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.MessageResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "registration complete", resp.Message)

	svc.AssertExpectations(t)
}

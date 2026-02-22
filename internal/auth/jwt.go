package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	PurposeRegistration      = "registration_setup"
	Purpose2FA               = "2fa_verification"
	PurposePinSetup          = "pin_setup"
	PurposeSession           = "session"
	PurposeAccess            = "access"
	PurposeRecoveryOTP       = "recovery_otp"
	PurposeRecoveryConfirmed = "recovery_confirmed"
)

type Claims struct {
	jwt.RegisteredClaims
	AccountID uuid.UUID `json:"account_id"`
	Purpose   string    `json:"purpose"`
}

func GenerateToken(secret []byte, accountID uuid.UUID, purpose string, ttl time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		AccountID: accountID,
		Purpose:   purpose,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(secret []byte, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

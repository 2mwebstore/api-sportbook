package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte { return []byte(os.Getenv("JWT_SECRET")) }

func accessTTL() time.Duration {
	v, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL_MINUTES"))
	if v == 0 {
		v = 15
	}
	return time.Duration(v) * time.Minute
}

func refreshTTL() time.Duration {
	v, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL_DAYS"))
	if v == 0 {
		v = 7
	}
	return time.Duration(v) * 24 * time.Hour
}

// GenerateTokenPair issues a short-lived access token and a long-lived refresh token.
func GenerateTokenPair(userID uint, role, email, phone string) (access, refresh string, expiresIn int64, err error) {
	now := time.Now()

	aTTL := accessTTL()
	aClaims := Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userID)),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(aTTL)),
		},
	}
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, aClaims)
	access, err = aToken.SignedString(jwtSecret())
	if err != nil {
		return
	}
	expiresIn = int64(aTTL.Seconds())

	rTTL := refreshTTL()
	rClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userID)),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(rTTL)),
		},
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	refresh, err = rToken.SignedString(jwtSecret())
	return
}

// ValidateToken parses and validates a JWT string, returning its claims.
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

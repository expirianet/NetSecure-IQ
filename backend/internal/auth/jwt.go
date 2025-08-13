package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         string  `json:"uid"`
	RoleID         int     `json:"rid"`
	OrganizationID *string `json:"org,omitempty"`
	jwt.RegisteredClaims
}

func IssueJWT(secret []byte, userID string, roleID int, organizationID *string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:         userID,
		RoleID:         roleID,
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseJWT(secret []byte, tokenString string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid_token")
	}
	return claims, nil
}

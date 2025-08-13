package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	CtxUserID         ctxKey = "user_id"
	CtxRoleName       ctxKey = "role"
	CtxRoleID         ctxKey = "role_id"
	CtxOrganizationID ctxKey = "organization_id"
)

func JWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header.Get("Authorization")
		if !strings.HasPrefix(ah, "Bearer ") {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		raw := strings.TrimPrefix(ah, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) == 0 {
			http.Error(w, "JWT secret not configured", http.StatusInternalServerError)
			return
		}
		token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		ctx := r.Context()
		if v, _ := claims["user_id"].(string); v != "" {
			ctx = context.WithValue(ctx, CtxUserID, v)
		}
		if v, _ := claims["role"].(string); v != "" {
			ctx = context.WithValue(ctx, CtxRoleName, v)
		}
		if v, ok := claims["role_id"].(float64); ok {
			ctx = context.WithValue(ctx, CtxRoleID, int(v))
		}
		if v, _ := claims["organization_id"].(string); v != "" {
			ctx = context.WithValue(ctx, CtxOrganizationID, v)
		}
		next(w, r.WithContext(ctx))
	}
}

func GetUserID(r *http.Request) string {
	if v := r.Context().Value(CtxUserID); v != nil { return v.(string) }
	return ""
}
func GetRoleName(r *http.Request) string {
	if v := r.Context().Value(CtxRoleName); v != nil { return v.(string) }
	return ""
}
func GetRoleID(r *http.Request) int {
	if v := r.Context().Value(CtxRoleID); v != nil { return v.(int) }
	return 0
}
func GetOrgID(r *http.Request) string {
	if v := r.Context().Value(CtxOrganizationID); v != nil { return v.(string) }
	return ""
}

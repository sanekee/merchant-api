package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/sanekee/merchant-api/backend/internal/log"
)

const (
	jwtHMACSecret = "my-super-secret"
)

func JWTAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			log.Error("auth header not found")
			writeForbidden(w)
			return
		}
		auths := strings.Split(authHeader, " ")
		if len(auths) != 2 && auths[0] != "Bearer" {
			log.Error("bearer token not found")
			writeForbidden(w)
			return
		}

		token, err := jwt.Parse(auths[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtHMACSecret), nil
		})
		if err != nil {
			log.Error("token parse error %s", err.Error())
			writeForbidden(w)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Error("invalid token token parse error")
			writeForbidden(w)
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok || len(userID) == 0 {
			log.Error("user id not found")
			writeForbidden(w)
			return
		}
		fn(w, r)
	}
}

func writeForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("forbidden area"))
}

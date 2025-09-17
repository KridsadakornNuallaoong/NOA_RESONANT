package sensitive

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userContextKey contextKey = "user"

// VerifyJWT verifies the JWT token
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte("User Login") // ใช้คีย์เดียวกับที่ใช้สร้าง JWT

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าใช้ SigningMethod ที่ถูกต้อง
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า Token ถูกต้องและไม่หมดอายุ
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AuthMiddleware is a middleware for validating JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// แยกคำว่า "Bearer" ออกจาก Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจสอบ JWT Token
		claims, err := VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// เพิ่มข้อมูล Claims ลงใน Context
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProtectedResource(w http.ResponseWriter, r *http.Request) {
	// ดึง Claims จาก Context
	claims, ok := r.Context().Value(userContextKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ใช้ข้อมูลจาก Claims
	username := claims["username"].(string)
	log.Println("Accessing protected resource for user:", username)

	response := map[string]string{
		"message": "Welcome " + username,
		"userID":  claims["userID"].(string),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

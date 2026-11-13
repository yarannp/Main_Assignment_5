package middlewares

import (
	"aaa_bekend/internal/auth"
	"aaa_bekend/internal/lmodel"
	"net/http"
	"strings"
)

func AuthN(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Formatnya adalah "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 3. Validasi token menggunakan fungsi dari package auth
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Pindahkan data dari token (claims) ke struct lmodel.UserContext
		userCtx := lmodel.UserContext{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		// 5. Titipkan data user yang berhasil login ke dalam r.Context()
		ctx := lmodel.ContextWithUser(r.Context(), userCtx)

		// 6. Request berikutnya (handler) dengan context yang sudah diisi
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

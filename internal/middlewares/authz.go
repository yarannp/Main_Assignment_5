package middlewares

import (
	"aaa_bekend/internal/lmodel"
	"net/http"
)

type RoleGuard struct {
	next         http.Handler
	allowedRoles []string
}

func (rg *RoleGuard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil data user yang sudah dititipkan oleh authn.go di dalam Context
	user, err := lmodel.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	// 2. Cek apakah Role milik user ini ada di dalam daftar allowedRoles
	for _, role := range rg.allowedRoles {
		if user.Role == role {
			// Jika rolenya cocok, izinkan request lanjut ke tujuan (handler)
			rg.next.ServeHTTP(w, r)
			return
		}
	}

	// 3. Jika looping selesai tapi tidak ada role yang cocok, blokir aksesnya!
	http.Error(w, "Forbidden: You don't have enough permission to access this resource", http.StatusForbidden)
}

func RequireRole(allowedRules ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &RoleGuard{
			next:         next,
			allowedRoles: allowedRules,
		}
	}
}

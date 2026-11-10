package handler

import (
	"aaa_bekend/internal/db"
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"net/http"
	"time"
)

// Struct tambahan untuk menangkap request body saat update data
type UpdateMyDataReq struct {
	Username string `json:"username"`
}

func GetMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// 1. Ambil data profil dari database (Hak Akses)[cite: 3]
	profile, err := db.GetUserByID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// 1. Parse JSON body yang dikirim oleh user
	var req UpdateMyDataReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Lakukan proses pembaruan data di database (Hak Perbaikan)[cite: 3]
	err := db.UpdateMyData(r.Context(), user.ID, req.Username)
	if err != nil {
		http.Error(w, "Failed to update data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data updated successfully",
		"username": req.Username,
	})
}

func GetMyAuditLog(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// 1. Ambil riwayat akses/log aktivitas dari database[cite: 3]
	logs, err := db.GetUserAuditLogs(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}

	// 2. Mencegah nilai log menjadi null jika masih kosong
	if logs == nil {
		logs = []lmodel.AuditLog{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user.Username,
		"total": len(logs),
		"logs":  logs,
	})
}

func ExportMyData(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// 1. Ambil seluruh entitas data milik user (Hak Portabilitas)[cite: 3]
	profile, _ := db.GetUserByID(r.Context(), user.ID)
	projects, _ := db.GetUserProjects(r.Context(), user.ID)
	logs, _ := db.GetUserAuditLogs(r.Context(), user.ID)

	if projects == nil {
		projects = []lmodel.Project{}
	}
	if logs == nil {
		logs = []lmodel.AuditLog{}
	}

	export := map[string]interface{}{
		"exported_at": time.Now(),
		"profile":     profile,
		"projects":    projects,
		"audit_logs":  logs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=my_data_export.json")
	json.NewEncoder(w).Encode(export)
}

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := lmodel.UserFromContext(r.Context())

	// 1. Hapus akun dan seluruh data terkait secara permanen (Hak Penghapusan)[cite: 3]
	err := db.DeleteUserAccount(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your account and all data have been permanently deleted",
	})
}

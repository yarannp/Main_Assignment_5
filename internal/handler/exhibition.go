package handler

import (
	"aaa_bekend/internal/db"
	"aaa_bekend/internal/lmodel"
	"encoding/json"

	// "errors"
	"net/http"
	"strconv"
	// "github.com/jackc/pgx/v5"
)

type SubmitProjectReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GradeProjectReq struct {
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

// POST /api/v1/projects

func SubmitProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())
	// 1. Parse JSON body dari request
	var req SubmitProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Simpan proyek ke database atas nama user yang sedang login[cite: 2]
	// Asumsi db.CreateProject mengeksekusi INSERT dan mengembalikan projectID
	projectID, err := db.CreateProject(r.Context(), req.Title, req.Description, userCtx.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 3. Kembalikan respons 201 beserta project_id[cite: 2]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"project_id": projectID})
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	// TODO
	// 1. Ambil seluruh daftar proyek dari database
	projects, err := db.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO
	// 2. Pastikan array tidak bernilai nil agar menjadi [] di JSON, bukan null
	if projects == nil {
		projects = []lmodel.Project{} // Sesuaikan dengan tipe struct Project Anda
	}

	// 3. Mengembalikan daftar proyek dalam format JSON[cite: 2]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&projects)
}

func GetAllProjectGrades(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil seluruh nilai proyek dari database
	grades, err := db.GetAllProjectGrades(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if grades == nil {
		grades = []lmodel.ProjectGrade{} // Tipe data disesuaikan dengan kembalian fungsi db.go
	}

	// 2. Mengembalikan nilai seluruh proyek dalam format JSON[cite: 2]
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&grades)
}

func GradeProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// 1. Ambil ID dari URL
	idParam := r.PathValue("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request: Invalid Project ID", http.StatusBadRequest)
		return
	}

	// 2. Parse skor dan feedback dari request body
	var req GradeProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
		return
	}

	// 3. Simpan nilai ke database[cite: 2]
	// Asumsi db.SubmitGrade mengeksekusi INSERT/UPDATE ke tabel grades
	err = db.CreateGrade(r.Context(), projectID, userCtx.ID, req.Score, req.Feedback)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. Respons 201 dengan pesan "Grade Submitted"[cite: 3]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Grade Submitted"})
}

func GetProjectGrades(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// TODO
	idParam := r.PathValue("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Bad Request: Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Not Found: Project does not exist", http.StatusNotFound)
		return
	}

	if userCtx.Role == "student" {
		if project.StudentID != userCtx.ID {
			http.Error(w, "Forbidden: you can only view grades for your own projects", http.StatusForbidden)
			return
		}
	}

	grades, err := db.GetGradesForProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
		return
	}

	if grades == nil {
		grades = []lmodel.Grade{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// TODO

// if project.StudentID != userCtx.ID {
// http.Error(w, "Forbidden: you can only view grades for your own projects", http.StatusForbidden)
// return
// }
// grades, err := db.GetGradesForProject(r.Context(), projectID)
// if err != nil {
// http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
// return
// }

// if grades == nil {
// grades = []lmodel.Grade{}
//}
//json.NewEncoder(w).Encode(grades)

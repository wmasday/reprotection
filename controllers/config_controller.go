package controllers

import (
	"net/http"
	"reprotection/config"
)

func StoreConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		workingProject := r.FormValue("working_project")
		_, err := config.DB.Exec("INSERT INTO config (working_project) VALUES (?) ON DUPLICATE KEY UPDATE working_project = VALUES(working_project)", workingProject)
		if err != nil {
			http.Error(w, "Failed to save working project", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
} 
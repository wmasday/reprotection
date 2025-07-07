package controllers

import (
	"net/http"
	"reprotection/config"
)

func StoreConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		workingProject := r.FormValue("working_project")
		_, err := config.DB.Exec("UPDATE config SET working_project = ? WHERE id = 1", workingProject)
		if err != nil {
			http.Error(w, "Failed to update working project", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetWorkingProject() (string, error) {
	var workingProject string
	err := config.DB.QueryRow("SELECT working_project FROM config WHERE id = 1").Scan(&workingProject)
	return workingProject, err
} 
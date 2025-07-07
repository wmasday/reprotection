package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"reprotection/config"
	"reprotection/models"
)

type ViewData struct {
	User           string
	Items          interface{}
	Timestamp      int64
	ChartData      string
	WorkingProject *models.Project
}

func Index(w http.ResponseWriter, r *http.Request) {
	items, err := models.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	workingProject, _ := models.GetWorkingProject()
	if workingProject == nil {
		workingProject = &models.Project{ID: 1, WorkingProject: "(not set)"}
	}

	session, _ := config.Store.Get(r, "session")
	username := ""
	if v, ok := session.Values["username"].(string); ok {
		username = v
	}

	viewData := ViewData{
		User:           username,
		Items:          items,
		Timestamp:      time.Now().Unix(),
		ChartData:      "",
		WorkingProject: workingProject,
	}

	render(w, r, "index.html", viewData)
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		models.Create(title)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	models.Delete(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func render(w http.ResponseWriter, r *http.Request, file string, data interface{}) {
	// session, _ := config.Store.Get(r, "session")

	tmpl := template.Must(template.ParseFiles(
		"views/layout/layout.html",
		"views/"+file,
	))

	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

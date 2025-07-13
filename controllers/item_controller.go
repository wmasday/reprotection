package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

type FileDetail struct {
	Filename string
	Filepath string
	Content  string
}

type DashboardItem struct {
	ID int
	Title string
	MaliciousCount int
}

func Index(w http.ResponseWriter, r *http.Request) {
	items, err := models.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var dashboardItems []DashboardItem
	for _, item := range items {
		count, _ := models.CountMaliciousByItemID(item.ID)
		dashboardItems = append(dashboardItems, DashboardItem{
			ID: item.ID,
			Title: item.Title,
			MaliciousCount: count,
		})
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

	chartDataJSON := ""
	data, _ := json.Marshal(items)
	chartDataJSON = string(data)

	viewData := ViewData{
		User:           username,
		Items:          dashboardItems,
		Timestamp:      time.Now().Unix(),
		ChartData:      chartDataJSON,
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

func Detail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	item, err := models.GetByID(id)
	if err != nil {
		http.Error(w, "Item not found", 404)
		return
	}

	workingProject, _ := models.GetWorkingProject()
	if workingProject == nil {
		http.Error(w, "Working project not set", 500)
		return
	}

	var files []FileDetail
	err = filepath.Walk(workingProject.WorkingProject, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(item.Title)) {
			content, _ := os.ReadFile(path)
			encodedContent := base64.StdEncoding.EncodeToString(content)
			files = append(files, FileDetail{
				Filename: info.Name(),
				Filepath: path,
				Content:  encodedContent,
			})
			_ = models.UpsertMalicious(item.ID, path)
		}
		return nil
	})
	if err != nil {
		http.Error(w, "Error scanning files", 500)
		return
	}

	session, _ := config.Store.Get(r, "session")
	username := ""
	if v, ok := session.Values["username"].(string); ok {
		username = v
	}

	viewData := struct {
		User           string
		Item           *models.Item
		Files          []FileDetail
		WorkingProject *models.Project
	}{
		User:           username,
		Item:           item,
		Files:          files,
		WorkingProject: workingProject,
	}

	render(w, r, "detail.html", viewData)
}

func Sync(w http.ResponseWriter, r *http.Request) {
	// Get all items
	items, err := models.GetAll()
	if err != nil {
		http.Error(w, "Failed to get items", 500)
		return
	}

	for _, item := range items {
		maliciousFiles, err := models.GetMaliciousByItemID(item.ID)
		if err != nil {
			continue
		}
		for _, m := range maliciousFiles {
			content, err := os.ReadFile(m.Filepath)
			if err != nil {
				// File not found, remove from table
				_ = models.DeleteMaliciousByID(m.ID)
				continue
			}
			if !strings.Contains(strings.ToLower(string(content)), strings.ToLower(item.Title)) {
				_ = models.DeleteMaliciousByID(m.ID)
			}
		}
	}
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

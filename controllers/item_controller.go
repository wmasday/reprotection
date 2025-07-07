package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"reprotection/config"
	"reprotection/models"
)

type ViewData struct {
	User      string
	Items     interface{}
	Timestamp int64
	ChartData string
}

func Index(w http.ResponseWriter, r *http.Request) {
	items, err := models.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	render(w, r, "index.html", items)
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
	session, _ := config.Store.Get(r, "session")

	username := ""
	if v, ok := session.Values["username"].(string); ok {
		username = v
	}

	tmpl := template.Must(template.ParseFiles(
		"views/layout/layout.html",
		"views/"+file,
	))

	// Convert items to JSON for chart data
	var chartDataJSON string
	if items, ok := data.([]models.Item); ok {
		chartData, _ := json.Marshal(items)
		chartDataJSON = string(chartData)
	}

	viewData := ViewData{
		User:      username,
		Items:     data,
		Timestamp: time.Now().Unix(),
		ChartData: chartDataJSON,
	}



	err := tmpl.Execute(w, viewData)
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

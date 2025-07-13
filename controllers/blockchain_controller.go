package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reprotection/config"
	"reprotection/models"
	"strings"
)

type BlacklistEntry struct {
	Keyword   string `json:"keyword"`
	CreatedBy string `json:"createdBy"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
	IsActive  bool   `json:"isActive"`
	IsApplied bool   `json:"isApplied"`
}

type Statistics struct {
	TotalApply          string `json:"totalApply"`
	TotalVerify         string `json:"totalVerify"`
	TotalBlacklistEntries string `json:"totalBlacklistEntries"`
}

type BlacklistViewData struct {
	User        string
	Entries     []BlacklistEntry
	Statistics  Statistics
	Message     string
	MessageType string
}

func BlacklistIndex(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")
	username := ""
	if v, ok := session.Values["username"].(string); ok {
		username = v
	}

	// Get blacklist entries from blockchain API
	entries, err := getBlacklistEntries()
	if err != nil {
		entries = []BlacklistEntry{}
	}

	// Get statistics from blockchain API
	stats, err := getStatistics()
	if err != nil {
		stats = Statistics{
			TotalApply:          "0",
			TotalVerify:         "0",
			TotalBlacklistEntries: "0",
		}
	}

	// Handle URL parameters for messages
	message := r.URL.Query().Get("message")
	messageType := r.URL.Query().Get("type")
	if messageType == "" {
		messageType = "info"
	}

	// Handle error and success messages
	if errorMsg := r.URL.Query().Get("error"); errorMsg != "" {
		message = errorMsg
		messageType = "danger"
	}
	if successMsg := r.URL.Query().Get("success"); successMsg != "" {
		message = successMsg
		messageType = "success"
	}

	viewData := BlacklistViewData{
		User:        username,
		Entries:     entries,
		Statistics:  stats,
		Message:     message,
		MessageType: messageType,
	}

	render(w, r, "blockchain/index.html", viewData)
}

func BlacklistAdd(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := config.Store.Get(r, "session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if user is authenticated
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		keyword := r.FormValue("keyword")

		if keyword == "" {
			http.Redirect(w, r, "/blockchain?error=Keyword is required", http.StatusSeeOther)
			return
		}

		// Call blockchain API to add keyword with authenticated user
		err := addBlacklistKeyword(keyword, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+err.Error(), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/blockchain?success=Keyword added successfully", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blockchain", http.StatusSeeOther)
}

func BlacklistRemove(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := config.Store.Get(r, "session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if user is authenticated
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		keyword := r.FormValue("keyword")

		if keyword == "" {
			http.Redirect(w, r, "/blockchain?error=Keyword is required", http.StatusSeeOther)
			return
		}

		// Check if user is the creator of the keyword
		isCreator, err := checkKeywordCreator(keyword, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+err.Error(), http.StatusSeeOther)
			return
		}

		if !isCreator {
			http.Redirect(w, r, "/blockchain?error=Only the creator can remove this keyword", http.StatusSeeOther)
			return
		}

		// Call blockchain API to remove keyword
		err = removeBlacklistKeyword(keyword, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+err.Error(), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/blockchain?success=Keyword removed successfully", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blockchain", http.StatusSeeOther)
}

func BlacklistCheck(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := config.Store.Get(r, "session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if user is authenticated
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		keyword := r.FormValue("keyword")

		if keyword == "" {
			http.Redirect(w, r, "/blockchain?error=Keyword is required", http.StatusSeeOther)
			return
		}

		// Check if keyword is blacklisted
		isBlacklisted, err := checkKeywordBlacklisted(keyword)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+err.Error(), http.StatusSeeOther)
			return
		}

		message := fmt.Sprintf("Keyword '%s' is %s", keyword, map[bool]string{true: "BLACKLISTED", false: "NOT blacklisted"}[isBlacklisted])
		messageType := map[bool]string{true: "danger", false: "success"}[isBlacklisted]

		http.Redirect(w, r, "/blockchain?message="+message+"&type="+messageType, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blockchain", http.StatusSeeOther)
}

// Helper functions to interact with blockchain API
func getBlacklistEntries() ([]BlacklistEntry, error) {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	resp, err := http.Get(apiURL + "/blacklist")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entries []BlacklistEntry
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		return nil, err
	}

	// Check which keywords are applied to items system
	for i := range entries {
		entries[i].IsApplied = isKeywordApplied(entries[i].Keyword)
	}

	return entries, err
}

func isKeywordApplied(keyword string) bool {
	// Check if keyword exists in items system
	items, err := models.GetAll()
	if err != nil {
		return false
	}

	for _, item := range items {
		if strings.ToLower(item.Title) == strings.ToLower(keyword) {
			return true
		}
	}

	return false
}

func getStatistics() (Statistics, error) {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	resp, err := http.Get(apiURL + "/blacklist/statistics")
	if err != nil {
		return Statistics{}, err
	}
	defer resp.Body.Close()

	var stats Statistics
	err = json.NewDecoder(resp.Body).Decode(&stats)
	return stats, err
}

func addBlacklistKeyword(keyword, username string) error {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	data := map[string]string{
		"keyword":  keyword,
		"username": username,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiURL+"/blacklist/add", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf(errorResp.Error)
	}

	return nil
}

func removeBlacklistKeyword(keyword, username string) error {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	data := map[string]string{
		"keyword":  keyword,
		"username": username,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiURL+"/blacklist/remove", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf(errorResp.Error)
	}

	return nil
}

func checkKeywordCreator(keyword, username string) (bool, error) {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	resp, err := http.Get(apiURL + "/blacklist/" + keyword + "/creator")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("keyword not found")
	}

	var result struct {
		Creator string `json:"creator"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.Creator == username, nil
}

func checkKeywordBlacklisted(keyword string) (bool, error) {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	resp, err := http.Get(apiURL + "/blacklist/" + keyword + "/isBlacklisted")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		Keyword      string `json:"keyword"`
		IsBlacklisted bool   `json:"isBlacklisted"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.IsBlacklisted, err
}

func BlacklistApply(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := config.Store.Get(r, "session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if user is authenticated
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		keyword := r.FormValue("keyword")

		if keyword == "" {
			http.Redirect(w, r, "/blockchain?error=Keyword is required", http.StatusSeeOther)
			return
		}

		// Check if keyword is blacklisted
		isBlacklisted, err := checkKeywordBlacklisted(keyword)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+err.Error(), http.StatusSeeOther)
			return
		}

		if !isBlacklisted {
			http.Redirect(w, r, "/blockchain?error=Keyword is not blacklisted", http.StatusSeeOther)
			return
		}

		// Apply keyword to items system
		err = applyKeywordToItems(keyword, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+err.Error(), http.StatusSeeOther)
			return
		}

		// Increment apply counter in blockchain
		err = incrementApplyCounter()
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to increment apply counter: %v\n", err)
		}

		http.Redirect(w, r, "/blockchain?success=Keyword applied to items system successfully", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blockchain", http.StatusSeeOther)
}

func applyKeywordToItems(keyword, username string) error {
	// Create item in the items system
	err := models.Create(keyword)
	if err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}

	// Log the application
	fmt.Printf("User %s applied blockchain keyword '%s' to items system\n", username, keyword)
	return nil
}

func incrementApplyCounter() error {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	resp, err := http.Post(apiURL+"/blacklist/apply", "application/json", strings.NewReader("{}"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to increment apply counter")
	}

	return nil
} 
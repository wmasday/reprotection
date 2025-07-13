package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reprotection/config"
	"reprotection/models"
	"strconv"
	"strings"
	"time"
)

type BlacklistEntry struct {
	Index     int    `json:"index"`
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

	// Get all blockchain keywords
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
			http.Redirect(w, r, "/blockchain?error=Keyword+is+required", http.StatusSeeOther)
			return
		}

		// Call blockchain API to add keyword
		err := addBlacklistKeyword(keyword, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+strings.ReplaceAll(err.Error(), " ", "+"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/blockchain?success=Keyword+added+to+blockchain+successfully", http.StatusSeeOther)
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
		indexStr := r.FormValue("index")

		if indexStr == "" {
			http.Redirect(w, r, "/blockchain?error=Index+is+required", http.StatusSeeOther)
			return
		}

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error=Invalid+index", http.StatusSeeOther)
			return
		}

		// Check if user is the creator of the keyword
		isCreator, err := checkKeywordCreator(index, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+strings.ReplaceAll(err.Error(), " ", "+"), http.StatusSeeOther)
			return
		}

		if !isCreator {
			http.Redirect(w, r, "/blockchain?error=Only+the+creator+can+remove+this+keyword", http.StatusSeeOther)
			return
		}

		// Call blockchain API to remove keyword
		err = removeBlacklistKeyword(index, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+strings.ReplaceAll(err.Error(), " ", "+"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/blockchain?success=Keyword+deactivated+successfully", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blockchain", http.StatusSeeOther)
}

func BlacklistToggle(w http.ResponseWriter, r *http.Request) {
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
		indexStr := r.FormValue("index")

		if indexStr == "" {
			http.Redirect(w, r, "/blockchain?error=Index+is+required", http.StatusSeeOther)
			return
		}

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error=Invalid+index", http.StatusSeeOther)
			return
		}

		// Check if user is the creator of the keyword
		isCreator, err := checkKeywordCreator(index, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+strings.ReplaceAll(err.Error(), " ", "+"), http.StatusSeeOther)
			return
		}

		if !isCreator {
			http.Redirect(w, r, "/blockchain?error=Only+the+creator+can+toggle+this+keyword", http.StatusSeeOther)
			return
		}

		// Call blockchain API to toggle keyword
		err = toggleBlacklistKeyword(index, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+strings.ReplaceAll(err.Error(), " ", "+"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/blockchain?success=Keyword+status+toggled+successfully", http.StatusSeeOther)
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

	fmt.Printf("Fetching blockchain keywords from: %s\n", apiURL)

	resp, err := http.Get(apiURL + "/blacklist")
	if err != nil {
		fmt.Printf("Error fetching blockchain keywords: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("Blockchain API response status: %d\n", resp.StatusCode)

	var entries []BlacklistEntry
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		fmt.Printf("Error decoding blockchain response: %v\n", err)
		return nil, err
	}

	fmt.Printf("Retrieved %d keywords from blockchain\n", len(entries))

	// Format timestamps to readable date/time
	for i := range entries {
		if entries[i].Timestamp != "" {
			// Convert Unix timestamp to formatted date
			timestamp, err := strconv.ParseInt(entries[i].Timestamp, 10, 64)
			if err == nil {
				// Convert seconds to time.Time
				t := time.Unix(timestamp, 0)
				entries[i].Timestamp = t.Format("Jan 02, 2006 15:04:05")
			}
		}
	}

	// Note: We're not tracking applied status anymore since we allow duplicates
	// Each user can apply their own keywords independently

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

func checkKeywordCreator(index int, username string) (bool, error) {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	resp, err := http.Get(fmt.Sprintf("%s/blacklist/%d/creator", apiURL, index))
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

func removeBlacklistKeyword(index int, username string) error {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	data := map[string]interface{}{
		"index":    index,
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

func toggleBlacklistKeyword(index int, username string) error {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	data := map[string]interface{}{
		"index":    index,
		"username": username,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiURL+"/blacklist/toggle", "application/json", strings.NewReader(string(jsonData)))
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

func checkKeywordExists(index int) (bool, error) {
	apiURL := os.Getenv("BLOCKCHAIN_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3001"
	}

	// Try to get the keyword entry
	resp, err := http.Get(fmt.Sprintf("%s/blacklist/%d", apiURL, index))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// If status is 200, keyword exists
	return resp.StatusCode == 200, nil
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
			http.Redirect(w, r, "/blockchain?error=Keyword+is+required", http.StatusSeeOther)
			return
		}

		// Since we allow duplicates, users can apply their own keywords
		// No need to check if already applied since each user's keyword is independent

		// Apply keyword to items system
		err = applyKeywordToItems(keyword, username)
		if err != nil {
			http.Redirect(w, r, "/blockchain?error="+strings.ReplaceAll(err.Error(), " ", "+"), http.StatusSeeOther)
			return
		}

		// Increment apply counter in blockchain
		err = incrementApplyCounter()
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to increment apply counter: %v\n", err)
		}

		http.Redirect(w, r, "/blockchain?success=Keyword+applied+to+items+system+successfully", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/blockchain", http.StatusSeeOther)
}

func applyKeywordToItems(keyword, username string) error {
	// Create a new item in the items system with just the keyword
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
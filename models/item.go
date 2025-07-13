package models

import (
	"os"
	"reprotection/config"
	"strings"
)

type Item struct {
	ID    int
	Title string
}

type Malicious struct {
	ID int
	ItemID int
	Filepath string
	Filename string
	Content string
}

func GetAll() ([]Item, error) {
	rows, err := config.DB.Query("SELECT id, title FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		rows.Scan(&i.ID, &i.Title)
		items = append(items, i)
	}
	return items, nil
}

func Create(title string) error {
	_, err := config.DB.Exec("INSERT INTO items (title) VALUES (?)", title)
	return err
}

func Delete(id int) error {
	// First delete all malicious records associated with this item
	_, err := config.DB.Exec("DELETE FROM malicious WHERE item_id = ?", id)
	if err != nil {
		return err
	}
	
	// Then delete the item
	_, err = config.DB.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

func GetMaliciousByItemID(itemID int) ([]Malicious, error) {
	rows, err := config.DB.Query("SELECT id, item_id, filepath FROM malicious WHERE item_id = ?", itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Malicious
	for rows.Next() {
		var m Malicious
		err := rows.Scan(&m.ID, &m.ItemID, &m.Filepath)
		if err != nil {
			continue
		}
		// Extract filename from filepath
		parts := strings.Split(m.Filepath, "/")
		m.Filename = parts[len(parts)-1]
		// Read file content (if exists)
		content, err := os.ReadFile(m.Filepath)
		if err == nil {
			m.Content = string(content)
		} else {
			m.Content = "(Could not read file)"
		}
		result = append(result, m)
	}
	return result, nil
}

func GetByID(id int) (*Item, error) {
	var i Item
	err := config.DB.QueryRow("SELECT id, title FROM items WHERE id = ?", id).Scan(&i.ID, &i.Title)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func GetByTitle(title string) (*Item, error) {
	var i Item
	err := config.DB.QueryRow("SELECT id, title FROM items WHERE title = ?", title).Scan(&i.ID, &i.Title)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

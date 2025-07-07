package models

import (
	"reprotection/config"
)

type Item struct {
	ID    int
	Title string
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
	_, err := config.DB.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

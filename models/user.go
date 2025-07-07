package models

import (
	"golang.org/x/crypto/bcrypt"
	"reprotection/config"
)

type User struct {
	ID                 int
	Username           string
	Password           string
	MustChangePassword bool
}

// Membuat user baru
func CreateUser(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = config.DB.Exec("INSERT INTO users (username, password, must_change_password) VALUES (?, ?, TRUE)", username, string(hashed))
	return err
}

// Login dan validasi password
func FindUser(username, password string) (*User, error) {
	u := User{}
	err := config.DB.QueryRow("SELECT id, username, password, must_change_password FROM users WHERE username = ?", username).
		Scan(&u.ID, &u.Username, &u.Password, &u.MustChangePassword)
	if err != nil {
		return nil, err
	}

	// Bandingkan hash password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Get user by ID
func GetUserByID(id int) (*User, error) {
	u := User{}
	err := config.DB.QueryRow("SELECT id, username, password, must_change_password FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Username, &u.Password, &u.MustChangePassword)
	if err != nil {
		return nil, err
	}
	return &u, nil
}


// Update password dan reset status wajib ganti
func UpdateUserPassword(id int, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = config.DB.Exec("UPDATE users SET password = ?, must_change_password = FALSE WHERE id = ?", string(hashed), id)
	return err
}

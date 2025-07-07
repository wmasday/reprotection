package controllers

import (
	"html/template"
	"net/http"
	"log"

	"golang.org/x/crypto/bcrypt"

	"reprotection/config"
	"reprotection/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")

	var errorMsg string
	if flashes := session.Flashes(); len(flashes) > 0 {
		errorMsg = flashes[0].(string)
	}
	session.Save(r, w)

	tmpl := template.Must(template.ParseFiles("views/auth/login.html"))
	data := map[string]interface{}{
		"Error": errorMsg,
	}
	tmpl.Execute(w, data)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := models.FindUser(username, password)
		session, _ := config.Store.Get(r, "session")

		if err != nil {
			session.AddFlash("Invalid username or password.")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session.Values["user_id"] = user.ID
		session.Values["username"] = user.Username

		// Jika user wajib ganti password
		if user.MustChangePassword {
			session.Values["must_change_password"] = true
			
			session.Save(r, w)
			http.Redirect(w, r, "/change-password", http.StatusSeeOther)
			return
		}

		// Login berhasil
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")
	delete(session.Values, "user_id")
	delete(session.Values, "username")
	delete(session.Values, "must_change_password")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ShowChangePassword menampilkan halaman ubah kata sandi untuk user yang diwajibkan mengganti password.
func ShowChangePassword(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "session")
	if err != nil {
		log.Printf("Gagal mendapatkan sesi di ShowChangePassword: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Cek apakah user sudah login
	if session.Values["user_id"] == nil {
		log.Println("[ShowChangePassword] Pengguna belum login, mengarahkan ke /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Pastikan user memang diwajibkan untuk mengganti kata sandi
	mustChange, ok := session.Values["must_change_password"].(bool)
	if !ok || !mustChange {
		log.Println("[ShowChangePassword] Pengguna tidak diwajibkan mengganti password, mengarahkan ke /")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Ambil flash message (pesan error)
	var errorMsg string
	flashes := session.Flashes()
	log.Printf("[ShowChangePassword] Flashes ditemukan: %+v", flashes)

	if len(flashes) > 0 {
		if msg, ok := flashes[0].(string); ok {
			errorMsg = msg
			log.Println("[ShowChangePassword] Pesan flash:", errorMsg)
		} else {
			log.Printf("[ShowChangePassword] Format flash bukan string: %T", flashes[0])
		}
	}

	// Wajib save setelah mengambil flash agar tidak muncul kembali
	if err := session.Save(r, w); err != nil {
		log.Printf("Gagal menyimpan sesi setelah mengambil flash: %v", err)
	}

	// Render template ubah password
	tmpl := template.Must(template.ParseFiles("views/auth/change_password.html"))
	if err := tmpl.Execute(w, map[string]interface{}{
		"Error": errorMsg,
	}); err != nil {
		log.Printf("Gagal mengeksekusi template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}


// ChangePassword menangani permintaan POST untuk mengubah kata sandi.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	session, err := config.Store.Get(r, "session")
	if err != nil {
		log.Printf("Gagal mendapatkan sesi di ChangePassword: %v", err)
		session.AddFlash("An error occurred with your session. Please try again.")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("[ChangePassword] ID user tidak ditemukan atau bukan integer.")
		session.AddFlash("Session expired or invalid user ID. Please log in again.")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Ambil nilai dari form
	oldPass := r.FormValue("old_password")
	newPass := r.FormValue("new_password")
	confirmPass := r.FormValue("confirm_password")

	log.Printf("[ChangePassword] Validasi input. UserID: %d", userID)

	// Validasi input tidak boleh kosong
	if oldPass == "" || newPass == "" || confirmPass == "" {
		session.AddFlash("All password fields must be filled!")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	// Validasi konfirmasi password baru
	if newPass != confirmPass {
		session.AddFlash("New password and confirmation do not match!")
		log.Println("[ChangePassword] Password baru dan konfirmasi tidak sama.")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	// Cek kekuatan password baru (misalnya minimal 8 karakter)
	if len(newPass) < 8 {
		session.AddFlash("New password must be at least 8 characters!")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	// Ambil user dari database
	user, err := models.GetUserByID(userID)
	if err != nil {
		log.Printf("[ChangePassword] Gagal mendapatkan user dari DB: %v", err)
		session.AddFlash("User not found or database error occurred.")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validasi password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPass)); err != nil {
		session.AddFlash("Incorrect current password!")
		log.Println("[ChangePassword] Password lama salah.")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	// Update password ke database
	if err := models.UpdateUserPassword(user.ID, newPass); err != nil {
		log.Printf("[ChangePassword] Gagal memperbarui password: %v", err)
		session.AddFlash("Failed to change password! Please try again.")
		if saveErr := session.Save(r, w); saveErr != nil {
			log.Printf("Gagal menyimpan sesi setelah flash: %v", saveErr)
		}
		http.Redirect(w, r, "/change-password", http.StatusSeeOther)
		return
	}

	// Sukses, hapus flag wajib ubah password
	delete(session.Values, "must_change_password")
	if err := session.Save(r, w); err != nil {
		log.Printf("Gagal menyimpan sesi setelah menghapus flag: %v", err)
	}

	log.Println("[ChangePassword] Password berhasil diperbarui, redirect ke /")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

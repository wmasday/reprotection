package main

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    // Daftar kata terlarang
    blockedKeywords := []string{
        "porno", "porn", "bokep", "sex", "seks", "18+", "nude", "nudity", "hot",
        "jav", "anal", "bdsm", "payudara", "pantat", "bugil", "mesum", "masturbasi",
        "bunuh", "pembunuhan", "tembak", "darah", "mati", "pembantaian", "tusuk",
        "bakar", "bom", "teroris", "ledakan", "pukulan", "gantung diri",
        "kafir", "babi", "anjing", "goblok", "bodoh", "jancok", "tolol", "kontol",
        "asu", "bangsat", "idiot", "penghinaan", "rasis", "fasis", "nazi",
        "gratis", "free", "diskon", "promo", "klik di sini", "penipuan", "bitcoin gratis",
        "investasi cepat", "cepat kaya", "pinjaman online", "uang cepat", "hacks", "cheat",
        "memek", "ngentot", "coli", "kampret", "taik", "tai", "peler", "titit",
        "setan", "iblis", "anjrit", "keparat",
    }

    // Ganti path sesuai kebutuhan
    targetDir := "/var/www/html/"

    // Walk semua file
    filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            fmt.Println("Error accessing:", path, "-", err)
            return nil
        }

        if d.IsDir() {
            return nil
        }

        // Baca file isi
        content, err := os.ReadFile(path)
        if err != nil {
            fmt.Println("Gagal baca:", path, "-", err)
            return nil
        }

        // Ubah isi ke lowercase string
        lowerContent := strings.ToLower(string(content))

        // Cek apakah mengandung kata blacklist
        for _, word := range blockedKeywords {
            if strings.Contains(lowerContent, word) {
                fmt.Println("Menghapus:", path)
                os.Remove(path)
                break
            }
        }

        return nil
    })
}

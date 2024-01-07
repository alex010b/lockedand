package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var eventChannel (chan VPS)

var emailTime = make(map[string]time.Time)

var MP string

// user json struct
type Person struct {
	Name       string `gorm:"primary_key;auto_increment"`
	Ip         string
	Email      string
	Password   []byte
	Keycheck   bool
	Emailcheck bool
}

type Logs struct {
	Ip    string
	Count int
}

type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Masterkey string `json:"masterkey"`
	Email     string `json:"email"`
	Vernumber string `json:"vernumber"`
}

type VPS struct {
	Hostname string `json:"hostname"`
	Label    string `json:"label"`
	Password string `json:"password"`
	Os       string `json:"os"`
	Username string `json:"username"`
}

func main() {

	var password string

	fmt.Print("password : ")

	fmt.Scanln(&password)

	ciphertext, err := os.ReadFile("encrypted.txt")
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}

	decrypted, err := decrypt([]byte(password), ciphertext)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	fmt.Println("Decryption successful. Decrypted data:", decrypted)

	MP = decrypted

	dbConnect()

	eventChannel = make(chan VPS)

	mux := http.NewServeMux()

	// Specify the directory containing your images
	imagesDir := "images/"

	// Create a file server handler for the images directory
	fileServer := http.FileServer(http.Dir(imagesDir))

	// Handle requests for the "/images/" path and strip the "/images/" prefix
	mux.Handle("/images/", http.StripPrefix("/images/", fileServer))
	mux.HandleFunc("/", mainPageHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/createkey", keyCreatingHandler)
	mux.HandleFunc("/verifyemail", verifyEmailHandler)
	mux.HandleFunc("/css/", cssHandler)
	mux.HandleFunc("/static/", StaticjsFileHandler)
	mux.HandleFunc("/user/", PersoHandler)
	mux.HandleFunc("/ws", wsHandler)

	fmt.Printf("Server listening on port 8000\n")
	err = http.ListenAndServe(":8000", generalHandler(mux))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func generalHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if rateLimit(r) {
			w.Write([]byte("Rate Limited!"))
			return
		}
		if len(r.UserAgent()) > 1000 {
			w.Write([]byte("invalid user agent"))
		}
		//requestHandler(w, r)
		next.ServeHTTP(w, r)
	})
}

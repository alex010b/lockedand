package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

// Getting the whole request and printing it to the console server
func requestHandler(w http.ResponseWriter, r *http.Request) {
	rawRequest, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error dumping the request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Print(string(rawRequest), time.Now(), "\n\n")
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var creds Credentials

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		if verifyVernumber(creds.Username, creds.Vernumber) {

			// setting the emailver flag to true
			err := db.Model(&Person{}).Where("name = ?", creds.Username).Update("emailcheck", true).Error
			if err != nil {
				fmt.Println("couldnt find user", err)
			}

			//authenticating the ip in the allip folder
			newIp := getClientIP(r) + " " + r.Header.Get("User-Agent")
			err = db.Model(&Person{}).Where("name = ?", creds.Username).Update("ip", newIp).Error
			if err != nil {
				fmt.Println("error with appending to ip record", err)
			}
			response := "Success! Redirecting"
			w.Write([]byte(response))
		} else {
			response := "Wrong numer"
			w.Write([]byte(response))
		}
	} else {
		http.ServeFile(w, r, "templates/verifyemail.html")
	}
}

func keyCreatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		//variable for putting the json into the golang struct
		var keycreating Credentials
		if err := json.NewDecoder(r.Body).Decode(&keycreating); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		//token like validation
		if isValid := CheckCredentials(keycreating.Username, keycreating.Password); isValid == "success" {
			//token like validation
			if isKeyValid := checkMasterKey(keycreating.Masterkey); isKeyValid == "success" {
				if createKey(keycreating.Username, keycreating.Masterkey) {
					response := "Successful! redirecting..."
					w.Write([]byte(response))
				} else {
					response := "couldn't create your key );"
					w.Write([]byte(response))
				}
			} else {
				response := isKeyValid
				w.Write([]byte(response))
			}
		} else {
			response := isValid
			w.Write([]byte(response))
		}
	} else {
		http.ServeFile(w, r, "templates/createkey.html")
	}
}

// Hosting the index.html
func mainPageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var login Credentials

		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}
		var responseData interface{}

		isValid := CheckCredentials(login.Username, login.Password)

		if isValid == "success" {

			var person Person

			if err := db.First(&person, "name = ?", login.Username).Error; err != nil {
				fmt.Println("couldnt find users info")
			}

			userAgent := r.Header.Get("User-Agent")

			wholeip := getClientIP(r) + " " + userAgent

			if checkIp(wholeip, person.Ip) && person.Emailcheck {

				token, err := userTokenGen(person.Name)
				if err != nil {
					fmt.Println("token gen err", err)
				}

				fmt.Println(token)

				responseData = map[string]interface{}{
					"status":  "login_success",
					"message": token,
				}

				jsonData, err := json.Marshal(responseData)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				w.Write(jsonData)

			} else {

				verMail(person.Name)

				responseData = map[string]interface{}{
					"status":  "error_mail",
					"message": person.Email,
				}

				jsonData, err := json.Marshal(responseData)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				w.Write(jsonData)
			}

		} else {
			responseData = map[string]interface{}{
				"status":  "login_error",
				"message": isValid,
			}
			jsonData, err := json.Marshal(responseData)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}

	} else {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		var register Credentials

		if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		if isValid := CheckRegister(register.Username, register.Password, register.Email); isValid == "success" {

			// create user data and folder
			createPrivatDataFolder(register.Username)
			registerUser(register.Password, register.Username, register.Email, r.Header.Get("User-Agent"), r)

			verMail(register.Username)

			tmpl, err := template.ParseFiles("templates/verifyemail.html")
			if err != nil {
				panic(err)
			}

			err = tmpl.Execute(w, register)
			if err != nil {
				panic(err)
			}

		} else {
			response := map[string]string{"Register": isValid}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		http.ServeFile(w, r, "templates/register.html")
	}

}

func serveCSSFile(w http.ResponseWriter, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/css")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to serve the file", http.StatusInternalServerError)
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	serveCSSFile(w, "css/webflow.css")
	serveCSSFile(w, "css/normalize.css")
	serveCSSFile(w, "css/lockedand-01e041.webflow.css")
}

func StaticjsFileHandler(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("static/webflow.js")
	if err != nil {

		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/javascript")

	_, err = io.Copy(w, file)
	if err != nil {

		http.Error(w, "Failed to serve the file", http.StatusInternalServerError)
	}
}

func PersoHandler(w http.ResponseWriter, r *http.Request) {

	usernameCookie, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		fmt.Println("No cookie found")
		http.Redirect(w, r, "https://www.lockedand.com", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Redirect(w, r, "https://www.lockedand.com", http.StatusSeeOther)
		return
	}
	username := usernameCookie.Value

	if _, err := validateUserToken(username, r); err != nil {
		http.Redirect(w, r, "https://www.lockedand.com", http.StatusSeeOther)
		return
	} else {

		if r.Method == http.MethodPost {

			var vps VPS

			if err := json.NewDecoder(r.Body).Decode(&vps); err != nil {
				http.Error(w, "Invalid JSON data", http.StatusBadRequest)
				return
			}

			if vps.Username == "" {
				vps.Username = username
			}

			fmt.Println("saddsaasd")

			eventChannel <- vps

			//gonna do sum config check

			w.Write([]byte("good"))

		} else {
			page := r.URL.Path[len("/user/"):]

			filename := "templates/" + page + ".html"

			// Serve the HTML file
			http.ServeFile(w, r, filename)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "https://www.lockedand.com"
	},
}
var test int

func wsHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Connection", "Upgrade")
	r.Header.Set("Upgrade", "websocket")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
}

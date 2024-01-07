package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

//
//
//

// checking for register

// main
func CheckRegister(username string, password string, email string) string {

	specialcharsToCheck := "~!@#$%^&*()_+`=-{}[];':/.,?\\"
	capCharToCheck := "QWERTYUIOPASDFGHJKLZXCVBNM"
	numbercharsToCheck := "1234567890"
	requiredSpecialChar := 5
	requiredNumberChar := 5
	requiredNumberCharfu := 1
	capRequired := 1

	if len(username) == 0 {
		return "you must provide a username cuh"
	}

	if len(password) == 0 {
		return "you must provide a password cuh"
	}

	if strings.Contains(username, " ") {
		return "username can't have spaces"
	}

	if strings.Contains(password, " ") {
		return "password can't have spaces"
	}

	for _, char := range specialcharsToCheck {
		if strings.ContainsRune(username, char) {
			errstr := "username cant contain special character " + string(char)
			return errstr
		}
	}

	var exists bool
	db.Model(&Person{}).Where("name = ?", username).Find(&Person{})

	if exists {
		return "username alredy exists try an other one"
	}

	numberCountfu := 0
	for _, char := range numbercharsToCheck {
		numberCountfu += strings.Count(username, string(char))
	}

	if numberCountfu < requiredNumberCharfu {
		return "The username must contain at least 1 numbers"
	}

	if len(username) < 5 {
		return "username muste have at least 5 characters"
	}

	if len(username) > 100 {
		return "ok too long fella"
	}

	if len(password) < 15 {
		return "password must have at least 15 characters"
	}

	if len(password) > 100 {
		return "ok too long fella"
	}

	specialCount := 0
	for _, char := range specialcharsToCheck {
		specialCount += strings.Count(password, string(char))
	}
	if specialCount < requiredSpecialChar {
		return "The password must contain at least 5 special characters ex: !@&"
	}

	numberCount := 0
	for _, char := range numbercharsToCheck {
		numberCount += strings.Count(password, string(char))
	}

	if numberCount < requiredNumberChar {
		return "The password must contain at least 5 numbers"
	}

	capCountp := 0
	for _, char := range capCharToCheck {
		capCountp += strings.Count(password, string(char))
	}

	if capCountp < capRequired {
		return "The password must contain at least 1 capital characters"
	}

	if len(password) == 0 {
		return "You need to provide an email"
	}

	if strings.Contains(email, " ") {
		return "Invalid email"
	}

	return "success"
}

//
//
//

// hashing password and storing it

// main

//
//
//

// creating user folders
func createPrivatDataFolder(username string) {
	folderName := "database/" + username

	err := os.Mkdir(folderName, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	fmt.Println("Directory created:", folderName)
}

//
//
//

// creating the user's folder

func registerUser(password string, username string, email string, userAgent string, r *http.Request) {

	keycheck := false

	emailcheck := false

	hashedPassword, _ := HashPassword(password)

	ip := getClientIP(r) + " " + userAgent

	person := Person{
		Name:       username,
		Ip:         ip,
		Email:      email,
		Password:   hashedPassword,
		Keycheck:   keycheck,
		Emailcheck: emailcheck,
	}

	db.Create(&person)

	var Person []Person
	db.Find(&Person)
}

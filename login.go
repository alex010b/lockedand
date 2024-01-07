package main

import (
	"fmt"
)

func CheckCredentials(username string, password string) string {
	var person Person
	if err := db.Select("password").Where("name = ?", username).First(&person).Error; err != nil {
		fmt.Println("couldnt find a password")
	}

	hashedPasswordFromFile := person.Password

	if err := comparePasswordWithHash(password, hashedPasswordFromFile); err != nil {
		return "Login password invalid"
	} else {
		return "success"
	}

}

func loadPersoLogin(username string) Person {
	var user Person
	if err := db.First(&user, "name = ?", username).Error; err != nil {
		fmt.Println("couldnt find users info")
	}

	return user

}

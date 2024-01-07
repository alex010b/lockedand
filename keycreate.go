package main

import (
	"fmt"
	"os"
)

func checkMasterKey(masterkey string) string {
	if len(masterkey) < 15 {
		return "Master key needs to be 15 characters at least"
	}
	return "success"
}

func createKey(username string, key string) bool {
	filePath := "database/" + username + "/chekkey.txt"

	ciphertext := encrypt(key, "test")

	os.WriteFile(filePath, ciphertext, 0644)

	fmt.Println("key created for : ", username)

	db.Model(&Person{}).Where("name = ?", username).Update("keycheck", "true")

	return true
}

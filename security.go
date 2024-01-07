package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ipRequestMap = make(map[string]int)
	count        int
	mutex        sync.Mutex
	neow         time.Time
)

func rateLimit(r *http.Request) bool {
	fullIp := getClientIP(r) + r.Header.Get("User-Agent")
	now := time.Now()
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := ipRequestMap[fullIp]
	if !exists || now.Sub(neow) >= 300*time.Millisecond {
		ipRequestMap[fullIp] = 0
		neow = time.Now()
	} else {
		ipRequestMap[fullIp]++
		neow = time.Now()
	}

	if ipRequestMap[fullIp] > 20 {
		fmt.Println("Rate Limited!")
		return true
	}
	return false
}

func checkIp(liveAddress string, storedAddress string) bool {

	fmt.Println(storedAddress)
	fmt.Println(liveAddress)

	if liveAddress == storedAddress || strings.Contains(string(storedAddress), liveAddress) {
		fmt.Println("was true")
		return true
	} else {
		fmt.Println("was false")
		return false
	}

}

func verifyVernumber(username string, vernumberProvided string) bool {

	folderPathNum := "database/" + username + "/emailCode.txt"

	vernumberFromFile, err := os.ReadFile(folderPathNum)
	if err != nil {
		fmt.Println("could no read json")
	}

	now := time.Now()
	fmt.Println(now.Sub(emailTime[username]))
	fmt.Println(emailTime)

	if err := comparePasswordWithHash(vernumberProvided, vernumberFromFile); err != nil && now.Sub(emailTime[username]) < 10*time.Minute {
		return false
	} else {
		os.Remove(folderPathNum)
		return true
	}

}

func verMail(username string) {
	var user Person

	user = loadPersoLogin(username)

	randomInt := rand.Intn(100000)

	from := "schyzis@lockedamail.com"
	password := MP
	to := user.Email
	subject := "Lockedand verification email"
	message := "Verification code: " + strconv.Itoa(randomInt)
	smtpServer := "mail.lockedamail.com"
	smtpPort := "587"
	smtpUsername := from

	email := []byte("From: " + from + "\r\n" + subject + "\r\n" + message)

	auth := smtp.PlainAuth("", smtpUsername, password, smtpServer)

	now := time.Now()
	emailTime[username] = now

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, email)
	if err != nil {
		os.WriteFile("database/errorlogs.txt", []byte(user.Email), 0644)
		fmt.Println("error sending email")
	}

	vernumber := strconv.Itoa(randomInt)

	hashevernumber, err := HashPassword(vernumber)
	if err != nil {
		fmt.Println("couldnt hash vernumber")
	}

	fmt.Println(vernumber)

	filepath := "database/" + username + "/emailCode.txt"

	err = os.WriteFile(filepath, hashevernumber, 0644)
	if err != nil {
		fmt.Println("could not write the verCode to the database")
	}

}

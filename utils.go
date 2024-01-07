package main

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

func detectWindows(w http.ResponseWriter, r *http.Request) bool {
	userAgent := r.UserAgent()
	return strings.Contains(userAgent, "Win64")
}

func redirectTowww(w http.ResponseWriter, r *http.Request) {

	host := r.Host

	if host == "lockedand.com" {
		http.Redirect(w, r, "https://www.lockedand.com", http.StatusSeeOther)

	}
}

func SendWs(conn *websocket.Conn, message string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println(err)
	}
}

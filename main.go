package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func init() {
	log.SetPrefix(fmt.Sprintf("[OPENVPN-AUTH] %s ", time.Now().Format("2025-07-17 00:04:05.000")))
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func main() {
	username := os.Getenv("username")
	password := os.Getenv("password")
	api := os.Getenv("ovpn_auth_api")

	if username == "" || password == "" || api == "" {
		log.Println("Missing required environment variables.")
		os.Exit(1)
	}

	resp, err := http.PostForm(api, url.Values{"username": {username}, "password": {password}})
	if err != nil {
		log.Println("Failed to send request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response:", err)
		os.Exit(1)
	}

	var data struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Failed to parse response:", err)
		os.Exit(1)
	}

	log.Printf("[%s] %s\n", username, data.Message)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Authentication failed for user %q", username)
		os.Exit(1)
	}

	os.Exit(0)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const SUPERSET_BASE = "http://localhost:8088"
const DASHBOARD_ID = "c0e94d84-82e6-4e8b-ba23-3e54987094cd"
const CLAUSE = "year_id > 1"
const DATASET = 2

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CSRFResponse struct {
	Result string `json:"result"`
}

type GuestTokenPayload struct {
	User      User               `json:"user"`
	Resources []Resource         `json:"resources"`
	RLS       []RowLevelSecurity `json:"rls"`
}

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type RowLevelSecurity struct {
	Clause  string `json:"clause"`
	Dataset int    `json:"dataset"`
}

type GuestTokenResponse struct {
	Token string `json:"token"`
}

func login() (string, string, error) {
	loginURL := fmt.Sprintf("%s/api/v1/security/login", SUPERSET_BASE)
	payload := map[string]interface{}{
		"password": os.Getenv("SUPERSET_PASS"),
		"provider": "db",
		"refresh":  true,
		"username": os.Getenv("SUPERSET_ADMIN"),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("login failed with status code %d", resp.StatusCode)
	}

	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return "", "", err
	}

	return loginResponse.AccessToken, loginResponse.RefreshToken, nil
}

func getCSRFToken(accessToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/security/csrf_token/", SUPERSET_BASE), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var csrfResponse CSRFResponse
	if err := json.NewDecoder(resp.Body).Decode(&csrfResponse); err != nil {
		return "", err
	}

	return csrfResponse.Result, nil
}

func createGuestToken() (string, error) {
	accessToken, _, err := login()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	csrfToken, err := getCSRFToken(accessToken)
	if err != nil {
		return "", err
	}

	guestTokenURL := fmt.Sprintf("%s/api/v1/security/guest_token/", SUPERSET_BASE)
	payload := GuestTokenPayload{
		User: User{
			Username:  "guest",
			FirstName: "Guest",
			LastName:  "User",
		},
		Resources: []Resource{
			{Type: "dashboard", ID: DASHBOARD_ID},
		},
		RLS: []RowLevelSecurity{
			{Clause: CLAUSE, Dataset: DATASET},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", guestTokenURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", guestTokenURL)
	req.Header.Set("X-CSRFToken", csrfToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var guestTokenResponse GuestTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&guestTokenResponse); err != nil {
		return "", err
	}
	return guestTokenResponse.Token, nil
}

func getGuestToken(w http.ResponseWriter, r *http.Request) {
	guestToken, err := createGuestToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]string{"guestToken": guestToken}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/guest-token", getGuestToken).Methods("GET")

	corsObj := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsObj)(r)))
}

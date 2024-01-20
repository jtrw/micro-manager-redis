package handler

import (
	"encoding/json"
	"io"
	//"github.com/go-chi/render"
	//"log"
	"crypto/sha256"
	"net/http"
)

type Auth struct {
	Username string
	Password string
}

func NewAuth(username, password string) Auth {
	return Auth{Username: username, Password: password}
}

func (a Auth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestData JSON

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = json.Unmarshal(b, &requestData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestData["username"] == nil || requestData["password"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.Username == requestData["username"] && a.Password == requestData["password"] {
		jsonResponse := JSON{"username": a.Username, "fullName": a.Username, "token": a.GetToken()}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
}

func (a Auth) GetToken() string {
	h := sha256.New()
	h.Write([]byte(a.Username + ":" + a.Password))
	return string(h.Sum(nil))
}
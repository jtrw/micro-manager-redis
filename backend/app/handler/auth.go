package handler

import (
	"encoding/json"
	"io"
	//"github.com/go-chi/render"
	//"log"
	"net/http"
)

type Auth struct {
	Username string
	Password string
}

func NewAuth(username, password string) Auth {
	return Auth{Username: username, Password: password}
}

func (a Auth) Auth(w http.ResponseWriter, r *http.Request) {
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
		jsonResponse := JSON{"id": 1, "username": a.Username, "fullName": a.Username, "token": "124"}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}
}

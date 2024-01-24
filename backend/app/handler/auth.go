package handler

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func (a Auth) GetToken() string {
	h := sha1.New()
	h.Write([]byte(a.Username + ":" + a.Password))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

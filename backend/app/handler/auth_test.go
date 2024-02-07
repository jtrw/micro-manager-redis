package handler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	auth := NewAuth("testuser", "testpassword")
	handler := http.HandlerFunc(auth.Login)

	requestBody := map[string]interface{}{
		"username": "testuser",
		"password": "testpassword",
	}

	reqBody, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	var jsonResponse JSON
	err = json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	assert.Nil(t, err)
	assert.Equal(t, "testuser", jsonResponse["username"])
	assert.Equal(t, "testuser", jsonResponse["fullName"])
	assert.NotNil(t, jsonResponse["token"])
}

func TestLogin_Unauthorized(t *testing.T) {
	auth := NewAuth("testuser", "testpassword")
	handler := http.HandlerFunc(auth.Login)

	requestBody := map[string]interface{}{
		"username": "wronguser",
		"password": "wrongpassword",
	}

	reqBody, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLogin_BadRequest(t *testing.T) {
	auth := NewAuth("testuser", "testpassword")
	handler := http.HandlerFunc(auth.Login)

	// Missing 'password' field in the request body
	requestBody := map[string]interface{}{
		"username": "testuser",
	}

	reqBody, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

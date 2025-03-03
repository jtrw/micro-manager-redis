package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware := Cors(handler)
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Control-Request, Content-Range, Request, Range, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Database", rr.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "GET, POST, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Range", rr.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "OK", rr.Body.String())
}

func Test_CorsMethodOptionsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	req, err := http.NewRequest("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middleware := Cors(handler)
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Authorized"))
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := "your_secret_token"
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	middleware := Auth(token)(handler)
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, "Authorized", rr.Body.String())
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Authorized"))
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	invalidToken := "invalid_token"
	req.Header.Set("Authorization", "Bearer "+invalidToken)

	rr := httptest.NewRecorder()
	middleware := Auth("your_secret_token")(handler)
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid token\n", rr.Body.String())
}

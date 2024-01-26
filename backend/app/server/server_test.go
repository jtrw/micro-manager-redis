package server

import (
	"context"
	//"fmt"
	//	"log"

	//"fmt"
	//	"embed"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	//"strings"
	"testing"
	"time"
)

func TestRest_Run(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	srv := Server{Listen: "localhost:54009", Version: "v1"}
	err := srv.Run(ctx)
	require.Error(t, err)
	assert.Equal(t, "http: Server closed", err.Error())
}

func TestRest_RobotsCheck(t *testing.T) {
	srv := Server{Listen: "localhost:54009", Version: "v1", Secret: "12345"}

	ts := httptest.NewServer(srv.routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/robots.txt")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "User-agent: *\nDisallow: /\n", string(body))
}

func TestRest_FileServerNotFound(t *testing.T) {
	srv := Server{Listen: "localhost:54009", Version: "v1", Secret: "12345", WebRoot: "./web"}

	ts := httptest.NewServer(srv.routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/web")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRest_FileServer(t *testing.T) {
	path, _ := os.Getwd()
	path = path + "/../../web"

	srv := Server{Listen: "localhost:54009", Version: "v1", Secret: "12345", WebRoot: path}

	ts := httptest.NewServer(srv.routes())
	defer ts.Close()

	testHTMLName := "index.html"

	resp, err := http.Get(ts.URL + "/web/" + testHTMLName)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Stubs for frontends\n", string(body))
}

func TestRest_FileServerNoFound(t *testing.T) {
	path, _ := os.Getwd()
	path = path + "/../../web"

	srv := Server{Listen: "localhost:54009", Version: "v1", Secret: "12345", WebRoot: path}

	ts := httptest.NewServer(srv.routes())
	defer ts.Close()

	testHTMLName := "index.html"

	resp, err := http.Get(ts.URL + "/web/test/" + testHTMLName)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

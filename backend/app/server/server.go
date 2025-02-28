package server

import (
	"context"
	"crypto/md5"
	"embed"
	"fmt"
	"io/fs"
	"log"
	manageHandler "micro-manager-redis/app/handler"
	repository "micro-manager-redis/app/repository"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jtrw/go-rest"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	//"fmt"
)

type Server struct {
	Listen         string
	PinSize        int
	MaxPinAttempts int
	MaxExpire      time.Duration
	WebRoot        string
	WebFS          embed.FS
	Secret         string
	Version        string
	Client         *redis.Client
	AuthLogin      string
	AuthPassword   string
	Context        context.Context
}

func (s Server) Run(ctx context.Context) error {
	log.Printf("[INFO] activate rest server")
	log.Printf("[INFO] Listen: %s", s.Listen)

	httpServer := &http.Server{
		Addr:              s.Listen,
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		if httpServer != nil {
			if clsErr := httpServer.Close(); clsErr != nil {
				log.Printf("[ERROR] failed to close proxy http server, %v", clsErr)
			}
		}
	}()

	err := httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)

	if err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}
	return err
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP)
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	router.Use(rest.AppInfo("Manag-RKeys", "Jrtw", s.Version), rest.Ping)
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))
	router.Use(middleware.Logger)

	redRep := repository.NewRedisRepository(s.Client)
	handler := manageHandler.NewHandler(redRep)
	authHandle := manageHandler.NewAuth(s.AuthLogin, s.AuthPassword)

	router.Route(
		"/api/v1", func(r chi.Router) {
			r.Use(Cors)
			r.Use(Auth(authHandle.GetToken()))
			r.Use(SetDatabase(&handler))
			r.Get("/keys", handler.AllKeys)
			r.Get("/keys/{key}", handler.GetKey)
			//r.Post("/keys", handler.CreateKey)
			//r.Put("/keys/{key}", handler.UpdateKey)
			r.Delete("/keys/{key}", handler.DeleteKey)
			r.Delete("/keys", handler.DeleteAllKeys)
			r.Get("/keys-group", handler.GroupKeys)
			r.Delete("/keys-group/{group}", handler.DeleteByGroup)
			r.Get("/keyspaces", handler.GetKeyspaces)
			r.Get("/databases", handler.GetDatabases)
			r.Post("/set-database", handler.SetDatabase)
		},
	)
	router.Route(
		"/auth", func(r chi.Router) {
			r.Use(Cors)
			r.Post("/", authHandle.Login)
		},
	)

	router.Get(
		"/robots.txt", func(w http.ResponseWriter, r *http.Request) {
			render.PlainText(w, r, "User-agent: *\nDisallow: /\n")
		},
	)

	addFileServer(router, s.WebFS, s.WebRoot, s.Version)

	return router
}

func addFileServer(r chi.Router, embedFS embed.FS, webRoot, version string) {
	var webFS http.Handler
	log.Printf("[INFO] webRoot: %s", webRoot)
	if _, err := os.Stat(webRoot); err == nil {
		log.Printf("[INFO] run file server from %s from the disk", webRoot)
		webFS = http.FileServer(http.Dir(webRoot))
	} else {
		log.Printf("[INFO] run file server, embedded")
		var contentFS, _ = fs.Sub(embedFS, "web")
		webFS = http.FileServer(http.FS(contentFS))
	}

	webFS = http.StripPrefix("/web", webFS)
	r.Get("/web", http.RedirectHandler("/web/", http.StatusMovedPermanently).ServeHTTP)

	r.With(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(20, nil)),
		middleware.Timeout(10*time.Second),
		cacheControl(time.Hour, version),
	).Get("/web/*", func(w http.ResponseWriter, r *http.Request) {
		// don't show dirs, just serve files
		if strings.HasSuffix(r.URL.Path, "/") && len(r.URL.Path) > 1 && r.URL.Path != ("/web/") {
			http.NotFound(w, r)
			return
		}
		webFS.ServeHTTP(w, r)
	})
}

func cacheControl(expiration time.Duration, version string) func(http.Handler) http.Handler {
	etag := func(r *http.Request, version string) string {
		s := version + ":" + r.URL.String()
		return fmt.Sprintf("%x", md5.Sum([]byte(s)))
	}

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			e := `"` + etag(r, version) + `"`
			w.Header().Set("Etag", e)
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, no-cache", int(expiration.Seconds())))

			if match := r.Header.Get("If-None-Match"); match != "" {
				if strings.Contains(match, e) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

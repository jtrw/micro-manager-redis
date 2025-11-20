package server

import (
	"context"
	"log"
	manageHandler "micro-manager-redis/app/handler"
	"net/http"
	"strconv"
	"strings"
)

type JSON map[string]interface{}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // change this later
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Control-Request, Content-Range, Request, Range, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Database")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Range")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Auth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")
			headerToken := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))

			if headerToken != token {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Database(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := r.Header.Get("X-Database")

		if database != "" {
			database = strings.TrimPrefix(database, "db")
			dbIndx, _ := strconv.Atoi(database)
			log.Printf("Database: %d", dbIndx)

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "database", dbIndx)))
			return
		}

		next.ServeHTTP(w, r)

	})
}

func SetDatabase(h *manageHandler.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbNumber := r.Header.Get("X-Database")
			if dbNumber != "" {
				dbNumber = strings.TrimPrefix(dbNumber, "db")
				index, err := strconv.Atoi(dbNumber)
				if err != nil {
					log.Printf("Invalid database number: %s", dbNumber)
					http.Error(w, "Invalid database number", http.StatusBadRequest)
					return
				}
				if index < 0 || index > 15 {
					log.Printf("Database number out of range: %d", index)
					http.Error(w, "Database number out of range (0-15)", http.StatusBadRequest)
					return
				}
				if err := h.SetRedisDatabase(index); err != nil {
					log.Printf("Error setting database: %v", err)
					http.Error(w, "Failed to set database", http.StatusInternalServerError)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

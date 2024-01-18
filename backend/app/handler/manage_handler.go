package handler

import (
	"context"
	"encoding/json"
	//"log"
	"net/http"

	//"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type JSON map[string]interface{}

type Handler struct {
	Database *redis.Client
}

const (
	StatusNew  = "new"
	StatusDone = "done"
)

func NewHandler(database *redis.Client) Handler {
	return Handler{Database: database}
}

type Keys struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Expire int    `json:"expire"`
}

func (h Handler) ShowTaskInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//uuidStr := chi.URLParam(r, "uuid")

	ctx := context.Background()

	iter := h.Database.Scan(ctx, 0, "*", 0).Iterator()
	allKeys := []Keys{}
	for iter.Next(ctx) {
		keys := Keys{
			Key:    iter.Val(),
			Value:  h.Database.Get(ctx, iter.Val()).Val(),
			Expire: int(h.Database.TTL(ctx, iter.Val()).Val().Seconds()),
		}
		allKeys = append(allKeys, keys)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(allKeys)
}

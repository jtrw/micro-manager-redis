package handler

import (
	"context"
	"encoding/json"
	//"log"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
)

type JSON map[string]interface{}

const (
	SEPARATOR = "::"
)

type Handler struct {
	Database *redis.Client
}

func NewHandler(database *redis.Client) Handler {
	return Handler{Database: database}
}

type Keys struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Expire int    `json:"expire"`
}

type SplitKeys struct {
	Key       string `json:"key"`
	Separator string `json:"separator"`
}

func (h Handler) AllKeys(w http.ResponseWriter, r *http.Request) {
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

func (h Handler) GroupKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//get separator
	separator := r.URL.Query().Get("separator")

	if separator == "" {
		separator = SEPARATOR
	}

	ctx := context.Background()

	iter := h.Database.Scan(ctx, 0, "*"+separator+"*", 0).Iterator()
	allKeys := []SplitKeys{}
	for iter.Next(ctx) {
		curentKey := iter.Val()
		splitKey := strings.Split(curentKey, "::")
		splitKeyLen := len(splitKey)
		if splitKeyLen > 1 {
			keys := SplitKeys{
				Key:       splitKey[0],
				Separator: separator,
			}
			allKeys = append(allKeys, keys)
		}
	}

	allKeys = removeDuplicate(allKeys)

	if err := iter.Err(); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(allKeys)
}

func removeDuplicate(keys []SplitKeys) []SplitKeys {
	result := []SplitKeys{}
	seen := map[string]string{}
	for _, val := range keys {
		if _, ok := seen[val.Key]; !ok {
			result = append(result, val)
			seen[val.Key] = val.Key
			seen[val.Separator] = val.Separator
		}
	}
	return result
}

func (h Handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key := chi.URLParam(r, "key")

	ctx := context.Background()

	h.Database.Del(ctx, key)

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) DeleteByGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group := chi.URLParam(r, "group")

	separator := r.URL.Query().Get("separator")

	if separator == "" {
		separator = SEPARATOR
	}

	ctx := context.Background()

	iter := h.Database.Scan(ctx, 0, group+separator+"*", 0).Iterator()
	for iter.Next(ctx) {
		h.Database.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) DeleteAllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//uuidStr := chi.URLParam(r, "uuid")

	ctx := context.Background()

	iter := h.Database.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		h.Database.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) GetKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key := chi.URLParam(r, "key")

	ctx := context.Background()

	value := h.Database.Get(ctx, key).Val()

	json.NewEncoder(w).Encode(JSON{"value": value})
}

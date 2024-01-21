package handler

import (
	"context"
	"encoding/json"
	"log"

	//"log"
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
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
	Key    string `json:"id"`
	Value  string `json:"value"`
	Expire int    `json:"expire"`
}

type SplitKeys struct {
	Key       string `json:"id"`
	Separator string `json:"separator"`
}

func (h Handler) getRange(r *http.Request) []string {
	ran := []string{}

	if r.URL.Query().Has("range") {
		rangePage := r.URL.Query().Get("range")
		rangePage = strings.ReplaceAll(rangePage, "[", "")
		rangePage = strings.ReplaceAll(rangePage, "]", "")

		ran = strings.Split(rangePage, ",")
	}

	return ran
}

func (h Handler) AllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ran := h.getRange(r)

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

	count := len(allKeys)
	log.Println(count)

	contentRange := "keys 0-" + strconv.Itoa(count) + "/" + strconv.Itoa(count)
	w.Header().Set("Content-Range", contentRange)

	if len(ran) > 0 {
		offset, _ := strconv.Atoi(ran[0])
		limit, _ := strconv.Atoi(ran[1])
		if limit > count {
			limit = count
		}
		allKeys = allKeys[offset:limit]
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

	ran := h.getRange(r)

	count := len(allKeys)
	if len(ran) > 0 {
		offset, _ := strconv.Atoi(ran[0])
		limit, _ := strconv.Atoi(ran[1])
		if limit > count {
			limit = count
		}
		allKeys = allKeys[offset:limit]
	}

	log.Println(count)
	contentRange := "keys-group 0-0/" + strconv.Itoa(count)
	w.Header().Set("Content-Range", contentRange)
	//w.Header().Set("Content-Range", string(count))
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
	log.Println(key)
	ctx := context.Background()

	h.Database.Del(ctx, key)
	h.Database.Expire(ctx, key, -1)

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
		h.Database.Expire(ctx, iter.Val(), -1)
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
	keyExpire := int(h.Database.TTL(ctx, key).Val().Seconds())
	keyCollection := Keys{
		Key:    key,
		Value:  value,
		Expire: keyExpire,
	}
	json.NewEncoder(w).Encode(keyCollection)
}

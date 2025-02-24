package handler

import (
	"encoding/json"
	"log"
	repository "micro-manager-redis/app/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type JSON map[string]interface{}

const (
	SEPARATOR = "::"
)

type Handler struct {
	RedisRepository repository.RedisRepositoryInterface
}

func NewHandler(rep repository.RedisRepositoryInterface) Handler {
	return Handler{RedisRepository: rep}
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

func getRange(r *http.Request) []string {
	ran := []string{}

	if r.URL.Query().Has("range") {
		rangePage := r.URL.Query().Get("range")
		rangePage = strings.ReplaceAll(rangePage, "[", "")
		rangePage = strings.ReplaceAll(rangePage, "]", "")

		ran = strings.Split(rangePage, ",")
	}

	return ran
}

func getFilter(r *http.Request) string {
	filter := r.URL.Query().Get("filter")

	if filter != "" {
		jsonFilter := JSON{}
		err := json.Unmarshal([]byte(filter), &jsonFilter)
		if err != nil {
			log.Println(err)
		}

		if jsonFilter["key"] != nil {
			filter = jsonFilter["key"].(string)
			return filter
		}
	}

	return ""
}

func (h Handler) AllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ran := getRange(r)

	pattern := "*"
	filter := getFilter(r)
	if filter != "" {
		pattern = pattern + filter + pattern
	}

	allKeys, err := h.RedisRepository.GetAllKeys(pattern)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	count := len(allKeys)

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

	pattern := "*" + separator + "*"
	filter := getFilter(r)
	if filter != "" {
		pattern = "*" + filter + "*" + separator + "*"
	}
	allKeys, err := h.RedisRepository.GroupKeys(pattern, separator)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	allKeys = removeDuplicate(allKeys)

	ran := getRange(r)

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

func removeDuplicate(keys []repository.SplitKeys) []repository.SplitKeys {
	result := []repository.SplitKeys{}
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
	h.RedisRepository.DeleteKey(key)

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) DeleteByGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	group := chi.URLParam(r, "group")

	separator := r.URL.Query().Get("separator")

	if separator == "" {
		separator = SEPARATOR
	}

	err := h.RedisRepository.DeleteByGroup(group + separator + "*")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) DeleteAllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//uuidStr := chi.URLParam(r, "uuid")
	h.RedisRepository.DeleteAllKeys()

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) GetKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key := chi.URLParam(r, "key")

	keyCollection, err := h.RedisRepository.GetKey(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(keyCollection)
}

func (h Handler) GetKeyspaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.RedisRepository.GetKeySpaces()

	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

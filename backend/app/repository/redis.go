package repository

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Keys struct {
	Key    string `json:"id"`
	Value  string `json:"value"`
	Expire int    `json:"expire"`
}

type SplitKeys struct {
	Key       string `json:"id"`
	Separator string `json:"separator"`
}

type RedisRepository struct {
	Database *redis.Client
}

type RedisRepositoryInterface interface {
	GetAllKeys(pattern string) ([]Keys, error)
	GroupKeys(pattern, separator string) ([]SplitKeys, error)
	GetKey(key string) (Keys, error)
	DeleteKey(key string) error
	DeleteAllKeys()
	DeleteByGroup(pattern string) error
	GetActiveKeySpaces() ([]int, error)
	GetCountDb() (int, error)
	SetActiveKeySpace(db int) error
}

func NewRedisRepository(database *redis.Client) RedisRepositoryInterface {
	return &RedisRepository{Database: database}
}

func (r *RedisRepository) GetAllKeys(pattern string) ([]Keys, error) {
	ctx := context.Background()

	iter := r.Database.Scan(ctx, 0, pattern, 0).Iterator()

	allKeys := []Keys{}

	if err := iter.Err(); err != nil {
		return allKeys, err
	}

	for iter.Next(ctx) {
		keys := Keys{
			Key:    iter.Val(),
			Value:  r.Database.Get(ctx, iter.Val()).Val(),
			Expire: int(r.Database.TTL(ctx, iter.Val()).Val().Seconds()),
		}
		allKeys = append(allKeys, keys)
	}

	return allKeys, nil
}

func (r *RedisRepository) GroupKeys(pattern, separator string) ([]SplitKeys, error) {
	ctx := context.Background()

	iter := r.Database.Scan(ctx, 0, pattern, 0).Iterator()

	allKeys := []SplitKeys{}

	if err := iter.Err(); err != nil {
		return allKeys, err
	}

	for iter.Next(ctx) {
		curentKey := iter.Val()
		splitKey := strings.Split(curentKey, separator)
		splitKeyLen := len(splitKey)
		if splitKeyLen > 1 {
			keys := SplitKeys{
				Key:       splitKey[0],
				Separator: separator,
			}
			allKeys = append(allKeys, keys)
		}
	}

	return allKeys, nil
}

func (r *RedisRepository) GetKey(key string) (Keys, error) {
	ctx := context.Background()

	value, err := r.Database.Get(ctx, key).Result()
	if err != nil {
		return Keys{}, err
	}

	keyExpire := int(r.Database.TTL(ctx, key).Val().Seconds())
	return Keys{
		Key:    key,
		Value:  value,
		Expire: keyExpire,
	}, nil
}

func (r *RedisRepository) DeleteKey(key string) error {
	ctx := context.Background()
	return r.Database.Del(ctx, key).Err()
}

func (r *RedisRepository) DeleteAllKeys() {
	ctx := context.Background()
	r.Database.FlushAll(ctx).Result()
}

func (r *RedisRepository) DeleteByGroup(pattern string) error {
	ctx := context.Background()

	iter := r.Database.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.Database.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) GetActiveKeySpaces() ([]int, error) {
	ctx := context.Background()

	value, err := r.Database.InfoMap(ctx, "keyspace").Result()
	if err != nil {
		return nil, err
	}

	dbsInt := []int{}

	for db, keyspace := range value {
		for key, _ := range keyspace {
			db = strings.Replace(key, "db", "", -1)
			dbInt, _ := strconv.Atoi(db)
			dbsInt = append(dbsInt, dbInt)
		}
	}

	return dbsInt, nil
}

func (r *RedisRepository) GetCountDb() (int, error) {
	ctx := context.Background()

	value, err := r.Database.ConfigGet(ctx, "databases").Result()
	if err != nil {
		return 0, err
	}

	if len(value) == 0 {
		return 0, nil
	}

	databases, err := strconv.Atoi(value["databases"])

	if err != nil {
		return 0, err
	}

	return databases, nil
}

func (r *RedisRepository) SetActiveKeySpace(db int) error {
	ctx := context.Background()
	log.Printf("SELECt DB %d", db)
	_, err := r.Database.Do(ctx, "SELECT", db).Result()

	if err != nil {
		return err
	}

	return nil
}

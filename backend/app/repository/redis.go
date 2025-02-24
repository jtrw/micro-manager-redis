package repository

import (
	"context"
	"fmt"
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
	DeleteKey(key string)
	DeleteAllKeys()
	DeleteByGroup(pattern string) error
	GetKeySpaces() ([]string, error)
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

func (r *RedisRepository) DeleteKey(key string) {
	ctx := context.Background()

	r.Database.Del(ctx, key).Result()
	r.Database.Expire(ctx, key, -1)
}

func (r *RedisRepository) DeleteAllKeys() {
	ctx := context.Background()

	r.Database.FlushAll(ctx).Result()

	// iter := r.Database.Scan(ctx, 0, "*", 0).Iterator()
	// for iter.Next(ctx) {
	// 	r.Database.Del(ctx, iter.Val())
	// }
}

func (r *RedisRepository) DeleteByGroup(pattern string) error {
	ctx := context.Background()

	iter := r.Database.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		r.Database.Del(ctx, iter.Val())
		r.Database.Expire(ctx, iter.Val(), -1)
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) GetKeySpaces() ([]int, error) {
	//redis-cli info keyspace
	ctx := context.Background()

	value, err := r.Database.InfoMap(ctx, "keyspace").Result()
	if err != nil {
		return nil, err
	}

	// config, err := r.Database.ConfigGet(ctx, "databases").Result()
	// if err != nil {
	// 	return nil, err
	// }
	dbs := []string{}
	dbsInt := []int{}
	for db, keyspace := range value {
		// get all keys from keyspace

		fmt.Printf("db%v: %v\n", db, keyspace)
		for key, _ := range keyspace {
			dbs = append(dbs, key)
			//replace db and convert to int
			db = strings.Replace(key, "db", "", -1)
			log.Printf("db: %v\n", db)
			dbInt, _ := strconv.Atoi(db)
			dbsInt = append(dbsInt, dbInt)
		}
	}

	fmt.Printf("dbs: %v\n", dbs)
	fmt.Printf("dbsInt: %v\n", dbsInt)
	for _, db := range dbsInt {
		fmt.Printf("db: %v\n", db)
	}

	// for db, keyspace := range value {
	// 	fmt.Printf("db%v: %v\n", db, keyspace)
	// }

	//fmt.Printf("value: %v\n", config)

	//db0:keys=7,expires=7,avg_ttl=51758148
	//db3:keys=576,expires=5,avg_ttl=849200
	//db4:keys=216,expires=216,avg_ttl=1466696

	// rows := strings.Split(value, "\n")
	// //count := len(rows)

	// var dbs []string

	// for _, row := range rows {
	// 	keyRow := strings.Split(row, ":")
	// 	if len(keyRow) > 1 {
	// 		value = keyRow[0]
	// 	}

	// 	dbs = append(dbs, value)

	// }

	// log.Printf("dbs: %v", dbs)

	return nil, nil
}

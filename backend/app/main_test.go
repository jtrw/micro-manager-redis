package main

import (
	"context"
	//"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	server "micro-manager-redis/app/server"
	"testing"
	"time"
)

func TestGetRedisConnection(t *testing.T) {
	// Mock options
	opts := Options{
		RedisUrl:     "localhost:6379",
		RedisPass:    "password",
		Database:     3,
		AuthLogin:    "admin",
		AuthPassword: "admin",
	}

	// Create a mock Redis client
	rdb, err := getRedisConnection(opts)

	// Assert that there is no error during connection
	assert.NoError(t, err)

	// Assert that the connection is not nil
	assert.NotNil(t, rdb)

	// Assert that we can successfully ping the Redis server
	result, err := rdb.Ping(context.Background()).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", result)

	// Close the Redis connection
	err = rdb.Close()
	assert.NoError(t, err)
}

func TestMainRun(t *testing.T) {
	// Mock options
	opts := Options{
		Listen:         ":8080",
		PinSize:        5,
		MaxExpire:      24 * time.Hour,
		MaxPinAttempts: 3,
		WebRoot:        "./web",
		Secret:         "123",
		RedisUrl:       "localhost:6379",
		Database:       3,
		RedisPass:      "password",
		AuthLogin:      "admin",
		AuthPassword:   "admin",
	}

	// Create a mock Redis client
	rdb, err := getRedisConnection(opts)
	assert.NoError(t, err)

	// Mock context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a mock server instance
	srv := server.Server{
		Listen:         opts.Listen,
		PinSize:        opts.PinSize,
		MaxExpire:      opts.MaxExpire,
		MaxPinAttempts: opts.MaxPinAttempts,
		WebRoot:        opts.WebRoot,
		WebFS:          webFS,
		Secret:         opts.Secret,
		Version:        revision,
		Client:         rdb,
		AuthLogin:      opts.AuthLogin,
		AuthPassword:   opts.AuthPassword,
	}

	// Run the server (this will block, so we run it in a goroutine)
	go func() {
		srv.Run(ctx)
		//assert.NoError(t, err)
	}()

	// Simulate an interrupt signal to stop the server
	cancel()
}

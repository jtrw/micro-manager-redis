package main

import (
	//"log"
	//"net/http"
	//"os"
	//"time"
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/redis/go-redis/v9"
	"log"
	server "micro-manager-redis/app/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	Listen         string        `short:"l" long:"listen" env:"LISTEN" default:":8080" description:"listen address"`
	Secret         string        `short:"s" long:"secret" env:"TASKS_SECRET_KEY" default:"123"`
	PinSize        int           `long:"pinszie" env:"PIN_SIZE" default:"5" description:"pin size"`
	MaxExpire      time.Duration `long:"expire" env:"MAX_EXPIRE" default:"24h" description:"max lifetime"`
	MaxPinAttempts int           `long:"pinattempts" env:"PIN_ATTEMPTS" default:"3" description:"max attempts to enter pin"`
	WebRoot        string        `long:"web" env:"WEB" default:"/" description:"web ui location"`
	RedisUrl       string        `long:"redis-url" env:"REDIS_URL" default:"localhost:6379" description:"redis url"`
	Database       string        `long:"redis-db" env:"REDIS_DATABASE" default:"3" description:"database name"`
	RedisPass      string        `long:"redis-pass" env:"REDIS_PASSWORD" default:"Y6zhcj769Fo1" description:"database password"`
}

var revision string

func main() {
	log.Printf("Micro Manager redis %s\n", revision)

	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {

		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic:\n%v", x)
			panic(x)
		}

		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	rdb, err := getRedisConnection(opts)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.Server{
		Listen:         opts.Listen,
		PinSize:        opts.PinSize,
		MaxExpire:      opts.MaxExpire,
		MaxPinAttempts: opts.MaxPinAttempts,
		WebRoot:        opts.WebRoot,
		Secret:         opts.Secret,
		Version:        revision,
		Client:         rdb,
	}
	if err := srv.Run(ctx); err != nil {
		log.Printf("[ERROR] failed, %+v", err)
	}
}

func getRedisConnection(opts Options) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Y6zhcj769Fo1", // no password set
		DB:       3,              // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()

	return rdb, err
}

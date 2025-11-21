package main

import (
	"context"
	"embed"
	"log"
	server "micro-manager-redis/app/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/redis/go-redis/v9"
)

var webFS embed.FS

type Options struct {
	Listen       string `short:"l" long:"listen" env:"LISTEN_SERVER" default:":8080" description:"listen address"`
	Secret       string `short:"s" long:"secret" env:"TASKS_SECRET_KEY" default:"123"`
	WebRoot      string `long:"web" env:"WEB" default:"./web" description:"web ui location"`
	RedisUrl     string `long:"redis-url" env:"REDIS_URL" default:"localhost:6379" description:"redis url"`
	Database     int    `long:"redis-db" env:"REDIS_DATABASE" default:"3" description:"database name"`
	RedisPass    string `long:"redis-pass" env:"REDIS_PASSWORD" default:"" description:"database password"`
	AuthLogin    string `long:"auth-login" env:"AUTH_LOGIN" default:"admin" description:"auth login"`
	AuthPassword string `long:"auth-password" env:"AUTH_PASSWORD" default:"admin" description:"auth password"`
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
		log.Printf("[ERROR] failed, %+v", err)
	}

	srv := server.Server{
		Listen:       opts.Listen,
		WebRoot:      opts.WebRoot,
		WebFS:        webFS,
		Secret:       opts.Secret,
		Version:      revision,
		Client:       rdb,
		AuthLogin:    opts.AuthLogin,
		AuthPassword: opts.AuthPassword,
		Context:      ctx,
	}
	if err := srv.Run(ctx); err != nil {
		log.Printf("[ERROR] failed, %+v", err)
	}
}

func getRedisConnection(opts Options) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.RedisUrl,
		Password: opts.RedisPass, // no password set
		DB:       opts.Database,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()

	return rdb, err
}

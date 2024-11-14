package main

import (
	"context"
	"log"

	"news-weeder/cmd"
	"news-weeder/internal/server/http"
	"news-weeder/internal/weeder/redis"
)

func main() {
	serviceConfig := cmd.Execute()

	redisService := redis.New(&serviceConfig.Redis)
	if err := redisService.Weeder.CreateSchema(); err != nil {
		log.Fatalln("failed to create schema: ", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	httpServer := http.Init(&serviceConfig.Server, redisService)
	go func() {
		err := httpServer.Server.Start(ctx)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	<-ctx.Done()
	cancel()
	if err := httpServer.Server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

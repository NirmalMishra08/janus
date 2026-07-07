package main

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/logger"
	"server/internal/redis"
	"server/internal/router"
)

func main() {

	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.InitRedis(cfg.RedisURL)
	r := router.New(cfg, rdb)


	r.Setup()
   
	addr := fmt.Sprintf(":%s",cfg.Server.Port)
	fmt.Print("Server is running on PORT:", cfg.Server.Port)
	http.ListenAndServe(addr, r.Handler())

}

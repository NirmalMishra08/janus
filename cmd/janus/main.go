package main

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/router"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := router.New(cfg)


	r.Setup()

   
	addr := fmt.Sprintf(":%s",cfg.Server.Port)
	http.ListenAndServe(addr, r.Handler())

}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal/config"
	"server/internal/logger"
	"server/internal/redis"
	"server/internal/router"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	logger.InitLogger()

	defer logger.GetLogger().Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.InitRedis(cfg.RedisURL)
	r := router.New(cfg, rdb)

	r.Setup()

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Print("Server is running on PORT:", cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r.Handler(),
	}

	go func() {
		logger.GetLogger().Info("API Gateway started", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.GetLogger().Info("Shutting down server...")

	// Give 30 seconds to finish ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.GetLogger().Error("Server forced to shutdown", zap.Error(err))
	}

	logger.GetLogger().Info("Server exited gracefully")

}

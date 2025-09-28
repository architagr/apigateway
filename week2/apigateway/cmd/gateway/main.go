package main

import (
	"apigateway/config"
	"apigateway/handler"
	"apigateway/pkg/forwarder"
	"apigateway/pkg/logger"
	"apigateway/pkg/middleware"
	"apigateway/pkg/router"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	//load config
	cfg, err := config.NewConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}
	fmt.Println(cfg)

	//build zap logger instance
	zapConfig := cfg.ZapConfig()
	loggerInstance, err := zapConfig.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer loggerInstance.Sync()

	//initialize loggers for application based and request based logs
	appLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	requestLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	//initialise grpc client
	fwd, err := forwarder.NewForwarder(cfg.Service1URL)
	if err != nil {
		appLogger.Error("Failed to initialize forwarder", zap.Error(err))
		os.Exit(1)
	}

	//initialize gin router
	r := gin.Default()
	r.Use(middleware.RequestLoggerMiddleware(requestLogger))
	h := handler.NewHandler(fwd, requestLogger)
	router.SetupRoutes(r, h)

	//start http server with graceful shutdown support
	srv := &http.Server{
		Addr:    ":" + cfg.GatewayPort,
		Handler: r,
	}

	go func() {
		appLogger.Info("Starting API Gateway", zap.String("port", cfg.GatewayPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed", zap.Error(err))
		}
	}()

	//wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	//graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		appLogger.Info("Server exited gracefully")
	}
}

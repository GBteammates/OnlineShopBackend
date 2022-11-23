package main

import (
	"OnlineShopBackend/cmd/app"
	"OnlineShopBackend/cmd/httpServer"
	"OnlineShopBackend/config"
	"OnlineShopBackend/internal/delivery"
	"OnlineShopBackend/internal/handlers"
	"OnlineShopBackend/internal/logger"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/usecase"
	"context"
	"log"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("can't initialize configuration")
	}
	logger := logger.NewLogger(cfg.LogLevel)
	l := logger.Logger
	store, err := repository.NewPgrepo(cfg.DSN, l)
	if err != nil {
		log.Fatalf("can't initalize storage: %v", err)
	}
	usecase := usecase.NewStorage(store, l)
	handlers := handlers.NewHandlers(usecase, l)
	delivery := delivery.NewDelivery(handlers, l)
	router := app.NewRouter(delivery, l)
	server := httpServer.NewServer(ctx, cfg.Port, router, l)
	err = server.Start(ctx)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
	var services []app.Service

	services = append(services, server)
	a := app.NewApp(services)
	log.Printf("Server started")
	a.Start()
}

package main

import (
	"OnlineShopBackend/config"
	"OnlineShopBackend/internal/app"
	"OnlineShopBackend/internal/app/logger"
	"OnlineShopBackend/internal/app/router"
	"OnlineShopBackend/internal/app/server"
	"OnlineShopBackend/internal/cash"
	"OnlineShopBackend/internal/delivery"
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/usecase"
	"context"
	"log"
	"os"
	"time"
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
		log.Fatalf("can't initialize storage: %v", err)
	}
	cash, err := cash.NewRedisCash(cfg.CashHost, cfg.CashPort, time.Duration(cfg.CashTTL), l)
	if err != nil {
		log.Fatalf("can't initialize cash: %v", err)
	}
	usecase := usecase.NewUsecase(store, store, cash, l)
	handlers := handlers.NewHandlers(usecase, l)
	filestorage := filestorage.NewInMemoryStorage(cfg.FsPath)
	delivery := delivery.NewDelivery(handlers, l, filestorage)
	router := router.NewRouter(delivery, l)
	server := server.NewServer(ctx, cfg.Port, router, l)
	err = server.Start(ctx)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
	var services []app.Service

	services = append(services, server)
	a := app.NewApp(l, services)
	log.Printf("Server started")
	a.Start()
}

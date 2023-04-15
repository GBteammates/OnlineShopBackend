package main

import (
	"OnlineShopBackend/config"
	"OnlineShopBackend/internal/delivery"
	"OnlineShopBackend/internal/delivery/router"
	"OnlineShopBackend/internal/delivery/server"
	"OnlineShopBackend/internal/delivery/user/password"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"OnlineShopBackend/internal/repository/filestorage"
	"OnlineShopBackend/internal/usecase"
	"OnlineShopBackend/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	log.Println("Start load configuration...")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("can't initialize configuration")
	}
	logger := logger.NewLogger(cfg.LogLevel)
	lsug := logger.Logger.Sugar()
	l := logger.Logger

	l.Info("Configuration sucessfully load")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	pgstore, err := repository.NewPgxStorage(ctx, lsug, cfg.DNS)
	if err != nil {
		log.Fatalf("can't initalize storage: %v", err)
	}
	itemStore := repository.NewItemRepo(pgstore, lsug)
	categoryStore := repository.NewCategoryRepo(pgstore, lsug)
	userStore := repository.NewUser(pgstore, lsug)

	setAdmin(userStore, cfg.AdminMail, cfg.AdminPass, l)
	setCustomerRights(userStore, l)

	cartStore := repository.NewCartStore(pgstore, lsug)
	orderStore := repository.NewOrderRepo(pgstore, lsug)

	redis, err := cash.NewRedisCash(cfg.CashHost, cfg.CashPort, time.Duration(cfg.CashTTL), l)
	if err != nil {
		log.Fatalf("can't initialize cash: %v", err)
	}
	itemsCash := cash.NewItemsCash(redis, l)
	categoriesCash := cash.NewCategoriesCash(redis, l)

	filestorage := filestorage.NewOnDiskLocalStorage(cfg.ServerURL, cfg.FsPath, l)

	itemUsecase := usecase.NewItemUsecase(itemStore, itemsCash, filestorage, l)
	categoryUsecase := usecase.NewCategoryUsecase(categoryStore, categoriesCash, l)
	userUsecase := usecase.NewUserUsecase(userStore, l)

	cartUsecase := usecase.NewCartUseCase(cartStore, l)
	orderUsecase := usecase.NewOrderUsecase(orderStore, lsug)

	delivery := delivery.NewDelivery(itemUsecase, userUsecase, categoryUsecase, cartUsecase, l, filestorage, orderUsecase)

	router := router.NewRouter(delivery, l)
	serverOptions := map[string]int{
		"ReadTimeout":       cfg.ReadTimeout,
		"WriteTimeout":      cfg.WriteTimeout,
		"ReadHeaderTimeout": cfg.ReadHeaderTimeout,
	}
	server := server.NewServer(cfg.Port, router, l, serverOptions)

	err = createCashOnStartService(ctx, categoryUsecase, itemUsecase, l)
	if err != nil {
		l.Sugar().Fatalf("error on create cash on start: %v", err)
	}

	server.Start()
	l.Info(fmt.Sprintf("Server start successful on port: %v", cfg.Port))

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			panic(err)
		}
	}()
	<-ctx.Done()

	err = pgstore.ShutDown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Database connection stopped sucessful")
	}

	err = redis.ShutDown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Cash connection stopped successful")
	}

	err = server.ShutDown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Server stopped successful")
	}

	cancel()
}

func createCashOnStartService(ctx context.Context, categoryUsecase usecase.ICategoryUsecase, itemUsecase usecase.IItemUsecase, l *zap.Logger) error {
	l.Debug("Enter in main createCashOnStartService")
	l.Debug("Start create cash...")
	categoryList, err := categoryUsecase.GetCategoryList(ctx)
	if err != nil {
		l.Sugar().Errorf("error on create category cash: %w", err)
		return err
	}
	l.Info("Category list cash create success")

	limitOptions := map[string]int{"offset": 0, "limit": 0}
	listOptions := []map[string]string{
		{"sortType": "name", "sortOrder": "asc"},
		{"sortType": "name", "sortOrder": "desc"},
		{"sortType": "price", "sortOrder": "asc"},
		{"sortType": "price", "sortOrder": "desc"},
	}
	for _, sortOptions := range listOptions {
		_, err = itemUsecase.ItemsList(ctx, limitOptions, sortOptions)
		if err != nil {
			l.Sugar().Errorf("error on create items list cash: %w", err)
			return err
		}
		l.Info("Items list cash create success")

		for _, category := range categoryList {
			_, err := itemUsecase.GetItemsByCategory(ctx, category.Name, limitOptions, sortOptions)
			if err != nil {
				l.Sugar().Errorf("error on create items list in category: %s cash: %w", category.Name, err)
				return err
			}
		}
	}
	l.Info("Items lists in categories cash create success")
	return nil
}

func setAdmin(userStore repository.UserStore, mail string, pass string, logger *zap.Logger) {
	logger.Debug("Enter in main setAdmin()")
	ctx := context.Background()
	exist, err := userStore.GetUserByEmail(ctx, mail)
	logger.Sugar().Debugf("existAdmin is: %v", exist)
	if err != nil {
		logger.Error(err.Error())
	}
	if exist.ID != uuid.Nil {
		logger.Info("User admin is already exists")
		return
	}

	adminRights := &models.Rights{}

	existAdminRights, err := userStore.GetRightsId(ctx, "Admin")
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Sugar().Debugf("ExistAdminRights: %v", existAdminRights)
	if existAdminRights.ID == uuid.Nil {
		adminRights.Name = "Admin"
		adminRights.Rules = []string{"Admin"}

		rightsId, err := userStore.CreateRights(ctx, adminRights)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		adminRights.ID = rightsId
	} else {
		logger.Info("rights admin is already exists")
	}
	newAdmin := &models.User{
		Firstname: "Admin",
		Lastname:  "Admin",
		Email:     mail,
		Password:  pass,
		Rights: models.Rights{
			ID: adminRights.ID,
		},
	}
	hash := password.GeneratePasswordHash(newAdmin.Password)
	newAdmin.Password = hash

	admin, err := userStore.Create(ctx, newAdmin)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if admin != nil {
		logger.Info("Set Admin success")
	} else {
		logger.Warn("Set Admin fail")
	}
}

func setCustomerRights(userStore repository.UserStore, logger *zap.Logger) {
	logger.Debug("Enter in setCustomerRights()")
	ctx := context.Background()

	existRights, err := userStore.GetRightsId(ctx, "Customer")
	if err != nil {
		logger.Error(err.Error())
	}
	if existRights.ID != uuid.Nil {
		logger.Sugar().Debugf("ExistCustomerRights: %v", existRights)
		return
	}
	customerRights := models.Rights{
		Name:  "Customer",
		Rules: []string{"Customer"},
	}
	rightsId, err := userStore.CreateRights(ctx, &customerRights)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Sugar().Infof("Customer rights with id: %v create success", rightsId)
}

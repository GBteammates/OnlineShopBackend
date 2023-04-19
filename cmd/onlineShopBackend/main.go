package main

import (
	"OnlineShopBackend/config"
	carts "OnlineShopBackend/internal/delivery/carts"
	categories "OnlineShopBackend/internal/delivery/categories"
	items "OnlineShopBackend/internal/delivery/items"
	orders "OnlineShopBackend/internal/delivery/orders"
	"OnlineShopBackend/internal/delivery/router"
	users "OnlineShopBackend/internal/delivery/users"
	"OnlineShopBackend/internal/delivery/users/user/password"
	"OnlineShopBackend/internal/models"
	redisCache "OnlineShopBackend/internal/repository/cache/redis"
	"OnlineShopBackend/internal/repository/db/postgres"
	"OnlineShopBackend/internal/repository/filestorage"
	cartUsecase "OnlineShopBackend/internal/usecase/carts"
	categoryUsecase "OnlineShopBackend/internal/usecase/categories"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	itemUsecase "OnlineShopBackend/internal/usecase/items"
	orderUsecase "OnlineShopBackend/internal/usecase/orders"
	userUsecase "OnlineShopBackend/internal/usecase/users"
	"OnlineShopBackend/logger"
	"OnlineShopBackend/server"
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

	pgstore, err := postgres.NewPgxStorage(ctx, lsug, cfg.DNS)
	if err != nil {
		log.Fatalf("can't initalize storage: %v", err)
	}
	itemStore := postgres.NewItemRepo(pgstore, lsug)
	categoryStore := postgres.NewCategoryRepo(pgstore, lsug)
	userStore := postgres.NewUser(pgstore, lsug)
	cartStore := postgres.NewCartStore(pgstore, lsug)
	orderStore := postgres.NewOrderRepo(pgstore, lsug)

	setAdmin(userStore, cfg.AdminMail, cfg.AdminPass, l)
	setCustomerRights(userStore, l)

	redis, err := redisCache.NewRedisCache(cfg.CashHost, cfg.CashPort, time.Duration(cfg.CashTTL), l)
	if err != nil {
		log.Fatalf("can't initialize cache: %v", err)
	}
	itemsCache := redisCache.NewItemsCache(redis, l)
	categoriesCache := redisCache.NewCategoriesСache(redis, l)

	filestorage := filestorage.NewFileStorage(cfg.ServerURL, cfg.FsPath, l)

	itemUsecase := itemUsecase.NewItemUsecase(itemStore, itemsCache, filestorage, l)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryStore, categoriesCache, l)
	userUsecase := userUsecase.NewUserUsecase(userStore, l)
	cartUsecase := cartUsecase.NewCartUsecase(cartStore, l)
	orderUsecase := orderUsecase.NewOrderUsecase(orderStore, lsug)

	itemDelivery := items.NewItemDelivery(itemUsecase, categoryUsecase, lsug)
	categoryDelivery := categories.NewCategoryDelivery(categoryUsecase, itemUsecase, lsug)
	cartDelivery := carts.NewCartDelivery(cartUsecase, lsug)
	orderDelivery := orders.NewOrderDelivery(orderUsecase, cartUsecase, lsug)
	userDelivery := users.NewUserDelivery(userUsecase, cartUsecase, lsug)

	router := router.NewRouter(itemDelivery, categoryDelivery, cartDelivery, orderDelivery, userDelivery, l)
	serverOptions := map[string]int{
		"ReadTimeout":       cfg.ReadTimeout,
		"WriteTimeout":      cfg.WriteTimeout,
		"ReadHeaderTimeout": cfg.ReadHeaderTimeout,
	}
	server := server.NewServer(cfg.Port, router, l, serverOptions)

	err = createCacheOnStartService(ctx, categoryUsecase, itemUsecase, l)
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

	err = pgstore.Shutdown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Database connection stopped sucessful")
	}

	err = redis.Shutdown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Cash connection stopped successful")
	}

	err = server.Shutdown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Server stopped successful")
	}

	cancel()
}

func createCacheOnStartService(ctx context.Context, categoryUsecase usecase.ICategoryUsecase, itemUsecase usecase.IItemUsecase, l *zap.Logger) error {
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

func setAdmin(userStore usecase.UserStore, mail string, pass string, logger *zap.Logger) {
	logger.Debug("Enter in main setAdmin()")
	ctx := context.Background()
	exist, err := userStore.GetUserByEmail(ctx, mail)
	logger.Sugar().Debugf("existAdmin is: %v", exist)
	if err != nil {
		logger.Error(err.Error())
	}
	if exist.Id != uuid.Nil {
		logger.Info("User admin is already exists")
		return
	}

	adminRights := &models.Rights{}

	existAdminRights, err := userStore.GetRightsId(ctx, "Admin")
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Sugar().Debugf("ExistAdminRights: %v", existAdminRights)
	if existAdminRights.Id == uuid.Nil {
		adminRights.Name = "Admin"
		adminRights.Rules = []string{"Admin"}

		rightsId, err := userStore.CreateRights(ctx, adminRights)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		adminRights.Id = rightsId
	} else {
		logger.Info("rights admin is already exists")
	}
	newAdmin := &models.User{
		Firstname: "Admin",
		Lastname:  "Admin",
		Email:     mail,
		Password:  pass,
		Rights: models.Rights{
			Id: adminRights.Id,
		},
	}
	hash := password.GeneratePasswordHash(newAdmin.Password)
	newAdmin.Password = hash

	adminId, err := userStore.CreateUser(ctx, newAdmin)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if adminId != uuid.Nil {
		logger.Info("Set Admin success")
	} else {
		logger.Warn("Set Admin fail")
	}
}

func setCustomerRights(userStore usecase.UserStore, logger *zap.Logger) {
	logger.Debug("Enter in setCustomerRights()")
	ctx := context.Background()

	existRights, err := userStore.GetRightsId(ctx, "Customer")
	if err != nil {
		logger.Error(err.Error())
	}
	if existRights.Id != uuid.Nil {
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

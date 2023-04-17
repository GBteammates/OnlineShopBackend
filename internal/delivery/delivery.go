package delivery

import (
	"OnlineShopBackend/internal/metrics"
	"OnlineShopBackend/internal/usecase"

	"go.uber.org/zap"
)

//	@title			Online Shop Backend Service
//	@version		1.0
//	@description	Backend service for online store
//	@license.name	MIT

//	@contact.url	https://github.com/GBteammates/OnlineShopBackend

//	@BasePath	/

type Delivery struct {
	itemUsecase     usecase.IItemUsecase
	categoryUsecase usecase.ICategoryUsecase
	userUsecase     usecase.IUserUsecase
	cartUsecase     usecase.ICartUsecase
	logger          *zap.Logger
	orderUsecase    usecase.IOrderUsecase
}

// NewDelivery initialize delivery layer
func NewDelivery(
	itemUsecase usecase.IItemUsecase,
	userUsecase usecase.IUserUsecase,
	categoryUsecase usecase.ICategoryUsecase,
	cartUsecase usecase.ICartUsecase,
	logger *zap.Logger,
	orderUsecase usecase.IOrderUsecase,
) *Delivery {
	logger.Debug("Enter in NewDelivery()")
	metrics.DeliveryMetrics.NewDeliveryTotal.Inc()

	return &Delivery{
		itemUsecase:     itemUsecase,
		categoryUsecase: categoryUsecase,
		cartUsecase:     cartUsecase,
		userUsecase:     userUsecase,
		logger:          logger,
		orderUsecase:    orderUsecase,
	}
}

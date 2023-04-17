package delivery

import (
	"OnlineShopBackend/internal/delivery/file"
	"OnlineShopBackend/internal/metrics"
	"OnlineShopBackend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
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

// Index is the index handler.
func (delivery *Delivery) Index(c *gin.Context) {
	delivery.logger.Debug("Enter in Index")
	c.String(http.StatusOK, "Hello World!")
}

// GetFileList returns list of files
//
//	@Summary		Get list of files
//	@Description	Method provides to get list of files.
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	file.FileListResponse	"List of files"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/images/list [get]
func (delivery *Delivery) GetFileList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetFileList()")
	fileInfos, err := delivery.filestorage.GetFileList()
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	var files file.FileListResponse
	files.Files = make([]file.FilesInfo, len(fileInfos))
	for i, info := range fileInfos {
		files.Files[i] = file.FilesInfo{
			Name:       info.Name,
			Path:       info.Path,
			CreateDate: info.CreateDate,
			ModifyDate: info.ModifyDate,
		}
	}

	c.JSON(http.StatusOK, files)
}

package postgres

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

var _ usecase.ItemStore = (*itemRepo)(nil)

type itemRepo struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

func NewItemRepo(storage *PGres, logger *zap.SugaredLogger) usecase.ItemStore {
	logger.Debug("Enter in repository NewItemRepo()")
	return &itemRepo{
		storage: storage,
		logger:  logger,
	}
}

// CreateItem insert new item in database
func (repo *itemRepo) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	repo.logger.Debugf("Enter in repository CreateItem() with args: ctx, item: %v", item)

	pool := repo.storage.GetPool()

	// Recording operations need transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return uuid.Nil, fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()
	var id uuid.UUID
	row := tx.QueryRow(ctx, `INSERT INTO items(name, category, description, price, vendor, pictures, deleted_at)
	values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		item.Title,
		item.Category.Id,
		item.Description,
		item.Price,
		item.Vendor,
		item.Images,
		nil,
	)
	err = row.Scan(&id)
	if err != nil {
		repo.logger.Errorf("can't create item %s", err)
		return uuid.Nil, fmt.Errorf("can't create item %w", err)
	}
	repo.logger.Info("Item create success")
	repo.logger.Debugf("id is %v\n", id)
	return id, nil
}

// UpdateItem сhanges the existing item
func (repo *itemRepo) UpdateItem(ctx context.Context, item *models.Item) error {
	repo.logger.Debugf("Enter in repository UpdateItem() with args: ctx, item: %v", item)

	pool := repo.storage.GetPool()

	// Recording operations need transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("Can't rollback %s", err)
			}
		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()

	_, err = tx.Exec(ctx, `UPDATE items SET name=$1, category=$2, description=$3, price=$4, vendor=$5, pictures = $6 WHERE id=$7`,
		item.Title,
		item.Category.Id,
		item.Description,
		item.Price,
		item.Vendor,
		item.Images,
		item.Id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error on update item %s: %s", item.Id, err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error on update item %s: %s", item.Id, err)
		return fmt.Errorf("error on update item %s: %w", item.Id, err)
	}
	repo.logger.Infof("Item %s successfully updated", item.Id)
	return nil
}

// GetItem returns *models.Item by id or error
func (repo *itemRepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	repo.logger.Debug("Enter in repository GetItem() with args: ctx, id: %v", id)

	pool := repo.storage.GetPool()

	item := models.Item{}
	row := pool.QueryRow(ctx, `
	SELECT 
	items.id, 
	items.name, 
	category, 
	categories.name, 
	categories.description, 
	categories.picture, 
	items.description, 
	price, 
	vendor, 
	pictures 
	FROM items 
	INNER JOIN categories 
	ON category=categories.id 
	AND items.id = $1 
	WHERE items.deleted_at is null 
	AND categories.deleted_at is null
	`, id)
	err := row.Scan(
		&item.Id,
		&item.Title,
		&item.Category.Id,
		&item.Category.Name,
		&item.Category.Description,
		&item.Category.Image,
		&item.Description,
		&item.Price,
		&item.Vendor,
		&item.Images,
	)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error in rows scan get item by id: %s", err)
		return &models.Item{}, models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error in rows scan get item by id: %s", err)
		return &models.Item{}, fmt.Errorf("error in rows scan get item by id: %w", err)
	}
	repo.logger.Info("Get item success")
	return &item, nil
}

// ItemsList reads all the items from the database and writes it to the
// output channel and returns this channel or error
func (repo *itemRepo) ItemsList(ctx context.Context) (chan models.Item, error) {
	repo.logger.Debug("Enter in repository ItemsList() with args: ctx")
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		pool := repo.storage.GetPool()

		item := &models.Item{}
		rows, err := pool.Query(ctx, `
		SELECT
		items.id,
		items.name, 
		category, 
		categories.name, 
		categories.description, 
		categories.picture, 
		items.description, 
		price, 
		vendor, 
		pictures 
		FROM items 
		INNER JOIN categories 
		ON category=categories.id 
		WHERE items.deleted_at is null 
		AND categories.deleted_at is null
		`)
		if err != nil {
			msg := fmt.Errorf("error on items list query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Description,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

// SearchLine allows to find all the items that satisfy the parameters from the search query and writes them to the output channel
func (repo *itemRepo) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	repo.logger.Debugf("Enter in repository SearchLine() with args: ctx, param: %s", param)

	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}
		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT 
		items.id, 
		items.name, 
		category, 
		categories.name, 
		categories.description,
		categories.picture, 
		items.description, 
		price, 
		vendor, 
		pictures 
		FROM items 
		INNER JOIN categories 
		ON category=categories.id 
		WHERE items.deleted_at is null 
		AND categories.deleted_at is null
		AND items.name ilike $1 
		OR items.description ilike $1 
		OR vendor ilike $1 
		OR categories.name ilike $1
		`, "%"+param+"%")
		if err != nil {
			msg := fmt.Errorf("error on search line query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Description,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			repo.logger.Info(fmt.Sprintf("find item: %v", item))
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

// GetItemsByCategory finds in the database all the items with a certain name of the category and writes them in the outgoing channel
func (repo *itemRepo) GetItemsByCategory(ctx context.Context, categoryName string) (chan models.Item, error) {
	repo.logger.Debugf("Enter in repository GetItemsByCategory() with args: ctx, categoryName: %s", categoryName)
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}
		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT items.id, 
		items.name, 
		category, 
		categories.name, 
		categories.description,
		categories.picture, 
		items.description, 
		price, 
		vendor, 
		pictures FROM items 
		INNER JOIN categories ON category=categories.id 
		WHERE items.deleted_at is null 
		AND categories.deleted_at is null 
		AND categories.name=$1
		`, categoryName)
		if err != nil {
			msg := fmt.Errorf("error on get items by category query context: %w", err)
			repo.logger.Error(msg.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Description,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

// DeleteItem changes the value of the deleted_at attribute in the deleted item for the current time
func (repo *itemRepo) DeleteItem(ctx context.Context, id uuid.UUID) error {
	repo.logger.Debugf("Enter in repository DeleteItem() with args: ctx, id: %v", id)
	pool := repo.storage.GetPool()

	// Removal operation is carried out in transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repo.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	repo.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			repo.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				repo.logger.Errorf("can't rollback %s", err)
			}

		} else {
			repo.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				repo.logger.Errorf("Can't commit %s", err)
			}
		}
	}()
	_, err = tx.Exec(ctx, `UPDATE items SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error on delete item %s: %s", id, err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("Error on delete item %s: %s", id, err)
		return fmt.Errorf("error on delete item %s: %w", id, err)
	}
	repo.logger.Infof("Item with id: %s successfully deleted from database", id)
	return nil
}

// AddFavouriteItem adds item to the list of favourites for a specific user
func (repo *itemRepo) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	repo.logger.Debug("Enter in repository AddFavouriteItem() with args: ctx, userid: %v, itemId: %v", userId, itemId)
	pool := repo.storage.GetPool()
	_, err := pool.Exec(ctx, `INSERT INTO favourite_items (user_id, item_id) VALUES ($1, $2)`, userId, itemId)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("can't add item to favourite_items: %s", err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("can't add item to favourite_items: %s", err)
		return fmt.Errorf("can't add item to favourite_items: %w", err)
	}
	return nil
}

// DeleteFavouriteItem deletes item from the list of favourites for a specific user
func (repo *itemRepo) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	repo.logger.Debug("Enter in repository DeleteFavouriteItem() with args: ctx, userid: %v, itemId: %v", userId, itemId)
	pool := repo.storage.GetPool()
	_, err := pool.Exec(ctx, `DELETE FROM favourite_items WHERE user_id=$1 AND item_id=$2`, userId, itemId)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("can't delete item from favourite: %s", err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		repo.logger.Errorf("can't delete item from favourite: %s", err)
		return fmt.Errorf("can't delete item from favourite: %w", err)
	}
	repo.logger.Info("Delete item from cart success")
	return nil
}

// GetItemsFavouriteItems finds in the database all the items in list of favourites for current user
// and writes them in the output channel
func (repo *itemRepo) GetFavouriteItems(ctx context.Context, userId uuid.UUID) (chan models.Item, error) {
	repo.logger.Debug("Enter in repository GetFavouriteItems() with args: ctx, userId: %v", userId)

	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		pool := repo.storage.GetPool()
		item := models.Item{}
		rows, err := pool.Query(ctx, `
		SELECT 	
		i.id, 
		i.name, 
		i.description, 
		i.category, 
		cat.name, 
		cat.description, 
		cat.picture, 
		i.price, 
		i.vendor, 
		i.pictures
		FROM favourite_items f, items i, categories cat
		WHERE f.user_id=$1 
		AND i.id = f.item_id 
		AND cat.id = i.category
		AND i.deleted_at IS NULL
		`, userId)
		if err != nil {
			repo.logger.Errorf("can't select items from favourite_items: %s", err)
			return
		}
		defer rows.Close()
		repo.logger.Debug("read info from db in pool.Query success")
		for rows.Next() {
			if err := rows.Scan(
				&item.Id,
				&item.Title,
				&item.Description,
				&item.Category.Id,
				&item.Category.Name,
				&item.Category.Description,
				&item.Category.Image,
				&item.Price,
				&item.Vendor,
				&item.Images,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			itemChan <- item
		}
	}()
	repo.logger.Info("Select items from favourites success")
	return itemChan, nil
}

// GetFavouriteItemsId returns list of identificators of favourite items for current user
func (repo *itemRepo) GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	repo.logger.Debug("Enter in repository GetFavouriteItemsId() with args: ctx, userId: %v", userId)

	pool := repo.storage.GetPool()

	result := make(map[uuid.UUID]uuid.UUID)
	item := models.Item{}
	rows, err := pool.Query(ctx, `
		SELECT 	i.id FROM favourite_items f, items i WHERE f.user_id=$1 and i.id = f.item_id`, userId)
	if err != nil {
		repo.logger.Errorf("can't select items from favourite_items: %s", err)
		return nil, err
	}
	defer rows.Close()
	repo.logger.Debug("read info from db in pool.Query success")
	for rows.Next() {
		err := rows.Scan(
			&item.Id,
		)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			repo.logger.Info("this user don't have favourite items")
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			repo.logger.Error(err.Error())
			return nil, err
		}
		result[item.Id] = userId
	}
	return &result, nil
}

// ItemsListQuantity returns quantity of all items or error
func (repo *itemRepo) ItemsListQuantity(ctx context.Context) (int, error) {
	repo.logger.Debug("Enter in repository ItemsListQuantity() with args: ctx")
	pool := repo.storage.GetPool()
	var quantity int
	row := pool.QueryRow(ctx, `SELECT COUNT(1) FROM items WHERE deleted_at IS NULL`)
	err := row.Scan(&quantity)
	if err != nil {
		repo.logger.Errorf("Error in row.Scan items list quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items list quantity: %w", err)
	}
	repo.logger.Info("Request for ItemsListQuantity success")
	return quantity, nil
}

// ItemsByCategoryQuantity returns quntity of items in category or error
func (repo *itemRepo) ItemsByCategoryQuantity(ctx context.Context, categoryName string) (int, error) {
	repo.logger.Debug("Enter in repository ItemsByCategoryQuantity() with args: ctx, categoryName: %s", categoryName)
	pool := repo.storage.GetPool()
	var quantity int
	row := pool.QueryRow(ctx, `
	SELECT COUNT(1) FROM items 
	INNER JOIN categories ON category=categories.id 
	WHERE items.deleted_at is null 
	AND categories.deleted_at is null 
	AND categories.name=$1
	`, categoryName)
	err := row.Scan(&quantity)
	if err != nil {
		repo.logger.Errorf("Error in row.Scan items by category quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items by category quantity: %w", err)
	}
	repo.logger.Info("Request for ItemsByCategoryQuantity success")
	return quantity, nil
}

// ItemsInSearchQuantity returns quantity of items in search results or error
func (repo *itemRepo) ItemsInSearchQuantity(ctx context.Context, searchRequest string) (int, error) {
	repo.logger.Debug("Enter in repository ItemsInSearchQuantity() with args: ctx, searchRequest: %s", searchRequest)
	pool := repo.storage.GetPool()
	var quantity int
	row := pool.QueryRow(ctx, `
		SELECT COUNT(1) 
		FROM items 
		INNER JOIN categories 
		ON category=categories.id 
		WHERE items.deleted_at is null 
		AND categories.deleted_at is null
		AND items.name ilike $1 
		OR items.description ilike $1 
		OR vendor ilike $1 
		OR categories.name ilike $1
		`, "%"+searchRequest+"%")
	err := row.Scan(&quantity)
	if err != nil {
		repo.logger.Errorf("Error in row.Scan items in search quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items in search quantity: %w", err)
	}
	repo.logger.Info("Request for ItemsInSearchQuantity success")
	return quantity, nil
}

// ItemsInFavouriteQuantity returns quantity or favourite items by user id or error
func (repo *itemRepo) ItemsInFavouriteQuantity(ctx context.Context, userId uuid.UUID) (int, error) {
	repo.logger.Debug("Enter in repository ItemsInFavouriteQuantity() with args: ctx, userId uuid.UUID: %v", userId)
	pool := repo.storage.GetPool()
	var quantity int
	row := pool.QueryRow(ctx, `
	SELECT COUNT(1) 
	FROM favourite_items f, items i
	WHERE f.user_id=$1 
	AND i.id = f.item_id
	AND i.deleted_at IS NULL
	`, userId)
	err := row.Scan(&quantity)
	if err != nil {
		repo.logger.Errorf("Error in row.Scan items in favourite quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items in favourite quantity: %w", err)
	}
	repo.logger.Info("Request for ItemsInFavouriteQuantity success")
	return quantity, nil
}

package postgres

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

var _ usecase.ItemStore = (*itemStore)(nil)

type itemStore struct {
	db     *PGres
	logger *zap.SugaredLogger
}

func NewItemStore(db *PGres, logger *zap.SugaredLogger) *itemStore {
	logger.Debug("Enter in repository NewItemRepo()")
	return &itemStore{
		db:     db,
		logger: logger,
	}
}

// CreateItem insert new item in database
func (r *itemStore) Create(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	r.logger.Debugf("Enter in items repository Create() with args: ctx, item: %v", item)

	pool := r.db.GetPool()

	// Recording operations need transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Can't create transaction: %s", err)
		return uuid.Nil, fmt.Errorf("can't create transaction: %w", err)
	}
	r.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			r.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			r.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				r.logger.Errorf("Can't commit %s", err)
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
		r.logger.Errorf("can't create item %s", err)
		return uuid.Nil, fmt.Errorf("can't create item %w", err)
	}
	r.logger.Info("Item create success")
	r.logger.Debugf("id is %v\n", id)
	return id, nil
}

// UpdateItem сhanges the existing item
func (r *itemStore) Update(ctx context.Context, item *models.Item) error {
	r.logger.Debugf("Enter in items repository Update() with args: ctx, item: %v", item)

	pool := r.db.GetPool()

	// Recording operations need transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	r.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			r.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Errorf("Can't rollback %s", err)
			}
		} else {
			r.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				r.logger.Errorf("Can't commit %s", err)
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
		r.logger.Errorf("Error on update item %s: %s", item.Id, err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error on update item %s: %s", item.Id, err)
		return fmt.Errorf("error on update item %s: %w", item.Id, err)
	}
	r.logger.Infof("Item %s successfully updated", item.Id)
	return nil
}

// GetItem returns *models.Item by id or error
func (r *itemStore) Get(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	r.logger.Debug("Enter in items repository Get() with args: ctx, id: %v", id)

	pool := r.db.GetPool()

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
		r.logger.Errorf("Error in rows scan get item by id: %s", err)
		return &models.Item{}, models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error in rows scan get item by id: %s", err)
		return &models.Item{}, fmt.Errorf("error in rows scan get item by id: %w", err)
	}
	r.logger.Info("Get item success")
	return &item, nil
}

// List reads all the items from the database and writes it to the
// output channel and returns this channel or error
func (r *itemStore) List(ctx context.Context, param string) (chan models.Item, error) {
	r.logger.Debug("Enter in repository List() with args: ctx")
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		pool := r.db.GetPool()

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
			r.logger.Error(msg.Error())
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
				r.logger.Error(err.Error())
				return
			}
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

// SearchLine allows to find all the items that satisfy the parameters from the search query and writes them to the output channel
func (r *itemStore) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	r.logger.Debugf("Enter in items repository SearchLine() with args: ctx, param: %s", param)

	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}
		pool := r.db.GetPool()
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
			r.logger.Error(msg.Error())
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
				r.logger.Error(err.Error())
				return
			}
			r.logger.Info(fmt.Sprintf("find item: %v", item))
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

// ListByCategory finds in the database all the items with a certain name of the category and writes them in the outgoing channel
func (r *itemStore) ListByCategory(ctx context.Context, param string) (chan models.Item, error) {
	r.logger.Debugf("Enter in repository ListByCategory() with args: ctx, param: %s", param)
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}
		pool := r.db.GetPool()
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
		`, param)
		if err != nil {
			msg := fmt.Errorf("error on get items by category query context: %w", err)
			r.logger.Error(msg.Error())
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
				r.logger.Error(err.Error())
				return
			}
			itemChan <- *item
		}
	}()
	return itemChan, nil
}

// Delete changes the value of the deleted_at attribute in the deleted item for the current time
func (r *itemStore) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debugf("Enter in repository Delete() with args: ctx, id: %v", id)
	pool := r.db.GetPool()

	// Removal operation is carried out in transaction
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Can't create transaction: %s", err)
		return fmt.Errorf("can't create transaction: %w", err)
	}
	r.logger.Debug("Transaction begin success")
	defer func() {
		if err != nil {
			r.logger.Errorf("Transaction rolled back")
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Errorf("can't rollback %s", err)
			}

		} else {
			r.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				r.logger.Errorf("Can't commit %s", err)
			}
		}
	}()
	_, err = tx.Exec(ctx, `UPDATE items SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error on delete item %s: %s", id, err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error on delete item %s: %s", id, err)
		return fmt.Errorf("error on delete item %s: %w", id, err)
	}
	r.logger.Infof("Item with id: %s successfully deleted from database", id)
	return nil
}

// AddFavouriteItem adds item to the list of favourites for a specific user
func (r *itemStore) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	r.logger.Debug("Enter in repository AddFavouriteItem() with args: ctx, userid: %v, itemId: %v", userId, itemId)
	pool := r.db.GetPool()
	_, err := pool.Exec(ctx, `INSERT INTO favourite_items (user_id, item_id) VALUES ($1, $2)`, userId, itemId)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("can't add item to favourite_items: %s", err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("can't add item to favourite_items: %s", err)
		return fmt.Errorf("can't add item to favourite_items: %w", err)
	}
	return nil
}

// DeleteFavouriteItem deletes item from the list of favourites for a specific user
func (r *itemStore) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	r.logger.Debug("Enter in repository DeleteFavouriteItem() with args: ctx, userid: %v, itemId: %v", userId, itemId)
	pool := r.db.GetPool()
	_, err := pool.Exec(ctx, `DELETE FROM favourite_items WHERE user_id=$1 AND item_id=$2`, userId, itemId)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("can't delete item from favourite: %s", err)
		return models.ErrorNotFound{}
	} else if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("can't delete item from favourite: %s", err)
		return fmt.Errorf("can't delete item from favourite: %w", err)
	}
	r.logger.Info("Delete item from cart success")
	return nil
}

// FavouriteItemsList finds in the database all the items in list of favourites for current user
// and writes them in the output channel
func (r *itemStore) ListFavouriteItems(ctx context.Context, userId string) (chan models.Item, error) {
	r.logger.Debug("Enter in repository ListFavouriteItems() with args: ctx, userId: %v", userId)

	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		pool := r.db.GetPool()
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
			r.logger.Errorf("can't select items from favourite_items: %s", err)
			return
		}
		defer rows.Close()
		r.logger.Debug("read info from db in pool.Query success")
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
				r.logger.Error(err.Error())
				return
			}
			itemChan <- item
		}
	}()
	r.logger.Info("Select items from favourites success")
	return itemChan, nil
}

// FavouriteItemsId returns list of identificators of favourite items for current user
func (r *itemStore) FavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	r.logger.Debug("Enter in repository FavouriteItemsId() with args: ctx, userId: %v", userId)

	pool := r.db.GetPool()

	result := make(map[uuid.UUID]uuid.UUID)
	item := models.Item{}
	rows, err := pool.Query(ctx, `
		SELECT 	i.id FROM favourite_items f, items i WHERE f.user_id=$1 and i.id = f.item_id`, userId)
	if err != nil {
		r.logger.Errorf("can't select items from favourite_items: %s", err)
		return nil, err
	}
	defer rows.Close()
	r.logger.Debug("read info from db in pool.Query success")
	for rows.Next() {
		err := rows.Scan(
			&item.Id,
		)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			r.logger.Info("this user don't have favourite items")
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			r.logger.Error(err.Error())
			return nil, err
		}
		result[item.Id] = userId
	}
	return &result, nil
}

// ListQuantity returns quantity of all items or error
func (r *itemStore) ListQuantity(ctx context.Context, param string) (int, error) {
	r.logger.Debug("Enter in repository ListQuantity() with args: ctx")
	pool := r.db.GetPool()
	var quantity int
	row := pool.QueryRow(ctx, `SELECT COUNT(1) FROM items WHERE deleted_at IS NULL`)
	err := row.Scan(&quantity)
	if err != nil {
		r.logger.Errorf("Error in row.Scan items list quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items list quantity: %w", err)
	}
	r.logger.Info("Request for ItemsListQuantity success")
	return quantity, nil
}

// ListByCategoryQuantity returns quntity of items in category or error
func (r *itemStore) ListByCategoryQuantity(ctx context.Context, param string) (int, error) {
	r.logger.Debug("Enter in repository ListByCategoryQuantity() with args: ctx, categoryName: %s", param)
	pool := r.db.GetPool()
	var quantity int
	row := pool.QueryRow(ctx, `
	SELECT COUNT(1) FROM items 
	INNER JOIN categories ON category=categories.id 
	WHERE items.deleted_at is null 
	AND categories.deleted_at is null 
	AND categories.name=$1
	`, param)
	err := row.Scan(&quantity)
	if err != nil {
		r.logger.Errorf("Error in row.Scan items by category quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items by category quantity: %w", err)
	}
	r.logger.Info("Request for ItemsByCategoryQuantity success")
	return quantity, nil
}

// InSearchQuantity returns quantity of items in search results or error
func (r *itemStore) InSearchQuantity(ctx context.Context, param string) (int, error) {
	r.logger.Debug("Enter in repository InSearchQuantity() with args: ctx, searchRequest: %s", param)
	pool := r.db.GetPool()
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
		`, "%"+param+"%")
	err := row.Scan(&quantity)
	if err != nil {
		r.logger.Errorf("Error in row.Scan items in search quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items in search quantity: %w", err)
	}
	r.logger.Info("Request for ItemsInSearchQuantity success")
	return quantity, nil
}

// InFavouriteQuantity returns quantity or favourite items by user id or error
func (r *itemStore) InFavouriteQuantity(ctx context.Context, userId string) (int, error) {
	r.logger.Debug("Enter in repository InFavouriteQuantity() with args: ctx, userId uuid.UUID: %v", userId)
	pool := r.db.GetPool()
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
		r.logger.Errorf("Error in row.Scan items in favourite quantity: %s", err)
		return -1, fmt.Errorf("error in row.Scan items in favourite quantity: %w", err)
	}
	r.logger.Info("Request for InFavouriteQuantity success")
	return quantity, nil
}

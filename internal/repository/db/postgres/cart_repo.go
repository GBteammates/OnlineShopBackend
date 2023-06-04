package postgres

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type cartStore struct {
	db     *PGres
	logger *zap.SugaredLogger
}

var _ usecase.CartStore = (*cartStore)(nil)

func NewCartRepo(db *PGres, logger *zap.SugaredLogger) *cartStore {
	return &cartStore{
		db:     db,
		logger: logger,
	}
}

// Create Shall we add items at the moment we create cart
func (r *cartStore) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	r.logger.Debugf("Enter in repository cart Create() with args: ctx, userId: %v", userId)
	select {
	case <-ctx.Done():
		return uuid.Nil, fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		var cartId uuid.UUID
		row := pool.QueryRow(ctx, `INSERT INTO carts (user_id) VALUES ($1) RETURNING id`,
			userId)
		err := row.Scan(&cartId)
		if err != nil {
			r.logger.Error(err)
			return uuid.Nil, fmt.Errorf("can't create cart object: %w", err)
		}
		r.logger.Info("Create cart success")
		return cartId, nil
	}
}

// AddItemToCart Maybe add to item
func (r *cartStore) AddItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	r.logger.Debugf("Enter in repository cart AddItemToCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	select {
	case <-ctx.Done():
		r.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		row := pool.QueryRow(ctx, `SELECT item_id from cart_items where item_id=$1 and cart_id=$2`, itemId, cartId)
		var checkId uuid.UUID
		err := row.Scan(&checkId)
		if err != nil {
			r.logger.Errorf("error on row.Scan: %s", err)
		}
		if checkId == uuid.Nil {
			_, err := pool.Exec(ctx, `INSERT INTO cart_items (cart_id, item_id, item_quantity) VALUES ($1, $2, $3)`, cartId, itemId, 1)
			if err != nil {
				r.logger.Errorf("can't add item to cart: %s", err)
				return fmt.Errorf("can't add item to cart: %w", err)
			}
		} else {
			_, err := pool.Exec(ctx, `UPDATE cart_items SET item_quantity = item_quantity + 1 WHERE cart_id=$1 and item_id=$2`, cartId, itemId)
			if err != nil {
				r.logger.Errorf("can't add item to cart: %s", err)
				return fmt.Errorf("can't add item to cart: %w", err)
			}
		}
	}
	return nil
}

func (r *cartStore) Delete(ctx context.Context, cartId uuid.UUID) error {
	r.logger.Debug("Enter in repository cart DeleteCart() with args: ctx, cartId: %v", cartId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		defer func() {
			if err != nil {
				r.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					r.logger.Errorf("can't rollback %s", err)
				}

			} else {
				r.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					r.logger.Errorf("can't commit %s", err)
				}
			}
		}()
		_, err = tx.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cartId)
		if err != nil {
			r.logger.Errorf("can't delete cart items from cart: %s", err)
			return fmt.Errorf("can't delete cart items from cart: %w", err)
		}
		_, err = tx.Exec(ctx, `DELETE FROM carts WHERE id=$1`, cartId)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			r.logger.Errorf("can't delete cart: %s", err)
			return models.ErrorNotFound{}
		}
		if err != nil {
			r.logger.Errorf("can't delete cart: %s", err)
			return fmt.Errorf("can't delete cart: %w", err)
		}
		r.logger.Info("Delete cart with id: %v from database success", cartId)
		return nil
	}
}

func (r *cartStore) DeleteItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	r.logger.Debug("Enter in repository cart DeleteItemFromCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		row := pool.QueryRow(ctx, `SELECT item_quantity from cart_items where item_id=$1 and cart_id=$2`, itemId, cartId)
		var quantity int
		err := row.Scan(&quantity)
		if err != nil {
			r.logger.Errorf("error on row.Scan: %s", err)
			return err
		}
		if quantity > 1 {
			_, err := pool.Exec(ctx, `UPDATE cart_items SET item_quantity = item_quantity - 1 WHERE cart_id=$1 and item_id=$2`, cartId, itemId)
			if err != nil {
				r.logger.Errorf("can't delete item from cart: %s", err)
				return fmt.Errorf("can't delete item from cart: %w", err)
			}
		} else if quantity == 1 {
			_, err := pool.Exec(ctx, `DELETE FROM cart_items WHERE item_id=$1 AND cart_id=$2`, itemId, cartId)
			if err != nil {
				r.logger.Errorf("can't delete item from cart: %s", err)
				return fmt.Errorf("can't delete item from cart: %w", err)
			}
		}
		r.logger.Info("Delete item from cart success")
		return nil
	}
}

func (r *cartStore) Get(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	r.logger.Debug("Enter in repository cart GetCart() with args: ctx, cartId: %v", cartId)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		var userId uuid.UUID
		row := pool.QueryRow(ctx, `SELECT user_id FROM carts WHERE id = $1`, cartId)
		err := row.Scan(&userId)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			r.logger.Error(err.Error())
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			r.logger.Error(err)
			return nil, fmt.Errorf("can't read user id: %w", err)
		}
		r.logger.Debug("read user id success: %v", userId)
		item := models.ItemWithQuantity{}
		rows, err := pool.Query(ctx, `
		SELECT 	i.id, i.name, i.description, i.category, cat.name, cat.description, cat.picture, i.price, i.vendor, i.pictures, r.item_quantity
		FROM cart_items c, items i, categories cat
		WHERE r.cart_id=$1 and i.id = r.item_id and cat.id = i.category`, cartId)
		if err != nil {
			r.logger.Errorf("can't select items from cart: %s", err)
			return nil, fmt.Errorf("can't select items from cart: %w", err)
		}
		defer rows.Close()
		r.logger.Debug("read info from db in pool.Query success")
		items := make([]models.ItemWithQuantity, 0, 100)
		for rows.Next() {
			err := rows.Scan(
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
				&item.Quantity,
			)
			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				r.logger.Error(err.Error())
				return nil, models.ErrorNotFound{}
			}
			if err != nil {
				r.logger.Error(err.Error())
				return nil, err
			}

			items = append(items, item)
		}
		r.logger.Info("Select items from cart success")
		r.logger.Info("Get cart success")
		return &models.Cart{
			Id:     cartId,
			UserId: userId,
			Items:  items,
		}, nil
	}
}

func (r *cartStore) CartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error) {
	r.logger.Debug("Enter in repository cart GetCartByUserId() with args: ctx, userId: %v", userId)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		var cartId uuid.UUID
		row := pool.QueryRow(ctx, `SELECT id FROM carts WHERE user_id = $1`, userId)
		err := row.Scan(&cartId)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			r.logger.Error(err.Error())
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			r.logger.Error(err)
			return nil, fmt.Errorf("can't read cart id: %w", err)
		}
		r.logger.Debug("read cart id success: %v", userId)
		item := models.ItemWithQuantity{}
		rows, err := pool.Query(ctx, `
		SELECT i.id, i.name, i.description, i.category, cat.name, cat.description, cat.picture, i.price, i.vendor, i.pictures, r.item_quantity
		FROM cart_items c, items i, categories cat
		WHERE r.cart_id=$1 and i.id = r.item_id and cat.id = i.category`, cartId)
		if err != nil {
			r.logger.Errorf("can't select items from cart: %s", err)
			return nil, fmt.Errorf("can't select items from cart: %w", err)
		}
		defer rows.Close()
		r.logger.Debug("read info from db in pool.Query success")
		items := make([]models.ItemWithQuantity, 0, 100)
		for rows.Next() {
			err := rows.Scan(
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
				&item.Quantity,
			)
			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				r.logger.Error(err.Error())
				return nil, models.ErrorNotFound{}
			}
			if err != nil {
				r.logger.Error(err.Error())
				return nil, err
			}

			items = append(items, item)
		}
		r.logger.Info("Select items from cart success")
		r.logger.Info("Get cart success")
		return &models.Cart{
			Id:     cartId,
			UserId: userId,
			Items:  items,
		}, nil
	}
}

package postgres

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type cart struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

var _ usecase.CartStore = (*cart)(nil)

func NewCartStore(storage *PGres, logger *zap.SugaredLogger) usecase.CartStore {
	return &cart{
		storage: storage,
		logger:  logger,
	}
}

// Create Shall we add items at the moment we create cart
func (c *cart) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	c.logger.Debugf("Enter in repository cart Create() with args: ctx, userId: %v", userId)
	select {
	case <-ctx.Done():
		return uuid.Nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		var cartId uuid.UUID
		row := pool.QueryRow(ctx, `INSERT INTO carts (user_id) VALUES ($1) RETURNING id`,
			userId)
		err := row.Scan(&cartId)
		if err != nil {
			c.logger.Error(err)
			return uuid.Nil, fmt.Errorf("can't create cart object: %w", err)
		}
		c.logger.Info("Create cart success")
		return cartId, nil
	}
}

// AddItemToCart Maybe add to item
func (c *cart) AddItemToCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Debugf("Enter in repository cart AddItemToCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	select {
	case <-ctx.Done():
		c.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT item_id from cart_items where item_id=$1 and cart_id=$2`, itemId, cartId)
		var checkId uuid.UUID
		err := row.Scan(&checkId)
		if err != nil {
			c.logger.Errorf("error on row.Scan: %s", err)
		}
		if checkId == uuid.Nil {
			_, err := pool.Exec(ctx, `INSERT INTO cart_items (cart_id, item_id, item_quantity) VALUES ($1, $2, $3)`, cartId, itemId, 1)
			if err != nil {
				c.logger.Errorf("can't add item to cart: %s", err)
				return fmt.Errorf("can't add item to cart: %w", err)
			}
		} else {
			_, err := pool.Exec(ctx, `UPDATE cart_items SET item_quantity = item_quantity + 1 WHERE cart_id=$1 and item_id=$2`, cartId, itemId)
			if err != nil {
				c.logger.Errorf("can't add item to cart: %s", err)
				return fmt.Errorf("can't add item to cart: %w", err)
			}
		}
	}
	return nil
}

func (c *cart) DeleteCart(ctx context.Context, cartId uuid.UUID) error {
	c.logger.Debug("Enter in repository cart DeleteCart() with args: ctx, cartId: %v", cartId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		defer func() {
			if err != nil {
				c.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					c.logger.Errorf("can't rollback %s", err)
				}

			} else {
				c.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					c.logger.Errorf("can't commit %s", err)
				}
			}
		}()
		_, err = tx.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cartId)
		if err != nil {
			c.logger.Errorf("can't delete cart items from cart: %s", err)
			return fmt.Errorf("can't delete cart items from cart: %w", err)
		}
		_, err = tx.Exec(ctx, `DELETE FROM carts WHERE id=$1`, cartId)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			c.logger.Errorf("can't delete cart: %s", err)
			return models.ErrorNotFound{}
		}
		if err != nil {
			c.logger.Errorf("can't delete cart: %s", err)
			return fmt.Errorf("can't delete cart: %w", err)
		}
		c.logger.Info("Delete cart with id: %v from database success", cartId)
		return nil
	}
}

func (c *cart) DeleteItemFromCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Debug("Enter in repository cart DeleteItemFromCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT item_quantity from cart_items where item_id=$1 and cart_id=$2`, itemId, cartId)
		var quantity int
		err := row.Scan(&quantity)
		if err != nil {
			c.logger.Errorf("error on row.Scan: %s", err)
			return err
		}
		if quantity > 1 {
			_, err := pool.Exec(ctx, `UPDATE cart_items SET item_quantity = item_quantity - 1 WHERE cart_id=$1 and item_id=$2`, cartId, itemId)
			if err != nil {
				c.logger.Errorf("can't delete item from cart: %s", err)
				return fmt.Errorf("can't delete item from cart: %w", err)
			}
		} else if quantity == 1 {
			_, err := pool.Exec(ctx, `DELETE FROM cart_items WHERE item_id=$1 AND cart_id=$2`, itemId, cartId)
			if err != nil {
				c.logger.Errorf("can't delete item from cart: %s", err)
				return fmt.Errorf("can't delete item from cart: %w", err)
			}
		}
		c.logger.Info("Delete item from cart success")
		return nil
	}
}

func (c *cart) GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	c.logger.Debug("Enter in repository cart GetCart() with args: ctx, cartId: %v", cartId)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		var userId uuid.UUID
		row := pool.QueryRow(ctx, `SELECT user_id FROM carts WHERE id = $1`, cartId)
		err := row.Scan(&userId)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			c.logger.Error(err.Error())
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			c.logger.Error(err)
			return nil, fmt.Errorf("can't read user id: %w", err)
		}
		c.logger.Debug("read user id success: %v", userId)
		item := models.ItemWithQuantity{}
		rows, err := pool.Query(ctx, `
		SELECT 	i.id, i.name, i.description, i.category, cat.name, cat.description, cat.picture, i.price, i.vendor, i.pictures, c.item_quantity
		FROM cart_items c, items i, categories cat
		WHERE c.cart_id=$1 and i.id = c.item_id and cat.id = i.category`, cartId)
		if err != nil {
			c.logger.Errorf("can't select items from cart: %s", err)
			return nil, fmt.Errorf("can't select items from cart: %w", err)
		}
		defer rows.Close()
		c.logger.Debug("read info from db in pool.Query success")
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
				c.logger.Error(err.Error())
				return nil, models.ErrorNotFound{}
			}
			if err != nil {
				c.logger.Error(err.Error())
				return nil, err
			}

			items = append(items, item)
		}
		c.logger.Info("Select items from cart success")
		c.logger.Info("Get cart success")
		return &models.Cart{
			Id:     cartId,
			UserId: userId,
			Items:  items,
		}, nil
	}
}

func (c *cart) GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error) {
	c.logger.Debug("Enter in repository cart GetCartByUserId() with args: ctx, userId: %v", userId)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		pool := c.storage.GetPool()
		var cartId uuid.UUID
		row := pool.QueryRow(ctx, `SELECT id FROM carts WHERE user_id = $1`, userId)
		err := row.Scan(&cartId)
		if err != nil && strings.Contains(err.Error(), "no rows in result set") {
			c.logger.Error(err.Error())
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			c.logger.Error(err)
			return nil, fmt.Errorf("can't read cart id: %w", err)
		}
		c.logger.Debug("read cart id success: %v", userId)
		item := models.ItemWithQuantity{}
		rows, err := pool.Query(ctx, `
		SELECT i.id, i.name, i.description, i.category, cat.name, cat.description, cat.picture, i.price, i.vendor, i.pictures, c.item_quantity
		FROM cart_items c, items i, categories cat
		WHERE c.cart_id=$1 and i.id = c.item_id and cat.id = i.category`, cartId)
		if err != nil {
			c.logger.Errorf("can't select items from cart: %s", err)
			return nil, fmt.Errorf("can't select items from cart: %w", err)
		}
		defer rows.Close()
		c.logger.Debug("read info from db in pool.Query success")
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
				c.logger.Error(err.Error())
				return nil, models.ErrorNotFound{}
			}
			if err != nil {
				c.logger.Error(err.Error())
				return nil, err
			}

			items = append(items, item)
		}
		c.logger.Info("Select items from cart success")
		c.logger.Info("Get cart success")
		return &models.Cart{
			Id:     cartId,
			UserId: userId,
			Items:  items,
		}, nil
	}
}

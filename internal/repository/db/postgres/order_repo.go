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

type orderStore struct {
	db     *PGres
	logger *zap.SugaredLogger
}

var _ usecase.OrderStore = (*orderStore)(nil)

func NeworderStore(db *PGres, log *zap.SugaredLogger) *orderStore {
	return &orderStore{
		db:     db,
		logger: log,
	}
}

func (r *orderStore) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	r.logger.Debug("Enter in rsitory order Create with args: ctx, order: %v", order)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("stopped with context")
	default:
		pool := r.db.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			r.logger.Errorf("can't create transaction: %s", err)
			return nil, fmt.Errorf("can't create transaction: %w", err)
		}
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
		row := tx.QueryRow(ctx, `INSERT INTO orders (created_at, shipment_time, user_id, status, address) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, order.CreatedAt, order.ShipmentTime, order.User.Id, order.Status,
			fmt.Sprintf("%s -> %s -> %s -> %s", order.Address.Zipcode, order.Address.Country, order.Address.City, order.Address.Street))
		err = row.Scan(&order.Id)
		if err != nil {
			r.logger.Errorf("can't add new order: %w", err)
			return nil, fmt.Errorf("can't add new order: %w", err)
		}
		query := `INSERT INTO order_items (order_id, item_id, item_quantity) VALUES`
		itemsString := ""
		for _, item := range order.Items {
			itemsString += fmt.Sprintf("('%s', '%s', '%d'),", order.Id.String(), item.Id.String(), item.Quantity)
		}
		itemsString = itemsString[:len(itemsString)-1]
		_, err = tx.Exec(ctx, fmt.Sprintf("%s %s;", query, itemsString))
		if err != nil {
			r.logger.Errorf("can't add items to order: %s", err)
			return nil, fmt.Errorf("can't add items to order: %w", err)
		}
		return order, nil
	}
}

func (r *orderStore) Delete(ctx context.Context, order *models.Order) error {
	r.logger.Debug("Enter in rsitory DeleteOrder() with args: ctx, order: %v", order)
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
		_, err = tx.Exec(ctx, `DELETE FROM order_items WHERE order_id=$1`, order.Id)
		if err != nil {
			r.logger.Errorf("can't delete order items from order: %s", err)
			return fmt.Errorf("can't delete order items from order: %w", err)
		}
		_, err = tx.Exec(ctx, `DELETE FROM orders WHERE id=$1`, order.Id)
		if err != nil {
			r.logger.Errorf("can't delete order: %s", err)
			return fmt.Errorf("can't delete order: %w", err)
		}
		return nil
	}
}
func (r *orderStore) ChangeAddress(ctx context.Context, order *models.Order) error {
	r.logger.Debug("Enter in rsitory order ChangeAddress() with args: ctx, order: %v", order)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		_, err := pool.Exec(ctx, `UPDATE orders SET address=$1 WHERE id=$2`,
			fmt.Sprintf("%s -> %s -> %s -> %s", order.Address.Zipcode, order.Address.Country, order.Address.City, order.Address.Street), order.Id)
		if err != nil {
			r.logger.Errorf("can't update address: %s", err)
			return fmt.Errorf("can't update address: %w", err)
		}
		return nil
	}
}
func (r *orderStore) ChangeStatus(ctx context.Context, order *models.Order) error {
	r.logger.Debug("Enter in rsitory order ChangeStatus() with args: ctx, order: %v", order)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		_, err := pool.Exec(ctx, `UPDATE orders SET status=$1 WHERE id=$2`, order.Status, order.Id)
		if err != nil {
			r.logger.Errorf("can't update status: %s", err)
			return fmt.Errorf("can't update status: %w", err)
		}
		return nil
	}
}
func (r *orderStore) Get(ctx context.Context, id uuid.UUID) (models.Order, error) {
	r.logger.Debug("Enter in rsitory GetOrderByID() with args: ctx, id: %v", id)
	select {
	case <-ctx.Done():
		r.logger.Errorf("context closed")
		return models.Order{}, fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		ordr := models.Order{
			Items: make([]models.ItemWithQuantity, 0),
		}
		rows, err := pool.Query(ctx, `SELECT items.id, items.name, categories.id, categories.name, categories.description, categories.picture,
				items.description, items.price, items.vendor, items.pictures, orders.id, orders.user_id, orders.status, orders.created_at, orders.shipment_time,
				orders.status, orders.address, order_items.item_quantity from items INNER JOIN categories ON categories.id=category  INNER JOIN order_items ON
				items.id=order_items.item_id INNER JOIN orders ON orders.id=order_items.order_id and orders.id = $1 ORDER BY order_id ASC`, id)
		if err != nil {
			r.logger.Errorf("can't get order from db: %s", err)
			return ordr, fmt.Errorf("can't get order from db: %w", err)
		}
		defer rows.Close()
		var address string
		for rows.Next() {
			item := models.ItemWithQuantity{}
			if err := rows.Scan(&item.Id, &item.Title, &item.Category.Id, &item.Category.Name, &item.Category.Description, &item.Category.Image,
				&item.Description, &item.Price, &item.Vendor, &item.Images, &ordr.Id, &ordr.User.Id, &ordr.Status, &ordr.CreatedAt, &ordr.ShipmentTime, &ordr.Status, &address, &item.Quantity); err != nil {
				r.logger.Errorf("can't scan data to order object: %w", err)
				return models.Order{}, err
			}
			ordr.Items = append(ordr.Items, item)
		}
		r.logger.Debug(address)
		splitted := strings.Split(address, " -> ")
		ordr.Address = models.UserAddress{
			Zipcode: splitted[0],
			Country: splitted[1],
			City:    splitted[2],
			Street:  splitted[3],
		}
		return ordr, nil
	}

}

func (r *orderStore) List(ctx context.Context, user *models.User) (chan models.Order, error) {
	r.logger.Debug("Enter in rsitory GetOrdersForUser() with args: ctx, user: %v", user)
	select {
	case <-ctx.Done():
		r.logger.Errorf("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		pool := r.db.GetPool()
		resChan := make(chan models.Order, 1)
		go func() {
			defer close(resChan)
			rows, err := pool.Query(ctx, `SELECT items.id, items.name, categories.id, categories.name, categories.description, categories.picture,
			items.description, items.price, items.vendor, items.pictures, orders.id, orders.user_id, orders.status, orders.created_at, orders.shipment_time,
			orders.status, orders.address, order_items.item_quantity from items INNER JOIN categories ON categories.id=category  INNER JOIN order_items ON
			items.id=order_items.item_id INNER JOIN orders ON orders.id=order_items.order_id and orders.user_id = $1 ORDER BY order_id ASC`, user.Id)
			if err != nil {
				r.logger.Errorf("can't get order from db: %s", err)
				return
			}
			defer rows.Close()
			prevOrder := models.Order{
				Items: make([]models.ItemWithQuantity, 0),
			}
			for rows.Next() {
				var address string
				item := models.ItemWithQuantity{}
				order := models.Order{}
				if err := rows.Scan(&item.Id, &item.Title, &item.Category.Id, &item.Category.Name, &item.Category.Description, &item.Category.Image,
					&item.Description, &item.Price, &item.Vendor, &item.Images, &order.Id, &order.User.Id, &order.Status, &order.CreatedAt, &order.ShipmentTime, &order.Status, &address, &item.Quantity); err != nil {
					r.logger.Errorf("can't scan data to order object: %w", err)
					return
				}
				if prevOrder.Id == uuid.Nil {
					prevOrder = order
				}
				if order.Id != prevOrder.Id {
					resChan <- prevOrder
					prevOrder = order
				}
				r.logger.Debug(address)
				splitted := strings.Split(address, " -> ")
				prevOrder.Address = models.UserAddress{
					Zipcode: splitted[0],
					Country: splitted[1],
					City:    splitted[2],
					Street:  splitted[3],
				}
				prevOrder.Items = append(prevOrder.Items, item)

			}
			resChan <- prevOrder
		}()
		return resChan, nil
	}
}

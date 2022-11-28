package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (repo *Pgrepo) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	repo.logger.Debug("Enter in repository CreateItem()")
	var id uuid.UUID
	_ = repo.db.QueryRowContext(ctx, `INSERT INTO items(name, category, description, price, vendor)
	values ($1, $2, $3, $4, $5) RETURNING id`,
		item.Title,
		item.Category,
		item.Description,
		item.Price,
		item.Vendor,
	).Scan(&id)

	repo.logger.Debug(fmt.Sprintf("id is %v\n", id))
	return id, nil
}

func (repo *Pgrepo) UpdateItem(ctx context.Context, item *models.Item) error {
	repo.logger.Debug("Enter in repository UpdateItem()")
	_, err := repo.db.ExecContext(ctx, `UPDATE items SET name=$1, category=$2, description=$3, price=$4, vendor=$5 WHERE id=$6`,
		item.Title,
		item.Category,
		item.Description,
		item.Price,
		item.Vendor,
		item.Id)
	if err != nil {
		return fmt.Errorf("error on update item: %w)", err)
	}
	return nil
}

func (repo *Pgrepo) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	repo.logger.Debug("Enter in repository GetItem()")
	item := models.Item{}
	rows, err := repo.db.QueryContext(ctx,
		`SELECT id, name, category, description, price, vendor FROM items WHERE id = $1`, id)
	if err != nil {
		return &models.Item{}, fmt.Errorf("error on get item: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Category,
			&item.Description,
			&item.Price,
			&item.Vendor,
		); err != nil {
			return &models.Item{}, fmt.Errorf("error in rows scan get item by id: %w", err)
		}
	}
	if uuid.UUID.String(item.Id) == "00000000-0000-0000-0000-000000000000" {
		err = fmt.Errorf("id not found")
		return &models.Item{}, err
	}

	return &item, nil
}

func (repo *Pgrepo) ItemsList(ctx context.Context) (chan models.Item, error) {
	repo.logger.Debug("Enter in repository ItemsList()")
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}

		rows, err := repo.db.QueryContext(ctx, `
		SELECT id, name, category, description, price, vendor FROM items`)
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
				&item.Category,
				&item.Description,
				&item.Price,
				&item.Vendor,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			itemChan <- *item
		}
	}()

	return itemChan, nil
}

func (repo *Pgrepo) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	repo.logger.Debug("Enter in repository SearchLine()")
	itemChan := make(chan models.Item, 100)
	go func() {
		defer close(itemChan)
		item := &models.Item{}

		rows, err := repo.db.QueryContext(ctx, `
		SELECT id, name, category, description, price, vendor FROM items WHERE name LIKE $1 OR description LIKE $1 OR vendor LIKE $1`,
			"%"+param+"%")
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
				&item.Category,
				&item.Description,
				&item.Price,
				&item.Vendor,
			); err != nil {
				repo.logger.Error(err.Error())
				return
			}
			fmt.Println(item)
			itemChan <- *item
		}
	}()

	return itemChan, nil
}

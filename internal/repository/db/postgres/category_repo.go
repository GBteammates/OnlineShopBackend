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

type categoryStore struct {
	db     *PGres
	logger *zap.SugaredLogger
}

var _ usecase.CategoryStore = (*categoryStore)(nil)

func NewCategoryStore(db *PGres, log *zap.SugaredLogger) *categoryStore {
	return &categoryStore{
		db:     db,
		logger: log,
	}
}

// CreateCategory create new category in database
func (r *categoryStore) Create(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	r.logger.Debugf("Enter in repository CreateCategory() with args: ctx, category: %v", category)

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
	// If name of created category = name of deleted category, update deleted category
	// and set deleted_at = null and return id of deleted category
	if id, ok := r.isDeleted(ctx, category.Name); ok {
		r.logger.Debug("Category with name: %s is deleted", category.Name)
		_, err := pool.Exec(ctx, `UPDATE categories SET description=$1, picture=$2, deleted_at=null WHERE name=$3`,
			category.Description,
			category.Image,
			category.Name)
		if err != nil {
			r.logger.Debug(err.Error())
			return uuid.Nil, err
		}
		r.logger.Debug("Category recreated from deleted category success")
		r.logger.Debugf("Category id is %v\n", id)
		return id, nil
	}
	var id uuid.UUID
	row := tx.QueryRow(ctx, `INSERT INTO categories(name, description, picture, deleted_at)
	values ($1, $2, $3, $4) RETURNING id`,
		category.Name,
		category.Description,
		category.Image,
		nil,
	)
	if err := row.Scan(&id); err != nil {
		r.logger.Errorf("Can't scan %s", err)
		return uuid.Nil, fmt.Errorf("can't scan %w", err)
	}
	r.logger.Debug("Category created success")
	r.logger.Debugf("Category id is %v\n", id)
	return id, nil
}

// isDeletedCategory check created category name and if it is a deleted category name, returns
// uid of deleted category and true
func (r *categoryStore) isDeleted(ctx context.Context, name string) (uuid.UUID, bool) {
	r.logger.Debug("Enter in repository is Deleted() with args: ctx, name: %s", name)
	pool := r.db.GetPool()
	category := models.Category{}
	row := pool.QueryRow(ctx,
		`SELECT id FROM categories WHERE deleted_at is not null AND name = $1`, name)
	err := row.Scan(
		&category.Id,
	)
	if err == nil && category.Id != uuid.Nil {
		return category.Id, true
	}
	r.logger.Error(err.Error())
	return uuid.Nil, false
}

// UpdateCategory сhanges the existing category
func (r *categoryStore) Update(ctx context.Context, category *models.Category) error {
	r.logger.Debugf("Enter in repository UpdateCategory() with args: ctx, category: %v", category)
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

	_, err = tx.Exec(ctx, `UPDATE categories SET name=$1, description=$2, picture=$3 WHERE id=$4`,
		category.Name,
		category.Description,
		category.Image,
		category.Id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error on update category %s: %s", category.Id, err)
		return models.ErrorNotFound{}
	}
	if err != nil {
		r.logger.Errorf("Error on update category %s: %s", category.Id, err)
		return fmt.Errorf("error on update category %s: %w", category.Id, err)
	}
	r.logger.Infof("Category with id %s successfully updated", category.Id)
	return nil
}

// GetCategory returns *models.Category by id or error
func (r *categoryStore) Get(ctx context.Context, param string) (*models.Category, error) {
	r.logger.Debugf("Enter in repository GetCategory() with args: ctx, param: %v", param)

	pool := r.db.GetPool()

	category := models.Category{}

	row := pool.QueryRow(ctx,
		`SELECT id, name, description, picture FROM categories WHERE deleted_at is null AND id = $1`, param)
	err := row.Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.Image,
	)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error in rows scan get category by id: %s", err)
		return nil, models.ErrorNotFound{}
	}
	if err != nil {
		r.logger.Errorf("Error in rows scan get category by id: %s", err)
		return nil, fmt.Errorf("error in rows scan get category by id: %w", err)
	}
	r.logger.Info("Get category success")
	return &category, nil
}

// CategoriesList reads all the categories from the database and writes it to the
// output channel and returns this channel or error
func (r *categoryStore) List(ctx context.Context) (chan models.Category, error) {
	r.logger.Debug("Enter in repository CategoriesList() with args: ctx")
	categoryChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryChan)
		category := &models.Category{}

		pool := r.db.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, description, picture FROM categories WHERE deleted_at is null`)
		if err != nil {
			r.logger.Error(fmt.Errorf("error on categories list query context: %w", err).Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&category.Id,
				&category.Name,
				&category.Description,
				&category.Image,
			); err != nil {
				r.logger.Error(err.Error())
				return
			}
			categoryChan <- *category
		}
	}()

	return categoryChan, nil
}

// Delete changes the value of the deleted_at attribute in the deleted category for the current time
func (r *categoryStore) Delete(ctx context.Context, id uuid.UUID) error {
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
				r.logger.Errorf("Can't rollback %s", err)
			}

		} else {
			r.logger.Info("Transaction commited")
			if err != tx.Commit(ctx) {
				r.logger.Errorf("Can't commit %s", err)
			}
		}
	}()
	_, err = tx.Exec(ctx, `UPDATE categories SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		r.logger.Errorf("Error on delete category %s: %s", id, err)
		return models.ErrorNotFound{}
	}
	if err != nil {
		r.logger.Errorf("Error on delete category %s: %s", id, err)
		return fmt.Errorf("error on delete category %s: %w", id, err)
	}
	r.logger.Infof("Category with id: %s successfully deleted from database", id)
	return nil
}

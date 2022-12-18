package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoryRepo struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

type Category struct {
	Id          uuid.UUID
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

var _ CategoryStore = (*categoryRepo)(nil)

func NewCategoryRepo(store *PGres, log *zap.SugaredLogger) CategoryStore {
	return &categoryRepo{
		storage: store,
		logger:  log,
	}
}

func (repo *categoryRepo) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	repo.logger.Debug("Enter in repository CreateCategory()")
	repoCategory := &Category{
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	var id uuid.UUID
	pool := repo.storage.GetPool()
	row := pool.QueryRow(ctx, `INSERT INTO categories(name, description, picture, created_at, updated_at, deleted_at)
	values ($1, $2, $3, $4, $5, $6) RETURNING id`,
		repoCategory.Name,
		repoCategory.Description,
		repoCategory.Image,
		repoCategory.CreatedAt,
		repoCategory.UpdatedAt,
		nil,
	)
	if err := row.Scan(&id); err != nil {
		repo.logger.Errorf("can't scan %s", err)
		return uuid.Nil, fmt.Errorf("can't scan %w", err)
	}

	repo.logger.Debugf("id is %v\n", id)
	return id, nil
}

func (repo *categoryRepo) UpdateCategory(ctx context.Context, category *models.Category) error {
	repo.logger.Debug("Enter in repository UpdateCategory()")
	pool := repo.storage.GetPool()
	_, err := pool.Exec(ctx, `UPDATE categories SET name=$1, description=$2, picture=$3, updated_at=$4 WHERE id=$4`,
		category.Name,
		category.Description,
		category.Image,
		category.Id,
		time.Now())
	if err != nil {
		repo.logger.Errorf("error on update category %s: %s", category.Id, err)
		return fmt.Errorf("error on update category %s: %w", category.Id, err)
	}
	repo.logger.Infof("category %s successfully updated", category.Id)
	return nil
}

func (repo *categoryRepo) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategory()")
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context closed")
	default:
		category := models.Category{}
		pool := repo.storage.GetPool()
		row := pool.QueryRow(ctx,
			`SELECT id, name, description, picture FROM categories WHERE id = $1`, id)
		err := row.Scan(
			&category.Id,
			&category.Name,
			&category.Description,
			&category.Image,
		)
		if err != nil {
			repo.logger.Errorf("error in rows scan get category by id: %s", err)
			return &models.Category{}, fmt.Errorf("error in rows scan get category by id: %w", err)
		}
		return &category, nil
	}
}

func (repo *categoryRepo) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	repo.logger.Debug("Enter in repository GetCategoryList()")
	categoryChan := make(chan models.Category, 100)
	go func() {
		defer close(categoryChan)
		category := &models.Category{}

		pool := repo.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, description, picture FROM categories WHERE deleted_at is null`)
		if err != nil {
			msg := fmt.Errorf("error on categories list query context: %w", err)
			repo.logger.Error(msg.Error())
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
				repo.logger.Error(err.Error())
				return
			}
			categoryChan <- *category
		}
	}()

	return categoryChan, nil
}

func (repo *categoryRepo) DeleteCategory(ctx context.Context, id uuid.UUID) (deletedCategoryName string, err error) {
	repo.logger.Debug("Enter in repository DeleteCategory()")
	pool := repo.storage.GetPool()
	repoCategory, err := repo.GetCategory(ctx, id)
	if err != nil {
		repo.logger.Errorf("error on get category: %v", err)
		return "", err
	}
	_, err = pool.Exec(ctx, `UPDATE categories SET deleted_at=$1 WHERE id=$2`,
		time.Now(), id)
	if err != nil {
		repo.logger.Errorf("error on delete category %s: %s", id, err)
		return "", fmt.Errorf("error on delete category %s: %w", id, err)
	}
	repo.logger.Infof("category %s successfully deleted", id)
	return repoCategory.Name, nil
}

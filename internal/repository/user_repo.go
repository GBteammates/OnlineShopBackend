package repository

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type user struct {
	storage *PGres
	logger  *zap.SugaredLogger
}

var _ UserStore = (*user)(nil)

func NewUser(storage *PGres, logger *zap.SugaredLogger) UserStore {
	return &user{
		storage: storage,
		logger:  logger,
	}
}

func (u *user) Create(ctx context.Context, user *models.User) (*models.User, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		// we create rights and address somewhere in usecase or get them from user
		row := pool.QueryRow(ctx, `INSERT INTO users 
		(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.ID,
			user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)
		var id uuid.UUID
		err := row.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("can't create user: %w", err)
		}
		user.ID = id
		return user, nil
	}
}

func (u *user) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u.logger.Debug("Enter in repository GetUserByEmail()")
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT users.id, users.name, lastname, password, email, rights.id, zipcode, country, city, street,
		rights.name, rights.rules FROM users INNER JOIN rights ON email=$1 and rights.id=users.rights`, email)
		var user = models.User{}
		err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.Rights.ID,
			&user.Address.Zipcode, &user.Address.Country, &user.Address.City, &user.Address.Street, &user.Rights.Name, &user.Rights.Rules)
		if err != nil {
			return &models.User{}, fmt.Errorf("can't get user from database: %w", err)
		}
		return &user, nil
	}
}

func (u *user) UpdateUserData(ctx context.Context, id uuid.UUID, user *models.User) (*models.User, error) {
	u.logger.Debug("Enter in repository UpdateUserData()")
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			u.logger.Errorf("can't create transaction: %s", err)
			return &models.User{}, fmt.Errorf("can't create transaction: %w", err)
		}
		u.logger.Debug("transaction begin success")
		defer func() {
			if err != nil {
				u.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					u.logger.Errorf("can't rollback %s", err)
				}

			} else {
				u.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					u.logger.Errorf("can't commit %s", err)
				}
			}
		}()

		_, err = tx.Exec(ctx, `UPDATE users SET name=$1, lastname=$2, country=$3, city=$4, street=$5, zipcode=$6 WHERE id=$7`,
			user.Firstname,
			user.Lastname,
			user.Address.Country,
			user.Address.City,
			user.Address.Street,
			user.Address.Zipcode,
			id)
		if err != nil {
			u.logger.Errorf("error on update user %s: %s", user.ID, err)
			return &models.User{}, fmt.Errorf("error on update item %s: %w", user.ID, err)
		}
		u.logger.Infof("item %s successfully updated %s", user.ID, user.Lastname)
		return user, nil
	}
}

func (u *user) UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error {
	u.logger.Debug("Enter in repository UpdateUserRole()")
	select {
	case <-ctx.Done():
		return fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			u.logger.Errorf("can't create transaction: %s", err)
			return fmt.Errorf("can't create transaction: %w", err)
		}
		u.logger.Debug("transaction begin success")
		defer func() {
			if err != nil {
				u.logger.Errorf("transaction rolled back")
				if err = tx.Rollback(ctx); err != nil {
					u.logger.Errorf("can't rollback %s", err)
				}

			} else {
				u.logger.Info("transaction commited")
				if err != tx.Commit(ctx) {
					u.logger.Errorf("can't commit %s", err)
				}
			}
		}()

		_, err = tx.Exec(ctx, `UPDATE users SET rights=$1 WHERE email=$2`,
			roleId,
			email)
		if err != nil {
			u.logger.Errorf("error on update user %s: %s", email, err)
			return fmt.Errorf("error on update user %s: %w", email, err)
		}
		u.logger.Infof("user role was successfully updated %s", email)
		return nil
	}
}

func (u *user) GetRightsId(ctx context.Context, name string) (models.Rights, error) {
	select {
	case <-ctx.Done():
		return models.Rights{}, fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		row := pool.QueryRow(ctx, `SELECT id, name, rules FROM rights WHERE name=$1`, name)
		var rights = models.Rights{}
		err := row.Scan(&rights.ID, &rights.Name, &rights.Rules)
		if err != nil {
			return models.Rights{}, fmt.Errorf("can't get rights from database: %w", err)
		}
		return rights, nil

	}
}

func (u *user) SaveSession(ctx context.Context, token string, t int64) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context is closed")
	default:
		pool := u.storage.GetPool()
		pool.QueryRow(ctx, `INSERT INTO session (token, timestamp) VALUES ($1, $2)`,
			token, t)
	}
	return nil
}

func (u *user) GetRightsList(ctx context.Context) (chan models.Rights, error) {
	u.logger.Debug("Enter in repository GetCategoryList() with args: ctx")
	rolesChan := make(chan models.Rights, 100)
	go func() {
		defer close(rolesChan)
		rights := &models.Rights{}

		pool := u.storage.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, rules FROM rights`) // WHERE deleted_at is null
		if err != nil {
			u.logger.Error(fmt.Errorf("error on rights list query context: %w", err).Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&rights.ID,
				&rights.Name,
				&rights.Rules,
			); err != nil {
				u.logger.Error(err.Error())
				return
			}
			rolesChan <- *rights
		}
	}()

	return rolesChan, nil
}
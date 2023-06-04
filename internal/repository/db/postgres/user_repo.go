package postgres

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type userStore struct {
	db     *PGres
	logger *zap.SugaredLogger
}

var _ usecase.UserStore = (*userStore)(nil)

func NewUser(db *PGres, logger *zap.SugaredLogger) *userStore {
	return &userStore{
		db:     db,
		logger: logger,
	}
}

func (r *userStore) Create(ctx context.Context, user *models.User) (uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return uuid.Nil, fmt.Errorf("context is closed")
	default:
		pool := r.db.GetPool()
		// we create rights and address somewhere in usecase or get them from user
		row := pool.QueryRow(ctx, `INSERT INTO users 
		(name, lastname, password, email, rights, zipcode, country, city, street) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
			user.Firstname, user.Lastname, user.Password, user.Email, user.Rights.Id,
			user.Address.Zipcode, user.Address.Country, user.Address.City, user.Address.Street)
		var id uuid.UUID
		err := row.Scan(&id)
		if err != nil {
			return uuid.Nil, fmt.Errorf("can't create user: %w", err)
		}
		return id, nil
	}
}

func (r *userStore) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.logger.Debug("Enter in rsitory GetUserByEmail()")
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := r.db.GetPool()
		row := pool.QueryRow(ctx, `SELECT users.id, users.name, lastname, password, email, rights.id, zipcode, country, city, street,
		rights.name, rights.rules FROM users INNER JOIN rights ON email=$1 and rights.id=users.rights`, email)
		var user = models.User{}
		err := row.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.Rights.Id,
			&user.Address.Zipcode, &user.Address.Country, &user.Address.City, &user.Address.Street, &user.Rights.Name, &user.Rights.Rules)
		if err != nil {
			return &models.User{}, fmt.Errorf("can't get user from database: %w", err)
		}
		return &user, nil
	}
}

func (r *userStore) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	r.logger.Debug("Enter in rsitory UpdateUserData()")
	select {
	case <-ctx.Done():
		return &models.User{}, fmt.Errorf("context is closed")
	default:
		pool := r.db.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			r.logger.Errorf("can't create transaction: %s", err)
			return &models.User{}, fmt.Errorf("can't create transaction: %w", err)
		}
		r.logger.Debug("transaction begin success")
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

		_, err = tx.Exec(ctx, `UPDATE users SET name=$1, lastname=$2, country=$3, city=$4, street=$5, zipcode=$6 WHERE id=$7`,
			user.Firstname,
			user.Lastname,
			user.Address.Country,
			user.Address.City,
			user.Address.Street,
			user.Address.Zipcode,
			user.Id)
		if err != nil {
			r.logger.Errorf("error on update user %s: %s", user.Id, err)
			return &models.User{}, fmt.Errorf("error on update item %s: %w", user.Id, err)
		}
		r.logger.Infof("item %s successfully updated %s", user.Id, user.Lastname)
		return user, nil
	}
}

func (r *userStore) UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error {
	r.logger.Debug("Enter in rsitory UpdateUserRole()")
	select {
	case <-ctx.Done():
		return fmt.Errorf("context is closed")
	default:
		pool := r.db.GetPool()
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			r.logger.Errorf("can't create transaction: %s", err)
			return fmt.Errorf("can't create transaction: %w", err)
		}
		r.logger.Debug("transaction begin success")
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

		_, err = tx.Exec(ctx, `UPDATE users SET rights=$1 WHERE email=$2`,
			roleId,
			email)
		if err != nil {
			r.logger.Errorf("error on update user %s: %s", email, err)
			return fmt.Errorf("error on update user %s: %w", email, err)
		}
		r.logger.Infof("user role was successfully updated %s", email)
		return nil
	}
}

func (r *userStore) GetRightsId(ctx context.Context, name string) (models.Rights, error) {
	select {
	case <-ctx.Done():
		return models.Rights{}, fmt.Errorf("context is closed")
	default:
		pool := r.db.GetPool()
		row := pool.QueryRow(ctx, `SELECT id, name, rules FROM rights WHERE name=$1`, name)
		var rights = models.Rights{}
		err := row.Scan(&rights.Id, &rights.Name, &rights.Rules)
		if err != nil {
			return models.Rights{}, fmt.Errorf("can't get rights from database: %w", err)
		}
		return rights, nil

	}
}

func (r *userStore) ListRights(ctx context.Context) (chan models.Rights, error) {
	r.logger.Debug("Enter in rsitory GetCategoryList() with args: ctx")
	rolesChan := make(chan models.Rights, 100)
	go func() {
		defer close(rolesChan)
		rights := &models.Rights{}

		pool := r.db.GetPool()
		rows, err := pool.Query(ctx, `
		SELECT id, name, rules FROM rights`) // WHERE deleted_at is null
		if err != nil {
			r.logger.Error(fmt.Errorf("error on rights list query context: %w", err).Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&rights.Id,
				&rights.Name,
				&rights.Rules,
			); err != nil {
				r.logger.Error(err.Error())
				return
			}
			rolesChan <- *rights
		}
	}()

	return rolesChan, nil
}

func (r *userStore) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	r.logger.Debugf("Enter in rsitory CreateRights() with args: ctx, rights: %v", rights)

	var id uuid.UUID
	pool := r.db.GetPool()

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
	row := tx.QueryRow(ctx, `INSERT INTO rights(name, rules) values ($1, $2) RETURNING id`,
		rights.Name,
		rights.Rules,
	)
	err = row.Scan(&id)
	if err != nil {
		r.logger.Errorf("can't create rights %s", err)
		return uuid.Nil, fmt.Errorf("can't create rights %w", err)
	}
	r.logger.Info("Rights create success")
	r.logger.Debugf("id is %v\n", id)
	return id, nil
}

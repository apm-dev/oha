package usersql

import (
	"context"
	"database/sql"

	"github.com/apm-dev/oha/src/domain"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) domain.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM users WHERE id=$1", id).
		Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Insert(ctx context.Context, user *domain.User) error {
	const op = "user.repo.sql.Insert"
	res, err := r.db.ExecContext(ctx, "INSERT INTO users (id,name) VALUES ($1,$2);", user.ID, user.Name)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.Wrap(domain.ErrDatabase, "failed to insert user into DB")
	}
	if rows > 1 {
		log.Warnf("%s expected to affect 1 row, affected %d", op, rows)
	}
	return nil
}

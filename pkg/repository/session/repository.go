package session

import (
	"context"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/jmoiron/sqlx"
	"time"
)

const timeOutQuery = 5 * time.Second

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Retrieve(u *user.Resource) (int64, error) {
	var id int64
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password_hash = $2 LIMIT 1", user.TableName)
	ctx, cancel := context.WithTimeout(context.Background(), timeOutQuery)
	defer cancel()

	err := r.db.QueryRowxContext(ctx, query, u.Email, u.Password).StructScan(&id)
	return id, err
}

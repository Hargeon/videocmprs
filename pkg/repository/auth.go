package repository

import (
	"context"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (auth AuthRepository) CreateUser(user *model.User) (int64, error) {
	var id int64
	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id", model.UserTableName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := auth.db.QueryRowxContext(ctx, query, user.Email, user.Password, user.CreatedAt)
	err := row.Scan(&id)
	return id, err
}

func (auth AuthRepository) GetUser(email, password string) (int64, error) {
	var id int64
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password_hash = $2 LIMIT 1", model.UserTableName)
	err := auth.db.QueryRowx(query, email, password).Scan(&id)
	return id, err
}

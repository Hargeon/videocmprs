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

	fmt.Println(*user)

	row := auth.db.QueryRowxContext(ctx, query, user.Email, user.Password, user.CreatedAt)
	err := row.Scan(&id)
	return id, err
}

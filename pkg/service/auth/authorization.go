// Package auth uses for authorization and authentication users
package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	"github.com/Hargeon/videocmprs/pkg/service/jwt"
	"github.com/google/jsonapi"
)

// Service ...
type Service struct {
	repo repository.Existable
}

// NewService ...
func NewService(repo repository.Existable) *Service {
	return &Service{repo: repo}
}

// GenerateToken jwt for user
func (srv *Service) GenerateToken(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	usr, ok := resource.(*user.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion in auth service")
	}

	hashPass := encryption.GenerateHash([]byte(usr.Password))
	id, err := srv.repo.Exists(ctx, usr.Email, fmt.Sprintf("%x", hashPass))
	if err != nil {
		return nil, err
	}

	token, err := jwt.SignedString(id)
	if err != nil {
		return nil, err
	}

	res := &user.Resource{
		ID:    id,
		Email: usr.Email,
		Token: token,
	}
	return res, nil
}

// Retrieve return user params
func (srv *Service) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	res, err := srv.repo.Retrieve(ctx, id)
	return res, err
}

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

// AuthorizationService ...
type AuthorizationService struct {
	repo repository.Existable
}

// NewService initialize AuthorizationService
func NewService(repo repository.Existable) *AuthorizationService {
	return &AuthorizationService{repo: repo}
}

// GenerateToken jwt for user
func (srv *AuthorizationService) GenerateToken(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
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

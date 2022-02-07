package user

import (
	"context"
	"fmt"

	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"

	"github.com/google/jsonapi"
)

type Service struct {
	repo repository.UserRepository
}

// NewService initialize Service
func NewService(repo repository.UserRepository) *Service {
	return &Service{repo: repo}
}

// Create function is hashing password and use repository to create user
func (srv *Service) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	usr, ok := resource.(*user.Resource)

	if !ok {
		return nil, service.ErrInvalidTypeAssertion
	}

	ok, err := srv.repo.Unique(ctx, usr.Email)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, service.ErrAlreadyExists
	}

	hashPass := encryption.GenerateHash([]byte(usr.Password))
	usr.Password = fmt.Sprintf("%x", hashPass)

	return srv.repo.Create(ctx, usr)
}

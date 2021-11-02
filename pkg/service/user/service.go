package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	"github.com/google/jsonapi"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// Create function is hashing password and use repository to create user
func (srv *Service) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	usr, ok := resource.(*user.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion in service")
	}
	hashPass := encryption.GenerateHash([]byte(usr.Password))
	usr.Password = fmt.Sprintf("%x", hashPass)
	return srv.repo.Create(ctx, usr)
}

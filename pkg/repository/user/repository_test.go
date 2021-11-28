package user

import (
	"github.com/google/jsonapi"
)

type notUser struct{}

func (n *notUser) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add",
	}
}

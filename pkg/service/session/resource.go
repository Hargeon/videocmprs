package session

import "github.com/google/jsonapi"

type Resource struct {
	Id    int64  `jsonapi:"primary,sessions"`
	Token string `jsonapi:"attr,token"`
}

func (r *Resource) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"details": "sessions meta information",
	}
}

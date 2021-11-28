package request

import (
	"github.com/google/jsonapi"
)

type invalidResource struct{}

func (r *invalidResource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "",
	}
}

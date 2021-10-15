package auth

import "github.com/google/jsonapi"

// Resource ...
type Resource struct {
	ID string `jsonapi:"primary,token"`
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}

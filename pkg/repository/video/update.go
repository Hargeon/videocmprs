package video

import (
	"context"

	"github.com/google/jsonapi"
)

// Update video TODO implement in future
func (r *Repository) Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error) {
	return nil, nil
}

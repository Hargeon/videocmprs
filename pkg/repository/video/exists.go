package video

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) RelationExists(ctx context.Context, userID, relationID int64) (int64, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64

	err := sq.Select("id").
		From(TableName).
		Where(sq.And{sq.Eq{"id": relationID}, sq.Eq{"user_id": userID}}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(c).
		Scan(&id)

	return id, err
}

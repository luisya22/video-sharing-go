package videos

import (
	"context"
	"database/sql"
)

type store interface {
	Insert(ctx context.Context, v *Video) error
	ReadById(ctx context.Context, videoId int64) error
}

type videoStore struct {
	db *sql.DB
}

func (v *videoStore) Insert(ctx context.Context, video *Video) error {
	query := `INSERT INTO videos (path) 
			VALUES ($1)
			RETURNING id, created_at, version`

	args := []any{video.Path}

	return v.db.QueryRowContext(ctx, query, args...).Scan(&video.ID, &video.CreatedAt, &video.Version)
}

func (v *videoStore) ReadById(ctx context.Context, videoId int64) error {

	return nil
}

// Initialize Store
func newStore(db *sql.DB) (*videoStore, error) {
	return &videoStore{
		db: db,
	}, nil
}

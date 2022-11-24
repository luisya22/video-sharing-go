package videos

import (
	"context"
	"database/sql"
	"errors"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/datastore"
	"time"
)

type store interface {
	Insert(ctx context.Context, v *Video) error
	Update(ctx context.Context, v *Video) error
	ReadById(ctx context.Context, videoId int64) (*Video, error)
}

type videoStore struct {
	db *sql.DB
}

func (v *videoStore) Insert(ctx context.Context, video *Video) error {
	query := `INSERT INTO videos (status) 
			VALUES ($1)
			RETURNING id, status, created_at, version`

	args := []any{"Uploading"}

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return v.db.QueryRowContext(dbCtx, query, args...).Scan(&video.ID, &video.Status, &video.CreatedAt, &video.Version)
}

func (v *videoStore) Update(ctx context.Context, video *Video) error {
	query := `UPDATE videos SET title = $1, description = $2, video_path = $3, thumbnail_path = $4, status = $5, 
                  published_at = $6, version = version + 1, updated_at = now()
              	  WHERE id = $7 AND version = $8
                  RETURNING version`

	args := []any{
		video.Title,
		video.Description,
		video.Path,
		video.ImgPath,
		video.Status,
		video.PublishedDate,
		video.ID,
		video.Version,
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := v.db.QueryRowContext(dbCtx, query, args...).Scan(&video.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return datastore.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (v *videoStore) ReadById(ctx context.Context, videoId int64) (*Video, error) {

	query := `SELECT id, title, description, video_path, thumbnail_path, status, published_at, version FROM videos 
			  WHERE id = $1`

	var video Video

	dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := v.db.QueryRowContext(dbCtx, query, videoId).Scan(
		&video.ID,
		&video.Description,
		&video.Title,
		&video.Path,
		&video.ImgPath,
		&video.Status,
		&video.PublishedDate,
		&video.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, datastore.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &video, nil
}

// Initialize Store
func newStore(db *sql.DB) (*videoStore, error) {
	return &videoStore{
		db: db,
	}, nil
}

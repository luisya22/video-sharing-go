package videos

import (
	"context"
	"database/sql"
	"errors"
	"luismatosgarcia.dev/video-sharing-go/internal/validator"
	"mime/multipart"
	"time"
)

var (
	VideoValidationError = errors.New("video data is not valid")
)

type Video struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title,omitempty"`
	Description   string    `json:"description,omitempty"`
	Path          string    `json:"path,omitempty"`
	ImgPath       string    `json:"img_path,omitempty"`
	PublishedDate time.Time `json:"published_date:omitempty"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	Version       int32     `json:"version"`
}

type Videos struct {
	store store
}

func ValidateVideo(v *validator.Validator, video *Video) {
	v.Check(video.Title != "", "title", "must be provided")
	v.Check(len(video.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(video.Description != "", "description", "must be provided")

	v.Check(video.PublishedDate.IsZero() || video.PublishedDate.After(time.Now()), "published_date", "must be in the future")
}

func (vs *Videos) UploadVideo(ctx context.Context, videoFile *multipart.File, fileHeader *multipart.FileHeader) (*Video, error, map[string]string) {
	// Upload to S3 Bucket

	//Save to database
	// Return Video
	return nil, nil, nil
}

func (vs *Videos) CreateVideo(ctx context.Context, video *Video) (*Video, error, map[string]string) {

	v := validator.New()

	if ValidateVideo(v, video); !v.Valid() {
		return nil, VideoValidationError, v.Errors
	}

	err := vs.store.Insert(ctx, video)
	if err != nil {
		return nil, err, nil
	}

	return video, nil, nil
}

// Initialize Video Service
func NewService(db *sql.DB) (*Videos, error) {
	vs, err := newStore(db)
	if err != nil {
		return nil, err
	}

	return &Videos{
		store: vs,
	}, nil
}

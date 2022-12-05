package videos

import (
	"context"
	"database/sql"
	"errors"
	"luismatosgarcia.dev/video-sharing-go/internal/background"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/filestore"
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
	Status        string    `json:"status,omitempty"`
	PublishedDate time.Time `json:"published_date,omitempty"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	Version       int32     `json:"version"`
}

type VideoInput struct {
	Title         *string    `json:"title"`
	Description   *string    `json:"description"`
	PublishedDate *time.Time `json:"published_date"`
}

type Videos struct {
	store      store
	filestore  filestore.FileStore
	background *background.Routine
}

func ValidateVideo(v *validator.Validator, video *Video) {
	v.Check(video.Title != "", "title", "must be provided")
	v.Check(len(video.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(video.Description != "", "description", "must be provided")

	v.Check(video.PublishedDate.IsZero() || video.PublishedDate.After(time.Now()), "published_date", "must be in the future")
}

func (vs *Videos) UploadVideo(ctx context.Context, videoFile *multipart.File, fileHeader *multipart.FileHeader) (*Video, error, map[string]string) {
	// Upload to S3 Bucket

	//TODO: Save to database return ID create background job with ID then update
	video := &Video{}

	err := vs.store.Insert(ctx, video)
	if err != nil {
		return nil, err, nil
	}

	args := []any{*video}

	vs.background.Dispatch(func(args []any) {
		file := *videoFile
		defer file.Close()

		//TODO: The video name should be the id
		var backgroundVideo = args[0].(Video)

		filepath, backgroundErr := vs.filestore.Set(backgroundVideo.ID, videoFile, fileHeader)
		if backgroundErr != nil {
			vs.background.Logger.PrintError(backgroundErr, nil)
		}

		vs.background.Logger.PrintInfo(filepath, nil)

		video.Path = filepath
		video.Status = "Uploaded"

		backgroundErr = vs.store.Update(ctx, video)
		if backgroundErr != nil {
			vs.background.Logger.PrintError(backgroundErr, nil)
		}
	}, args)

	// Return Video
	return video, nil, nil
}

func (vs *Videos) CreateVideo(ctx context.Context, video *Video) (*Video, error, map[string]string) {

	validator := validator.New()

	if ValidateVideo(validator, video); !validator.Valid() {
		return nil, VideoValidationError, validator.Errors
	}

	err := vs.store.Insert(ctx, video)
	if err != nil {
		return nil, err, nil
	}

	return video, nil, nil
}

func (vs *Videos) ReadVideo(ctx context.Context, videoId int64) (*Video, error, map[string]string) {

	video, err := vs.store.ReadById(ctx, videoId)
	if err != nil {
		return nil, err, nil
	}

	return video, nil, nil
}

func (vs *Videos) UpdateVideo(ctx context.Context, videoId int64, videoInput *VideoInput) (*Video, error, map[string]string) {

	video, err := vs.store.ReadById(ctx, videoId)
	if err != nil {
		return nil, err, nil
	}

	if videoInput.Title != nil {
		video.Title = *videoInput.Title
	}

	if videoInput.Description != nil {
		video.Description = *videoInput.Description
	}

	if videoInput.PublishedDate != nil {
		video.PublishedDate = *videoInput.PublishedDate
	}

	validator := validator.New()

	if ValidateVideo(validator, video); !validator.Valid() {
		return nil, VideoValidationError, validator.Errors
	}

	err = vs.store.Update(ctx, video)
	if err != nil {
		return nil, err, nil
	}

	return video, nil, nil
}

// Initialize Video Service
func NewService(db *sql.DB, fs filestore.FileStore, bg *background.Routine) (*Videos, error) {
	vs, err := newStore(db)
	if err != nil {
		return nil, err
	}

	return &Videos{
		store:      vs,
		filestore:  fs,
		background: bg,
	}, nil
}

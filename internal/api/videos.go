package api

import (
	"context"
	"luismatosgarcia.dev/video-sharing-go/internal/videos"
	"mime/multipart"
)

func (api *API) CreateVideo(ctx context.Context, v *videos.Video) (*videos.Video, error, map[string]string) {
	v, err, validationErrors := api.videos.CreateVideo(ctx, v)
	if err != nil {
		api.Logger.PrintError(err, validationErrors)
		return nil, err, validationErrors
	}

	return v, nil, nil
}

func (api *API) UploadVideo(ctx context.Context, videoFile *multipart.File, fileHeader *multipart.FileHeader) (*videos.Video, error, map[string]string) {
	v, err, validationErrors := api.videos.UploadVideo(ctx, videoFile, fileHeader)
	if err != nil {
		api.Logger.PrintError(err, validationErrors)
		return nil, err, validationErrors
	}

	return v, nil, nil
}

func (api *API) ReadVideo(ctx context.Context, videoId int64) (*videos.Video, error, map[string]string) {
	v, err, validationErrors := api.videos.ReadVideo(ctx, videoId)
	if err != nil {
		api.Logger.PrintError(err, validationErrors)
		return nil, err, validationErrors
	}

	return v, nil, nil
}

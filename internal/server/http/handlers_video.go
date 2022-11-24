package http

import (
	"context"
	"errors"
	"fmt"
	"luismatosgarcia.dev/video-sharing-go/internal/videos"
	"net/http"
	"time"
)

func (h *Handlers) UploadVideo(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.errorHandler.badRequestResponse(w, r, err)
	}

	file, fileHeader, err := r.FormFile("video")
	if err != nil {
		h.errorHandler.badRequestResponse(w, r, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//TODO: Setup localstack and file upload interface

	videoId, err, validationErrors := h.api.UploadVideo(ctx, &file, fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, videos.VideoValidationError):
			h.errorHandler.failedValidationResponse(w, r, validationErrors)
		default:
			h.errorHandler.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"videoId": videoId,
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/videos/%d", videoId))

	err = h.httpHelper.writeJSON(w, http.StatusOK, data, headers)
	if err != nil {
		h.errorHandler.serverErrorResponse(w, r, err)
	}
}

func (h *Handlers) ReadVideo(w http.ResponseWriter, r *http.Request) {
	id, err := h.httpHelper.readIDParam(r)
	if err != nil {
		h.errorHandler.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	video, err, validationErrors := h.api.ReadVideo(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, videos.VideoValidationError):
			h.errorHandler.failedValidationResponse(w, r, validationErrors)
		default:
			h.errorHandler.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"video": video,
	}

	err = h.httpHelper.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		h.errorHandler.serverErrorResponse(w, r, err)
	}
}

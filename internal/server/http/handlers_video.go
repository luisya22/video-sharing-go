package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"luismatosgarcia.dev/video-sharing-go/internal/videos"
	"net/http"
	"time"
)

func (h *Handlers) UploadVideo(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.errorHandler.badRequestResponse(w, r, err)
	}

	f, fileHeader, err := r.FormFile("video")
	if err != nil {
		h.errorHandler.badRequestResponse(w, r, err)
	}

	file := io.Reader(f)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

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
	headers.Set("Location", fmt.Sprintf("/v1/videos/%v", videoId))

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

func (h *Handlers) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	id, err := h.httpHelper.readIDParam(r)
	if err != nil {
		h.errorHandler.notFoundResponse(w, r)
		return
	}

	var input = videos.VideoInput{}

	err = h.httpHelper.readJSON(w, r, &input)
	if err != nil {
		h.errorHandler.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	video, err, validatorErrors := h.api.UpdateVideo(ctx, id, &input)
	if err != nil {
		switch {
		case errors.Is(err, videos.VideoValidationError):
			h.errorHandler.failedValidationResponse(w, r, validatorErrors)
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

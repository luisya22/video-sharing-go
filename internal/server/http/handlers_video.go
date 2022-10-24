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

	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()

	//TODO: Setup localstack and file upload interface

	videoId, err, validationErrors := h.api.UploadVideo(ctx, &file, fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, videos.VideoValidationError):
			h.errorHandler.failedValidationResponse(w, r, validationErrors)
		default:
			h.errorHandler.serverErrorResponse(w, r, err)
		}
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

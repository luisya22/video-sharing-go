package http

import (
	"github.com/julienschmidt/httprouter"
	"luismatosgarcia.dev/video-sharing-go/internal/api"
	"net/http"
)

type Handlers struct {
	api          *api.API
	httpHelper   *Helper
	errorHandler *ErrorHandler
}

func (h *Handlers) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(h.errorHandler.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(h.errorHandler.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.healthCheckHandler)

	// Video Routes
	router.HandlerFunc(http.MethodPost, "/v1/videos", h.UploadVideo)
	router.HandlerFunc(http.MethodGet, "/v1/videos/:id", h.ReadVideo)
	router.HandlerFunc(http.MethodPatch, "/v1/videos/:id", h.UpdateVideo)

	return router
}

func (h *Handlers) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "Testing",
			"version":     "1.0.0",
		},
	}

	err := h.httpHelper.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		h.errorHandler.serverErrorResponse(w, r, err)
	}
}

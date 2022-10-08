package http

import (
	"github.com/julienschmidt/httprouter"
	"luismatosgarcia.dev/video-sharing-go/internal/api"
	"luismatosgarcia.dev/video-sharing-go/storage"
	"net/http"
)

type Handlers struct {
	api          *api.API
	httpHelper   *Helper
	errorHandler *ErrorHandler
}

func (handler *Handlers) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(handler.errorHandler.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(handler.errorHandler.methodNotAllowedResponse)

	fileServer := http.FileServer(http.FS(storage.Files))
	router.Handler(http.MethodGet, "/videos/*filepath", fileServer)
	router.Handler(http.MethodGet, "/html/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.healthCheckHandler)

	return router
}

func (handler *Handlers) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": "Testing",
			"version":     "1.0.0",
		},
	}

	err := handler.httpHelper.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		handler.errorHandler.serverErrorResponse(w, r, err)
	}
}

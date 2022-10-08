package http

import (
	"fmt"
	"luismatosgarcia.dev/video-sharing-go/internal/api"
	"net/http"
)

type ErrorHandler struct {
	api        *api.API
	httpHelper *Helper
}

func (e *ErrorHandler) logError(r *http.Request, err error) {
	e.api.Logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (e *ErrorHandler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := e.httpHelper.writeJSON(w, status, env, nil)
	if err != nil {
		e.logError(r, err)
		w.WriteHeader(500)
	}
}

func (e *ErrorHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	e.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (e *ErrorHandler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	e.errorResponse(w, r, http.StatusNotFound, message)
}

func (e *ErrorHandler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	e.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (e *ErrorHandler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (e *ErrorHandler) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	e.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (e *ErrorHandler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to and edit conflict, please try again"
	e.errorResponse(w, r, http.StatusConflict, message)
}

func (e *ErrorHandler) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	e.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (e *ErrorHandler) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *ErrorHandler) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *ErrorHandler) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *ErrorHandler) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	e.errorResponse(w, r, http.StatusForbidden, message)
}

func (e *ErrorHandler) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	e.errorResponse(w, r, http.StatusForbidden, message)
}

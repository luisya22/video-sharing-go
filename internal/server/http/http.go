package http

import (
	"context"
	"errors"
	"fmt"
	"luismatosgarcia.dev/video-sharing-go/internal/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTP struct {
	server *http.Server
	cfg    *Config
	api    *api.API
}

// TODO: Move config where is used DB, cache, etc...
type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	Cors struct {
		TrustedOrigins []string
	}
}

func (h *HTTP) Start() {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		h.api.Logger.PrintInfo("shutting down server", map[string]string{"signal": s.String()})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := h.server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		h.api.Logger.PrintInfo("completing background tasks", map[string]string{
			"addr": h.server.Addr,
		})

		h.api.Wg.Wait()
		shutdownError <- nil

	}()

	h.api.Logger.PrintInfo("starting server", map[string]string{
		"addr": h.server.Addr,
		"env":  h.cfg.Env,
	})

	err := h.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		h.api.Logger.PrintFatal(err, nil)
	}

	err = <-shutdownError
	if err != nil {
		h.api.Logger.PrintFatal(err, nil)
	}

	h.api.Logger.PrintInfo("stopped server", map[string]string{
		"addr": h.server.Addr,
	})
}

func NewService(cfg *Config, a *api.API) (*HTTP, error) {
	helper := &Helper{
		api: a,
	}

	errHandler := &ErrorHandler{
		api:        a,
		httpHelper: helper,
	}

	h := &Handlers{
		api:          a,
		httpHelper:   helper,
		errorHandler: errHandler,
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      h.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &HTTP{
		server: httpServer,
		cfg:    cfg,
		api:    a,
	}, nil

}

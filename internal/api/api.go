package api

import (
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/jsonlog"
	"sync"
	"time"
)

var (
	now = time.Now()
)

type API struct {
	Logger *jsonlog.Logger
	Wg     sync.WaitGroup
}

func NewService(l *jsonlog.Logger) (*API, error) {
	return &API{
		Logger: l,
	}, nil
}

package api

import (
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/jsonlog"
	"luismatosgarcia.dev/video-sharing-go/internal/videos"
	"sync"
	"time"
)

var (
	now = time.Now()
)

type API struct {
	Logger *jsonlog.Logger
	Wg     sync.WaitGroup
	videos *videos.Videos
}

func NewService(l *jsonlog.Logger, v *videos.Videos) (*API, error) {
	return &API{
		Logger: l,
		videos: v,
	}, nil
}

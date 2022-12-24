package api

import (
	"luismatosgarcia.dev/video-sharing-go/internal/background"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/jsonlog"
	"luismatosgarcia.dev/video-sharing-go/internal/videos"
	"time"
)

var (
	now = time.Now()
)

type API struct {
	Logger            *jsonlog.Logger
	BackgroundRoutine background.Routine
	videos            videos.Videos
}

func NewService(l *jsonlog.Logger, bg background.Routine, v videos.Videos) (*API, error) {
	return &API{
		Logger:            l,
		videos:            v,
		BackgroundRoutine: bg,
	}, nil
}

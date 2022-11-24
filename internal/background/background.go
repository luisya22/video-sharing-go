package background

import (
	"errors"
	"fmt"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/jsonlog"
	"sync"
)

var (
	ArgsTypeMismatchError = errors.New("background args types mistmatched")
)

type Routine struct {
	Wg     sync.WaitGroup
	Logger *jsonlog.Logger
}

func (br *Routine) Dispatch(fn func(args []any), args []any) {
	br.Wg.Add(1)

	go func(args []any) {
		defer br.Wg.Done()

		defer func() {
			if err := recover(); err != nil {
				br.Logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn(args)
	}(args)
}

func NewService(l *jsonlog.Logger) (*Routine, error) {
	return &Routine{
		Logger: l,
	}, nil
}

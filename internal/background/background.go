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

type Routine interface {
	Dispatch(fn func(args []any), args []any)
	Wait()
	PrintInfo(message string, properties map[string]string)
	PrintError(err error, properties map[string]string)
}

type RoutineImpl struct {
	Wg     sync.WaitGroup
	Logger *jsonlog.Logger
}

func (br *RoutineImpl) Dispatch(fn func(args []any), args []any) {
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

func (br *RoutineImpl) PrintInfo(message string, properties map[string]string) {
	br.Logger.PrintInfo(message, properties)
}

func (br *RoutineImpl) PrintError(err error, properties map[string]string) {
	br.Logger.PrintError(err, properties)
}

func (br *RoutineImpl) Wait() {
	br.Wg.Wait()
}

func NewService(l *jsonlog.Logger) (Routine, error) {
	return &RoutineImpl{
		Logger: l,
	}, nil
}

// RoutineMock mock for testing purposes
type RoutineMock struct{}

func (r *RoutineMock) Dispatch(fn func(args []any), args []any) {}

func (r *RoutineMock) Wait() {}

func (r *RoutineMock) PrintInfo(message string, properties map[string]string) {}

func (r *RoutineMock) PrintError(err error, properties map[string]string) {}

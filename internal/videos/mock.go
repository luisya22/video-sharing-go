package videos

import (
	"context"
	"io"
	"luismatosgarcia.dev/video-sharing-go/internal/tests"
	"mime/multipart"
)

type Mock struct {
	Video     *Video
	Err       error
	ErrorsMap map[string]string
}

func (m Mock) UploadVideo(ctx context.Context, videoFileReader *io.Reader, fileHeader *multipart.FileHeader) (*Video, error, map[string]string) {
	return m.Video, m.Err, m.ErrorsMap
}

func (m Mock) uploadVideoBackground(ctx context.Context, video *Video, videoFileReader *io.Reader, fileHeader *multipart.FileHeader) {
}

func (m Mock) CreateVideo(ctx context.Context, video *Video) (*Video, error, map[string]string) {
	return m.Video, m.Err, m.ErrorsMap
}

func (m Mock) ReadVideo(ctx context.Context, videoId int64) (*Video, error, map[string]string) {
	return m.Video, m.Err, m.ErrorsMap
}

func (m Mock) UpdateVideo(ctx context.Context, videoId int64, videoInput *VideoInput) (*Video, error, map[string]string) {
	return m.Video, m.Err, m.ErrorsMap
}

// Store

type storeMock struct {
	fnCalls map[string]int
	video   *Video
	err     map[string]error
}

func (s storeMock) Insert(ctx context.Context, v *Video) error {
	tests.Called(s.fnCalls, "Insert")
	return s.err["Insert"]
}

func (s storeMock) Update(ctx context.Context, v *Video) error {
	tests.Called(s.fnCalls, "Update")
	return s.err["Update"]
}

func (s storeMock) ReadById(ctx context.Context, videoId int64) (*Video, error) {
	tests.Called(s.fnCalls, "ReadById")
	return s.video, s.err["ReadById"]
}

func (s storeMock) GetFnCalls(fnName string) int {
	value, exists := s.fnCalls[fnName]

	if !exists {
		return 0
	}

	return value
}

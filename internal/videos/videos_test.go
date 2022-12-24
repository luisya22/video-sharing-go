package videos

import (
	"bytes"
	"context"
	"io"
	"luismatosgarcia.dev/video-sharing-go/internal/background"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/filestore"
	"luismatosgarcia.dev/video-sharing-go/internal/tests/assert"
	"mime/multipart"
	"testing"
)

func TestVideos_UploadVideo(t *testing.T) {
	testsMap := []struct {
		name           string
		videoFile      io.Reader
		fileHeader     multipart.FileHeader
		wants          Video
		storeMock      store
		filestoreMock  filestore.FileStore
		backgroundMock background.Routine
	}{
		{
			name:      "Can Upload",
			videoFile: &bytes.Buffer{},
			fileHeader: multipart.FileHeader{
				Filename: "Video.mp4",
				Header:   nil,
				Size:     0,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     nil,
			},
			filestoreMock:  &filestore.Mock{},
			backgroundMock: &background.RoutineMock{},
		},
	}

	for _, tt := range testsMap {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				store:      tt.storeMock,
				filestore:  tt.filestoreMock,
				background: tt.backgroundMock,
			}

			_, err, _ := service.UploadVideo(context.Background(), &tt.videoFile, &tt.fileHeader)

			assert.NilError(t, err)

			fs := tt.storeMock.(storeMock)

			assert.Equal(t, fs.GetFnCalls("Insert"), 1)

			//TODO: What to assert
			//TODO: Should I make sure that it calls other services(mocks in this case)
		})
	}

}

// Mock

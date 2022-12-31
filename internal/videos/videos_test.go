package videos

import (
	"bytes"
	"context"
	"io"
	"luismatosgarcia.dev/video-sharing-go/internal/background"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/datastore"
	"luismatosgarcia.dev/video-sharing-go/internal/pkg/filestore"
	"luismatosgarcia.dev/video-sharing-go/internal/tests/assert"
	"mime/multipart"
	"testing"
	"time"
)

type testResult struct {
	video          Video
	fnCalls        map[string]int
	shouldError    bool
	validateFields bool
}

func TestService_UploadVideo(t *testing.T) {
	testsMap := []struct {
		name           string
		videoFile      io.Reader
		fileHeader     multipart.FileHeader
		wants          testResult
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
			filestoreMock: filestore.Mock{
				FnCalls: make(map[string]int),
				Str:     "/videos/Video.mp4",
				Err:     nil,
			},
			backgroundMock: &background.RoutineMock{},
			wants: testResult{
				video: Video{},
				fnCalls: map[string]int{
					"vsInsert": 1,
					"fsSet":    1,
				},
				shouldError:    false,
				validateFields: false,
			},
		},
		{
			name:      "Store Returns Error",
			videoFile: &bytes.Buffer{},
			fileHeader: multipart.FileHeader{
				Filename: "Video.mp4",
				Header:   nil,
				Size:     0,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     map[string]error{"Insert": datastore.ErrEditConflict},
			},
			filestoreMock: filestore.Mock{
				FnCalls: make(map[string]int),
				Str:     "/videos/Video.mp4",
				Err:     nil,
			},
			backgroundMock: &background.RoutineMock{},
			wants: testResult{
				video: Video{},
				fnCalls: map[string]int{
					"vsInsert": 1,
					"fsSet":    0,
				},
				shouldError: true,
			},
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
			tt.backgroundMock.Wait()

			if !tt.wants.shouldError {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err)
			}

			vs := tt.storeMock.(storeMock)
			assert.Equal(t, vs.GetFnCalls("Insert"), tt.wants.fnCalls["vsInsert"])

			fs := tt.filestoreMock.(filestore.Mock)
			assert.Equal(t, fs.GetFnCalls("Set"), tt.wants.fnCalls["fsSet"])
		})
	}

}

func TestService_CreateVideo(t *testing.T) {

	newTitle := "New Video Title"
	newDescription := "New Video Description"
	newPublishedDate := time.Now().AddDate(0, 0, 1)
	pastPublishedDate := time.Now().AddDate(0, 0, -1)

	testsMap := []struct {
		name       string
		id         int64
		videoInput *VideoInput
		wants      testResult
		storeMock  store
	}{
		{
			name: "Can Update",
			id:   1,
			videoInput: &VideoInput{
				Title:         &newTitle,
				Description:   &newDescription,
				PublishedDate: &newPublishedDate,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video: &Video{
					Title:         "Random Title",
					Description:   "Random Description",
					PublishedDate: time.Now(),
				},
				err: nil,
			},
			wants: testResult{
				video: Video{
					Title:         newTitle,
					Description:   newDescription,
					PublishedDate: newPublishedDate,
				},
				fnCalls: map[string]int{
					"vsReadById": 1,
					"vsUpdate":   1,
				},
				shouldError:    false,
				validateFields: true,
			},
		},
		{
			name: "Validate Title",
			id:   1,
			videoInput: &VideoInput{
				Title:         nil,
				Description:   &newDescription,
				PublishedDate: &newPublishedDate,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     nil,
			},
			wants: testResult{
				video: Video{},
				fnCalls: map[string]int{
					"vsReadById": 1,
					"vsUpdate":   0,
				},
				shouldError:    true,
				validateFields: false,
			},
		},
		{
			name: "Validate Description",
			id:   1,
			videoInput: &VideoInput{
				Title:         &newTitle,
				Description:   nil,
				PublishedDate: &newPublishedDate,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     nil,
			},
			wants: testResult{
				video: Video{},
				fnCalls: map[string]int{
					"vsReadById": 1,
					"vsUpdate":   0,
				},
				shouldError:    true,
				validateFields: false,
			},
		},
		{
			name: "Validate Published Date",
			id:   1,
			videoInput: &VideoInput{
				Title:         &newTitle,
				Description:   &newDescription,
				PublishedDate: &pastPublishedDate,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     nil,
			},
			wants: testResult{
				video: Video{},
				fnCalls: map[string]int{
					"vsReadById": 1,
					"vsUpdate":   0,
				},
				shouldError:    true,
				validateFields: false,
			},
		},
		{
			name: "Store Update Error",
			id:   1,
			videoInput: &VideoInput{
				Title:         &newTitle,
				Description:   &newDescription,
				PublishedDate: &newPublishedDate,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     map[string]error{"Update": datastore.ErrEditConflict},
			},
			wants: testResult{
				video: Video{
					Title:         newTitle,
					Description:   newDescription,
					PublishedDate: pastPublishedDate,
				},
				fnCalls: map[string]int{
					"vsReadById": 1,
					"vsUpdate":   1,
				},
				shouldError:    true,
				validateFields: false,
			},
		},
		{
			name: "Store ReadById Error",
			id:   1,
			videoInput: &VideoInput{
				Title:         &newTitle,
				Description:   &newDescription,
				PublishedDate: &newPublishedDate,
			},
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video:   &Video{},
				err:     map[string]error{"ReadById": datastore.ErrRecordNotFound, "Update": nil},
			},
			wants: testResult{
				video: Video{
					Title:         newTitle,
					Description:   newDescription,
					PublishedDate: newPublishedDate,
				},
				fnCalls: map[string]int{
					"vsReadById": 1,
					"vsUpdate":   0,
				},
				shouldError:    true,
				validateFields: false,
			},
		},
	}

	for _, tt := range testsMap {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				store: tt.storeMock,
			}

			v, err, _ := service.UpdateVideo(context.Background(), tt.id, tt.videoInput)

			if !tt.wants.shouldError {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err)
			}

			if tt.wants.validateFields {
				assert.Equal(t, v.Title, tt.wants.video.Title)
				assert.Equal(t, v.Description, tt.wants.video.Description)
				assert.Equal(t, v.PublishedDate, tt.wants.video.PublishedDate)
			}

			vs := tt.storeMock.(storeMock)
			assert.Equal(t, vs.GetFnCalls("ReadById"), tt.wants.fnCalls["vsReadById"])
			assert.Equal(t, vs.GetFnCalls("Update"), tt.wants.fnCalls["vsUpdate"])
		})
	}
}

func TestService_ReadVideo(t *testing.T) {
	testMaps := []struct {
		name      string
		id        int64
		wants     testResult
		storeMock store
	}{
		{
			name: "Can Read",
			id:   1,
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video: &Video{
					ID:            1,
					Title:         "Video Title",
					Description:   "Video Description",
					Path:          "/videos/1",
					ImgPath:       "/videosImg/1",
					Status:        "Published",
					PublishedDate: time.Now(),
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
					Version:       1,
				},
				err: nil,
			},
			wants: testResult{
				video: Video{
					ID:            1,
					Title:         "Video Title",
					Description:   "Video Description",
					Path:          "/videos/1",
					ImgPath:       "/videosImg/1",
					Status:        "Published",
					PublishedDate: time.Now(),
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
					Version:       1,
				},
				fnCalls: map[string]int{
					"vsReadById": 1,
				},
				shouldError:    false,
				validateFields: true,
			},
		},
		{
			name: "Store ReadById Error",
			id:   1,
			storeMock: storeMock{
				fnCalls: make(map[string]int),
				video: &Video{
					ID:            1,
					Title:         "Video Title",
					Description:   "Video Description",
					Path:          "/videos/1",
					ImgPath:       "/videosImg/1",
					Status:        "Published",
					PublishedDate: time.Now(),
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
					Version:       1,
				},
				err: map[string]error{"ReadById": datastore.ErrRecordNotFound},
			},
			wants: testResult{
				fnCalls: map[string]int{
					"vsReadById": 1,
				},
				shouldError:    true,
				validateFields: false,
			},
		},
	}

	for _, tt := range testMaps {
		t.Run(tt.name, func(t *testing.T) {
			service := Service{
				store: tt.storeMock,
			}

			v, err, _ := service.ReadVideo(context.Background(), tt.id)

			if !tt.wants.shouldError {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err)
			}

			if tt.wants.validateFields {
				assert.Equal(t, v.ID, tt.wants.video.ID)
				assert.Equal(t, v.Title, tt.wants.video.Title)
				assert.Equal(t, v.Description, tt.wants.video.Description)
				assert.Equal(t, v.Path, tt.wants.video.Path)
				assert.Equal(t, v.ImgPath, tt.wants.video.ImgPath)
				assert.Equal(t, v.Status, tt.wants.video.Status)
				assert.Equal(t, v.PublishedDate, tt.wants.video.PublishedDate)
				assert.Equal(t, v.CreatedAt, tt.wants.video.CreatedAt)
				assert.Equal(t, v.UpdatedAt, tt.wants.video.UpdatedAt)
				assert.Equal(t, v.Version, tt.wants.video.Version)
			}

			vs := tt.storeMock.(storeMock)
			assert.Equal(t, vs.GetFnCalls("ReadById"), tt.wants.fnCalls["vsReadById"])
		})
	}
}

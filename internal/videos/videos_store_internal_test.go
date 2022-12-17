package videos

import (
	"context"
	"luismatosgarcia.dev/video-sharing-go/internal/tests"
	"luismatosgarcia.dev/video-sharing-go/internal/tests/assert"
	"testing"
	"time"
)

func TestVideos_InsertVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("store: skipping integration test")
	}

	testsMap := []struct {
		name  string
		video *Video
	}{
		{
			name:  "Can Insert",
			video: &Video{},
		},
		{
			name: "Populated Video does not insert data",
			video: &Video{
				ID:            0,
				Title:         "Video Title",
				Description:   "Video Description",
				Path:          "/videos/path",
				ImgPath:       "/videos/path",
				Status:        "Published",
				PublishedDate: time.Time{},
				CreatedAt:     time.Time{},
				UpdatedAt:     time.Time{},
				Version:       17,
			},
		},
	}

	db := tests.NewTestDB(t)

	for _, tt := range testsMap {
		t.Run(tt.name, func(t *testing.T) {
			store := videoStore{db: db}

			err := store.Insert(context.Background(), tt.video)

			assert.Equal(t, tt.video.Status, "Uploading")

			assert.NilError(t, err)

			//Assert value from db
			query := `SELECT id, status FROM videos
                      WHERE id = $1`

			var video Video

			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = db.QueryRowContext(dbCtx, query, tt.video.ID).Scan(
				&video.ID,
				&video.Status,
			)
			assert.NilError(t, err)

			assert.Equal(t, tt.video.ID, video.ID)
			assert.Equal(t, tt.video.Status, video.Status)
		})
	}

}

func TestVideos_UpdateVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("store: skipping integration test")
	}

	testsMap := []struct {
		name  string
		video *Video
		wants *Video
	}{
		{
			name: "Can Update",
			video: &Video{
				ID:            1,
				Title:         "Video #1",
				Description:   "Video #1 description",
				Path:          "http://example.com/video.mp4",
				ImgPath:       "http://example.com/img.jgp",
				Status:        "Processed",
				PublishedDate: time.Now(),
				CreatedAt:     time.Time{},
				UpdatedAt:     time.Time{},
				Version:       1,
			},
			wants: &Video{
				ID:            1,
				Title:         "Video #1",
				Description:   "Video #1 description",
				Path:          "http://example.com/video.mp4",
				ImgPath:       "http://example.com/img.jgp",
				Status:        "Processed",
				PublishedDate: time.Now().UTC(),
				CreatedAt:     time.Time{},
				UpdatedAt:     time.Time{},
				Version:       2,
			},
		},
	}

	db := tests.NewTestDB(t)

	for _, tt := range testsMap {
		t.Run(tt.name, func(t *testing.T) {
			store := videoStore{db: db}

			// Run function
			err := store.Update(context.Background(), tt.video)
			assert.NilError(t, err)

			query := `SELECT id, title, description, video_path, thumbnail_path, status, published_at, version FROM videos
                      WHERE id = $1`

			var video Video

			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err = db.QueryRowContext(dbCtx, query, tt.video.ID).Scan(
				&video.ID,
				&video.Title,
				&video.Description,
				&video.Path,
				&video.ImgPath,
				&video.Status,
				&video.PublishedDate,
				&video.Version,
			)

			assert.Equal(t, tt.wants.ID, video.ID)
			assert.Equal(t, tt.wants.Title, video.Title)
			assert.Equal(t, tt.wants.Description, video.Description)
			assert.Equal(t, tt.wants.Path, video.Path)
			assert.Equal(t, tt.wants.ImgPath, video.ImgPath)
			assert.Equal(t, tt.wants.Status, video.Status)
			assert.Equal(t, tt.wants.PublishedDate.UTC().Format("2022-12-01"), video.PublishedDate.UTC().Format("2022-12-01"))
			assert.Equal(t, tt.wants.Version, video.Version)

		})
	}
}

func TestVideoStore_ReadByIdVideo(t *testing.T) {
	testsMap := []struct {
		name    string
		videoId int64
		wants   Video
	}{
		{
			name:    "Can Read",
			videoId: 1,
			wants: Video{
				ID:          1,
				Title:       "Video #0",
				Description: "Video Description",
				Path:        "No Path",
				ImgPath:     "No Thumbnail",
				Status:      "No Status",
			},
		},
	}

	db := tests.NewTestDB(t)
	for _, tt := range testsMap {
		t.Run(tt.name, func(t *testing.T) {

			store := videoStore{db: db}

			video, err := store.ReadById(context.Background(), tt.videoId)

			assert.NilError(t, err)

			assert.Equal(t, video.ID, tt.videoId)
			assert.Equal(t, video.Description, tt.wants.Description)
			assert.Equal(t, video.Path, tt.wants.Path)
			assert.Equal(t, video.ImgPath, tt.wants.ImgPath)
			assert.Equal(t, video.Status, tt.wants.Status)

		})
	}
}

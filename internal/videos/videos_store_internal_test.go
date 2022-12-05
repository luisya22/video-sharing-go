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

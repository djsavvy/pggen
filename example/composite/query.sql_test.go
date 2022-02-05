package composite

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/djsavvy/pggen/internal/errs"
	"github.com/djsavvy/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewQuerier_SearchScreenshots(t *testing.T) {
	conn, cleanup := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanup()

	q := NewQuerier(conn)
	screenshotID := 99
	screenshot1 := insertScreenshotBlock(t, q, screenshotID, "body1")
	screenshot2 := insertScreenshotBlock(t, q, screenshotID, "body2")
	want := []SearchScreenshotsRow{
		{
			ID: screenshotID,
			Blocks: []Blocks{
				{
					ID:           screenshot1.ID,
					ScreenshotID: screenshotID,
					Body:         screenshot1.Body,
				},
				{
					ID:           screenshot2.ID,
					ScreenshotID: screenshotID,
					Body:         screenshot2.Body,
				},
			},
		},
	}

	t.Run("SearchScreenshots", func(t *testing.T) {
		rows, err := q.SearchScreenshots(context.Background(), SearchScreenshotsParams{
			Body:   "body",
			Limit:  5,
			Offset: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, want, rows)
	})

	t.Run("SearchScreenshotsBatch", func(t *testing.T) {
		batch := &pgx.Batch{}
		q.SearchScreenshotsBatch(batch, SearchScreenshotsParams{
			Body:   "body",
			Limit:  5,
			Offset: 0,
		})
		results := conn.SendBatch(context.Background(), batch)
		defer errs.CaptureT(t, results.Close, "close batch results")
		rows, err := q.SearchScreenshotsScan(results)
		require.NoError(t, err)
		assert.Equal(t, want, rows)
	})

	t.Run("SearchScreenshotsOneCol", func(t *testing.T) {
		rows, err := q.SearchScreenshotsOneCol(context.Background(), SearchScreenshotsOneColParams{
			Body:   "body",
			Limit:  5,
			Offset: 0,
		})
		require.NoError(t, err)
		assert.Equal(t, [][]Blocks{want[0].Blocks}, rows)
	})

	t.Run("SearchScreenshotsOneColBatch", func(t *testing.T) {
		batch := &pgx.Batch{}
		q.SearchScreenshotsOneColBatch(batch, SearchScreenshotsOneColParams{
			Body:   "body",
			Limit:  5,
			Offset: 0,
		})
		results := conn.SendBatch(context.Background(), batch)
		defer errs.CaptureT(t, results.Close, "close batch results")
		rows, err := q.SearchScreenshotsOneColScan(results)
		require.NoError(t, err)
		assert.Equal(t, [][]Blocks{want[0].Blocks}, rows)
	})
}

func insertScreenshotBlock(t *testing.T, q *DBQuerier, screenID int, body string) InsertScreenshotBlocksRow {
	t.Helper()
	row, err := q.InsertScreenshotBlocks(context.Background(), screenID, body)
	require.NoError(t, err, "insert screenshot blocks")
	return row
}

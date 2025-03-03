package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lelouchhh/friendly-basketball-media/domain"
	"github.com/lelouchhh/friendly-basketball-media/internal/logger"
	"go.uber.org/zap"
)

type Video struct {
	conn   *sqlx.DB
	logger logger.Logger
}

func (r Video) GetVideo(ctx context.Context, videoID int) (domain.Video, error) {
	r.logger.Debug("get video...", zap.String("layer", "repo"))

	query := `
		select
			video_id, event_id, title, file_path, upload_date, duration, resolution, size, url
		from
			media.videos
		where
			video_id = $1;	
	`
	var video domain.Video
	err := r.conn.QueryRowContext(ctx, query, videoID).
		Scan(&video.VideoID, &video.EventID, &video.Title, &video.FilePath, &video.UploadDate, &video.Duration, &video.Resolution, &video.Size, &video.URL)
	if err != nil {
		r.logger.Error("failed to query video", zap.Error(err))
		return domain.Video{}, err
	}
	return video, nil
}

func (r Video) GetVideoList(ctx context.Context, eventID int) ([]domain.Video, error) {
	r.logger.Debug("get video...", zap.String("layer", "repo"))

	query := `
		select
			video_id, event_id, title, file_path, upload_date, duration, resolution, size, url
		from
			media.videos
		where
			event_id = $1;	
	`
	rows, err := r.conn.QueryContext(ctx, query, eventID)
	if err != nil {
		r.logger.Error("failed to query video", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var videos []domain.Video
	for rows.Next() {
		var video domain.Video
		err := rows.Scan(&video.VideoID, &video.EventID, &video.Title, &video.FilePath, &video.UploadDate, &video.Duration, &video.Resolution, &video.Size, &video.URL)
		if err != nil {
			r.logger.Error("failed to scan video", zap.Error(err))
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func NewVideo(conn *sqlx.DB, logger logger.Logger) Video {
	return Video{
		conn:   conn,
		logger: logger,
	}
}

func (r Video) Upload(ctx context.Context, videos []domain.Video) ([]domain.Video, error) {
	r.logger.Debug("Uploading videos...", zap.String("layer", "repo"))
	tx, err := r.conn.Beginx()
	var resultVideos []domain.Video
	if err != nil {
		r.logger.Error("Can't begin transaction", zap.Error(err))

		return nil, err
	}
	defer func() {
		if err != nil {
			r.logger.Error("can't commit", zap.Error(err))
			_ = tx.Rollback()
		} else {
			r.logger.Debug("commit")

			_ = tx.Commit()
		}
	}()

	for _, video := range videos {
		var v domain.Video
		err := tx.QueryRowContext(ctx, `
            INSERT INTO media.videos (event_id, title, file_path, upload_date, size, url)
            VALUES ($1, $2, $3, $4, $5, $6) returning video_id, event_id, title, file_path, upload_date, size, url;
        `, video.EventID, video.Title, video.FilePath, video.UploadDate, video.Size, video.URL).Scan(&v.VideoID, &v.EventID, &v.Title, &v.FilePath, &v.UploadDate, &v.Size, &v.URL)
		if err != nil {
			return nil, err
		}
		resultVideos = append(resultVideos, v)
	}

	return resultVideos, nil
}

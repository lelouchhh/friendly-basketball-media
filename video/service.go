package video

import (
	"context"
	"errors"
	"github.com/lelouchhh/friendly-basketball-media/domain"
	"github.com/lelouchhh/friendly-basketball-media/internal/logger"
	"go.uber.org/zap"
	"time"
)

type VideoRepository interface {
	Upload(ctx context.Context, videos []domain.Video) ([]domain.Video, error)
	GetVideo(ctx context.Context, videoID int) (domain.Video, error)
	GetVideoList(ctx context.Context, eventID int) ([]domain.Video, error)
}

type Service struct {
	repo   VideoRepository
	logger logger.Logger
}

func (s Service) GetVideo(ctx context.Context, videoID int) (domain.Video, error) {
	s.logger.Debug("get video...", zap.String("layer", "service"))

	return s.repo.GetVideo(ctx, videoID)
}

func (s Service) GetVideoList(ctx context.Context, eventID int) ([]domain.Video, error) {
	s.logger.Debug("get video list...", zap.String("layer", "service"))

	return s.repo.GetVideoList(ctx, eventID)
}

func NewService(repo VideoRepository, logger logger.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s Service) Upload(ctx context.Context, eventID int, videos []domain.Video) ([]domain.Video, error) {
	s.logger.Debug("Uploading videos...", zap.String("layer", "service"))
	if eventID <= 0 {
		s.logger.Error("eventID must be positive")
		return nil, errors.New("invalid event_id")
	}

	for i := range videos {
		videos[i].EventID = eventID
		videos[i].UploadDate = time.Now()
	}

	return s.repo.Upload(ctx, videos)
}

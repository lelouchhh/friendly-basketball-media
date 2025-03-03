package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lelouchhh/friendly-basketball-media/domain"
	"github.com/lelouchhh/friendly-basketball-media/internal/delivery/http/middleware"
	"github.com/lelouchhh/friendly-basketball-media/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type VideoService interface {
	Upload(ctx context.Context, eventID int, videos []domain.Video) ([]domain.Video, error)
	GetVideo(ctx context.Context, videoID int) (domain.Video, error)
	GetVideoList(ctx context.Context, eventID int) ([]domain.Video, error)
}

type VideoHandler struct {
	videoService VideoService
	jwt          string
	logger       logger.Logger
}

func NewVideoHandler(e *echo.Echo, secretKey string, service VideoService, logger logger.Logger) {
	handler := &VideoHandler{
		videoService: service,
		jwt:          secretKey,
		logger:       logger,
	}
	g := e.Group("/api/v1/video")
	g.Use(middleware.JWTMiddleware(secretKey))
	g.POST("/upload/:event_id", handler.Upload)
	g.GET("/:video_id", handler.GetVideo)
	g.GET("/event/:event_id", handler.GetVideos)

}

func (h *VideoHandler) Upload(c echo.Context) error {
	eventIDParam := c.Param("event_id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil || eventID <= 0 {
		h.logger.Error("Invalid event_id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event_id"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("Failed to parse multipart form", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
	}

	files := form.File["videos"]
	if len(files) == 0 {
		h.logger.Warn("No files uploaded")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No files uploaded"})
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			h.logger.Error("Failed to create upload directory", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
		}
	}

	var videos []domain.Video
	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			h.logger.Error("Failed to open file", zap.String("filename", fileHeader.Filename), zap.Error(err))
			continue
		}
		defer src.Close()

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
		filePath := filepath.Join(uploadDir, fileName)

		dst, err := os.Create(filePath)
		if err != nil {
			h.logger.Error("Failed to create file", zap.String("filename", fileName), zap.Error(err))
			continue
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			h.logger.Error("Failed to copy file", zap.String("filename", fileName), zap.Error(err))
			continue
		}

		// Формируем URL для доступа к файлу
		fileURL := fmt.Sprintf("http://%s/videos/%s", c.Request().Host, fileName)

		videos = append(videos, domain.Video{
			Title:    fileHeader.Filename,
			FilePath: filePath,
			Size:     int(fileHeader.Size),
			URL:      fileURL, // Добавляем URL для доступа к файлу
		})
	}

	if len(videos) == 0 {
		h.logger.Warn("No valid videos uploaded")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No valid videos uploaded"})
	}

	// Сохраняем видео в базу данных через сервис
	if videos, err = h.videoService.Upload(c.Request().Context(), eventID, videos); err != nil {
		h.logger.Error("Failed to upload videos", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload videos"})
	}

	h.logger.Info("Videos uploaded successfully", zap.Int("count", len(videos)))
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Videos uploaded successfully",
		"videos":  videos,
	})
}

func (h *VideoHandler) GetVideo(c echo.Context) error {
	videoIDparam := c.Param("video_id")
	videoID, err := strconv.Atoi(videoIDparam)
	if err != nil || videoID <= 0 {
		h.logger.Error("Invalid event_id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event_id"})
	}
	video, err := h.videoService.GetVideo(c.Request().Context(), videoID)
	if err != nil {
		h.logger.Error("Failed to get video", zap.Error(err))
		return err
	}
	return c.JSON(http.StatusOK, video)
}
func (h *VideoHandler) GetVideos(c echo.Context) error {
	eventIDparam := c.Param("event_id")
	eventID, err := strconv.Atoi(eventIDparam)
	if err != nil || eventID <= 0 {
		h.logger.Error("Invalid event_id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event_id"})
	}
	video, err := h.videoService.GetVideoList(c.Request().Context(), eventID)
	if err != nil {
		h.logger.Error("Failed to get video", zap.Error(err))
		return err
	}
	return c.JSON(http.StatusOK, video)
}

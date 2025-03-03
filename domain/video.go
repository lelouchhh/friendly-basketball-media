package domain

import "time"

type Video struct {
	VideoID    int       `db:"video_id" json:"video_id"`
	EventID    int       `db:"event_id" json:"event_id"`
	Title      string    `db:"title" json:"title"`
	FilePath   string    `db:"file_path" json:"file_path"`
	UploadDate time.Time `db:"upload_date" json:"upload_date"`
	Duration   *string   `db:"duration" json:"duration"`
	Resolution *string   `db:"resolution" json:"resolution"`
	Size       int       `db:"size" json:"size"`
	URL        string    `db:"url" json:"url"`
}

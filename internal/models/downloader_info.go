package models

type InfoResponse struct {
	Title     string         `json:"title"`
	Author    string         `json:"author"`
	Thumbnail string         `json:"thumbnail"`
	Duration  string         `json:"duration"`
	Formats   []FormatOption `json:"formats"`
}

type FormatOption struct {
	FormatID       string `json:"format_id"`
	Quality        string `json:"quality"`
	MimeType       string `json:"mime_type"`
	HasVideo       bool   `json:"has_video"`
	HasAudio       bool   `json:"has_audio"`
	AudioChannels  int    `json:"audio_channels"`
	Filesize       *int64 `json:"filesize"`
	FilesizeApprox *int64 `json:"filesize_approx"`
}

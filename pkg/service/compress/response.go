package compress

import "github.com/Hargeon/videocmprs/pkg/repository/video"

type Response struct {
	RequestID      int64           `json:"request_id"`
	OriginalVideo  *video.Resource `json:"original_video,omitempty"`
	ConvertedVideo *video.Resource `json:"converted_video,omitempty"`
	Error          string          `json:"error,omitempty"`
}

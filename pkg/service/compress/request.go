package compress

import (
	"fmt"

	"github.com/Hargeon/videocmprs/pkg/repository/request"
)

type Request struct {
	RequestID      int64  `json:"request_id"`
	Bitrate        int64  `json:"bitrate"`
	Resolution     string `json:"resolution"`
	Ratio          string `json:"ratio"`
	VideoID        int64  `json:"video_id"`
	VideoServiceID string `json:"video_service_id"`
}

// NewRequest initialize *Request. Need use *request.Resource with *OriginalVideo
func NewRequest(r *request.Resource) *Request {
	req := new(Request)
	req.RequestID = r.ID
	req.Bitrate = r.Bitrate

	if r.ResolutionX != 0 || r.ResolutionY != 0 {
		req.Resolution = fmt.Sprintf("%d:%d",
			r.ResolutionX, r.ResolutionY)
	}

	if r.RatioX != 0 || r.RatioY != 0 {
		req.Ratio = fmt.Sprintf("%d:%d",
			r.RatioX, r.RatioY)
	}

	req.VideoID = r.OriginalVideo.ID
	req.VideoServiceID = r.OriginalVideo.ServiceID

	return req
}

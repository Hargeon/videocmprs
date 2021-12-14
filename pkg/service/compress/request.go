// Package compress uses for generation request to compress worker, updating
// request and original video in db, creating converted video in db
package compress

import (
	"fmt"

	"github.com/Hargeon/videocmprs/pkg/repository/request"
)

// Request for compress worker
type Request struct {
	RequestID      int64  `json:"request_id"`
	Bitrate        int64  `json:"bitrate"`
	Resolution     string `json:"resolution"`
	Ratio          string `json:"ratio"`
	VideoID        int64  `json:"video_id"`
	UserID         int64  `json:"user_id"`
	VideoServiceID string `json:"video_service_id"`
}

// NewRequest initialize *Request. Need use *request.Resource with *OriginalVideo
func NewRequest(r *request.Resource) *Request {
	req := &Request{
		RequestID: r.ID,
		Bitrate:   r.Bitrate,
		UserID:    r.UserID,
	}

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

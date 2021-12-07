package compress

import (
	"reflect"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name string
		args *request.Resource
		want *Request
	}{
		{
			name: "With resolution",
			args: &request.Resource{
				ID:          1,
				UserID:      1,
				ResolutionX: 800,
				ResolutionY: 600,
				OriginalVideo: &video.Resource{
					ID:        1,
					ServiceID: "test_video",
				},
			},
			want: &Request{
				RequestID:      1,
				UserID:         1,
				Resolution:     "800:600",
				Ratio:          "",
				VideoID:        1,
				VideoServiceID: "test_video",
			},
		},
		{
			name: "With ratio",
			args: &request.Resource{
				ID:     1,
				UserID: 1,
				RatioX: 4,
				RatioY: 3,
				OriginalVideo: &video.Resource{
					ID:        1,
					ServiceID: "test_video",
				},
			},
			want: &Request{
				RequestID:      1,
				UserID:         1,
				Ratio:          "4:3",
				VideoID:        1,
				VideoServiceID: "test_video",
			},
		},
		{
			name: "With bitrate",
			args: &request.Resource{
				ID:      1,
				UserID:  1,
				Bitrate: 64000,
				OriginalVideo: &video.Resource{
					ID:        1,
					ServiceID: "test_video",
				},
			},
			want: &Request{
				RequestID:      1,
				UserID:         1,
				Bitrate:        64000,
				VideoID:        1,
				VideoServiceID: "test_video",
			},
		},
		{
			name: "With all fields",
			args: &request.Resource{
				ID:          1,
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				OriginalVideo: &video.Resource{
					ID:        1,
					ServiceID: "test_video",
				},
			},
			want: &Request{
				RequestID:      1,
				UserID:         1,
				Bitrate:        64000,
				Resolution:     "800:600",
				Ratio:          "4:3",
				VideoID:        1,
				VideoServiceID: "test_video",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequest(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

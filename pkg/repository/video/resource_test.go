package video

import (
	"testing"
)

func TestBuildFields(t *testing.T) {
	cases := []struct {
		name           string
		video          *Resource
		expectedFields map[string]interface{}
	}{
		{
			name:  "With name",
			video: &Resource{Name: "test_video"},
			expectedFields: map[string]interface{}{
				"name": "test_video",
			},
		},
		{
			name:  "With Size",
			video: &Resource{Size: 50000},
			expectedFields: map[string]interface{}{
				"size": 50000,
			},
		},
		{
			name:  "With Bitrate",
			video: &Resource{Bitrate: 64000},
			expectedFields: map[string]interface{}{
				"bitrate": 64000,
			},
		},
		{
			name:  "With ResolutionX",
			video: &Resource{ResolutionX: 600},
			expectedFields: map[string]interface{}{
				"resolution_x": 600,
			},
		},
		{
			name:  "With ResolutionY",
			video: &Resource{ResolutionY: 400},
			expectedFields: map[string]interface{}{
				"resolution_y": 400,
			},
		},
		{
			name:  "With RatioX",
			video: &Resource{RatioX: 4},
			expectedFields: map[string]interface{}{
				"ratio_x": 4,
			},
		},
		{
			name:  "With RatioY",
			video: &Resource{RatioY: 3},
			expectedFields: map[string]interface{}{
				"ratio_y": 3,
			},
		},
		{
			name:  "With ServiceID",
			video: &Resource{ServiceID: "check_video"},
			expectedFields: map[string]interface{}{
				"service_id": "check_video",
			},
		},
		{
			name: "With all fields",
			video: &Resource{
				Name:        "test_video",
				Size:        50000,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				ServiceID:   "check_video",
			},
			expectedFields: map[string]interface{}{
				"name":         "test_video",
				"size":         50000,
				"bitrate":      64000,
				"resolution_x": 800,
				"resolution_y": 600,
				"ratio_x":      4,
				"ratio_y":      3,
				"service_id":   "check_video",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fields := testCase.video.BuildFields()
			if len(fields) != len(testCase.expectedFields) {
				t.Fatalf("Invalid field Number")
			}

			for expectedKey, expectedValue := range testCase.expectedFields {
				var found bool
				for key, value := range fields {
					if key == expectedKey {
						valueInt, ok := value.(int64)
						if ok {
							if int(valueInt) == expectedValue {
								found = true

								break
							}
						} else {
							if expectedValue == value {
								found = true

								break
							}
						}
					}
				}

				if !found {
					t.Errorf("Field %s with value %v in not present in fields",
						expectedKey, expectedValue)
				}
			}
		})
	}
}

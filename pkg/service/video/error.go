package video

import "errors"

var (
	// ErrVideoNotPresent returns if video doesn't exists
	ErrVideoNotPresent = errors.New("video does not exists")
	// ErrVideoAssertion returns if jsonapi.Linkable after Retrieve can't convert to *video.Resource
	ErrVideoAssertion = errors.New("invalid type assertion *video.Resource")
	// ErrVideoNotInCloud returns if video doesn't exists in cloud
	ErrVideoNotInCloud = errors.New("video doesn't exists in cloud")
)

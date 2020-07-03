package media_items

import "github.com/duffpl/google-photos-api-client/common"

type mediaItemsResponse struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

type batchGetMediaItemsResponse struct {
	MediaItemResults []MediaItemWithStatus `json:"mediaItemResults"`
}

type MediaItem struct {
	ID              string          `json:"id"`
	Description     string          `json:"description"`
	ProductURL      string          `json:"productUrl"`
	BaseURL         string          `json:"baseUrl"`
	MimeType        string          `json:"mimeType"`
	MediaMetadata   MediaMetadata   `json:"mediaMetaData"`
	ContributorInfo ContributorInfo `json:"contributorInfo"`
	Filename        string          `json:"filename"`
}

type MediaItemWithStatus struct {
	MediaItem MediaItem        `json:"mediaItem"`
	Status    common.APIStatus `json:"status"`
}

type MediaMetadata struct {
	CreationTime  string         `json:"creationTime"`
	Width         string         `json:"width"`
	Height        string         `json:"height"`
	PhotoMetadata *PhotoMetadata `json:"photo,omitempty"`
	VideoMetadata *VideoMetadata `json:"video,omitempty"`
}

type PhotoMetadata struct {
	CameraMake      string  `json:"cameraMake"`
	CameraModel     string  `json:"cameraModel"`
	FocalLength     float32 `json:"focalLength"`
	ApertureFNumber float32 `json:"apertureFNumber"`
	IsoEquivalent   int     `json:"isoEquivalent"`
	ExposureTime    string  `json:"exposureTime"`
}

type VideoMetadata struct {
	CameraMake  string  `json:"cameraMake"`
	CameraModel string  `json:"cameraModel"`
	Fps         float32 `json:"fps"`
	Status      string  `json:"status"`
}

type ContributorInfo struct {
	ProfilePictureBaseURL string `json:"profilePictureBaseUrl"`
	DisplayName           string `json:"displayName"`
}

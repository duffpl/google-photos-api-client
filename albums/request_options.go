package albums

type AlbumPosition struct {
	Position                 AlbumPositionType `json:"position"`
	RelativeMediaItemId      string            `json:"relativeMediaItemId,omitempty"`
	RelativeEnrichmentItemId string            `json:"relativeEnrichmentItemId,omitempty"`
}

type TextEnrichment struct {
	Text string `json:"text"`
}

type LocationEnrichment struct {
	Location Location `json:"location"`
}

type MapEnrichment struct {
	Origin      Location `json:"origin"`
	Destination Location `json:"destination"`
}

type LatLng struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Location struct {
	LocationName string `json:"locationName"`
	LatLng       LatLng `json:"latLng"`
}

type NewEnrichmentItem struct {
	TextEnrichment     TextEnrichment     `json:"textEnrichment,omitempty"`
	LocationEnrichment LocationEnrichment `json:"locationEnrichment, omitempty"`
	MapEnrichment      MapEnrichment      `json:"mapEnrichment,omitempty"`
}

type SharedAlbumRequestOptions struct {
	IsCollaborative bool `json:"isCollaborative"`
	IsCommentable   bool `json:"isCommentable"`
}

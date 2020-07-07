package albums

type Album struct {
	ID                    string         `json:"id,omitempty"`
	Title                 string         `json:"title"`
	ProductURL            string         `json:"productUrl,omitempty"`
	IsWriteable           bool           `json:"isWriteable,omitempty"`
	ShareInfo             AlbumShareInfo `json:"shareInfo,omitempty"`
	MediaItemsCount       string         `json:"mediaItemsCount,omitempty"`
	CoverPhotoBaseURL     string         `json:"coverPhotoBaseUrl,omitempty"`
	CoverPhotoMediaItemID string         `json:"coverPhotoMediaItemId,omitempty"`
}

type AlbumShareInfo struct {
	SharedAlbumOptions SharedAlbumOptions `json:"sharedAlbumOptions"`
	ShareableURL       string             `json:"shareableUrl"`
	ShareToken         string             `json:"shareToken"`
	IsJoined           bool               `json:"isJoined"`
	IsOwned            bool               `json:"isOwned"`
}

type SharedAlbumOptions struct {
	IsCollaborative bool `json:"isCollaborative"`
	IsCommentable   bool `json:"isCommentable"`
}

type EnrichmentItem struct {
	Id string `json:"id"`
}

package albums

type Album struct {
	ID                    string         `json:"id"`
	Title                 string         `json:"title"`
	ProductURL            string         `json:"productUrl"`
	IsWriteable           bool           `json:"isWriteable"`
	ShareInfo             AlbumShareInfo `json:"shareInfo"`
	MediaItemsCount       string         `json:"mediaItemsCount"`
	CoverPhotoBaseURL     string         `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemID string         `json:"coverPhotoMediaItemId"`
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

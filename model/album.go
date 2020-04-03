package model

type Album struct {
	ID                    string    `json:"id"`
	Title                 string    `json:"title"`
	ProductURL            string    `json:"productUrl"`
	IsWriteable           bool      `json:"isWriteable"`
	ShareInfo             ShareInfo `json:"shareInfo"`
	MediaItemsCount       string    `json:"mediaItemsCount"`
	CoverPhotoBaseURL     string    `json:"coverPhotoBaseUrl"`
	CoverPhotoMediaItemID string    `json:"coverPhotoMediaItemId"`
}
type SharedAlbumOptions struct {
	IsCollaborative bool `json:"isCollaborative"`
	IsCommentable   bool `json:"isCommentable"`
}
type ShareInfo struct {
	SharedAlbumOptions SharedAlbumOptions `json:"sharedAlbumOptions"`
	ShareableURL       string             `json:"shareableUrl"`
	ShareToken         string             `json:"shareToken"`
	IsJoined           bool               `json:"isJoined"`
	IsOwned            bool               `json:"isOwned"`
}

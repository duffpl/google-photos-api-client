package shared_albums

import "github.com/duffpl/google-photos-api-client/albums"

type singleAlbumResponse struct {
	Album albums.Album `json:"album"`
}

type multipleAlbumsResponse struct {
	SharedAlbums  []albums.Album `json:"sharedAlbums"`
	NextPageToken string         `json:"nextPageToken"`
}

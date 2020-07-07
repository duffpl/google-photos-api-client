package google_photos_api_client

import (
	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/duffpl/google-photos-api-client/media_items"
	"github.com/duffpl/google-photos-api-client/uploader"
	"net/http"
)

type ApiClient struct {
	Albums     albums.AlbumsService
	MediaItems media_items.MediaItemsService
}

// Creates new client with all resource services
func NewApiClient(authenticatedClient *http.Client) ApiClient {
	httpUploader := uploader.NewHttpMediaUploader(authenticatedClient)
	return ApiClient{
		Albums:     albums.NewHttpAlbumsService(authenticatedClient),
		MediaItems: media_items.NewHttpMediaItemsService(authenticatedClient, httpUploader),
	}
}

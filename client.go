package google_photos_api_client

import (
	"net/http"
)

type ApiClient struct {
	Albums     AlbumsService
	MediaItems MediaItemsService
}
// Creates new client with all resource services
func NewApiClient(authenticatedClient *http.Client) ApiClient {
	return ApiClient{
		Albums:     NewHttpAlbumsService(authenticatedClient),
		MediaItems: NewHttpMediaItemsService(authenticatedClient),
	}
}

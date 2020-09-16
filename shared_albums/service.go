package shared_albums

import (
	"context"
	"fmt"
	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/duffpl/google-photos-api-client/internal"
	"github.com/imdario/mergo"
	"net/http"
)

// Interface for https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums resource
type SharedAlbumsService interface {
	Get(shareToken string, ctx context.Context) (*albums.Album, error)
	Join(shareToken string, ctx context.Context) (*albums.Album, error)
	Leave(shareToken string, ctx context.Context) error
	List(options *ListOptions, pageToken string, ctx context.Context) (result []albums.Album, nextPageToken string, err error)
	ListAll(options *ListOptions, ctx context.Context) ([]albums.Album, error)
	ListAllAsync(options *ListOptions, ctx context.Context) (<-chan albums.Album, <-chan error)
}

type HttpSharedAlbumsService struct {
	c    *internal.HttpClient
	path string
}

// Fetches album based on specified shareToken
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/get
func (s HttpSharedAlbumsService) Get(shareToken string, ctx context.Context) (*albums.Album, error) {
	responseModel := &albums.Album{}
	err := s.c.FetchWithGet(s.path+"/"+shareToken, nil, responseModel, nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch shared album: %w", err)
	}
	return responseModel, nil

}

// Joins a shared album on behalf of the Google Photos user.
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/join
func (s HttpSharedAlbumsService) Join(shareToken string, ctx context.Context) (*albums.Album, error) {
	responseModel := &singleAlbumResponse{}
	body := shareTokenBody{
		ShareToken: shareToken,
	}
	err := s.c.PostJSON(s.path+":join", nil, body, responseModel, nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot join shared album: %w", err)
	}
	return &responseModel.Album, nil
}

// Leaves a previously-joined shared album on behalf of the Google Photos user. The user must not own this album.
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/leave
func (s HttpSharedAlbumsService) Leave(shareToken string, ctx context.Context) error {
	body := shareTokenBody{
		ShareToken: shareToken,
	}
	err := s.c.PostJSON(s.path+":leave", nil, body, nil, nil, ctx)
	if err != nil {
		return fmt.Errorf("cannot leave shared album: %w", err)
	}
	return nil
}

// Lists all shared albums available in the Sharing tab of the user's Google Photos app.
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/list
func (s HttpSharedAlbumsService) List(options *ListOptions, pageToken string, ctx context.Context) (result []albums.Album, nextPageToken string, err error) {
	requestOptions := ListOptions{
		PageSize: 50,
	}
	if options != nil {
		_ = mergo.Merge(&requestOptions, options, mergo.WithOverride)
	}
	optionsWithToken := struct {
		ListOptions
		PageToken string `url:"pageToken,omitEmpty"`
	}{
		requestOptions,
		pageToken,
	}
	responseModel := &multipleAlbumsResponse{}
	err = s.c.FetchWithGet(s.path, optionsWithToken, responseModel, nil, ctx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot list shared albums: %w", err)
	}
	return responseModel.SharedAlbums, responseModel.NextPageToken, nil
}

// Synchronous wrapper for ListAllAsync
func (s HttpSharedAlbumsService) ListAll(options *ListOptions, ctx context.Context) ([]albums.Album, error) {
	albumsC, errorsC := s.ListAllAsync(options, ctx)
	result := make([]albums.Album, 0)
	for {
		select {
		case item, ok := <-albumsC:
			if !ok {
				return result, nil
			}
			result = append(result, item)
		case err := <-errorsC:
			return nil, err
		}
	}
}

// Asynchronous wrapper for List that takes care of pagination. Returned channel has buffer size of 50
func (s HttpSharedAlbumsService) ListAllAsync(options *ListOptions, ctx context.Context) (<-chan albums.Album, <-chan error) {
	albumsC := make(chan albums.Album, 50)
	errorsC := make(chan error)
	pageToken := ""
	go func() {
		defer close(albumsC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			items, token, err := s.List(options, pageToken, ctx)
			if err != nil {
				errorsC <- err
				return
			}
			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				case albumsC <- item:
				}
			}
			if token == "" {
				return
			}
			pageToken = token
		}
	}()
	return albumsC, errorsC
}

func NewHttpSharedAlbumsService(authenticatedClient *http.Client) HttpSharedAlbumsService {
	return HttpSharedAlbumsService{
		c:    internal.NewHttpClient(authenticatedClient),
		path: "v1/sharedAlbums",
	}
}

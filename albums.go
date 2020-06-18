package google_photos_api_client

import (
	"context"
	"fmt"
	"github.com/duffpl/google-photos-api-client/internal"
	"github.com/duffpl/google-photos-api-client/model"
	"github.com/imdario/mergo"
	"net/http"
)

type AlbumsService interface {
	List(options *AlbumsListOptions, pageToken string, ctx context.Context) (result []model.Album, nextPageToken string, err error)
	ListAll(options *AlbumsListOptions, ctx context.Context) ([]model.Album, error)
	ListAllAsync(options *AlbumsListOptions, ctx context.Context) (<-chan model.Album, <-chan error)
}

type AlbumsListOptions struct {
	PageSize int `url:"pageSize"`
	ExcludeNonAppCreatedData bool `url:"excludeNonAppCreatedData"`
}

type getAlbumsResponse struct {
	Albums        []model.Album `json:"albums"`
	NextPageToken string        `json:"nextPageToken"`
}

type HttpAlbumsService struct {
	c *internal.HttpClient
	path string
}

func (s HttpAlbumsService) List(options *AlbumsListOptions, pageToken string, ctx context.Context) (result []model.Album, nextPageToken string, err error) {
	requestOptions := AlbumsListOptions{
		PageSize: 50,
	}
	if options != nil {
		_ = mergo.Merge(&requestOptions, options, mergo.WithOverride)
	}
	optionsWithToken := struct {
		AlbumsListOptions
		PageToken string `url:"pageToken,omitEmpty"`
	}{
		requestOptions,
		pageToken,
	}
	responseModel := &getAlbumsResponse{}
	err = s.c.FetchWithGet(s.path, optionsWithToken, responseModel, ctx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot process request: %w", err)
	}
	return responseModel.Albums, responseModel.NextPageToken, nil
}

func (s HttpAlbumsService) ListAll(options *AlbumsListOptions, ctx context.Context) ([]model.Album, error) {
	albumsC, errorsC := s.ListAllAsync(options, ctx)
	result := make([]model.Album, 0)
	for {
		select {
		case item, ok := <- albumsC:
			if !ok {
				return result, nil
			}
			result = append(result, item)
		case err := <-errorsC:
			return nil, err
		}
	}
}

func (s HttpAlbumsService) ListAllAsync(options *AlbumsListOptions, ctx context.Context) (<-chan model.Album, <-chan error) {
	albumsC := make(chan model.Album)
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

func NewHttpAlbumsService(authenticatedClient *http.Client) AlbumsService {
	return HttpAlbumsService{
		c:    internal.NewHttpClient(authenticatedClient),
		path: "v1/albums",
	}
}

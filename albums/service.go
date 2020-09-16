package albums

import (
	"context"
	"errors"
	"fmt"
	"github.com/duffpl/google-photos-api-client/internal"
	"github.com/imdario/mergo"
	"math"
	"net/http"
	"net/url"
	"strings"
)

// Interface for https://developers.google.com/photos/library/reference/rest/v1/albums resource
type AlbumsService interface {
	AddEnrichment(albumId string, enrichment NewEnrichmentItem, ctx context.Context) (*EnrichmentItem, error)
	BatchAddMediaItems(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItems(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error
	Create(title string, ctx context.Context) (*Album, error)
	Get(id string, ctx context.Context) (*Album, error)
	List(options *AlbumsListOptions, pageToken string, ctx context.Context) (result []Album, nextPageToken string, err error)
	ListAll(options *AlbumsListOptions, ctx context.Context) ([]Album, error)
	ListAllAsync(options *AlbumsListOptions, ctx context.Context) (<-chan Album, <-chan error)
	Patch(album Album, fieldMask []Field, ctx context.Context) (*Album, error)
	Share(id string, options SharedAlbumOptions, ctx context.Context) (*AlbumShareInfo, error)
	Unshare(id string, ctx context.Context) error
}

type AlbumsListOptions struct {
	PageSize                 int  `url:"pageSize"`
	ExcludeNonAppCreatedData bool `url:"excludeNonAppCreatedData"`
}

type getAlbumsResponse struct {
	Albums        []Album `json:"albums"`
	NextPageToken string  `json:"nextPageToken"`
}

type mediaItemsRequestBody struct {
	MediaItemIds []string `json:"mediaItemIds"`
}

type HttpAlbumsService struct {
	c    *internal.HttpClient
	path string
}

// Adds enrichment item to album specified by id
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/addEnrichment
func (s HttpAlbumsService) AddEnrichment(albumId string, enrichment NewEnrichmentItem, ctx context.Context) (*EnrichmentItem, error) {
	responseModel := &EnrichmentItem{}
	err := s.c.PostJSON(s.path+"/"+albumId+":addEnrichment", nil, enrichment, responseModel, nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot add enrichment: %w", err)
	}
	return responseModel, nil
}

// Removes multiple media items (max 50) from album specified by id
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/batchRemoveMediaItems
func (s HttpAlbumsService) BatchRemoveMediaItems(albumId string, mediaItemIds []string, ctx context.Context) error {
	if len(mediaItemIds) > 50 {
		return errors.New("maximum allowed IDs is 50")
	}
	body := mediaItemsRequestBody{mediaItemIds}
	err := s.c.PostJSON(s.path+"/"+albumId+":batchRemoveMediaItems", nil, body, nil, nil, ctx)
	if err != nil {
		return fmt.Errorf("cannot batch remove media items: %w", err)
	}
	return nil
}

// Removes multiple media items (no limit) using multiple BatchRemoveMediaItems requests
func (s HttpAlbumsService) BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	itemsPerRequest := 50
	requestCount := int(math.Ceil(float64(len(mediaItemIds)) / float64(itemsPerRequest)))
	idsCount := len(mediaItemIds)
	for i := 0; i < requestCount; i++ {
		startOffset := i * itemsPerRequest
		endOffset := startOffset + internal.Min(idsCount-startOffset, itemsPerRequest)
		err := s.BatchRemoveMediaItems(albumId, mediaItemIds[startOffset:endOffset], ctx)
		if err != nil {
			return fmt.Errorf("cannot batch remove all media items: %w", err)
		}
	}
	return nil
}

// Adds multiple media items (no limit) using multiple BatchAddMediaItems requests
func (s HttpAlbumsService) BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error {
	itemsPerRequest := 50
	requestCount := int(math.Ceil(float64(len(mediaItemIds)) / float64(itemsPerRequest)))
	idsCount := len(mediaItemIds)
	for i := 0; i < requestCount; i++ {
		startOffset := i * itemsPerRequest
		endOffset := startOffset + internal.Min(idsCount-startOffset, itemsPerRequest)
		err := s.BatchAddMediaItems(albumId, mediaItemIds[startOffset:endOffset], ctx)
		if err != nil {
			return fmt.Errorf("cannot batch add all media items: %w", err)
		}
	}
	return nil
}

// Adds multiple media items (max 50) to album specified by id
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/batchAddMediaItems
func (s HttpAlbumsService) BatchAddMediaItems(albumId string, mediaItemIds []string, ctx context.Context) error {
	if len(mediaItemIds) > 50 {
		return errors.New("maximum allowed IDs is 50")
	}
	body := mediaItemsRequestBody{mediaItemIds}
	err := s.c.PostJSON(s.path+"/"+albumId+":batchAddMediaItems", nil, body, nil, nil, ctx)
	if err != nil {
		return fmt.Errorf("cannot batch add media items: %w", err)
	}
	return nil
}

// Unshares album specified by id
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/unshare
func (s HttpAlbumsService) Unshare(id string, ctx context.Context) error {
	err := s.c.PostJSON(s.path+"/"+id+":unshare", nil, nil, nil, nil, ctx)
	if err != nil {
		return fmt.Errorf("cannot unshare album: %w", err)
	}
	return nil
}

// Shares album specified by id
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/share
func (s HttpAlbumsService) Share(id string, options SharedAlbumOptions, ctx context.Context) (*AlbumShareInfo, error) {
	responseModel := &AlbumShareInfo{}
	err := s.c.PostJSON(s.path+"/"+id+":share", nil, options, responseModel, nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot share album: %w", err)
	}
	return responseModel, nil
}

// Create new album
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/create
func (s HttpAlbumsService) Create(title string, ctx context.Context) (*Album, error) {
	responseModel := &Album{}
	err := s.c.PostJSON(s.path, nil, createAlbumInput{
		Album: Album{
			Title: title,
		},
	}, responseModel, nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create album: %w", err)
	}
	return responseModel, nil
}

// Fetch album by id
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/get
func (s HttpAlbumsService) Get(id string, ctx context.Context) (*Album, error) {
	responseModel := &Album{}
	err := s.c.FetchWithGet(s.path+"/"+id, nil, responseModel, nil, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch album: %w", err)
	}
	return responseModel, nil
}

// Lists all albums
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/list
func (s HttpAlbumsService) List(options *AlbumsListOptions, pageToken string, ctx context.Context) (result []Album, nextPageToken string, err error) {
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
	err = s.c.FetchWithGet(s.path, optionsWithToken, responseModel, nil, ctx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot process request: %w", err)
	}
	return responseModel.Albums, responseModel.NextPageToken, nil
}

// Synchronous wrapper for ListAllAsync
func (s HttpAlbumsService) ListAll(options *AlbumsListOptions, ctx context.Context) ([]Album, error) {
	albumsC, errorsC := s.ListAllAsync(options, ctx)
	result := make([]Album, 0)
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
func (s HttpAlbumsService) ListAllAsync(options *AlbumsListOptions, ctx context.Context) (<-chan Album, <-chan error) {
	albumsC := make(chan Album, 50)
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

// Patches album. updateMask argument can be used to update only selected fields. Currently only id, title
// and coverPhotoMediaItemId are read
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/albums/patch
func (s HttpAlbumsService) Patch(album Album, updateMask []Field, ctx context.Context) (*Album, error) {
	responseModel := &Album{}
	queryValues := url.Values{}
	if len(updateMask) > 0 {
		fields := []string{}
		for i := range updateMask {
			fields = append(fields, string(updateMask[i]))
		}
		queryValues["updateMask"] = []string{strings.Join(fields, ",")}
	}
	err := s.c.PatchJSON(s.path+"/"+album.ID, queryValues, album, responseModel, nil, ctx)
	if err != nil {
		return nil, err
	}
	return responseModel, nil
}

func NewHttpAlbumsService(authenticatedClient *http.Client) HttpAlbumsService {
	return HttpAlbumsService{
		c:    internal.NewHttpClient(authenticatedClient),
		path: "v1/albums",
	}
}

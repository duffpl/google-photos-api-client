package google_photos_api_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/duffpl/google-photos-api-client/internal"
	"github.com/duffpl/google-photos-api-client/model"
	"github.com/imdario/mergo"
	"net/http"
	"net/url"
)

// Interface for https://developers.google.com/photos/library/reference/rest/v1/mediaItems resource
type MediaItemsService interface {
	Get(itemId string, ctx context.Context) (mediaItem *model.MediaItem, err error)
	List(options *MediaItemsListOptions, pageToken string, ctx context.Context) (mediaItems []model.MediaItem, nextPageToken string, err error)
	ListAll(options *MediaItemsListOptions, ctx context.Context) ([]model.MediaItem, error)
	ListAllAsync(options *MediaItemsListOptions, ctx context.Context) (<-chan model.MediaItem, <-chan error)
	Search(options *MediaItemsSearchOptions, pageToken string, ctx context.Context) (mediaItems []model.MediaItem, nextPageToken string, err error)
	SearchAll(options *MediaItemsSearchOptions, ctx context.Context) ([]model.MediaItem, error)
	SearchAllAsync(options *MediaItemsSearchOptions, ctx context.Context) (<-chan model.MediaItem, <-chan error)
	BatchGetItems(ids []string, ctx context.Context) (mediaItems []model.MediaItemWithStatus, err error)
	BatchGetItemsAll(ids []string, ctx context.Context) ([]model.MediaItemWithStatus, error)
	BatchGetItemsAllAsync(ids []string, ctx context.Context) (<-chan model.MediaItemWithStatus, <-chan error)
}

type MediaItemsListOptions struct {
	PageSize  int    `url:"pageSize"`
}

type MediaItemsSearchOptions struct {
	PageSize  int                     `json:"pageSize"`
	AlbumId   string                  `json:"albumId,omitEmpty"`
	Filters   *MediaItemSearchFilters `json:"filters,omitEmpty"`
}

type MediaItemSearchFilters struct {
	FeatureFilter            *MediaItemFeatureFilter   `json:"featureFilter,omitEmpty"`
	DateFilter               *MediaItemDateFilter      `json:"dateFilter, omitEmpty"`
	ContentFilter            *MediaItemContentFilter   `json:"contentFilter, omitEmpty"`
	MediaTypeFilter          *MediaItemMediaTypeFilter `json:"mediaTypeFilter, omitEmpty"`
	IncludeArchivedMedia     bool                      `json:"includeArchivedMedia"`
	ExcludeNonAppCreatedData bool                      `json:"excludeNonAppCreatedData"`
}

type MediaItemContentFilter struct {
	IncludedContentCategories []ContentCategory `json:"includedContentCategories"`
	ExcludedContentCategories []ContentCategory `json:"excludedContentCategories"`
}

type MediaItemMediaTypeFilter struct {
	MediaTypes []MediaType `json:"mediaTypes"`
}

type MediaItemDateFilter struct {
	Dates  []MediaItemDateFilterDateItem  `json:"dates"`
	Ranges []MediaItemDateFilterRangeItem `json:"ranges"`
}

type MediaItemDateFilterDateItem struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type MediaItemDateFilterRangeItem struct {
	StartDate MediaItemDateFilterDateItem `json:"startDate"`
	EndDate   MediaItemDateFilterDateItem `json:"endDate"`
}

type MediaItemFeatureFilter struct {
	IncludedFeatures []Feature `json:"includedFeatures,omitEmpty"`
}

type batchGetMediaItemsResponse struct {
	MediaItemResults []model.MediaItemWithStatus `json:"mediaItemResults"`
}

type mediaItemsResponse struct {
	MediaItems    []model.MediaItem `json:"mediaItems"`
	NextPageToken string            `json:"nextPageToken"`
}

type HttpMediaItemsService struct {
	c    *internal.HttpClient
	path string
}

// Fetches media item specified by ID
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/get
func (s HttpMediaItemsService) Get(itemId string, ctx context.Context) (mediaItem *model.MediaItem, err error) {
	q := url.Values{"mediaItemId": []string{itemId}}
	responseModel := &model.MediaItem{}
	err = s.c.FetchWithGet(s.path, q, responseModel, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot complete request: %w", err)
	}
	return responseModel, nil
}

// Fetches multiple media items (max 50)
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchGet
func (s HttpMediaItemsService) BatchGetItems(ids []string, ctx context.Context) (mediaItems []model.MediaItemWithStatus, err error) {
	if len(ids) > 50 {
		return nil, errors.New("max 50 ids allowed")
	}
	q := url.Values{"mediaItemIds": ids}
	responseModel := &batchGetMediaItemsResponse{}
	err = s.c.FetchWithGet(s.path+":batchGet", q, responseModel, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot complete request: %w", err)
	}
	return responseModel.MediaItemResults, nil
}

// Synchronous wrapper for BatchGetItemsAllAsync
func (s HttpMediaItemsService) BatchGetItemsAll(ids []string, ctx context.Context) ([]model.MediaItemWithStatus, error) {
	itemsC, errorsC := s.BatchGetItemsAllAsync(ids, ctx)
	result := make([]model.MediaItemWithStatus, 0)
	for {
		select {
		case item, ok := <-itemsC:
			if !ok {
				return result, nil
			}
			result = append(result, item)
		case err := <-errorsC:
			return nil, err
		}
	}
}

// Asynchronous wrapper for BatchGetItems
// Fetches any number of media items in 50 items chunks.
func (s HttpMediaItemsService) BatchGetItemsAllAsync(ids []string, ctx context.Context) (<-chan model.MediaItemWithStatus, <-chan error) {
	itemsC := make(chan model.MediaItemWithStatus)
	errC := make(chan error)
	go func() {
		defer close(itemsC)
		itemsPerChunk := 50
		chunkCount := (len(ids) / itemsPerChunk) + 1
		for i := 0; i < chunkCount; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			sliceStart := i * itemsPerChunk
			sliceEnd := min((i+1)*itemsPerChunk, len(ids))
			idsChunk := ids[sliceStart:sliceEnd]
			mediaItems, err := s.BatchGetItems(idsChunk, ctx)
			if err != nil {
				errC <- err
				return
			}
			for _, item := range mediaItems {
				select {
				case <-ctx.Done():
					return
				case itemsC <- item:
				}
			}
		}
	}()
	return itemsC, errC
}

// Fetches all media items. Default page size is 50
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/list
func (s HttpMediaItemsService) List(options *MediaItemsListOptions, pageToken string, ctx context.Context) (mediaItems []model.MediaItem, nextPageToken string, err error) {
	responseModel := &mediaItemsResponse{}
	requestOptions := MediaItemsListOptions{
		PageSize: 50,
	}
	if options != nil {
		_ = mergo.Merge(&requestOptions, options, mergo.WithOverride)
	}
	optionsWithToken := struct {
		MediaItemsListOptions
		PageToken string `url:"pageToken,omitEmpty"`
	}{
		requestOptions,
		pageToken,
	}
	err = s.c.FetchWithGet(s.path, optionsWithToken, responseModel, ctx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot complete request: %w", err)
	}
	return responseModel.MediaItems, responseModel.NextPageToken, nil
}

// Synchronous wrapper for ListAllAsync
func (s HttpMediaItemsService) ListAll(options *MediaItemsListOptions, ctx context.Context) ([]model.MediaItem, error) {
	itemsC, errorsC := s.ListAllAsync(options, ctx)
	result := make([]model.MediaItem, 0)
	for {
		select {
		case item, ok := <-itemsC:
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
func (s HttpMediaItemsService) ListAllAsync(options *MediaItemsListOptions, ctx context.Context) (<-chan model.MediaItem, <-chan error) {
	itemsC := make(chan model.MediaItem, 50)
	errorsC := make(chan error)
	pageToken := ""
	go func() {
		defer close(itemsC)
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
				case itemsC <- item:
				}
			}
			if token == "" {
				return
			}
			pageToken = token
		}
	}()
	return itemsC, errorsC
}

// Fetches all media items based on search criteria. Default page size is 50
//
// Doc: https://developers.google.com/photos/library/reference/rest/v1/mediaItems/search
func (s HttpMediaItemsService) Search(options *MediaItemsSearchOptions, pageToken string, ctx context.Context) (mediaItems []model.MediaItem, nextPageToken string, err error) {
	requestOptions := MediaItemsSearchOptions{
		PageSize: 50,
	}
	if options != nil {
		_ = mergo.Merge(&requestOptions, options, mergo.WithOverride)
	}
	responseModel := &mediaItemsResponse{}
	optionsWithToken := struct {
		MediaItemsSearchOptions
		PageToken string `json:"pageToken,omitEmpty"`
	}{
		requestOptions,
		pageToken,
	}
	err = s.c.FetchWithPost(s.path+":search", nil, optionsWithToken, responseModel, ctx)
	if err != nil {
		return nil, "", fmt.Errorf("cannot complete request: %w", err)
	}
	return responseModel.MediaItems, responseModel.NextPageToken, nil
}

// Synchronous wrapper for SearchAllAsync
func (s HttpMediaItemsService) SearchAll(options *MediaItemsSearchOptions, ctx context.Context) ([]model.MediaItem, error) {
	itemsC, errorsC := s.SearchAllAsync(options, ctx)
	result := make([]model.MediaItem, 0)
	for {
		select {
		case item, ok := <-itemsC:
			if !ok {
				return result, nil
			}
			result = append(result, item)
		case err := <-errorsC:
			return nil, err
		}
	}
}

// Asynchronous wrapper for Search that takes care of pagination. Returned channel has buffer size of 50
func (s HttpMediaItemsService) SearchAllAsync(options *MediaItemsSearchOptions, ctx context.Context) (<-chan model.MediaItem, <-chan error) {
	itemsC := make(chan model.MediaItem, 50)
	errorsC := make(chan error)
	pageToken := ""
	go func() {
		defer close(itemsC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			items, token, err := s.Search(options, pageToken, ctx)

			if err != nil {
				errorsC <- fmt.Errorf("cannot perform search: %w", err)
				return
			}
			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				case itemsC <- item:
				}
			}
			if token == "" {
				return
			}
			pageToken = token
		}
	}()
	return itemsC, errorsC
}

func NewHttpMediaItemsService(httpClient *http.Client) HttpMediaItemsService {
	return HttpMediaItemsService{
		c:    internal.NewHttpClient(httpClient),
		path: "v1/mediaItems",
	}
}

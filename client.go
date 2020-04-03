package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/duffpl/google-photos-api-client/model"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
)

type Client interface {
	GetAlbums(...RequestOption) (result []model.Album, nextPageToken string, err error)
	GetAllAlbumsC(context.Context) (<-chan model.Album, <-chan error)
	GetAllAlbums() ([]model.Album, error)
	GetMediaItems(...RequestOption) (result []model.MediaItem, nextPageToken string, err error)
	GetAllMediaItemsC(ctx context.Context) (<-chan model.MediaItem, <-chan error)
	SearchMediaItems(params SearchParams) ([]model.MediaItem, string, error)
	SearchAllMediaItemsC(params SearchParams, ctx context.Context) (<-chan model.MediaItem, <-chan error)
	SearchAllMediaItems(params SearchParams) ([]model.MediaItem, error)
	BatchGetItems(ids []string) ([]model.MediaItem, error)
}

type getAlbumsResponse struct {
	Albums        []model.Album `json:"albums"`
	NextPageToken string        `json:"nextPageToken"`
}

type getMediaItemsResponse struct {
	MediaItems    []model.MediaItem `json:"mediaItems"`
	NextPageToken string            `json:"nextPageToken"`
}

type WrappedMediaItem struct {
	MediaItem model.MediaItem `json:"mediaItem"`
}

type batchGetMediaItemsResponse struct {
	MediaItemResults []WrappedMediaItem `json:"mediaItemResults"`
}

type httpClient struct {
	*http.Client
	logrus.StdLogger
	LibraryCallCount int
}

func (h httpClient) batch(ids []string) ([]model.MediaItem, error) {
	q := url.Values{}
	q["mediaItemIds"] = ids
	u := &url.URL{
		Scheme:   "https",
		Opaque:   "",
		Host:     "photoslibrary.googleapis.com",
		Path:     "v1/mediaItems:batchGet",
		RawQuery: q.Encode(),
	}
	response, _ := h.doGetRequest(u.String())
	responseModel := &batchGetMediaItemsResponse{}
	err := json.Unmarshal(response, responseModel)
	if err != nil {
		return nil, err
	}
	mediaItems := make([]model.MediaItem, 0)
	for _, wrappedItem := range responseModel.MediaItemResults {
		mediaItems = append(mediaItems, wrappedItem.MediaItem)
	}
	return mediaItems, nil

}
func (h httpClient) BatchGetItems(ids []string) ([]model.MediaItem, error) {
	pageCount := (len(ids) / 50) + 1
	mediaItems := make([]model.MediaItem, 0)
	for i := 0; i< pageCount;i++ {

		idsSlice := ids[i*50:int(math.Min(float64(i*50+50),float64(len(ids))))]
		result, _ := h.batch(idsSlice)
		mediaItems = append(mediaItems, result...)
	}
	return mediaItems, nil
}

func (h httpClient) SearchAllMediaItems(params SearchParams) ([]model.MediaItem, error) {
	items := make([]model.MediaItem, 0)
	iC, eC := h.SearchAllMediaItemsC(params, context.Background())
	err := func() error {
		for {
			select {
			case item, ok := <-iC:
				if !ok {
					return nil
				}
				items = append(items, item)
			case err := <-eC:
				if err != nil {
					return err
				}
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (h httpClient) GetAllAlbums() ([]model.Album, error) {
	albums := make([]model.Album, 0)
	aC, eC := h.GetAllAlbumsC(context.Background())
	err := func() error {
		for {
			select {
			case album, ok := <-aC:
				if !ok {
					return nil
				}
				albums = append(albums, album)
			case err := <-eC:
				if err != nil {
					return err
				}
			}
		}
	}()
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (h httpClient) GetAlbums(options ...RequestOption) (result []model.Album, nextPageToken string, err error) {
	defaultOptions := []RequestOption{
		PageSizeOption{"50"},
	}
	options = append(defaultOptions, options...)
	q := url.Values{}
	for _, o := range options {
		q.Set(o.getKey(), o.getValue())
	}
	u := &url.URL{
		Scheme:   "https",
		Opaque:   "",
		Host:     "photoslibrary.googleapis.com",
		Path:     "v1/albums",
		RawQuery: q.Encode(),
	}


	response, _ := h.doGetRequest(u.String())
	responseModel := &getAlbumsResponse{}
	err = json.Unmarshal(response, responseModel)
	if err != nil {
		return nil, "", err
	}
	return responseModel.Albums, responseModel.NextPageToken, nil
}

func (h httpClient) GetAllMediaItemsC(ctx context.Context) (<-chan model.MediaItem, <-chan error) {
	itemsC := make(chan model.MediaItem, 50)
	errorsC := make(chan error)
	var opts []RequestOption
	go func() {
		defer close(itemsC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			items, token, err := h.GetMediaItems(opts...)
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
			opts = []RequestOption{PageTokenOption{token}}
		}
	}()
	return itemsC, errorsC
}

func (h httpClient) doGetRequest(url string) ([]byte, error) {
	response, _ := h.Client.Get(url)
	bytes_, _ := ioutil.ReadAll(response.Body)
	return bytes_, nil
}

func (h httpClient) doPostRequest(url string, contentType string, body io.Reader) ([]byte, error) {
	response, _ := h.Client.Post(url, contentType, body)
	bytes_, _ := ioutil.ReadAll(response.Body)
	return bytes_, nil
}

func (h httpClient) GetAllAlbumsC(ctx context.Context) (<-chan model.Album, <-chan error) {

	albumsChan := make(chan model.Album)
	errorsChan := make(chan error)
	var opts []RequestOption
	go func() {
		defer close(albumsChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			items, token, err := h.GetAlbums(opts...)
			if err != nil {
				errorsChan <- err
				return
			}
			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				case albumsChan <- item:
				}
			}
			if token == "" {
				return
			}
			opts = []RequestOption{PageTokenOption{token}}
		}
	}()
	return albumsChan, errorsChan
}

type RequestOption interface {
	getKey() string
	getValue() string
}

type PageSizeOption struct {
	value string
}

func (p PageSizeOption) getKey() string {
	return "pageSize"
}

func (p PageSizeOption) getValue() string {
	return p.value
}

type AlbumIDOption struct {
	value string
}

func (p AlbumIDOption) getKey() string {
	return "albumId"
}

func (p AlbumIDOption) getValue() string {
	return p.value
}

type FiltersOption struct {
	value *filters
}

func (p FiltersOption) getKey() string {
	return "filters"
}

func (p FiltersOption) getValue() string {
	marshalledFilters, _ := json.Marshal(p.value)
	return string(marshalledFilters)
}

type PageTokenOption struct {
	value string
}

func (p PageTokenOption) getKey() string {
	return "pageToken"
}

func (p PageTokenOption) getValue() string {
	return p.value
}

func (h httpClient) GetMediaItems(options ...RequestOption) (result []model.MediaItem, nextPageToken string, err error) {
	defaultOptions := []RequestOption{
		PageSizeOption{"100"},
	}
	options = append(defaultOptions, options...)
	q := url.Values{}
	for _, o := range options {
		q.Set(o.getKey(), o.getValue())
	}
	u := &url.URL{
		Scheme:   "https",
		Opaque:   "",
		Host:     "photoslibrary.googleapis.com",
		Path:     "v1/mediaItems",
		RawQuery: q.Encode(),
	}
	h.LibraryCallCount++
	h.StdLogger.Println("library call count: ", h.LibraryCallCount)
	response, _ := h.doGetRequest(u.String())
	responseModel := &getMediaItemsResponse{}
	err = json.Unmarshal(response, responseModel)
	if err != nil {
		return nil, "", err
	}
	return responseModel.MediaItems, responseModel.NextPageToken, nil
}

func (h httpClient) SearchMediaItems(params SearchParams) ([]model.MediaItem, string, error) {
	defaultParams := SearchParams{
		PageSize: 100,
	}
	err := mergo.Merge(&params, defaultParams, mergo.WithOverride)
	searchBody, _ := json.Marshal(params)
	fmt.Println(string(searchBody))
	response, _ := h.doPostRequest("https://photoslibrary.googleapis.com/v1/mediaItems:search", "application/javascript", bytes.NewReader(searchBody))
	responseModel := &getMediaItemsResponse{}
	err = json.Unmarshal(response, responseModel)
	if err != nil {
		return nil, "", err
	}
	return responseModel.MediaItems, responseModel.NextPageToken, nil
}

type SearchParams struct {
	PageSize  int      `json:"pageSize"`
	AlbumId   string   `json:"albumId,omitEmpty"`
	Filters   *filters `json:"filters,omitEmpty"`
	PageToken string   `json:"pageToken"`
}
type filters struct {
	FeatureFilter *featureFilter `json:"featureFilter,omitEmpty"`
}

type featureFilter struct {
	IncludedFeatures []string `json:"includedFeatures,omitEmpty"`
}

func (h httpClient) SearchAllMediaItemsC(params SearchParams, ctx context.Context) (<-chan model.MediaItem, <-chan error) {
	itemsC := make(chan model.MediaItem)
	errorsC := make(chan error)
	go func() {
		defer close(itemsC)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			items, token, err := h.SearchMediaItems(params)
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
			params.PageToken = token
		}
	}()
	return itemsC, errorsC
}

func NewClient(_httpClient *http.Client) Client {
	l := logrus.New()
	return httpClient{Client: _httpClient, StdLogger: l}
}

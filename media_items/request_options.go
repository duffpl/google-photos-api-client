package media_items

import "github.com/duffpl/google-photos-api-client/albums"

type ListOptions struct {
	PageSize int `url:"pageSize"`
}

type SearchOptions struct {
	PageSize int            `json:"pageSize"`
	AlbumId  string         `json:"albumId,omitempty"`
	Filters  *SearchFilters `json:"filters,omitempty"`
}

type SearchFilters struct {
	FeatureFilter            *FeatureFilter   `json:"featureFilter,omitempty"`
	DateFilter               *DateFilter      `json:"dateFilter,omitempty"`
	ContentFilter            *ContentFilter   `json:"contentFilter,omitempty"`
	MediaTypeFilter          *MediaTypeFilter `json:"mediaTypeFilter,omitempty"`
	IncludeArchivedMedia     bool             `json:"includeArchivedMedia"`
	ExcludeNonAppCreatedData bool             `json:"excludeNonAppCreatedData"`
}

type ContentFilter struct {
	IncludedContentCategories []ContentCategory `json:"includedContentCategories"`
	ExcludedContentCategories []ContentCategory `json:"excludedContentCategories"`
}

type MediaTypeFilter struct {
	MediaTypes []MediaType `json:"mediaTypes"`
}

type DateFilter struct {
	Dates  []DateFilterDateItem  `json:"dates"`
	Ranges []DateFilterRangeItem `json:"ranges"`
}

type DateFilterDateItem struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type DateFilterRangeItem struct {
	StartDate DateFilterDateItem `json:"startDate"`
	EndDate   DateFilterDateItem `json:"endDate"`
}

type FeatureFilter struct {
	IncludedFeatures []Feature `json:"includedFeatures,omitempty"`
}

type BatchCreateOptions struct {
	AlbumId       string               `json:"albumId"`
	AlbumPosition albums.AlbumPosition `json:"albumPosition"`
	NewMediaItems []NewMediaItem       `json:"newMediaItems"`
}

type SimpleMediaItem struct {
	UploadToken string `json:"uploadToken"`
	FileName    string `json:"fileName"`
}

type NewMediaItem struct {
	Description     string          `json:"description"`
	SimpleMediaItem SimpleMediaItem `json:"simpleMediaItem"`
}

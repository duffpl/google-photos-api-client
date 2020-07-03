package media_items

type ListOptions struct {
	PageSize int `url:"pageSize"`
}

type SearchOptions struct {
	PageSize int            `json:"pageSize"`
	AlbumId  string         `json:"albumId,omitEmpty"`
	Filters  *SearchFilters `json:"filters,omitEmpty"`
}

type SearchFilters struct {
	FeatureFilter            *FeatureFilter   `json:"featureFilter,omitEmpty"`
	DateFilter               *DateFilter      `json:"dateFilter, omitEmpty"`
	ContentFilter            *ContentFilter   `json:"contentFilter, omitEmpty"`
	MediaTypeFilter          *MediaTypeFilter `json:"mediaTypeFilter, omitEmpty"`
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
	IncludedFeatures []Feature `json:"includedFeatures,omitEmpty"`
}

package media_items

// Content category used in search filter
type ContentCategory string

const (
	ContentCategoryAnimals      ContentCategory = "ANIMALS"
	ContentCategoryArts         ContentCategory = "ARTS"
	ContentCategoryBirthdays    ContentCategory = "BIRTHDAYS"
	ContentCategoryCityscapes   ContentCategory = "CITYSCAPES"
	ContentCategoryCrafts       ContentCategory = "CRAFTS"
	ContentCategoryDocuments    ContentCategory = "DOCUMENTS"
	ContentCategoryFashion      ContentCategory = "FASHION"
	ContentCategoryFlowers      ContentCategory = "FLOWERS"
	ContentCategoryFood         ContentCategory = "FOOD"
	ContentCategoryGardens      ContentCategory = "GARDENS"
	ContentCategoryHolidays     ContentCategory = "HOLIDAYS"
	ContentCategoryHouses       ContentCategory = "HOUSES"
	ContentCategoryLandmarks    ContentCategory = "LANDMARKS"
	ContentCategoryLandscapes   ContentCategory = "LANDSCAPES"
	ContentCategoryNight        ContentCategory = "NIGHT"
	ContentCategoryPeople       ContentCategory = "PEOPLE"
	ContentCategoryPerformances ContentCategory = "PERFORMANCES"
	ContentCategoryPets         ContentCategory = "PETS"
	ContentCategoryReceipts     ContentCategory = "RECEIPTS"
	ContentCategoryScreenshots  ContentCategory = "SCREENSHOTS"
	ContentCategorySelfies      ContentCategory = "SELFIES"
	ContentCategorySport        ContentCategory = "SPORT"
	ContentCategoryTravel       ContentCategory = "TRAVEL"
	ContentCategoryUtility      ContentCategory = "UTILITY"
	ContentCategoryWeddings     ContentCategory = "WEDDINGS"
	ContentCategoryWhiteboards  ContentCategory = "WHITEBOARDS"
)

// Feature used in search filter
type Feature string

const (
	FeatureNone      Feature = "NONE"
	FeatureFavorites Feature = "FAVORITES"
)

// Media type used in search filter
type MediaType string

const (
	MediaTypeFilterPhoto    MediaType = "PHOTO"
	MediaTypeFilterVideo    MediaType = "VIDEO"
	MediaTypeFilterAllMedia MediaType = "ALL_MEDIA"
)

package albums

type AlbumPositionType string

const (
	AlbumPositionTypeUnspecified         AlbumPositionType = "POSITION_TYPE_UNSPECIFIED"
	AlbumPositionTypeFirstInAlbum        AlbumPositionType = "FIRST_IN_ALBUM"
	AlbumPositionTypeLastInAlbum         AlbumPositionType = "LAST_IN_ALBUM"
	AlbumPositionTypeAfterMediaItem      AlbumPositionType = "AFTER_MEDIA_ITEM"
	AlbumPositionTypeAfterEnrichmentItem AlbumPositionType = "AFTER_ENRICHMENT_ITEM"
)

// Used for updateMask attribute in patch method
type Field string

const (
	AlbumFieldTitle                 Field = "title"
	AlbumFieldCoverPhotoMediaItemId Field = "coverPhotoMediaItemId"
)

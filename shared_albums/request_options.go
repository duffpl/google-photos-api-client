package shared_albums

type ListOptions struct {
	PageSize                 int  `url:"pageSize"`
	ExcludeNonAppCreatedData bool `url:"excludeNonAppCreatedData"`
}

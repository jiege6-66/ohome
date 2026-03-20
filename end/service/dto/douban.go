package dto

type DoubanCategory struct {
	Category string `json:"category"`
	Selected bool   `json:"selected"`
	Type     string `json:"type"`
	Title    string `json:"title"`
}

type DoubanRecentHotResp struct {
	Items      []any            `json:"items"`
	Subjects   []any            `json:"subjects"`
	Categories []DoubanCategory `json:"categories"`
	Total      int              `json:"total"`
	Notice     string           `json:"notice"`
}

type DoubanCategoryMappingItem struct {
	Category string `json:"category"`
	Type     string `json:"type"`
}

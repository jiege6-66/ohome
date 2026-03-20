package quarksearch

import "time"

type Settings struct {
	HTTPProxy      string
	HTTPSProxy     string
	Channels       []string
	EnabledPlugins []string
}

type FilterConfig struct {
	Include []string
	Exclude []string
}

type Request struct {
	Keyword      string
	Channels     []string
	Concurrency  int
	ForceRefresh bool
	ResultType   string
	SourceType   string
	Plugins      []string
	Ext          map[string]any
	Filter       *FilterConfig
}

type Link struct {
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	Password  string    `json:"password"`
	Datetime  time.Time `json:"datetime,omitempty"`
	WorkTitle string    `json:"work_title,omitempty"`
}

type SearchResult struct {
	MessageID string    `json:"message_id"`
	UniqueID  string    `json:"unique_id"`
	Channel   string    `json:"channel"`
	Datetime  time.Time `json:"datetime"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Links     []Link    `json:"links"`
	Images    []string  `json:"images,omitempty"`
}

type MergedLink struct {
	URL      string    `json:"url"`
	Password string    `json:"password"`
	Note     string    `json:"note"`
	Datetime time.Time `json:"datetime"`
	Source   string    `json:"source,omitempty"`
	Images   []string  `json:"images,omitempty"`
}

type SearchResponse struct {
	Total        int                     `json:"total"`
	Results      []SearchResult          `json:"results,omitempty"`
	MergedByType map[string][]MergedLink `json:"merged_by_type,omitempty"`
}

package dto

import "io"

type DoubanImage struct {
	ContentType   string
	ContentLength int64
	Body          io.ReadCloser
	CacheControl  string
	ETag          string
	LastModified  string
}

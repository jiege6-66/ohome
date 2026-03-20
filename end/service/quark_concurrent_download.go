package service

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	// defaults used when config values are missing or invalid
	defaultQuarkConcurrency     = 5
	defaultQuarkPartSizeMB      = 10
	defaultQuarkChunkMaxRetries = 3
)

// getQuarkConcurrency returns the configured concurrent download connections count.
func getQuarkConcurrency() int {
	v := viper.GetInt("quarkFs.concurrency")
	if v <= 0 {
		return defaultQuarkConcurrency
	}
	return v
}

// getQuarkPartSizeBytes returns the configured download chunk size in bytes.
func getQuarkPartSizeBytes() int64 {
	v := viper.GetInt("quarkFs.partSizeMB")
	if v <= 0 {
		return int64(defaultQuarkPartSizeMB) * 1024 * 1024
	}
	return int64(v) * 1024 * 1024
}

// getQuarkChunkMaxRetries returns the configured max retries per chunk.
func getQuarkChunkMaxRetries() int {
	v := viper.GetInt("quarkFs.chunkMaxRetries")
	if v <= 0 {
		return defaultQuarkChunkMaxRetries
	}
	return v
}

// QuarkStreamResult holds the result of a concurrent stream operation.
type QuarkStreamResult struct {
	Body          io.ReadCloser
	ContentLength int64
	ContentType   string
	StatusCode    int
	ContentRange  string
}

type QuarkProxyResponseMeta struct {
	ContentLength int64
	StatusCode    int
	ContentRange  string
}

// cancelReadCloser wraps an io.ReadCloser and calls a cancel function on Close.
type cancelReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *cancelReadCloser) Close() error {
	c.cancel()
	return c.ReadCloser.Close()
}

// parseRangeHeader parses an HTTP Range header and returns the byte offset and length.
// Returns (0, fileSize, nil) for empty or invalid headers (i.e. request the full file).
func parseRangeHeader(rangeHeader string, fileSize int64) (start, length int64, err error) {
	if strings.TrimSpace(rangeHeader) == "" {
		return 0, fileSize, nil
	}

	const prefix = "bytes="
	if !strings.HasPrefix(rangeHeader, prefix) {
		return 0, fileSize, nil
	}

	spec := strings.TrimPrefix(rangeHeader, prefix)
	// Only handle the first range in a multi-range header
	if idx := strings.IndexByte(spec, ','); idx >= 0 {
		spec = spec[:idx]
	}
	spec = strings.TrimSpace(spec)

	// Suffix range: "-500" means last 500 bytes
	if strings.HasPrefix(spec, "-") {
		suffix, parseErr := strconv.ParseInt(spec[1:], 10, 64)
		if parseErr != nil || suffix <= 0 {
			return 0, fileSize, nil
		}
		start = fileSize - suffix
		if start < 0 {
			start = 0
		}
		return start, fileSize - start, nil
	}

	dashIdx := strings.IndexByte(spec, '-')
	if dashIdx < 0 {
		return 0, fileSize, nil
	}

	start, err = strconv.ParseInt(spec[:dashIdx], 10, 64)
	if err != nil || start < 0 {
		return 0, fileSize, nil
	}
	if start >= fileSize {
		return 0, 0, fmt.Errorf("range start %d >= file size %d", start, fileSize)
	}

	endStr := spec[dashIdx+1:]
	if endStr == "" {
		// "bytes=100-" means from 100 to end
		return start, fileSize - start, nil
	}

	end, parseErr := strconv.ParseInt(endStr, 10, 64)
	if parseErr != nil || end < start {
		return 0, fileSize, nil
	}
	if end >= fileSize {
		end = fileSize - 1
	}
	return start, end - start + 1, nil
}

func BuildQuarkProxyResponseMeta(fileSize int64, rangeHeader string) (QuarkProxyResponseMeta, error) {
	if fileSize == 0 {
		return QuarkProxyResponseMeta{
			ContentLength: 0,
			StatusCode:    http.StatusOK,
		}, nil
	}
	start, length, err := parseRangeHeader(rangeHeader, fileSize)
	if err != nil {
		return QuarkProxyResponseMeta{}, err
	}
	if length <= 0 {
		return QuarkProxyResponseMeta{}, fmt.Errorf("invalid range: start=%d length=%d fileSize=%d", start, length, fileSize)
	}

	meta := QuarkProxyResponseMeta{
		ContentLength: length,
		StatusCode:    http.StatusOK,
	}
	if !(start == 0 && length == fileSize) {
		meta.StatusCode = http.StatusPartialContent
		meta.ContentRange = fmt.Sprintf("bytes %d-%d/%d", start, start+length-1, fileSize)
	}
	return meta, nil
}

// openConcurrentStream downloads from the given URL using multiple concurrent
// HTTP connections (similar to openlist's multi-threaded downloader).
// It splits the requested range into chunks and downloads them in parallel,
// then returns a single ordered io.ReadCloser.
func (c *quarkClient) openConcurrentStream(
	ctx context.Context,
	downloadURL string,
	fileSize int64,
	rangeHeader string,
	filename string,
) (*QuarkStreamResult, error) {

	reqStart, reqLength, err := parseRangeHeader(rangeHeader, fileSize)
	if err != nil {
		return nil, err
	}
	if reqLength <= 0 {
		return nil, fmt.Errorf("invalid range: start=%d length=%d fileSize=%d", reqStart, reqLength, fileSize)
	}

	isPartial := !(reqStart == 0 && reqLength == fileSize)

	// Determine content type from filename extension
	contentType := mime.TypeByExtension(path.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	buildResult := func(body io.ReadCloser, ct string) *QuarkStreamResult {
		if ct != "" {
			contentType = ct
		}
		result := &QuarkStreamResult{
			Body:          body,
			ContentLength: reqLength,
			ContentType:   contentType,
		}
		if isPartial {
			result.StatusCode = http.StatusPartialContent
			result.ContentRange = fmt.Sprintf("bytes %d-%d/%d", reqStart, reqStart+reqLength-1, fileSize)
		} else {
			result.StatusCode = http.StatusOK
		}
		return result
	}

	// ---- Small request: single connection is sufficient ----
	partSizeBytes := getQuarkPartSizeBytes()
	if reqLength <= partSizeBytes {
		rh := ""
		if isPartial {
			rh = fmt.Sprintf("bytes=%d-%d", reqStart, reqStart+reqLength-1)
		}
		resp, err := c.openDownloadStream(ctx, downloadURL, rh)
		if err != nil {
			return nil, err
		}
		upstreamCT := resp.Header.Get("Content-Type")
		return buildResult(resp.Body, upstreamCT), nil
	}

	// ---- Large request: concurrent multi-connection download ----
	partSize := partSizeBytes
	numParts := int((reqLength + partSize - 1) / partSize)
	concurrency := getQuarkConcurrency()
	if concurrency > numParts {
		concurrency = numParts
	}

	logQuarkWarnf("[quarkFs:stream] concurrent download start url_len=%d fileSize=%d range=%d-%d parts=%d concurrency=%d",
		len(downloadURL), fileSize, reqStart, reqStart+reqLength-1, numParts, concurrency)

	// chunkData carries either an open response body (for streaming) or an error.
	// Unlike the previous buffered approach, the response body is NOT fully read
	// into memory; instead it is streamed directly to the pipe, which gives the
	// client data immediately (fast TTFB) and reduces memory usage.
	type chunkData struct {
		body io.ReadCloser
		err  error
	}

	// Create result channels for ordered retrieval (one per chunk)
	channels := make([]chan chunkData, numParts)
	for i := range channels {
		channels[i] = make(chan chunkData, 1)
	}

	// Cancellable context for all download goroutines
	dlCtx, dlCancel := context.WithCancel(ctx)

	// downloadChunk opens an HTTP connection for a single chunk.
	// On success it sends the response body (un-read) via channel;
	// the consumer goroutine will stream from it.
	downloadChunk := func(idx int) {
		chunkStart := reqStart + int64(idx)*partSize
		chunkEnd := chunkStart + partSize - 1
		if chunkEnd >= reqStart+reqLength {
			chunkEnd = reqStart + reqLength - 1
		}

		rh := fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)
		var lastErr error

		maxRetries := getQuarkChunkMaxRetries()
		for retry := 0; retry < maxRetries; retry++ {
			if dlCtx.Err() != nil {
				channels[idx] <- chunkData{err: dlCtx.Err()}
				return
			}

			resp, err := c.openDownloadStream(dlCtx, downloadURL, rh)
			if err != nil {
				lastErr = err
				logQuarkWarnf("[quarkFs:stream] chunk %d/%d connect failed (attempt %d): %v", idx, numParts, retry+1, err)
				continue
			}

			// Successfully connected — send the body for streaming
			channels[idx] <- chunkData{body: resp.Body}
			return
		}

		channels[idx] <- chunkData{err: lastErr}
	}

	// Launch initial batch of goroutines (up to concurrency)
	for i := 0; i < concurrency && i < numParts; i++ {
		go downloadChunk(i)
	}

	// Pipe to deliver chunks in order to the caller
	pr, pw := io.Pipe()

	go func() {
		consumedUpTo := -1
		defer func() {
			dlCancel()
			// Drain and close any response bodies that were not consumed
			for i := consumedUpTo + 1; i < numParts; i++ {
				select {
				case cd := <-channels[i]:
					if cd.body != nil {
						cd.body.Close()
					}
				default:
				}
			}
		}()

		for i := 0; i < numParts; i++ {
			// Wait for chunk i or context cancellation
			var cd chunkData
			select {
			case cd = <-channels[i]:
			case <-ctx.Done():
				pw.CloseWithError(ctx.Err())
				return
			}
			consumedUpTo = i

			if cd.err != nil {
				logQuarkWarnf("[quarkFs:stream] chunk %d/%d failed: %v", i, numParts, cd.err)
				pw.CloseWithError(cd.err)
				return
			}

			// Launch next download goroutine (sliding window)
			next := i + concurrency
			if next < numParts {
				go downloadChunk(next)
			}

			// Stream chunk body directly to pipe — data flows to the client
			// immediately without buffering the entire chunk in memory.
			_, copyErr := io.Copy(pw, cd.body)
			cd.body.Close()
			if copyErr != nil {
				logQuarkWarnf("[quarkFs:stream] chunk %d/%d stream failed: %v", i, numParts, copyErr)
				pw.CloseWithError(copyErr)
				return
			}
		}
		logQuarkWarnf("[quarkFs:stream] concurrent download finished parts=%d", numParts)
		pw.Close()
	}()

	// Wrap with cancelReadCloser so all goroutines stop when the caller closes the body
	body := &cancelReadCloser{ReadCloser: pr, cancel: dlCancel}
	return buildResult(body, ""), nil
}

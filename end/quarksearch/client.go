package quarksearch

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func newHTTPClient(settings Settings) *http.Client {
	transport := &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		DialContext: (&net.Dialer{
			Timeout:   20 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	httpProxy := strings.TrimSpace(settings.HTTPProxy)
	httpsProxy := strings.TrimSpace(settings.HTTPSProxy)
	if httpProxy != "" || httpsProxy != "" {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			raw := ""
			switch strings.ToLower(req.URL.Scheme) {
			case "https":
				raw = httpsProxy
				if raw == "" {
					raw = httpProxy
				}
			default:
				raw = httpProxy
			}
			if strings.TrimSpace(raw) == "" {
				return nil, nil
			}
			return url.Parse(raw)
		}
	}

	return &http.Client{
		Timeout:   20 * time.Second,
		Transport: transport,
	}
}

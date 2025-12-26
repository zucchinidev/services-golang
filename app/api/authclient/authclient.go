package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// This provides a default client configuration, but it is recommended
// this is replaced by the user with application specific settings
// using the WithClient function at the time a AuthAPI is constructed.
var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 0,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

type Logger func(ctx context.Context, format string, args ...any)

type Client struct {
	url  string
	log  Logger
	http *http.Client
}

// New construct a Client that can be used to talk with the auth service
func New(url string, log Logger, options ...func(cln *Client)) *Client {
	cln := Client{
		url:  url,
		log:  log,
		http: &defaultClient,
	}

	for _, option := range options {
		option(&cln)
	}

	return &cln
}

// WithClient add custom HTTP client for processing requests.
// It is recommended to not use the default client and provide your own.
func WithClient(http *http.Client) func(cln *Client) {
	return func(cln *Client) {
		cln.http = http
	}
}

func (cln *Client) Authenticate(ctx context.Context, authorization string) (AuthenticateResp, error) {
	endpoint := fmt.Sprintf("%s/auth/authenticate", cln.url)

	header := map[string]string{
		"authorization": authorization,
	}

	var res AuthenticateResp
	if err := cln.rawRequest(ctx, http.MethodGet, endpoint, header, nil, &res); err != nil {
		return AuthenticateResp{}, err
	}

	return res, nil
}

func (cln *Client) Authorize(ctx context.Context, auth Authorize) error {
	endpoint := fmt.Sprintf("%s/auth/authorize", cln.url)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(auth); err != nil {
		return fmt.Errorf("encode auth request: %w", err)
	}

	if err := cln.rawRequest(ctx, http.MethodPost, endpoint, nil, &b, nil); err != nil {
		return fmt.Errorf("authorize request: %w", err)
	}

	return nil
}

func (cln *Client) rawRequest(ctx context.Context, method string, url string, headers map[string]string, r io.Reader, v any) error {
	cln.log(ctx, "authClient rawRequest: started:", "method", method, "url", url)
	defer cln.log(ctx, "authClient rawRequest: completed:", "method", method, url)

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	for k, v := range headers {
		cln.log(ctx, "rawRequest: header:", "key", k, "value", v)
		req.Header.Set(k, v)
	}

	resp, err := cln.http.Do(req)
	if err != nil {
		return fmt.Errorf("http: do: error: %w", err)
	}
	defer resp.Body.Close()

	cln.log(ctx, "rawRequest: client do:", "method", method, "url", url, "status", resp.StatusCode)

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("failed: response: %s, decodeing error: %w", string(data), err)
		}

		return nil
	case http.StatusUnauthorized:
		var e Error
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("failed: response: %s, decodeing error: %w", string(data), err)
		}
		return e
	}

	return nil
}

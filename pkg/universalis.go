package universalis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://universalis.app/api/v2/"
)

var errNonNilContext = errors.New("context must be non-nil")

type Client struct {
	client  *http.Client
	limiter *rate.Limiter

	// Base URL with trailing slash, defaults to "https://universalis.app/api/v2/"
	BaseUrl *url.URL

	common service

	Listings *ListingService
	History  *HistoryService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseUrl, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseUrl: baseUrl}
	c.limiter = rate.NewLimiter(rate.Every(time.Second), 20)
	c.common.client = c
	c.Listings = (*ListingService)(&c.common)
	c.History = (*HistoryService)(&c.common)

	return c
}

func NewClientWithCustomLimiter(httpClient *http.Client, limiter *rate.Limiter) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseUrl, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseUrl: baseUrl}
	c.limiter = limiter
	c.common.client = c
	c.Listings = (*ListingService)(&c.common)
	c.History = (*HistoryService)(&c.common)

	return c
}

func (c *Client) NewRequest(method, urlString string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseUrl.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseUrl)
	}

	u, err := c.BaseUrl.Parse(urlString)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	err := c.limiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	bareResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	resp := newResponse(bareResp)
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

type Response struct {
	*http.Response

	NextPage int
	PrevPage int
	CurPage  int
}

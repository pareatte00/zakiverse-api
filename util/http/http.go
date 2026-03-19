package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/zakiverse/zakiverse-api/core/cst"
	"github.com/zakiverse/zakiverse-api/logger"
)

type Client struct {
	baseUrl *url.URL
	client  *http.Client
}

func New(baseUrl string) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return &Client{
		baseUrl: u,
		client:  client,
	}, nil
}

func MustNew(baseUrl string) *Client {
	client, err := New(baseUrl)
	if err != nil {
		logger.Fatal("Failed to init outbound client", logger.Field(cst.KeyError, err), logger.Field(cst.KeyBaseUrl, baseUrl))
	}

	return client
}

func (c *Client) Disconnect() {
	c.client.CloseIdleConnections()
}

func (c *Client) newRequest(ctx context.Context, method string, p string, body any, bodyType BodyType, headers map[string]string) (*http.Request, error) {
	if bodyType == 0 {
		bodyType = BodyJson
	}

	u := *c.baseUrl
	u.Path = path.Join(c.baseUrl.Path, p)

	var bodyReader io.Reader

	switch bodyType {
	case BodyForm:
		if body != nil {
			switch b := body.(type) {
			case url.Values:
				bodyReader = strings.NewReader(b.Encode())
			case map[string]string:
				values := url.Values{}
				for k, v := range b {
					values.Set(k, v)
				}
				bodyReader = strings.NewReader(values.Encode())
			case map[string]any:
				values := url.Values{}
				for k, v := range b {
					values.Set(k, fmt.Sprint(v))
				}
				bodyReader = strings.NewReader(values.Encode())
			default:
				return nil, fmt.Errorf("form body must be url.Values or map[string]string or map[string]any")
			}
		}
	case BodyJson:
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	switch bodyType {
	case BodyForm:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case BodyJson:
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

type BodyType int

const (
	BodyJson BodyType = iota
	BodyForm
)

type Response struct {
	StatusCode int
	Headers    http.Header
	BodyType   BodyType
}

func (c *Client) do(req *http.Request, out any) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return &Response{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
		}, fmt.Errorf("http %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if out != nil {
		contentType := resp.Header.Get("Content-Type")
		switch v := out.(type) {

		case map[string]any, *map[string]any:
			if err := json.Unmarshal(bodyBytes, out); err != nil {
				return nil, err
			}
		case *string:
			*v = string(bodyBytes)
		case *[]byte:
			*v = bodyBytes
		default:
			if strings.Contains(contentType, "application/json") {
				if err := json.Unmarshal(bodyBytes, out); err != nil {
					return nil, err
				}
			}
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}

type RequestParam struct {
	Path     string
	Header   map[string]string
	Body     any
	BodyType BodyType
}

func (c *Client) Get(ctx context.Context, out any, param RequestParam) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, param.Path, param.Body, param.BodyType, param.Header)
	if err != nil {
		return nil, err
	}
	return c.do(req, out)
}

func (c *Client) Post(ctx context.Context, out any, param RequestParam) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, param.Path, param.Body, param.BodyType, param.Header)
	if err != nil {
		return nil, err
	}
	return c.do(req, out)
}

func (c *Client) Put(ctx context.Context, out any, param RequestParam) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, param.Path, param.Body, param.BodyType, param.Header)
	if err != nil {
		return nil, err
	}
	return c.do(req, out)
}

func (c *Client) Patch(ctx context.Context, out any, param RequestParam) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPatch, param.Path, param.Body, param.BodyType, param.Header)
	if err != nil {
		return nil, err
	}
	return c.do(req, out)
}

func (c *Client) Delete(ctx context.Context, out any, param RequestParam) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, param.Path, param.Body, param.BodyType, param.Header)
	if err != nil {
		return nil, err
	}
	return c.do(req, out)
}

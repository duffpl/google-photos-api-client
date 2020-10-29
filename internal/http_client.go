package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	c *http.Client
}

func NewHttpClient(c *http.Client) *HttpClient {
	return &HttpClient{
		c: c,
	}
}

func (c *HttpClient) FetchWithGet(path string, queryValues interface{}, responseModel interface{}, reqCb func(req *http.Request), ctx context.Context) error {
	req, err := prepareGetRequest(path, queryValues, ctx)
	if err != nil {
		return fmt.Errorf("cannot prepare request: %w", err)
	}
	if reqCb != nil {
		reqCb(req)
	}
	return c.fetchRequest(req, responseModel)
}

func (c *HttpClient) PostJSON(path string, queryValues interface{}, body interface{}, responseModel interface{}, reqCb func(req *http.Request), ctx context.Context) error {
	return c.doJSONRequest(path, queryValues, body, http.MethodPost, responseModel, reqCb, ctx)
}

func (c *HttpClient) PatchJSON(path string, queryValues interface{}, body interface{}, responseModel interface{}, reqCb func(req *http.Request), ctx context.Context) error {
	return c.doJSONRequest(path, queryValues, body, http.MethodPatch, responseModel, reqCb, ctx)
}

func (c *HttpClient) PostFile(path string, queryValues interface{}, file io.Reader, responseModel interface{}, reqCb func(req *http.Request), ctx context.Context) error {
	req, err := prepareFilePostRequest(path, queryValues, file, ctx)
	if err != nil {
		return fmt.Errorf("cannot prepare request: %w", err)
	}
	if reqCb != nil {
		reqCb(req)
	}
	return c.fetchRequest(req, responseModel)
}

func (c *HttpClient) doJSONRequest(path string, queryValues interface{}, body interface{}, method string, responseModel interface{}, reqCb func(req *http.Request), ctx context.Context) error {
	req, err := prepareJsonRequest(path, queryValues, body, method, ctx)
	if err != nil {
		return fmt.Errorf("cannot prepare request: %w", err)
	}
	if reqCb != nil {
		reqCb(req)
	}
	return c.fetchRequest(req, responseModel)
}

func (c *HttpClient) fetchRequest(req *http.Request, responseModel interface{}) error {
	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("cannot fetch response: %w", err)
	}
	err = GetErrorFromResponse(res)
	if err != nil {
		return fmt.Errorf("invalid response: %w", err)
	}
	err = UnmarshalResponse(res, responseModel)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return nil
}

func prepareFilePostRequest(path string, queryValues interface{}, file io.Reader, ctx context.Context) (*http.Request, error) {
	reqUrl, err := prepareRequestURL(path, queryValues)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request url: %w", err)
	}
	return http.NewRequestWithContext(ctx, http.MethodPost, reqUrl.String(), file)

}
func prepareJsonRequest(path string, queryValues interface{}, body interface{}, method string, ctx context.Context) (*http.Request, error) {
	reqUrl, err := prepareRequestURL(path, queryValues)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request url: %w", err)
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal body: %w", err)
	}
	bodyReader := bytes.NewReader(jsonBody)
	return http.NewRequestWithContext(ctx, method, reqUrl.String(), bodyReader)
}

func prepareRequestURL(path string, queryValues interface{}) (*url.URL, error) {
	var err error
	qValues, ok := queryValues.(url.Values)
	if !ok {
		qValues, err = query.Values(queryValues)
		if err != nil {
			return nil, fmt.Errorf("cannot get query values: %w", err)
		}
	}
	reqUrl := url.URL{
		Scheme:   "https",
		Host:     "photoslibrary.googleapis.com",
		Path:     path,
		RawQuery: qValues.Encode(),
	}
	return &reqUrl, nil
}

func prepareGetRequest(path string, queryValues interface{}, ctx context.Context) (*http.Request, error) {
	reqUrl, err := prepareRequestURL(path, queryValues)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request url: %w", err)
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
}

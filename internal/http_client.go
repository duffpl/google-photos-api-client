package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duffpl/google-photos-api-client/common"
	"github.com/google/go-querystring/query"
	"io/ioutil"
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

func (c *HttpClient) FetchWithGet(path string, queryValues interface{}, responseModel interface{}, ctx context.Context) error {
	req, err := prepareGetRequest(path, queryValues, ctx)
	if err != nil {
		return fmt.Errorf("cannot prepare request: %w", err)
	}
	return c.fetchRequest(err, req, responseModel)
}

func (c *HttpClient) FetchWithPost(path string, queryValues interface{}, body interface{}, responseModel interface{}, ctx context.Context) error {
	req, err := preparePostRequest(path, queryValues, body, ctx)
	if err != nil {
		return fmt.Errorf("cannot prepare request: %w", err)
	}
	return c.fetchRequest(err, req, responseModel)
}

func (c *HttpClient) fetchRequest(err error, req *http.Request, responseModel interface{}) error {
	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("cannot fetch response: %w", err)
	}
	errorReturned := false
	if res.StatusCode >= 400 {
		if res.StatusCode == 404 {
			return errors.New("url not found")
		}
		responseModel = &common.ErrorResponse{}
		errorReturned = true
	}
	err = unmarshalResponse(res, responseModel)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response: %w", err)
	}
	if errorReturned {
		return responseModel.(*common.ErrorResponse).Error
	}
	return nil
}

func unmarshalResponse(res *http.Response, dst interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read body bytes: %w", err)
	}
	err = json.Unmarshal(b, dst)
	if err != nil {
		return fmt.Errorf("json unmarshal common: %w", err)
	}
	return nil
}

func preparePostRequest(path string, queryValues interface{}, body interface{}, ctx context.Context) (*http.Request, error) {
	reqUrl, err := prepareRequestURL(path, queryValues)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request url: %w", err)
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal body: %w", err)
	}
	return http.NewRequestWithContext(ctx, http.MethodPost, reqUrl.String(), bytes.NewReader(jsonBody))
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

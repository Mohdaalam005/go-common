package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client interface {
	Post(ctx context.Context, url string, header http.Header, body interface{}, response interface{}) error
	Get(ctx context.Context, url string, header http.Header, response interface{}) error
	Put(ctx context.Context, url string, header http.Header, body interface{}, response interface{}) error
	Delete(ctx context.Context, url string, header http.Header, body interface{}, response interface{}) error
}

type clientImp struct {
	client *http.Client
}

func (c clientImp) Post(ctx context.Context, url string, header http.Header, body interface{}, response interface{}) error {
	return c.doRequest(ctx, http.MethodPost, url, header, body, response)

}

func (c clientImp) Get(ctx context.Context, url string, header http.Header, response interface{}) error {
	return c.doRequest(ctx, http.MethodGet, url, header, nil, response)
}

func (c clientImp) Put(ctx context.Context, url string, header http.Header, body interface{}, response interface{}) error {
	return c.doRequest(ctx, http.MethodPut, url, header, body, response)
}

func (c clientImp) Delete(ctx context.Context, url string, header http.Header, body interface{}, response interface{}) error {
	return c.doRequest(ctx, http.MethodDelete, url, header, body, response)
}
func (c *clientImp) doRequest(ctx context.Context, method string, url string, header http.Header, body interface{}, response interface{}) error {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	var bodyBufferReader io.Reader
	if len(bodyBytes) != 0 {
		bodyBufferReader = bytes.NewBuffer(bodyBytes)
	}

	r, err := http.NewRequestWithContext(ctx, method, url, bodyBufferReader)
	if err != nil {
		log.Println("failed to create request")
		return err
	}

	r.Header = header

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return err
	}

	return nil
}

func NewClient(client *http.Client) Client {
	return &clientImp{
		client: client,
	}
}

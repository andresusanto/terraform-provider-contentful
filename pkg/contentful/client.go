package contentful

import (
	"context"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	client         *http.Client
	baseURL        string
	token          string
	organisationID string
	envID          string

	ContentType IContentTypeService
}

func NewClient(token string, organisationID string, envID string) *Client {
	c := &Client{
		client:         &http.Client{},
		token:          token,
		organisationID: organisationID,
		envID:          envID,
		baseURL:        "https://api.contentful.com",
	}
	c.ContentType = NewContentTypeService(c)

	return c
}

func (c *Client) createRequest(ctx context.Context, method string, path string, version int, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/vnd.contentful.delivery.v1+json")
	req.Header.Set("X-Contentful-User-Agent", "contentful-go/1.0.0")

	if version != 0 {
		req.Header.Set("X-Contentful-Version", strconv.Itoa(version))
	}

	return req, nil
}

func (c *Client) getEnv(env string) string {
	envID := env
	if envID == "" {
		envID = c.envID
	}
	return envID
}

package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type IContentTypeService interface {
	Activate(ctx context.Context, spaceID string, env string, id string, version int) (map[string]interface{}, error)
	Read(ctx context.Context, spaceID string, env string, id string) (map[string]interface{}, error)
	Put(ctx context.Context, spaceID string, env string, id string, version int, body map[string]interface{}) (map[string]interface{}, error)
}

type contentTypeService struct {
	c *Client
}

func NewContentTypeService(c *Client) IContentTypeService {
	return &contentTypeService{c: c}
}

func (s *contentTypeService) Activate(ctx context.Context, spaceID string, env string, id string, version int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/published", spaceID, s.c.getEnv(env), id)
	req, err := s.c.createRequest(ctx, "PUT", path, version, nil)

	if err != nil {
		return nil, err
	}

	res, err := s.c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("contentful-api: received http status code %d when activating content_type\n\n%s", res.StatusCode, string(body))
	}

	body := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *contentTypeService) Read(ctx context.Context, spaceID string, env string, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s", spaceID, s.c.getEnv(env), id)
	req, err := s.c.createRequest(ctx, "GET", path, 0, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("contentful-api: received http status code %d when reading content_type\n\n%s", res.StatusCode, string(body))
	}

	body := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *contentTypeService) Put(ctx context.Context, spaceID string, env string, id string, version int, body map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s", spaceID, s.c.getEnv(env), id)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := s.c.createRequest(ctx, "PUT", path, version, bytes.NewReader(bodyBytes))

	if err != nil {
		return nil, err
	}

	res, err := s.c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("contentful-api: received http status code %d when updating content_type\n\n%s", res.StatusCode, string(body))
	}

	resBody := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

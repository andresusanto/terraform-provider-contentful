package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type IContentTypeEditorService interface {
	Read(ctx context.Context, spaceID string, env string, id string) (map[string]interface{}, error)
	Put(ctx context.Context, spaceID string, env string, id string, version int, body map[string]interface{}) (map[string]interface{}, error)
}

type contentTypeEditorService struct {
	c *Client
}

func NewContentTypeEditorService(c *Client) IContentTypeEditorService {
    return &contentTypeEditorService{c: c}
}

func (s *contentTypeEditorService) Read(ctx context.Context, spaceID string, env string, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/editor_interface", spaceID, s.c.getEnv(env), id)
	res, err := s.c.do(ctx, "GET", path, 0, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("contentful-api: received http status code %d when reading content_type editor interface\n\n%s", res.StatusCode, string(body))
	}

	body := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *contentTypeEditorService) Put(ctx context.Context, spaceID string, env string, id string, version int, body map[string]interface{}) (map[string]interface{}, error) {
    path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/editor_interface", spaceID, s.c.getEnv(env), id)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := s.c.do(ctx, "PUT", path, version, bytes.NewReader(bodyBytes))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("contentful-api: received http status code %d when updating content_type editor interface\n\n%s", res.StatusCode, string(body))
	}

	resBody := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

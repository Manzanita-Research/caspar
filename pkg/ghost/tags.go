package ghost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type tagsResponse struct {
	Tags []Tag     `json:"tags"`
	Meta *metaWrap `json:"meta,omitempty"`
}

func (c *Client) ListTags(params ListParams) ([]Tag, *Pagination, error) {
	data, err := c.Get("/tags/", params.toValues())
	if err != nil {
		return nil, nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil, fmt.Errorf("parsing tags: %w", err)
	}

	var pag *Pagination
	if resp.Meta != nil {
		pag = &resp.Meta.Pagination
	}
	return resp.Tags, pag, nil
}

func (c *Client) GetTag(idOrSlug string, params ListParams) (*Tag, error) {
	var path string
	if IsID(idOrSlug) {
		path = "/tags/" + idOrSlug + "/"
	} else {
		path = "/tags/slug/" + idOrSlug + "/"
	}

	v := url.Values{}
	if params.Fields != "" {
		v.Set("fields", params.Fields)
	}
	if params.Include != "" {
		v.Set("include", params.Include)
	}

	data, err := c.Get(path, v)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing tag: %w", err)
	}
	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("tag not found: %s", idOrSlug)
	}
	return &resp.Tags[0], nil
}

type CreateTagInput struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

func (c *Client) CreateTag(input CreateTagInput) (*Tag, error) {
	body, err := json.Marshal(map[string]any{"tags": []CreateTagInput{input}})
	if err != nil {
		return nil, fmt.Errorf("encoding tag: %w", err)
	}

	data, err := c.Post("/tags/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing created tag: %w", err)
	}
	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("no tag returned after create")
	}
	return &resp.Tags[0], nil
}

type UpdateTagInput struct {
	Name        *string `json:"-"`
	Slug        *string `json:"-"`
	Description *string `json:"-"`
	Visibility  *string `json:"-"`
}

func (c *Client) UpdateTag(id string, input UpdateTagInput) (*Tag, error) {
	tag := map[string]any{}
	if input.Name != nil {
		tag["name"] = *input.Name
	}
	if input.Slug != nil {
		tag["slug"] = *input.Slug
	}
	if input.Description != nil {
		tag["description"] = *input.Description
	}
	if input.Visibility != nil {
		tag["visibility"] = *input.Visibility
	}

	body, err := json.Marshal(map[string]any{"tags": []map[string]any{tag}})
	if err != nil {
		return nil, fmt.Errorf("encoding tag update: %w", err)
	}

	data, err := c.Put("/tags/"+id+"/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing updated tag: %w", err)
	}
	if len(resp.Tags) == 0 {
		return nil, fmt.Errorf("no tag returned after update")
	}
	return &resp.Tags[0], nil
}

func (c *Client) DeleteTag(id string) error {
	_, err := c.Delete("/tags/" + id + "/")
	return err
}

package ghost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type pagesResponse struct {
	Pages []Page    `json:"pages"`
	Meta  *metaWrap `json:"meta,omitempty"`
}

// ListPages returns a list of pages.
func (c *Client) ListPages(params ListParams) ([]Page, *Pagination, error) {
	data, err := c.Get("/pages/", params.toValues())
	if err != nil {
		return nil, nil, err
	}

	var resp pagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil, fmt.Errorf("parsing pages: %w", err)
	}

	var pag *Pagination
	if resp.Meta != nil {
		pag = &resp.Meta.Pagination
	}
	return resp.Pages, pag, nil
}

// GetPage fetches a single page by ID or slug.
func (c *Client) GetPage(idOrSlug string, params ListParams) (*Page, error) {
	var path string
	if IsID(idOrSlug) {
		path = "/pages/" + idOrSlug + "/"
	} else {
		path = "/pages/slug/" + idOrSlug + "/"
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

	var resp pagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing page: %w", err)
	}
	if len(resp.Pages) == 0 {
		return nil, fmt.Errorf("page not found: %s", idOrSlug)
	}
	return &resp.Pages[0], nil
}

// CreatePage creates a new page.
func (c *Client) CreatePage(input CreatePostInput, useHTML bool) (*Page, error) {
	post := map[string]any{
		"title": input.Title,
	}

	if input.HTML != "" {
		post["html"] = input.HTML
	}
	if input.Lexical != "" {
		post["lexical"] = input.Lexical
	}
	if input.Status != "" {
		post["status"] = input.Status
	}
	if input.Slug != "" {
		post["slug"] = input.Slug
	}
	if input.Featured {
		post["featured"] = true
	}
	if input.PublishedAt != "" {
		post["published_at"] = input.PublishedAt
	}
	if len(input.Tags) > 0 {
		tags := make([]map[string]string, len(input.Tags))
		for i, t := range input.Tags {
			tags[i] = map[string]string{"name": t}
		}
		post["tags"] = tags
	}

	body, err := json.Marshal(map[string]any{"pages": []map[string]any{post}})
	if err != nil {
		return nil, fmt.Errorf("encoding page: %w", err)
	}

	path := "/pages/"
	if useHTML {
		path += "?source=html"
	}

	data, err := c.Post(path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp pagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing created page: %w", err)
	}
	if len(resp.Pages) == 0 {
		return nil, fmt.Errorf("no page returned after create")
	}
	return &resp.Pages[0], nil
}

// UpdatePage updates an existing page.
func (c *Client) UpdatePage(id string, input UpdatePostInput, useHTML bool) (*Page, error) {
	page := map[string]any{
		"updated_at": input.UpdatedAt,
	}

	if input.Title != nil {
		page["title"] = *input.Title
	}
	if input.HTML != nil {
		page["html"] = *input.HTML
	}
	if input.Lexical != nil {
		page["lexical"] = *input.Lexical
	}
	if input.Status != nil {
		page["status"] = *input.Status
	}
	if input.Slug != nil {
		page["slug"] = *input.Slug
	}
	if input.Featured != nil {
		page["featured"] = *input.Featured
	}
	if input.PublishedAt != nil {
		page["published_at"] = *input.PublishedAt
	}
	if len(input.Tags) > 0 {
		tags := make([]map[string]string, len(input.Tags))
		for i, t := range input.Tags {
			tags[i] = map[string]string{"name": t}
		}
		page["tags"] = tags
	}

	body, err := json.Marshal(map[string]any{"pages": []map[string]any{page}})
	if err != nil {
		return nil, fmt.Errorf("encoding page update: %w", err)
	}

	path := "/pages/" + id + "/"
	if useHTML {
		path += "?source=html"
	}

	data, err := c.Put(path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp pagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing updated page: %w", err)
	}
	if len(resp.Pages) == 0 {
		return nil, fmt.Errorf("no page returned after update")
	}
	return &resp.Pages[0], nil
}

// DeletePage deletes a page by ID.
func (c *Client) DeletePage(id string) error {
	_, err := c.Delete("/pages/" + id + "/")
	return err
}

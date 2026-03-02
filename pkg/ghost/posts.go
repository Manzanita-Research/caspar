package ghost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{24}$`)

// IsID returns true if the string looks like a Ghost object ID (24-char hex).
func IsID(s string) bool {
	return uuidRegex.MatchString(s)
}

type postsResponse struct {
	Posts []Post     `json:"posts"`
	Meta  *metaWrap  `json:"meta,omitempty"`
}

type metaWrap struct {
	Pagination Pagination `json:"pagination"`
}

func (p ListParams) toValues() url.Values {
	v := url.Values{}
	if p.Limit > 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Page > 0 {
		v.Set("page", strconv.Itoa(p.Page))
	}
	if p.Filter != "" {
		v.Set("filter", p.Filter)
	}
	if p.Order != "" {
		v.Set("order", p.Order)
	}
	if p.Fields != "" {
		v.Set("fields", p.Fields)
	}
	if p.Include != "" {
		v.Set("include", p.Include)
	}
	if p.Formats != "" {
		v.Set("formats", p.Formats)
	}
	return v
}

// ListPosts returns a list of posts.
func (c *Client) ListPosts(params ListParams) ([]Post, *Pagination, error) {
	data, err := c.Get("/posts/", params.toValues())
	if err != nil {
		return nil, nil, err
	}

	var resp postsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil, fmt.Errorf("parsing posts: %w", err)
	}

	var pag *Pagination
	if resp.Meta != nil {
		pag = &resp.Meta.Pagination
	}
	return resp.Posts, pag, nil
}

// GetPostsByIDs fetches multiple posts by ID in a single request.
func (c *Client) GetPostsByIDs(ids []string, params ListParams) ([]Post, error) {
	params.Filter = "id:[" + strings.Join(ids, ",") + "]"
	params.Limit = len(ids)
	posts, _, err := c.ListPosts(params)
	return posts, err
}

// GetPost fetches a single post by ID or slug.
func (c *Client) GetPost(idOrSlug string, params ListParams) (*Post, error) {
	var path string
	if IsID(idOrSlug) {
		path = "/posts/" + idOrSlug + "/"
	} else {
		path = "/posts/slug/" + idOrSlug + "/"
	}

	v := url.Values{}
	if params.Fields != "" {
		v.Set("fields", params.Fields)
	}
	if params.Include != "" {
		v.Set("include", params.Include)
	}
	if params.Formats != "" {
		v.Set("formats", params.Formats)
	}

	data, err := c.Get(path, v)
	if err != nil {
		return nil, err
	}

	var resp postsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing post: %w", err)
	}
	if len(resp.Posts) == 0 {
		return nil, fmt.Errorf("post not found: %s", idOrSlug)
	}
	return &resp.Posts[0], nil
}

// CreatePostInput holds fields for creating a post.
type CreatePostInput struct {
	Title       string   `json:"title"`
	HTML        string   `json:"html,omitempty"`
	Lexical     string   `json:"lexical,omitempty"`
	Status      string   `json:"status,omitempty"`
	Slug        string   `json:"slug,omitempty"`
	Tags        []string `json:"-"`
	Featured    bool     `json:"featured,omitempty"`
	PublishedAt string   `json:"-"` // ISO 8601 datetime
	Visibility  string   `json:"-"`
}

type createPostRequest struct {
	Posts []map[string]any `json:"posts"`
}

// CreatePost creates a new post.
func (c *Client) CreatePost(input CreatePostInput, useHTML bool) (*Post, error) {
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
	if input.Visibility != "" {
		post["visibility"] = input.Visibility
	}
	if len(input.Tags) > 0 {
		tags := make([]map[string]string, len(input.Tags))
		for i, t := range input.Tags {
			tags[i] = map[string]string{"name": t}
		}
		post["tags"] = tags
	}

	body, err := json.Marshal(createPostRequest{Posts: []map[string]any{post}})
	if err != nil {
		return nil, fmt.Errorf("encoding post: %w", err)
	}

	path := "/posts/"
	if useHTML {
		path += "?source=html"
	}

	data, err := c.Post(path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp postsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing created post: %w", err)
	}
	if len(resp.Posts) == 0 {
		return nil, fmt.Errorf("no post returned after create")
	}
	return &resp.Posts[0], nil
}

// UpdatePostInput holds fields for updating a post.
type UpdatePostInput struct {
	Title     *string  `json:"-"`
	HTML      *string  `json:"-"`
	Lexical   *string  `json:"-"`
	Status    *string  `json:"-"`
	Slug      *string  `json:"-"`
	Tags        []string `json:"-"`
	Featured      *bool    `json:"-"`
	PublishedAt   *string  `json:"-"` // ISO 8601 datetime
	Visibility    *string  `json:"-"`
	CustomExcerpt *string  `json:"-"`
	UpdatedAt     string   `json:"-"` // required for conflict resolution
}

// UpdatePost updates an existing post. Requires the current updated_at for conflict detection.
func (c *Client) UpdatePost(id string, input UpdatePostInput, useHTML bool) (*Post, error) {
	post := map[string]any{
		"updated_at": input.UpdatedAt,
	}

	if input.Title != nil {
		post["title"] = *input.Title
	}
	if input.HTML != nil {
		post["html"] = *input.HTML
	}
	if input.Lexical != nil {
		post["lexical"] = *input.Lexical
	}
	if input.Status != nil {
		post["status"] = *input.Status
	}
	if input.Slug != nil {
		post["slug"] = *input.Slug
	}
	if input.Featured != nil {
		post["featured"] = *input.Featured
	}
	if input.PublishedAt != nil {
		post["published_at"] = *input.PublishedAt
	}
	if input.Visibility != nil {
		post["visibility"] = *input.Visibility
	}
	if input.CustomExcerpt != nil {
		post["custom_excerpt"] = *input.CustomExcerpt
	}
	if len(input.Tags) > 0 {
		tags := make([]map[string]string, len(input.Tags))
		for i, t := range input.Tags {
			tags[i] = map[string]string{"name": t}
		}
		post["tags"] = tags
	}

	body, err := json.Marshal(createPostRequest{Posts: []map[string]any{post}})
	if err != nil {
		return nil, fmt.Errorf("encoding post update: %w", err)
	}

	path := "/posts/" + id + "/"
	if useHTML {
		path += "?source=html"
	}

	data, err := c.Put(path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp postsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing updated post: %w", err)
	}
	if len(resp.Posts) == 0 {
		return nil, fmt.Errorf("no post returned after update")
	}
	return &resp.Posts[0], nil
}

// DeletePost deletes a post by ID.
func (c *Client) DeletePost(id string) error {
	_, err := c.Delete("/posts/" + id + "/")
	return err
}

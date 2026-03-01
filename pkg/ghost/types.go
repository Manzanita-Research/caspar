package ghost

import "time"

// Post represents a Ghost post.
type Post struct {
	ID          string     `json:"id"`
	UUID        string     `json:"uuid,omitempty"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	HTML        string     `json:"html,omitempty"`
	Lexical     string     `json:"lexical,omitempty"`
	Status      string     `json:"status"`
	Featured    bool       `json:"featured"`
	Excerpt     string     `json:"excerpt,omitempty"`
	CustomExcerpt string   `json:"custom_excerpt,omitempty"`
	FeatureImage  string   `json:"feature_image,omitempty"`
	URL         string     `json:"url,omitempty"`
	Tags        []Tag      `json:"tags,omitempty"`
	Authors     []Author   `json:"authors,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// Page has the same shape as a Post in the Ghost API.
type Page = Post

// Tag represents a Ghost tag.
type Tag struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

// Author represents a Ghost author/user.
type Author struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Slug  string `json:"slug,omitempty"`
	Email string `json:"email,omitempty"`
}

// Member represents a Ghost member.
type Member struct {
	ID               string     `json:"id"`
	UUID             string     `json:"uuid,omitempty"`
	Email            string     `json:"email"`
	Name             string     `json:"name,omitempty"`
	Status           string     `json:"status,omitempty"`
	Labels           []Label    `json:"labels,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// Label represents a member label.
type Label struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// Newsletter represents a Ghost newsletter.
type Newsletter struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

// Pagination holds Ghost API pagination metadata.
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
	Total int `json:"total"`
	Next  *int `json:"next"`
	Prev  *int `json:"prev"`
}

// ImageUpload represents the response from uploading an image.
type ImageUpload struct {
	URL string `json:"url"`
}

// ListParams holds common query parameters for list endpoints.
type ListParams struct {
	Limit   int
	Page    int
	Filter  string
	Order   string
	Fields  string
	Include string
}

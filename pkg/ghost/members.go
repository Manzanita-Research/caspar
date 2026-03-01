package ghost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type membersResponse struct {
	Members []Member  `json:"members"`
	Meta    *metaWrap `json:"meta,omitempty"`
}

func (c *Client) ListMembers(params ListParams) ([]Member, *Pagination, error) {
	data, err := c.Get("/members/", params.toValues())
	if err != nil {
		return nil, nil, err
	}

	var resp membersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil, fmt.Errorf("parsing members: %w", err)
	}

	var pag *Pagination
	if resp.Meta != nil {
		pag = &resp.Meta.Pagination
	}
	return resp.Members, pag, nil
}

func (c *Client) GetMember(id string, params ListParams) (*Member, error) {
	v := url.Values{}
	if params.Fields != "" {
		v.Set("fields", params.Fields)
	}

	data, err := c.Get("/members/"+id+"/", v)
	if err != nil {
		return nil, err
	}

	var resp membersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing member: %w", err)
	}
	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("member not found: %s", id)
	}
	return &resp.Members[0], nil
}

type CreateMemberInput struct {
	Email  string   `json:"email"`
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"-"`
}

func (c *Client) CreateMember(input CreateMemberInput) (*Member, error) {
	member := map[string]any{
		"email": input.Email,
	}
	if input.Name != "" {
		member["name"] = input.Name
	}
	if len(input.Labels) > 0 {
		labels := make([]map[string]string, len(input.Labels))
		for i, l := range input.Labels {
			labels[i] = map[string]string{"name": l}
		}
		member["labels"] = labels
	}

	body, err := json.Marshal(map[string]any{"members": []map[string]any{member}})
	if err != nil {
		return nil, fmt.Errorf("encoding member: %w", err)
	}

	data, err := c.Post("/members/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp membersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing created member: %w", err)
	}
	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("no member returned after create")
	}
	return &resp.Members[0], nil
}

type UpdateMemberInput struct {
	Email  *string  `json:"-"`
	Name   *string  `json:"-"`
	Labels []string `json:"-"`
}

func (c *Client) UpdateMember(id string, input UpdateMemberInput) (*Member, error) {
	member := map[string]any{}
	if input.Email != nil {
		member["email"] = *input.Email
	}
	if input.Name != nil {
		member["name"] = *input.Name
	}
	if len(input.Labels) > 0 {
		labels := make([]map[string]string, len(input.Labels))
		for i, l := range input.Labels {
			labels[i] = map[string]string{"name": l}
		}
		member["labels"] = labels
	}

	body, err := json.Marshal(map[string]any{"members": []map[string]any{member}})
	if err != nil {
		return nil, fmt.Errorf("encoding member update: %w", err)
	}

	data, err := c.Put("/members/"+id+"/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var resp membersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing updated member: %w", err)
	}
	if len(resp.Members) == 0 {
		return nil, fmt.Errorf("no member returned after update")
	}
	return &resp.Members[0], nil
}

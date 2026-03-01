package ghost

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type newslettersResponse struct {
	Newsletters []Newsletter `json:"newsletters"`
	Meta        *metaWrap    `json:"meta,omitempty"`
}

func (c *Client) ListNewsletters(params ListParams) ([]Newsletter, *Pagination, error) {
	data, err := c.Get("/newsletters/", params.toValues())
	if err != nil {
		return nil, nil, err
	}

	var resp newslettersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil, fmt.Errorf("parsing newsletters: %w", err)
	}

	var pag *Pagination
	if resp.Meta != nil {
		pag = &resp.Meta.Pagination
	}
	return resp.Newsletters, pag, nil
}

func (c *Client) GetNewsletter(id string, params ListParams) (*Newsletter, error) {
	v := url.Values{}
	if params.Fields != "" {
		v.Set("fields", params.Fields)
	}

	data, err := c.Get("/newsletters/"+id+"/", v)
	if err != nil {
		return nil, err
	}

	var resp newslettersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing newsletter: %w", err)
	}
	if len(resp.Newsletters) == 0 {
		return nil, fmt.Errorf("newsletter not found: %s", id)
	}
	return &resp.Newsletters[0], nil
}

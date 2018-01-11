package fastly

import (
	"fmt"
	"sort"
	"time"
)

// VCLSnippetType defines the location in generated VCL where the snippet should be placed.
type VCLSnippetType string

const (
	// Place the snippet above all subroutines.
	VCLSnippetTypeInit VCLSnippetType = "init"

	// Place the snippet within the vcl_recv subroutine but below the boilerplate VCL and above any objects.
	VCLSnippetTypeRecv VCLSnippetType = "recv"

	// Place the snippet within the vcl_hit subroutine.
	VCLSnippetTypeHit VCLSnippetType = "hit"

	// Place the snippet within the vcl_miss subroutine.
	VCLSnippetTypeMiss VCLSnippetType = "miss"

	// Place the snippet within the vcl_pass subroutine.
	VCLSnippetTypePass VCLSnippetType = "pass"

	// Place the snippet within the vcl_fetch subroutine.
	VCLSnippetTypeFetch VCLSnippetType = "fetch"

	// Place the snippet within the vcl_error subroutine.
	VCLSnippetTypeError VCLSnippetType = "error"

	// Place the snippet within the vcl_deliver subroutine.
	VCLSnippetTypeDeliver VCLSnippetType = "deliver"

	// Place the snippet within the vcl_log subroutine.
	VCLSnippetTypeLog VCLSnippetType = "log"

	// Don't render the snippet in VCL so it can be manually included in custom VCL.
	VCLSnippetTypeNone VCLSnippetType = "none"
)

// VCLSnippet represents a VCL Snippet response from the Fastly API.
type VCLSnippet struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Content   string         `mapstructure:"content"`
	Dynamic   bool           `mapstructure:"dynamic"`
	ID        string         `mapstructure:"id"`
	Name      string         `mapstructure:"name"`
	Priority  uint           `mapstructure:"priority"`
	Type      VCLSnippetType `mapstructure:"type"`
	CreatedAt *time.Time     `mapstructure:"created_at"`
	UpdatedAt *time.Time     `mapstructure:"updated_at"`
	DeletedAt *time.Time     `mapstructure:"deleted_at"`
}

// snippetsByName is a sortable list of VCL Snippets.
type snippetsByName []*VCLSnippet

// Len, Swap, and Less implement the sortable interface.
func (s snippetsByName) Len() int      { return len(s) }
func (s snippetsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s snippetsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListVCLSnippetsInput is used as input to the ListVCLSnippets function.
type ListVCLSnippetsInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListVCLSnippets returns the list of VCL Snippets for the configuration version.
func (c *Client) ListVCLSnippets(i *ListVCLSnippetsInput) ([]*VCLSnippet, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var bs []*VCLSnippet
	if err := decodeJSON(&bs, resp.Body); err != nil {
		return nil, err
	}
	sort.Stable(snippetsByName(bs))
	return bs, nil
}

// CreateVCLSnippetInput is used as input to the CreateVCLSnippet function.
type CreateVCLSnippetInput struct {
	ServiceID string `form:"service_id"`
	Version   int    `form:"version"`

	Content   string         `form:"content"`
	Dynamic   *Compatibool           `form:"dynamic"`
	Name      string         `form:"name"`
	Priority  uint           `form:"priority"`
	Type      VCLSnippetType `form:"type"`
}

// CreateVCLSnippet creates a new snippet.
func (c *Client) CreateVCLSnippet(i *CreateVCLSnippetInput) (*VCLSnippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet", i.ServiceID, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var v *VCLSnippet
	if err := decodeJSON(&v, resp.Body); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateVCLSnippetInput is used as input to the UpdateVCLSnippet function.
type UpdateVCLSnippetInput struct {
	ServiceID string `form:"service_id"`
	Version   int    `form:"version"`

	Content   string         `form:"content"`
	Name      string         `form:"name"`
	Priority  uint           `form:"priority"`
	Type      VCLSnippetType `form:"type"`
}

// UpdateVCLSnippet updates a specific snippet.
func (c *Client) UpdateVCLSnippet(i *UpdateVCLSnippetInput) (*VCLSnippet, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.ServiceID, i.Version, i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var v *VCLSnippet
	if err := decodeJSON(&v, resp.Body); err != nil {
		return nil, err
	}
	return v, nil
}

// DeleteVCLSnippetInput is the input parameter to DeleteVCLSnippet.
type DeleteVCLSnippetInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the snippet to delete (required).
	Name string
}

// DeleteVCLSnippet deletes the VCL snippet with the given name.
func (c *Client) DeleteVCLSnippet(i *DeleteVCLSnippetInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/snippet/%s", i.Service, i.Version, i.Name)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeJSON(&r, resp.Body); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("Not Ok: %s", r.Msg)
	}
	return nil
}

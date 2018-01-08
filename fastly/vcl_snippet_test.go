package fastly

import (
	"testing"
)

func TestClient_VCLSnippets(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "vcl_snippets/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// List
	var ss []*VCLSnippet
	record(t, "vcl_snippets/list", func(c *Client) {
		ss, err = c.ListVCLSnippets(&ListVCLSnippetsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("Expected to receive at least one snippet: %v", ss)
	}
}

func TestClient_ListVCLSnippet_validation(t *testing.T) {
	var err error
	_, err = testClient.ListVCLSnippets(&ListVCLSnippetsInput{})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}
}

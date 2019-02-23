package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"intel/isecl/tdservice/types"
	"net/http"
	"net/url"
)

// Client for consuming the TDService API
type Client struct {
	BaseURL string
	// Username used to authenticate with the KMS. Username is only used for obtaining an authorization token, which is automatically used for requests.
	Username string
	// Password to supply for the Username
	Password string
	// A reference to the underlying http Client.
	// If the value is nil, a default client will be created and used.
	HTTPClient *http.Client
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	return c.HTTPClient
}

func (c *Client) resolvePath(path string) (string, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	rel := baseURL.ResolveReference(&url.URL{Path: path})
	return rel.String(), err
}

func (c *Client) dispatchRequest(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.Username, c.Password)
	return c.httpClient().Do(req)
}

func (c *Client) AddHost(h types.HostInfo) (*types.Host, error) {
	hosts, err := c.resolvePath("/hosts")
	if err != nil {
		return nil, err
	}

	// to fix: you can stream the marshalling directly to HTTP using a buffer
	payload, err := json.Marshal(&h)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(http.MethodPost, hosts, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create host: HTTP Code: %d", rsp.StatusCode)
	}
	// parse it and return
	var created types.Host
	err = json.NewDecoder(rsp.Body).Decode(&created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) AddReportRaw(report string) (*types.Report, error) {
	reports, err := c.resolvePath("/reports")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, reports, bytes.NewBufferString(report))
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.dispatchRequest(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusCreated {
		return nil, err
	}
	var created types.Report
	err = json.NewDecoder(rsp.Body).Decode(&created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) AddReport(r types.Report) (*types.Report, error) {
	// fix this as well
	payload, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	return c.AddReportRaw(string(payload))
}

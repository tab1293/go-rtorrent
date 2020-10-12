package xmlrpc

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Client implements a basic XMLRPC client
type Client struct {
	addr       string
	httpClient *http.Client
}

// NewClient returns a new instance of Client
// Pass in a true value for `insecure` to turn off certificate verification
func NewClient(addr string, insecure bool) *Client {
	transport := &http.Transport{}
	if insecure {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	httpClient := &http.Client{Transport: transport}

	return &Client{
		addr:       addr,
		httpClient: httpClient,
	}
}

// NewClientWithHTTPClient returns a new instance of Client.
// This allows you to use a custom http.Client setup for your needs.
func NewClientWithHTTPClient(addr string, client *http.Client) *Client {
	return &Client{
		addr:       addr,
		httpClient: client,
	}
}

// Call calls the method with "name" with the given args
// Returns the result, and an error for communication errors
func (c *Client) Call(name string, args ...interface{}) (interface{}, error) {
	req := bytes.NewBuffer(nil)
	if err := Marshal(req, name, args...); err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}
	fmt.Printf("req %s\n", req)
	resp, err := c.httpClient.Post(c.addr, "text/xml", req)
	if err != nil {
		return nil, errors.Wrap(err, "POST failed")
	}
	defer resp.Body.Close()

	_, val, fault, err := Unmarshal(resp.Body)
	if fault != nil {
		err = errors.Errorf("Error: %v: %v", err, fault)
	}
	return val, err
}

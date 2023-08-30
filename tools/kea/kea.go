package kea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// nolint: gosec
const (
	packageName = "kea"
	envKEAUSER  = "KEA_USERNAME"
	envKEAPASS  = "KEA_PASSWORD"
)

type (
	// Client : Stored memory objects for the CradlePoint client.
	Client struct {
		client *http.Client
		log    *logrus.Logger
		auth   auth
		remote string
	}
	// Response : Similar response returned for all Kea queries.
	Response struct {
		Result    int             `json:"result"`
		Text      string          `json:"text"`
		Arguments json.RawMessage `json:"arguments"`
	}

	// Request : Generic Kea request payload.
	Request struct {
		Command   string         `json:"command"`
		Arguments map[string]any `json:"arguments,omitempty"`
		Service   []string       `json:"service,omitempty"`
	}

	auth struct {
		username, password string
	}

	// Metadata : Metadata returned from Kea.
	Metadata struct {
		ServerTags []string `json:"server-tags"`
	}
)

// New : Function used to create a new CradlePoint client data type.
func New(opts ...Option) *Client {
	client := new(Client)
	client.processOptions(opts...)
	return client
}

// make : creates an API request; a relative URI should be provided for uri
// and should not have a leading slash for a proper url.Parse() merge
//
//nolint:unparam
func (c *Client) make(method, url string, body interface{}, queryParameters *url.Values) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	if !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("https://%s", url)
	}
	if !strings.HasSuffix(url, "/") {
		url = fmt.Sprintf("%s/", url)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.auth.username, c.auth.password)

	if queryParameters != nil && len(*queryParameters) > 0 {
		req.URL.RawQuery = queryParameters.Encode()
	}

	if c.log.Level == logrus.DebugLevel {
		c.log.WithFields(logrus.Fields{
			"package":       packageName,
			"method":        "newRequest",
			"requestMethod": req.Method,
			"requestHost":   req.Host,
			"requestPath":   req.URL.Path,
			"requestQuery":  req.URL.RawQuery,
		}).Debug("HTTP request debugging")
	}

	return req, nil
}

// checkResponse : Kea returns a 200 OK on every response, even if there is a failure. Instead,
// it returns result codes in the json payload. We will need to check for these
// codes to determine an error.
//
//	0 - The command has been processed successfully.
//	1 - A general error or failure has occurred during the command processing.
//	2 - Specified command is unsupported by the server receiving it.
//	3 - The requested operation has been completed but the requested resource was not found.
//	    This status code is returned when a command returns no resources or affects no resources.
//	4 - The well-formed command has been processed but the requested changes could not be applied because
//	    they were in conflict with the server state or its notion of the configuration.
func checkResponse(resp *http.Response) (*Response, error) {
	e := make([]Response, 0)

	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, fmt.Errorf("checkResponse.json.Unmarshal(%w)", err)
	}
	if len(e) == 0 {
		return nil, fmt.Errorf("checkResponse(Unable to find response in %v)", e)
	}
	if e[0].Result == 0 {
		return &e[0], nil
	}

	return nil, fmt.Errorf("result:%d(%s)", e[0].Result, e[0].Text)
}

// do : sends an API request and JSON-decodes the API response.
// The response is stored in the value pointed to by v
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	c.log.WithFields(logrus.Fields{
		"package":    packageName,
		"method":     "rawDo",
		"statusCode": resp.StatusCode,
		"status":     resp.Status,
	}).Debug("HTTP response debugging")

	// Run through an error check process to determine if the response received was
	// an error.
	res, err := checkResponse(resp)
	if err != nil {
		// Before any other logic, execute a deferred close if the
		// body is not nil. This will prevent the body from not being closed
		// if there is some other error in the succeeding logic.
		if resp != nil {
			defer func(b io.ReadCloser) {
				if e := b.Close(); e != nil {
					err = fmt.Errorf("%w:%s", err, e)
				}
			}(resp.Body)
		}
		return nil, err
	}
	if v == nil {
		return res, nil
	}
	if e := json.NewDecoder(bytes.NewReader(res.Arguments)).Decode(v); e != nil && !errors.Is(e, io.EOF) {
		err = fmt.Errorf("failure decoding payload: %w", e)
	}
	return res, err
}

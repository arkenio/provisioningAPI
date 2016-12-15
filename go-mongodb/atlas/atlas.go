package atlas

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpDigestAuth "github.com/pteich/http-digest-auth-client"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1.0"
	defaultBaseURL = "https://cloud.mongodb.com/api/atlas/v1.0/"
	userAgent      = "go-mongodb/" + libraryVersion
	mediaType      = "application/json"
	format         = "json"
)

type Client struct {
	client *http.Client

	BaseURL *url.URL

	Username string
	APIKey   string

	UserAgent string

	// Services
	GroupWhiteList GroupWhiteListService
}

type Response struct {
	*http.Response
}

type ErrorResponse struct {
	Response *http.Response
	Detail   string `json:"detail"`
}

func NewClient(username string, apiKey string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: http.DefaultClient, BaseURL: baseURL, UserAgent: userAgent, Username: username, APIKey: apiKey}

	c.GroupWhiteList = &GroupWhiteListServiceOp{client: c}
	return c
}

type ClientOpt func(*Client) error

func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s+%s", ua, c.UserAgent)
		return nil
	}
}

func SetAuth(username string, apiKey string) ClientOpt {
	return func(c *Client) error {
		c.Username = username
		c.APIKey = apiKey
		return nil
	}
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	path := fmt.Sprintf("%s%s", c.BaseURL, urlStr)

	log.Printf("Setting PATH: %s", path)

	buf := new(bytes.Buffer)

	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, path, buf)

	if err != nil {
		return nil, err
	}

	DigestAuth := &httpDigestAuth.DigestHeaders{}
	DigestAuth, _ = DigestAuth.Auth(c.Username, c.APIKey, defaultBaseURL)
	DigestAuth.ApplyAuth(req)

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	return &response
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Detail)
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)

	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err := json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

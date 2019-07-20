package timeular

import (
	"bytes"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON = "application/json"

	defaultBaseURL = "https://api.timeular.com/api/v2"
)

type service struct {
	c *Client
}

type Client struct {
	c http.Client

	BaseURL string
	token string

	common service

	DeveloperService *developerService
	ServiceLessService *serviceLessService
	TrackingService *trackingService
}

func NewClient(APIKey, APISecret string) (*Client, error) {
	c := &Client{}
	c.BaseURL = defaultBaseURL
	c.common.c = c

	c.DeveloperService = (*developerService)(&c.common)
	c.ServiceLessService = (*serviceLessService)(&c.common)
	c.TrackingService = (*trackingService)(&c.common)

	token, err := c.DeveloperService.SignIn(APIKey, APISecret)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.token = token

	return c, nil
}

type message struct {
	Message string `json:"message"`
}

func checkResponseStatusCode(res *http.Response) error {
	if res.StatusCode != 200 {
		var msg message

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "http request failed with status code %d. Could not parse body.", res.StatusCode)
		}

		body := bytes.NewReader(b)

		err = json.NewDecoder(body).Decode(&msg)
		if err != nil {
			return errors.Wrapf(err, "http request failed with status code %d. Could not parse body.", res.StatusCode)
		}

		res.Body = ioutil.NopCloser(bytes.NewReader(b))

		return errors.WithStack(fmt.Errorf("http request failed with status code %d and message: %s", res.StatusCode, msg.Message))
	}

	return nil
}

func (c *Client) setAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
}

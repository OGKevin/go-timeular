package timeular

import (
	"bytes"
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

type developerService struct {
	c *Client
}

func (d *developerService) SignInContext(ctx context.Context, APIKey, APISecret string) (string, error) {
	return d.signIn(ctx, APIKey, APISecret)
}

func (d *developerService) SignIn(APIKey, APISecret string) (string, error) {
	return d.signIn(context.Background(), APIKey, APISecret)
}

func (d *developerService) signIn(ctx context.Context, APIKey, APISecret string) (string, error) {
	type body struct {
		APIKey    string `json:"apiKey"`
		APISecret string `json:"apiSecret"`
	}

	buf := bytes.NewBufferString("")

	err := json.NewEncoder(buf).Encode(&body{APIKey: APIKey, APISecret: APISecret})
	if err != nil {
		return "", errors.Wrap(err, "could not encode json body.")
	}

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/developer/sign-in", d.c.BaseURL), buf)
	if err != nil {
		return "nil", errors.Wrap(err, "could not create request for sign in.")
	}

	r.Header.Set(headerContentType, contentTypeJSON)

	res, err := d.c.c.Do(r)
	if err != nil {
		return "", errors.Wrap(err, "could not execute request for sign in.")
	}

	err = checkResponseStatusCode(res)
	if err != nil {
		return "", errors.Wrap(err, "sign in failed.")
	}

	type tokenBody struct {
		Token string `json:"token"`
	}

	resBody := tokenBody{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return "", errors.Wrap(err, "could not parse response body")
	}

	return resBody.Token, nil
}

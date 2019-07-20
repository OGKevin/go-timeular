package timeular

import (
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

type serviceLessService struct {
	c *Client
}

type TagsAndMentions struct {
	Tags []struct {
		Key   string `json:"key"`
		Label string `json:"label"`
	} `json:"tags"`
	Mentions []struct {
		Key   string `json:"key"`
		Label string `json:"label"`
	} `json:"mentions"`
}

func (s *serviceLessService) FetchTagsAndMentions() (*TagsAndMentions, error) {
	return s.fetchTagsAndMentions(context.Background())
}

func (s *serviceLessService) FetchTagsAndMentionsContext(ctx context.Context) (*TagsAndMentions, error) {
	return s.fetchTagsAndMentions(ctx)
}

func (s *serviceLessService) fetchTagsAndMentions(ctx context.Context) (*TagsAndMentions, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tags-and-mentions", s.c.BaseURL), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request for tags and mentions.")
	}

	r = r.WithContext(ctx)
	s.c.setAuthHeader(r)

	res, err := s.c.c.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request to fetch tags and mentions.")
	}

	err = checkResponseStatusCode(res)
	if err != nil {
		return nil,  errors.Wrap(err, "request to fetch tags and mentions failed.")
	}

	var resBody TagsAndMentions
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse body for tags and mentions response.")
	}

	return &resBody, nil
}

package timeular

import (
	"bytes"
	"context"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	//  YYYY-MM-DDTHH:mm:ss.SSS
	timeFormat = "2006-01-02T15:04:05.000"
)

type TimeularTime struct {
	time.Time
}

func (t *TimeularTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	return []byte(strconv.Quote(t.Format(timeFormat))), nil
}

// returns time.Now() no matter what!
func (t *TimeularTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if s == "nil" {
		t.Time = time.Time{}
		return nil
	}

	tt, err := time.Parse(timeFormat, s)
	if err != nil {
		return errors.WithStack(err)
	}

	t.Time = tt

	return nil
}

type trackingService struct {
	c *Client
}

type Note struct {
	Text string `json:"text"`
	TagsAndMentions
}

type CurrentTracking struct {
	CurrentTracking struct {
		Activity struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Color       string `json:"color"`
			Integration string `json:"integration"`
		} `json:"activity"`
		StartedAt TimeularTime `json:"startedAt"`
		Note      Note         `json:"note"`
	} `json:"currentTracking"`
}

func (s *trackingService) ShowCurrentTrackingContext(ctx context.Context) (*CurrentTracking, error) {
	return s.showCurrentTracking(ctx)
}

func (s *trackingService) ShowCurrentTracking() (*CurrentTracking, error) {
	return s.showCurrentTracking(context.Background())
}

func (s *trackingService) showCurrentTracking(ctx context.Context) (*CurrentTracking, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tracking", s.c.BaseURL), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request for show current tracking.")
	}

	r = r.WithContext(ctx)
	s.c.setAuthHeader(r)

	res, err := s.c.c.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request to fetch show current tracking.")
	}

	err = checkResponseStatusCode(res)
	if err != nil {
		return nil, errors.Wrap(err, "request to fetch show current tracking failed.")
	}

	var resBody CurrentTracking
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse body for show current tracking response.")
	}

	return &resBody, nil
}

type EditCurrentTracking struct {
	StartedAt  TimeularTime `json:"startedAt,omitempty"`
	ActivityId string       `json:"activityId,omitempty"`
	Note       NoteEdit     `json:"note,omitempty"`
}

type NoteEdit struct {
	Text     string   `json:"text"`
	Tags     []Entity `json:"tags"`
	Mentions []Entity `json:"mentions"`
}

type Entity struct {
	Indices []int  `json:"indices"`
	Key     string `json:"key"`
}

func (s *trackingService) EditCurrentTracking(activityID string, body EditCurrentTracking) error {
	return s.editCurrentTracking(context.Background(), activityID, body)
}

func (s *trackingService) EditCurrentTrackingContext(ctx context.Context, activityID string, body EditCurrentTracking) error {
	return s.editCurrentTracking(ctx, activityID, body)
}

func (s *trackingService) editCurrentTracking(ctx context.Context, activityID string, body EditCurrentTracking) error {
	buf := bytes.NewBufferString("")
	err := json.NewEncoder(buf).Encode(&body)
	if err != nil {
		return errors.Wrap(err, "could not encode request body.")
	}

	r, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/tracking/%s", s.c.BaseURL, activityID), buf)
	if err != nil {
		return errors.Wrap(err, "could not create request for edit current tracking.")
	}

	r = r.WithContext(ctx)
	r.Header.Set(headerContentType, contentTypeJSON)
	s.c.setAuthHeader(r)

	res, err := s.c.c.Do(r)
	if err != nil {
		return errors.Wrap(err, "could not do request to fetch edit current tracking.")
	}

	err = checkResponseStatusCode(res)
	if err != nil {
		return errors.Wrap(err, "request to fetch edit current tracking failed.")
	}

	var resBody CurrentTracking
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return errors.Wrap(err, "could not parse body for edit current tracking response.")
	}

	return nil
}

type TimeEntryList struct {
	TimeEntries []TimeEntry `json:"timeEntries"`
}

type TimeEntry struct {
	ID       string `json:"id"`
	Activity struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Color       string `json:"color"`
		Integration string `json:"integration"`
	} `json:"activity"`
	Duration struct {
		StartedAt TimeularTime `json:"startedAt"`
		StoppedAt TimeularTime `json:"stoppedAt"`
	} `json:"duration"`
	Note struct {
		Text string `json:"text"`
		Tags []struct {
			Indices []int  `json:"indices"`
			Key     string `json:"key"`
		} `json:"tags"`
		Mentions []struct {
			Indices []int  `json:"indices"`
			Key     string `json:"key"`
		} `json:"mentions"`
	} `json:"note"`
}

func (s *trackingService) GetTimeEntriesByRangeContext(ctx context.Context, startedBefore, stopppedAfter time.Time) (*TimeEntryList, error) {
	return s.getTimeEntriesByRange(ctx, startedBefore, stopppedAfter)
}

func (s *trackingService) GetTimeEntriesByRange(startedBefore, stoppedAfter time.Time) (*TimeEntryList, error) {
	return s.getTimeEntriesByRange(context.Background(), startedBefore, stoppedAfter)
}

func (s *trackingService) getTimeEntriesByRange(ctx context.Context, startedBefore, stoppedAfter time.Time) (*TimeEntryList, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/time-entries/%s/%s", s.c.BaseURL, stoppedAfter.Format(timeFormat), startedBefore.Format(timeFormat)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request for getting time entries by range.")
	}

	r = r.WithContext(ctx)
	r.Header.Set(headerContentType, contentTypeJSON)
	s.c.setAuthHeader(r)

	res, err := s.c.c.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request to fetch time entries by range.")
	}

	err = checkResponseStatusCode(res)
	if err != nil {
		return nil, errors.Wrap(err, "request to fetch time entries by range failed.")
	}

	var resBody TimeEntryList
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse body for time entries by range response.")
	}

	return &resBody, nil
}

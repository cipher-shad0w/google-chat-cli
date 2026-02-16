package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// EventsService provides methods for interacting with the Google Chat
// Space Events API (spaces.spaceEvents).
type EventsService struct {
	client *Client
}

// NewEventsService creates a new EventsService backed by the given API client.
func NewEventsService(client *Client) *EventsService {
	return &EventsService{client: client}
}

// List returns a paginated list of events from a space.
// GET /v1/{parent}/spaceEvents
// parent is a space name or ID (normalized with "spaces/" prefix).
func (s *EventsService) List(ctx context.Context, parent string, filter string, pageSize int, pageToken string) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/spaceEvents", parent)

	params := url.Values{}
	AddQueryParam(params, "filter", filter)
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)

	return s.client.Get(ctx, path, params)
}

// Get retrieves a single space event by its full resource name.
// GET /v1/{name}
// Name format: spaces/{space}/spaceEvents/{spaceEvent}
func (s *EventsService) Get(ctx context.Context, name string) (json.RawMessage, error) {
	return s.client.Get(ctx, name, nil)
}

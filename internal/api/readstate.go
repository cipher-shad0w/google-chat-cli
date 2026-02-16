package api

import (
	"context"
	"encoding/json"
	"net/url"
)

// ReadStateService provides methods for interacting with the Google Chat
// Read State API (users.spaces.spaceReadState and
// users.spaces.threads.threadReadState).
type ReadStateService struct {
	client *Client
}

// NewReadStateService creates a new ReadStateService backed by the given API client.
func NewReadStateService(client *Client) *ReadStateService {
	return &ReadStateService{client: client}
}

// GetSpaceReadState returns the read state of a space for the calling user.
// GET /v1/{name}
// Name format: users/{user}/spaces/{space}/spaceReadState
func (s *ReadStateService) GetSpaceReadState(ctx context.Context, name string) (json.RawMessage, error) {
	return s.client.Get(ctx, name, nil)
}

// UpdateSpaceReadState updates the read state of a space for the calling user.
// PATCH /v1/{name}
// Name format: users/{user}/spaces/{space}/spaceReadState
func (s *ReadStateService) UpdateSpaceReadState(ctx context.Context, name string, state map[string]interface{}, updateMask string) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "updateMask", updateMask)

	return s.client.Patch(ctx, name, params, state)
}

// GetThreadReadState returns the read state of a thread for the calling user.
// GET /v1/{name}
// Name format: users/{user}/spaces/{space}/threads/{thread}/threadReadState
func (s *ReadStateService) GetThreadReadState(ctx context.Context, name string) (json.RawMessage, error) {
	return s.client.Get(ctx, name, nil)
}

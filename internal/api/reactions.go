package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ReactionsService provides methods for interacting with the Google Chat
// Reactions API (spaces.messages.reactions).
type ReactionsService struct {
	client *Client
}

// NewReactionsService creates a new ReactionsService backed by the given client.
func NewReactionsService(client *Client) *ReactionsService {
	return &ReactionsService{client: client}
}

// List returns reactions on a message.
// parent is the message resource name, e.g. "spaces/{space}/messages/{message}".
func (s *ReactionsService) List(ctx context.Context, parent string, pageSize int, pageToken, filter string) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/reactions", parent)

	params := url.Values{}
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)
	AddQueryParam(params, "filter", filter)

	return s.client.Get(ctx, path, params)
}

// Create adds a reaction to a message.
// parent is the message resource name, e.g. "spaces/{space}/messages/{message}".
// reaction is the request body describing the reaction to create.
func (s *ReactionsService) Create(ctx context.Context, parent string, reaction map[string]interface{}) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/reactions", parent)

	return s.client.Post(ctx, path, nil, reaction)
}

// Delete removes a reaction.
// name is the full reaction resource name,
// e.g. "spaces/{space}/messages/{message}/reactions/{reaction}".
func (s *ReactionsService) Delete(ctx context.Context, name string) (json.RawMessage, error) {
	name = NormalizeName(name, "spaces/")
	return s.client.Delete(ctx, name, nil)
}

package api

import (
	"context"
	"encoding/json"
	"net/url"
)

// EmojiService handles custom emoji operations on the Google Chat API.
type EmojiService struct {
	client *Client
}

// NewEmojiService creates a new EmojiService bound to the given API client.
func NewEmojiService(client *Client) *EmojiService {
	return &EmojiService{client: client}
}

// List retrieves a paginated list of custom emojis.
// GET /v1/customEmojis
func (s *EmojiService) List(ctx context.Context, filter string, pageSize int, pageToken string) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "filter", filter)
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)

	return s.client.Get(ctx, "customEmojis", params)
}

// Get retrieves a single custom emoji by name or ID.
// GET /v1/{name}
func (s *EmojiService) Get(ctx context.Context, name string) (json.RawMessage, error) {
	name = NormalizeName(name, "customEmojis/")
	return s.client.Get(ctx, name, nil)
}

// Create creates a new custom emoji.
// POST /v1/customEmojis
func (s *EmojiService) Create(ctx context.Context, emoji map[string]interface{}) (json.RawMessage, error) {
	return s.client.Post(ctx, "customEmojis", nil, emoji)
}

// Delete deletes a custom emoji by name or ID.
// DELETE /v1/{name}
func (s *EmojiService) Delete(ctx context.Context, name string) (json.RawMessage, error) {
	name = NormalizeName(name, "customEmojis/")
	return s.client.Delete(ctx, name, nil)
}

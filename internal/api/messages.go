package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// MessagesService provides methods for interacting with the Google Chat
// Messages API resource.
type MessagesService struct {
	client *Client
}

// NewMessagesService creates a new MessagesService backed by the given client.
func NewMessagesService(client *Client) *MessagesService {
	return &MessagesService{client: client}
}

// List retrieves messages from a space.
// GET /v1/{parent}/messages
func (s *MessagesService) List(ctx context.Context, parent string, pageSize int, pageToken, filter, orderBy string, showDeleted bool) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/messages", parent)

	params := url.Values{}
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)
	AddQueryParam(params, "filter", filter)
	AddQueryParam(params, "orderBy", orderBy)
	AddQueryParamBool(params, "showDeleted", showDeleted)

	return s.client.Get(ctx, path, params)
}

// Get retrieves a single message by its full resource name.
// GET /v1/{name}
// Name format: spaces/{space}/messages/{message}
func (s *MessagesService) Get(ctx context.Context, name string) (json.RawMessage, error) {
	return s.client.Get(ctx, name, nil)
}

// Create sends a new message to a space.
// POST /v1/{parent}/messages
func (s *MessagesService) Create(ctx context.Context, parent string, message map[string]interface{}, threadKey, requestID, messageID, messageReplyOption string) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/messages", parent)

	params := url.Values{}
	AddQueryParam(params, "threadKey", threadKey)
	AddQueryParam(params, "requestId", requestID)
	AddQueryParam(params, "messageId", messageID)
	AddQueryParam(params, "messageReplyOption", messageReplyOption)

	return s.client.Post(ctx, path, params, message)
}

// Patch partially updates a message.
// PATCH /v1/{name}
func (s *MessagesService) Patch(ctx context.Context, name string, message map[string]interface{}, updateMask string, allowMissing bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "updateMask", updateMask)
	AddQueryParamBool(params, "allowMissing", allowMissing)

	return s.client.Patch(ctx, name, params, message)
}

// Update fully replaces a message.
// PUT /v1/{name}
func (s *MessagesService) Update(ctx context.Context, name string, message map[string]interface{}, updateMask string, allowMissing bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "updateMask", updateMask)
	AddQueryParamBool(params, "allowMissing", allowMissing)

	return s.client.Put(ctx, name, params, message)
}

// Delete removes a message.
// DELETE /v1/{name}
// If force is true, the force query parameter is set, which also deletes
// threaded replies to the message.
func (s *MessagesService) Delete(ctx context.Context, name string, force bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParamBool(params, "force", force)

	return s.client.Delete(ctx, name, params)
}

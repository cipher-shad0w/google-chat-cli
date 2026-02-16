package api

import (
	"context"
	"encoding/json"
)

// AttachmentsService provides methods for interacting with the Google Chat
// Attachments API (spaces.messages.attachments).
type AttachmentsService struct {
	client *Client
}

// NewAttachmentsService creates a new AttachmentsService backed by the given client.
func NewAttachmentsService(client *Client) *AttachmentsService {
	return &AttachmentsService{client: client}
}

// Get returns metadata for a message attachment.
// name is the full attachment resource name,
// e.g. "spaces/{space}/messages/{message}/attachments/{attachment}".
func (s *AttachmentsService) Get(ctx context.Context, name string) (json.RawMessage, error) {
	name = NormalizeName(name, "spaces/")
	return s.client.Get(ctx, name, nil)
}

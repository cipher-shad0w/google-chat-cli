package api

import (
	"context"
	"encoding/json"
	"net/url"
)

// NotificationsService provides methods for interacting with the Google Chat
// Space Notification Settings API (users.spaces.spaceNotificationSetting).
type NotificationsService struct {
	client *Client
}

// NewNotificationsService creates a new NotificationsService backed by the given API client.
func NewNotificationsService(client *Client) *NotificationsService {
	return &NotificationsService{client: client}
}

// Get returns the notification setting for a space.
// GET /v1/{name}
// Name format: users/{user}/spaces/{space}/spaceNotificationSetting
func (s *NotificationsService) Get(ctx context.Context, name string) (json.RawMessage, error) {
	return s.client.Get(ctx, name, nil)
}

// Patch updates the notification setting for a space.
// PATCH /v1/{name}
// Name format: users/{user}/spaces/{space}/spaceNotificationSetting
func (s *NotificationsService) Patch(ctx context.Context, name string, setting map[string]interface{}, updateMask string) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "updateMask", updateMask)

	return s.client.Patch(ctx, name, params, setting)
}

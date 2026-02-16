package api

import (
	"context"
	"encoding/json"
	"net/url"
)

// SpacesService provides methods for interacting with the Google Chat
// Spaces API endpoints.
type SpacesService struct {
	client *Client
}

// NewSpacesService creates a new SpacesService backed by the given API client.
func NewSpacesService(client *Client) *SpacesService {
	return &SpacesService{client: client}
}

// List returns a paginated list of spaces the caller is a member of.
// GET /v1/spaces
func (s *SpacesService) List(ctx context.Context, filter string, pageSize int, pageToken string) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "filter", filter)
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)

	return s.client.Get(ctx, "spaces", params)
}

// Get returns a single space by name.
// GET /v1/{name}
func (s *SpacesService) Get(ctx context.Context, name string, useAdminAccess bool) (json.RawMessage, error) {
	name = NormalizeName(name, "spaces/")

	params := url.Values{}
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Get(ctx, name, params)
}

// Create creates a new space.
// POST /v1/spaces
func (s *SpacesService) Create(ctx context.Context, space map[string]interface{}, requestID string) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "requestId", requestID)

	return s.client.Post(ctx, "spaces", params, space)
}

// Patch updates an existing space.
// PATCH /v1/{name}
func (s *SpacesService) Patch(ctx context.Context, name string, space map[string]interface{}, updateMask string, useAdminAccess bool) (json.RawMessage, error) {
	name = NormalizeName(name, "spaces/")

	params := url.Values{}
	AddQueryParam(params, "updateMask", updateMask)
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Patch(ctx, name, params, space)
}

// Delete deletes a space.
// DELETE /v1/{name}
func (s *SpacesService) Delete(ctx context.Context, name string, useAdminAccess bool) (json.RawMessage, error) {
	name = NormalizeName(name, "spaces/")

	params := url.Values{}
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Delete(ctx, name, params)
}

// Search searches for spaces visible to the caller.
// GET /v1/spaces:search
func (s *SpacesService) Search(ctx context.Context, query string, pageSize int, pageToken, orderBy string, useAdminAccess bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "query", query)
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)
	AddQueryParam(params, "orderBy", orderBy)
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Get(ctx, "spaces:search", params)
}

// Setup creates a space and adds specified users to it.
// POST /v1/spaces:setup
func (s *SpacesService) Setup(ctx context.Context, request map[string]interface{}) (json.RawMessage, error) {
	return s.client.Post(ctx, "spaces:setup", nil, request)
}

// FindDirectMessage finds a direct message space with the specified user.
// GET /v1/spaces:findDirectMessage
func (s *SpacesService) FindDirectMessage(ctx context.Context, userName string) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "name", userName)

	return s.client.Get(ctx, "spaces:findDirectMessage", params)
}

// CompleteImport completes the import process for a space, making it visible
// to users and allowing new messages.
// POST /v1/{name}:completeImport
func (s *SpacesService) CompleteImport(ctx context.Context, name string) (json.RawMessage, error) {
	name = NormalizeName(name, "spaces/")

	return s.client.Post(ctx, name+":completeImport", nil, map[string]interface{}{})
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// MembersService handles membership operations for Google Chat spaces.
type MembersService struct {
	client *Client
}

// NewMembersService creates a new MembersService backed by the given Client.
func NewMembersService(client *Client) *MembersService {
	return &MembersService{client: client}
}

// List retrieves members of a space.
// parent is the space resource name (e.g. "spaces/AAAA" or just "AAAA").
func (s *MembersService) List(ctx context.Context, parent string, pageSize int, pageToken, filter string, showInvited, showGroups, useAdminAccess bool) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/members", parent)

	params := url.Values{}
	AddQueryParamInt(params, "pageSize", pageSize)
	AddQueryParam(params, "pageToken", pageToken)
	AddQueryParam(params, "filter", filter)
	AddQueryParamBool(params, "showInvited", showInvited)
	AddQueryParamBool(params, "showGroups", showGroups)
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Get(ctx, path, params)
}

// Get retrieves a single membership by its resource name.
// name is the full resource name (e.g. "spaces/AAAA/members/123456").
func (s *MembersService) Get(ctx context.Context, name string, useAdminAccess bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Get(ctx, name, params)
}

// Create adds a new member to a space.
// parent is the space resource name (e.g. "spaces/AAAA" or just "AAAA").
// membership is the membership resource body to create.
func (s *MembersService) Create(ctx context.Context, parent string, membership map[string]interface{}, useAdminAccess bool) (json.RawMessage, error) {
	parent = NormalizeName(parent, "spaces/")
	path := fmt.Sprintf("%s/members", parent)

	params := url.Values{}
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Post(ctx, path, params, membership)
}

// Patch updates an existing membership.
// name is the full resource name (e.g. "spaces/AAAA/members/123456").
// membership is the membership resource body with updated fields.
// updateMask specifies which fields to update (comma-separated field paths).
func (s *MembersService) Patch(ctx context.Context, name string, membership map[string]interface{}, updateMask string, useAdminAccess bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParam(params, "updateMask", updateMask)
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Patch(ctx, name, params, membership)
}

// Delete removes a membership from a space.
// name is the full resource name (e.g. "spaces/AAAA/members/123456").
func (s *MembersService) Delete(ctx context.Context, name string, useAdminAccess bool) (json.RawMessage, error) {
	params := url.Values{}
	AddQueryParamBool(params, "useAdminAccess", useAdminAccess)

	return s.client.Delete(ctx, name, params)
}

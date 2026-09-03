package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// Hash provides Hash Map operations scoped to a specific namespace.
type Hash struct {
	client    *Client
	namespace string
}

// Hash returns a Hash Map sub-client bound to the specified namespace.
func (c *Client) Hash(namespace string) *Hash {
	return &Hash{
		client:    c,
		namespace: namespace,
	}
}

// Set creates or resets an empty hash table under key with an optional TTL in seconds.
// Returns an error if the operation fails on the server.
func (h *Hash) Set(ctx context.Context, key string, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HSET,
		Namespace: h.namespace,
		Key:       key,
		Ttl:       ttl,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Get retrieves all field-value pairs stored in the hash table as a map[string]string.
// Returns an error if the hash does not exist or has expired.
func (h *Hash) Get(ctx context.Context, key string) (map[string]string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HGET,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.GetSuccess() {
		return nil, newServerError(resp.GetError())
	}
	return resp.GetMapValue(), nil
}

// Del removes the entire hash table under key.
// Returns an error if the operation fails on the server.
func (h *Hash) Del(ctx context.Context, key string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HDEL,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Expire sets or updates the time-to-live for the hash table in seconds.
// Returns an error if the hash does not exist or has expired.
func (h *Hash) Expire(ctx context.Context, key string, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HEXPIRE,
		Namespace: h.namespace,
		Key:       key,
		Ttl:       ttl,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// TTL returns the remaining time-to-live of the hash table in seconds.
// Returns -1 if the hash exists without an expiration time.
// Returns -2 or an error if the hash does not exist.
func (h *Hash) TTL(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HTTL,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Exists checks whether an unexpired hash table exists under key.
// Returns true if the hash table exists, false otherwise.
func (h *Hash) Exists(ctx context.Context, key string) (bool, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HEXIST,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return false, err
	}
	if resp.GetError() != "" {
		return false, newServerError(resp.GetError())
	}
	return resp.GetSuccess(), nil
}

// Len returns the number of fields currently stored in the hash table.
// Returns an error if the hash does not exist or has expired.
func (h *Hash) Len(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HLEN,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Keys returns a slice containing all field names in the hash table.
// Returns an error if the hash does not exist or has expired.
func (h *Hash) Keys(ctx context.Context, key string) ([]string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HKEYS,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.GetSuccess() {
		return nil, newServerError(resp.GetError())
	}
	return resp.GetSliceValue(), nil
}

// Values returns a slice containing all field values in the hash table.
// Returns an error if the hash does not exist or has expired.
func (h *Hash) Values(ctx context.Context, key string) ([]string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HVALS,
		Namespace: h.namespace,
		Key:       key,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.GetSuccess() {
		return nil, newServerError(resp.GetError())
	}
	return resp.GetSliceValue(), nil
}

// FSet sets or updates the value of a specific field within the hash table under key.
// Creates the hash table automatically if it does not already exist.
func (h *Hash) FSet(ctx context.Context, key, field, value string) error {
	req := &memapv1.Request{
		Command:     memapv1.CommandType_HFSET,
		Namespace:   h.namespace,
		Key:         key,
		Field:       field,
		StringValue: value,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// FGet retrieves the value of a specific field from the hash table under key.
// Returns an error if the hash or the specified field does not exist.
func (h *Hash) FGet(ctx context.Context, key, field string) (string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HFGET,
		Namespace: h.namespace,
		Key:       key,
		Field:     field,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

// FDel removes a specific field from the hash table under key.
// Returns an error if the operation fails on the server.
func (h *Hash) FDel(ctx context.Context, key, field string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_HFDEL,
		Namespace: h.namespace,
		Key:       key,
		Field:     field,
	}
	resp, err := h.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

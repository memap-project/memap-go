package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// Set provides Set operations scoped to a specific namespace.
type Set struct {
	client    *Client
	namespace string
}

// Set returns a Set sub-client bound to the specified namespace.
func (c *Client) Set(namespace string) *Set {
	return &Set{
		client:    c,
		namespace: namespace,
	}
}

// Add adds a member to the set under key in the namespace.
// If the set does not exist, it is created automatically.
// Returns nil on success.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) Add(ctx context.Context, key, value string) error {
	req := &memapv1.Request{
		Command:     memapv1.CommandType_SADD,
		Namespace:   s.namespace,
		Key:         key,
		StringValue: value,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Remove removes a member from the set under key in the namespace.
// Removing a non-existent member or from a non-existent set succeeds without error.
// Returns nil on success.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) Remove(ctx context.Context, key, value string) error {
	req := &memapv1.Request{
		Command:     memapv1.CommandType_SREMOVE,
		Namespace:   s.namespace,
		Key:         key,
		StringValue: value,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// IsMember checks whether value is a member of the set under key in the namespace.
// Returns true if the value is a member of the set, false if it is not or if the set does not exist.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) IsMember(ctx context.Context, key, value string) (bool, error) {
	req := &memapv1.Request{
		Command:     memapv1.CommandType_SISMEMBER,
		Namespace:   s.namespace,
		Key:         key,
		StringValue: value,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return false, err
	}
	if !resp.GetSuccess() {
		return false, newServerError(resp.GetError())
	}
	return resp.GetIntValue() == 1, nil
}

// Card returns the cardinality (number of members) of the set under key in the namespace.
// Returns the count of members on success.
// Returns [ErrKeyNotFound] if the set does not exist or has expired.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) Card(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_SCARD,
		Namespace: s.namespace,
		Key:       key,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Members returns all members of the set under key in the namespace.
// Returns a slice of member strings on success.
// Returns [ErrKeyNotFound] if the set does not exist or has expired.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) Members(ctx context.Context, key string) ([]string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_SMEMBERS,
		Namespace: s.namespace,
		Key:       key,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.GetSuccess() {
		return nil, newServerError(resp.GetError())
	}
	return resp.GetSliceValue(), nil
}

// Expire sets or updates the time-to-live for the set under key in seconds.
// Returns nil on success.
// Returns [ErrKeyNotFound] if the set does not exist or has expired.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) Expire(ctx context.Context, key string, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_SEXPIRE,
		Namespace: s.namespace,
		Key:       key,
		Ttl:       ttl,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// TTL returns the remaining time-to-live of the set under key in seconds.
// Returns -1 if the set exists without an expiration time.
// Returns [ErrKeyNotFound] if the set does not exist or has expired.
// Returns [ErrNamespaceNotFound] if the namespace does not exist.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
func (s *Set) TTL(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_STTL,
		Namespace: s.namespace,
		Key:       key,
	}
	resp, err := s.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// Counter provides Counter operations scoped to a specific namespace.
type Counter struct {
	client    *Client
	namespace string
}

// Counter returns a Counter sub-client bound to the specified namespace.
func (c *Client) Counter(namespace string) *Counter {
	return &Counter{
		client:    c,
		namespace: namespace,
	}
}

// Init initializes a new counter under key with an upper limit and optional TTL in seconds.
// Returns an error if the counter already exists or initialization fails.
func (cnt *Counter) Init(ctx context.Context, key string, limit, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CINIT,
		Namespace: cnt.namespace,
		Key:       key,
		Limit:     limit,
		Ttl:       ttl,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// SetLimit updates the upper limit for the counter under key.
// Returns an error if the counter does not exist or has expired.
func (cnt *Counter) SetLimit(ctx context.Context, key string, limit int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CSLIMIT,
		Namespace: cnt.namespace,
		Key:       key,
		Limit:     limit,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// GetLimit returns the configured upper limit of the counter under key.
// Returns an error if the counter does not exist or has expired.
func (cnt *Counter) GetLimit(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CGLIMIT,
		Namespace: cnt.namespace,
		Key:       key,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Get returns the current integer value of the counter under key.
// Returns an error if the counter does not exist or has expired.
func (cnt *Counter) Get(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CGET,
		Namespace: cnt.namespace,
		Key:       key,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Del removes the counter under key.
// Returns an error if the operation fails on the server.
func (cnt *Counter) Del(ctx context.Context, key string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CDEL,
		Namespace: cnt.namespace,
		Key:       key,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Expire sets or updates the time-to-live for the counter in seconds.
// Returns an error if the counter does not exist or has expired.
func (cnt *Counter) Expire(ctx context.Context, key string, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CEXPIRE,
		Namespace: cnt.namespace,
		Key:       key,
		Ttl:       ttl,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// TTL returns the remaining time-to-live of the counter in seconds.
// Returns -1 if the counter exists without an expiration time.
// Returns -2 or an error if the counter does not exist.
func (cnt *Counter) TTL(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CTTL,
		Namespace: cnt.namespace,
		Key:       key,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// IncrBy increments the counter under key by delta and returns the new value.
// Returns an error if the counter does not exist, has expired, or if the increment exceeds the limit.
func (cnt *Counter) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CINCRBY,
		Namespace: cnt.namespace,
		Key:       key,
		IntValue:  delta,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// DecrBy decrements the counter under key by delta and returns the new value.
// Returns an error if the counter does not exist, has expired, or if the decrement would drop below zero.
func (cnt *Counter) DecrBy(ctx context.Context, key string, delta int64) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CDECRBY,
		Namespace: cnt.namespace,
		Key:       key,
		IntValue:  delta,
	}
	resp, err := cnt.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// KV provides Key-Value operations scoped to a specific namespace.
type KV struct {
	client    *Client
	namespace string
}

// KV returns a Key-Value sub-client bound to the specified namespace.
func (c *Client) KV(namespace string) *KV {
	return &KV{
		client:    c,
		namespace: namespace,
	}
}

// Set stores a string value under key with an optional TTL in seconds (0 = no expiration).
// Returns an error if the operation fails on the server.
func (kv *KV) Set(ctx context.Context, key, value string, ttl int64) error {
	req := &memapv1.Request{
		Command:     memapv1.CommandType_SET,
		Namespace:   kv.namespace,
		Key:         key,
		StringValue: value,
		Ttl:         ttl,
	}
	resp, err := kv.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Get retrieves the string value associated with key.
// Returns an error if the key does not exist or has expired.
func (kv *KV) Get(ctx context.Context, key string) (string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_GET,
		Namespace: kv.namespace,
		Key:       key,
	}
	resp, err := kv.client.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

// Del removes the specified key and its associated value.
// Returns an error if the operation fails on the server.
func (kv *KV) Del(ctx context.Context, key string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_DEL,
		Namespace: kv.namespace,
		Key:       key,
	}
	resp, err := kv.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Expire sets or updates the time-to-live for key in seconds.
// Returns an error if the key does not exist or has expired.
func (kv *KV) Expire(ctx context.Context, key string, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_EXPIRE,
		Namespace: kv.namespace,
		Key:       key,
		Ttl:       ttl,
	}
	resp, err := kv.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// TTL returns the remaining time-to-live of key in seconds.
// Returns -1 if the key exists without an expiration time.
// Returns -2 or an error if the key does not exist.
func (kv *KV) TTL(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_TTL,
		Namespace: kv.namespace,
		Key:       key,
	}
	resp, err := kv.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

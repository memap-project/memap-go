package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// RBuffer provides RingBuffer operations scoped to a specific namespace.
type RBuffer struct {
	client    *Client
	namespace string
}

// RBuffer returns a RingBuffer sub-client bound to the specified namespace.
func (c *Client) RBuffer(namespace string) *RBuffer {
	return &RBuffer{
		client:    c,
		namespace: namespace,
	}
}

// Init initializes a new ring buffer under key with the given capacity and optional TTL in seconds.
// Returns an error if the ring buffer already exists or initialization fails.
func (b *RBuffer) Init(ctx context.Context, key string, capacity, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BINIT,
		Namespace: b.namespace,
		Key:       key,
		Limit:     capacity,
		Ttl:       ttl,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Push appends a value to the ring buffer under key, overwriting the oldest entry if full.
// Returns an error if the ring buffer does not exist or has expired.
func (b *RBuffer) Push(ctx context.Context, key, value string) error {
	req := &memapv1.Request{
		Command:     memapv1.CommandType_BPUSH,
		Namespace:   b.namespace,
		Key:         key,
		StringValue: value,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Pop removes and returns the oldest entry (head) from the ring buffer.
// Returns an error if the ring buffer does not exist, has expired, or is empty.
func (b *RBuffer) Pop(ctx context.Context, key string) (string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BPOP,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

// At retrieves an entry by its logical index (0 = oldest) without removing it.
// Returns an error if the ring buffer does not exist, has expired, or the index is out of bounds.
func (b *RBuffer) At(ctx context.Context, key string, index int64) (string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BAT,
		Namespace: b.namespace,
		Key:       key,
		IntValue:  index,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

// Slice returns all elements currently in the ring buffer in chronological order (oldest to newest).
// Returns an error if the ring buffer does not exist or has expired.
func (b *RBuffer) Slice(ctx context.Context, key string) ([]string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BSLICE,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.GetSuccess() {
		return nil, newServerError(resp.GetError())
	}
	return resp.GetSliceValue(), nil
}

// Peek returns the oldest entry (head) from the ring buffer without removing it.
// Returns an error if the ring buffer does not exist, has expired, or is empty.
func (b *RBuffer) Peek(ctx context.Context, key string) (string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BPEEK,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

// Back returns the newest entry (tail) from the ring buffer without removing it.
// Returns an error if the ring buffer does not exist, has expired, or is empty.
func (b *RBuffer) Back(ctx context.Context, key string) (string, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BBACK,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

// Cap returns the configured capacity of the ring buffer under key.
// Returns an error if the ring buffer does not exist or has expired.
func (b *RBuffer) Cap(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BCAP,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Len returns the current number of elements in the ring buffer under key.
// Returns an error if the ring buffer does not exist or has expired.
func (b *RBuffer) Len(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BLEN,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

// Reset clears all elements from the ring buffer under key, resetting its length to 0.
// Returns an error if the ring buffer does not exist or has expired.
func (b *RBuffer) Reset(ctx context.Context, key string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BRESET,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Del removes the ring buffer under key.
// Returns an error if the operation fails on the server.
func (b *RBuffer) Del(ctx context.Context, key string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BDEL,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Expire sets or updates the time-to-live for the ring buffer in seconds.
// Returns an error if the ring buffer does not exist or has expired.
func (b *RBuffer) Expire(ctx context.Context, key string, ttl int64) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BEXPIRE,
		Namespace: b.namespace,
		Key:       key,
		Ttl:       ttl,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// TTL returns the remaining time-to-live of the ring buffer in seconds.
// Returns -1 if the ring buffer exists without an expiration time.
// Returns -2 or an error if the ring buffer does not exist.
func (b *RBuffer) TTL(ctx context.Context, key string) (int64, error) {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_BTTL,
		Namespace: b.namespace,
		Key:       key,
	}
	resp, err := b.client.do(ctx, req)
	if err != nil {
		return 0, err
	}
	if !resp.GetSuccess() {
		return 0, newServerError(resp.GetError())
	}
	return resp.GetIntValue(), nil
}

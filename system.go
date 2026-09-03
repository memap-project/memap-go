package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// Create creates a new namespace on the server.
// Returns an error if the namespace already exists or creation fails.
func (c *Client) Create(ctx context.Context, namespace string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_CREATE,
		Namespace: namespace,
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Drop deletes the namespace and all its keys on the server.
// Returns an error if the namespace does not exist or deletion fails.
func (c *Client) Drop(ctx context.Context, namespace string) error {
	req := &memapv1.Request{
		Command:   memapv1.CommandType_DROP,
		Namespace: namespace,
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Erase drops all custom namespaces and flushes the default namespace.
// Returns an error if the operation fails on the server.
func (c *Client) Erase(ctx context.Context) error {
	req := &memapv1.Request{
		Command: memapv1.CommandType_ERASE,
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

// Flush removes all keys across all namespaces while preserving the namespaces themselves.
// Returns an error if the operation fails on the server.
func (c *Client) Flush(ctx context.Context) error {
	req := &memapv1.Request{
		Command: memapv1.CommandType_FLUSH,
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.GetSuccess() {
		return newServerError(resp.GetError())
	}
	return nil
}

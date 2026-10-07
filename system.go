package memap

import (
	"context"

	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// Create creates a new namespace on the server.
// Returns nil on success.
// Returns [ErrNamespaceAlreadyExists] if the namespace already exists.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
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

// Drop deletes the namespace and all its stored keys from the server.
// Deleting a non-existent namespace succeeds without error.
// Returns nil on success.
// Returns a [*ServerError] if deletion fails on the server.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
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

// Erase flushes all data in the default namespace and drops all custom namespaces from the server.
// Returns nil on success.
// Returns a [*ServerError] if the operation fails on the server.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
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

// Flush removes all keys across all namespaces (including the default namespace) while preserving the namespaces themselves.
// Returns nil on success.
// Returns a [*ServerError] if the operation fails on the server.
// Returns [ErrClosed] if the client connection is closed, or a network/context error.
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

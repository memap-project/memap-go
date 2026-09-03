package memap

import "context"

// NamespaceClient provides an API scoped to operations within a single namespace.
type NamespaceClient struct {
	client    *Client
	namespace string

	kv      *KV
	hash    *Hash
	counter *Counter
	rbuffer *RBuffer
}

// Namespace returns a NamespaceClient bound to the specified namespace name.
func (c *Client) Namespace(name string) *NamespaceClient {
	return &NamespaceClient{
		client:    c,
		namespace: name,
		kv:        c.KV(name),
		hash:      c.Hash(name),
		counter:   c.Counter(name),
		rbuffer:   c.RBuffer(name),
	}
}

// Name returns the name of the namespace associated with this client.
func (n *NamespaceClient) Name() string {
	return n.namespace
}

// Create creates this namespace on the server.
// Returns an error if the namespace already exists or creation fails.
func (n *NamespaceClient) Create(ctx context.Context) error {
	return n.client.Create(ctx, n.namespace)
}

// Drop deletes this namespace and all its stored data from the server.
// Returns an error if the namespace does not exist or deletion fails.
func (n *NamespaceClient) Drop(ctx context.Context) error {
	return n.client.Drop(ctx, n.namespace)
}

// KV returns the Key-Value sub-client bound to this namespace.
func (n *NamespaceClient) KV() *KV {
	return n.kv
}

// Hash returns the Hash Map sub-client bound to this namespace.
func (n *NamespaceClient) Hash() *Hash {
	return n.hash
}

// Counter returns the Counter sub-client bound to this namespace.
func (n *NamespaceClient) Counter() *Counter {
	return n.counter
}

// RBuffer returns the RingBuffer sub-client bound to this namespace.
func (n *NamespaceClient) RBuffer() *RBuffer {
	return n.rbuffer
}

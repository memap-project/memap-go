package memap

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/dmi3midd/protorw"
	memapv1 "github.com/memap-project/memap-proto/gen/memapv1/go"
)

// Client is a thread-safe client for interacting with the Memap server over TCP.
type Client struct {
	cfg    Config
	mu     sync.Mutex
	conn   net.Conn
	closed bool
}

// New creates and connects a new Memap client to the specified address using functional options.
// If addr is empty, the default address ("localhost:2118") is used.
// Returns the connected [*Client] on success, or an error if the initial TCP dial fails.
func New(addr string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	if addr != "" {
		cfg.Addr = addr
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return NewClient(cfg)
}

// NewClient creates and connects a new Memap client based on the provided [Config].
// Returns the connected [*Client] on success, or an error if the initial TCP dial fails.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Addr == "" {
		cfg.Addr = defaultAddr
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = defaultDialTimeout
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}

	c := &Client{
		cfg: cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if _, err := c.getConn(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to memap server at %s: %w", cfg.Addr, err)
	}

	return c, nil
}

func (c *Client) getConn(ctx context.Context) (net.Conn, error) {
	if c.closed {
		return nil, ErrClosed
	}
	if c.conn != nil {
		return c.conn, nil
	}

	dialer := &net.Dialer{
		Timeout: c.cfg.DialTimeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", c.cfg.Addr)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	return c.conn, nil
}

// Close closes the underlying network connection to the Memap server.
// Subsequent operations on this client will return [ErrClosed].
// Returns nil if the client has already been closed or closes cleanly.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

// do sends a Protobuf request and reads the response over the TCP connection.
func (c *Client) do(ctx context.Context, req *memapv1.Request) (*memapv1.Response, error) {
	if req == nil {
		return nil, errNilRequest
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}

	var deadline time.Time
	if ctxDeadline, ok := ctx.Deadline(); ok {
		deadline = ctxDeadline
	} else if c.cfg.Timeout > 0 {
		deadline = time.Now().Add(c.cfg.Timeout)
	}
	_ = conn.SetDeadline(deadline)

	if err := protorw.WriteMsg(conn, req); err != nil {
		_ = conn.Close()
		c.conn = nil
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var resp memapv1.Response
	if err := protorw.ReadMsg(conn, &resp); err != nil {
		_ = conn.Close()
		c.conn = nil
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	_ = conn.SetDeadline(time.Time{})

	return &resp, nil
}

// Ping sends a PING command to verify server availability and connection health.
// Returns "PONG" on success.
// Returns a [*ServerError] if the server rejects the request,
// [ErrClosed] if the client connection is closed, or a network/context error.
func (c *Client) Ping(ctx context.Context) (string, error) {
	req := &memapv1.Request{
		Command: memapv1.CommandType_PING,
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.GetSuccess() {
		return "", newServerError(resp.GetError())
	}
	return resp.GetStringValue(), nil
}

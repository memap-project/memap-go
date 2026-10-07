package memap

import "time"

const (
	defaultAddr        = "localhost:2118"
	defaultDialTimeout = 5 * time.Second
	defaultTimeout     = 5 * time.Second
)

// Config defines connection and timeout settings for the [Client].
type Config struct {
	Addr        string        // Server network address (host:port)
	DialTimeout time.Duration // Maximum time allowed to establish a connection
	Timeout     time.Duration // Maximum time allowed for individual read/write operations
}

// Option represents a functional configuration option for modifying [Config] settings.
type Option func(*Config)

// WithAddr returns an [Option] that sets the target server network address.
func WithAddr(addr string) Option {
	return func(c *Config) {
		c.Addr = addr
	}
}

// WithDialTimeout returns an [Option] that sets the connection establishment timeout.
func WithDialTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.DialTimeout = timeout
	}
}

// WithTimeout returns an [Option] that sets the per-operation read/write deadline.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func defaultConfig() Config {
	return Config{
		Addr:        defaultAddr,
		DialTimeout: defaultDialTimeout,
		Timeout:     defaultTimeout,
	}
}

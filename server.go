package quizme

import (
	"log/slog"
	"net/http"

	"github.com/JDinABox/quiz-me/internal/dev"
)

type Config struct {
	Listen  string
	Logging bool
	Dev     bool
}

type Option func(*Config) error

// WithListenAddr sets the server's listen address
func WithListenAddr(addr string) Option {
	return func(c *Config) error {
		c.Listen = addr
		return nil
	}
}

// WithLogging sets the logging state
func WithLogging(logging bool) Option {
	return func(c *Config) error {
		c.Logging = logging
		return nil
	}
}

// NewConfig creates a Config instance with optional settings
func NewConfig(opts ...Option) (*Config, error) {
	// Set defaults
	cfg := &Config{
		Listen:  "127.0.0.1:8080",
		Logging: true,
		Dev:     dev.IsDev,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func Start(opts ...Option) error {
	// Create config from options
	conf, err := NewConfig(opts...)
	if err != nil {
		return err
	}

	if conf.Dev {
		slog.Warn("Running in development mode!")
	}
	// Initialize HTTP router
	router, err := newApp(conf)
	if err != nil {
		return err
	}

	// Start HTTP server
	server := &http.Server{
		Addr:    conf.Listen,
		Handler: router,
	}

	slog.Info("Listening on", "address", conf.Listen)
	return server.ListenAndServe()
}

package config

import (
	"crypto/tls"
	"fmt"
	"time"
)

type Config struct {
	Host         string
	Port         int
	Timeout      time.Duration
	MaxHeaderLen int
	JUnitReport  string
	Strict       bool
	DryRun       bool
	TLS          bool
	Insecure     bool
	Verbose      bool
	Sections     []string
	Targets      map[string]bool
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) TLSConfig() *tls.Config {
	if !c.TLS {
		return nil
	}

	config := tls.Config{
		InsecureSkipVerify: c.Insecure,
	}

	if config.NextProtos == nil {
		config.NextProtos = append(config.NextProtos, "h2", "h2-16")
	}

	return &config
}

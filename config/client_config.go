package config

import (
	"crypto/tls"
	"fmt"
	"time"
)

// Config represents the configuration of h2spec.
type ClientSpecConfig struct {
	Host         string
	Port         int
	Timeout      time.Duration
	MaxHeaderLen int
	JUnitReport  string
	Strict       bool
	DryRun       bool
	TLS          bool
	Verbose      bool
	Sections     []string
	CertFile     string
	CertKeyFile  string
	Exec         string
}

// Addr returns the string concatinated with hostname and port number.
func (c *ClientSpecConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ClientSpecConfig) Scheme() string {
	if c.TLS {
		return "https"
	} else {
		return "http"
	}
}

// TLSConfig returns a tls.Config based on the configuration of h2spec.
func (c *ClientSpecConfig) TLSConfig() (*tls.Config, error) {
	if !c.TLS {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(c.CertFile, c.CertKeyFile)
	if err != nil {
		return nil, err
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "h2-16"},
	}

	return &config, nil
}

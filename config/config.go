package config

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"
)

const (
	RunModeAll = iota
	RunModeGroup
	RunModeNone
)

// Config represents the configuration of h2spec.
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
	targetMap    map[string]bool
	CertFile     string
	CertKeyFile  string
	Exec         string
	FromPort     int
}

// Addr returns the string concatinated with hostname and port number.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) Scheme() string {
	if c.TLS {
		return "https"
	} else {
		return "http"
	}
}

// TLSConfig returns a tls.Config based on the configuration of h2spec.
func (c *Config) TLSConfig() (*tls.Config, error) {
	if !c.TLS {
		return nil, nil
	}

	config := tls.Config{
		InsecureSkipVerify: c.Insecure,
	}

	if config.NextProtos == nil {
		config.NextProtos = append(config.NextProtos, "h2", "h2-16")
	}

	if c.CertFile != "" && c.CertKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.CertKeyFile)
		if err != nil {
			return nil, err
		}
		config.Certificates = []tls.Certificate{cert}
	}

	return &config, nil
}

// RunMode returns a run mode of specified the section number.
// This is used to decide whether to run test cases.
func (c *Config) RunMode(section string) int {
	if c.targetMap == nil {
		c.buildTargetMap()
	}

	if len(c.targetMap) == 0 {
		return RunModeAll
	}

	comps := strings.Split(section, "/")
	compLen := len(comps)

	if compLen == 0 || compLen > 3 {
		return RunModeNone
	}

	keys := []string{comps[0]}

	if compLen > 1 {
		nums := strings.Split(comps[1], ".")
		for i, _ := range nums {
			key := fmt.Sprintf("%s/%s", comps[0], strings.Join(nums[:i+1], "."))
			keys = append(keys, key)
		}
	}

	if compLen > 2 {
		keys = append(keys, section)
	}

	var result int
	for _, key := range keys {
		val, ok := c.targetMap[key]
		if ok {
			if val {
				return RunModeAll
			}
			result = RunModeGroup
		} else {
			result = RunModeNone
		}
	}

	return result
}

func (c *Config) buildTargetMap() {
	c.targetMap = map[string]bool{}

	for _, section := range c.Sections {
		comps := strings.Split(section, "/")
		compLen := len(comps)

		// Validate the format of the section string.
		if compLen == 0 || compLen > 3 {
			fmt.Printf("Invalid section: %s", section)
			continue
		}

		// Check the section string is root section or not.
		if compLen == 1 {
			c.targetMap[comps[0]] = true
			continue
		}

		_, ok := c.targetMap[comps[0]]
		if !ok {
			c.targetMap[comps[0]] = false
		}

		// The parent group of the test case associated with this section
		// must only run test cases included in the group.
		nums := strings.Split(comps[1], ".")
		for i, _ := range nums {
			key := fmt.Sprintf("%s/%s", comps[0], strings.Join(nums[:i+1], "."))
			c.targetMap[key] = false
		}

		// The test case associated with this section string must be run.
		c.targetMap[section] = true
	}
}

func (c *Config) IsBrowserMode() bool {
	return c.Exec == ""
}

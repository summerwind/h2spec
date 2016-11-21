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

		// Invalid section
		if compLen == 0 || compLen > 3 {
			fmt.Println("Invalid section: %s", section)
			continue
		}

		// Root section
		if compLen == 1 {
			c.targetMap[comps[0]] = true
			continue
		}

		_, ok := c.targetMap[comps[0]]
		if !ok {
			c.targetMap[comps[0]] = false
		}

		nums := strings.Split(comps[1], ".")
		for i, _ := range nums {
			key := fmt.Sprintf("%s/%s", comps[0], strings.Join(nums[:i+1], "."))
			c.targetMap[key] = false
		}

		c.targetMap[section] = true
	}
}

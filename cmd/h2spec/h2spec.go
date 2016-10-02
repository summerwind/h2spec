package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/summerwind/h2spec"
	"github.com/summerwind/h2spec/config"
)

var (
	VERSION string = "0.0.0"
	COMMIT  string = "(none)"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "h2spec [section...]",
		Short: "Conformance testing tool for HTTP/2 implementation",
		Long:  "Conformance testing tool for HTTP/2 implementation.",
		RunE:  run,
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	flags := cmd.Flags()
	flags.StringP("host", "h", "127.0.0.1", "Target host")
	flags.IntP("port", "p", 0, "Target port")
	flags.IntP("timeout", "o", 2, "Maximum time allowed for test")
	flags.Int("max-header-length", 4000, "Maximum header length")
	flags.String("junit-report", "", "Generate JUnit test reports")
	flags.Bool("strict", false, "Strict mode")
	flags.Bool("dryrun", false, "Dry-run mode")
	flags.BoolP("tls", "t", false, "Connect over TLS")
	flags.BoolP("insecure", "k", false, "Don't verify server's certificate")
	flags.BoolP("verbose", "v", false, "Output verbose log")
	flags.Bool("version", false, "Display version information and exit")
	flags.Bool("help", false, "Display this help and exit")

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	v, err := flags.GetBool("version")
	if err != nil {
		return err
	}

	if v {
		version()
		return nil
	}

	host, err := flags.GetString("host")
	if err != nil {
		return err
	}

	port, err := flags.GetInt("port")
	if err != nil {
		return err
	}

	timeout, err := flags.GetInt("timeout")
	if err != nil {
		return err
	}

	maxHeaderLen, err := flags.GetInt("max-header-length")
	if err != nil {
		return err
	}

	junit, err := flags.GetString("junit-report")
	if err != nil {
		return err
	}

	strict, err := flags.GetBool("strict")
	if err != nil {
		return err
	}

	dryRun, err := flags.GetBool("dryrun")
	if err != nil {
		return err
	}

	tls, err := flags.GetBool("tls")
	if err != nil {
		return err
	}

	insecure, err := flags.GetBool("insecure")
	if err != nil {
		return err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return err
	}

	if port == 0 {
		if tls {
			port = 443
		} else {
			port = 80
		}
	}

	c := &config.Config{
		Host:         host,
		Port:         port,
		Timeout:      time.Duration(timeout) * time.Second,
		MaxHeaderLen: maxHeaderLen,
		JUnit:        junit,
		Strict:       strict,
		DryRun:       dryRun,
		TLS:          tls,
		Insecure:     insecure,
		Verbose:      verbose,
		Sections:     args,
		Targets:      buildTargets(args),
	}

	return h2spec.Run(c)
}

func buildTargets(sections []string) map[string]bool {
	targets := map[string]bool{}

	for _, section := range sections {
		comps := strings.Split(section, "/")
		compLen := len(comps)

		// Invalid section
		if compLen == 0 || compLen > 3 {
			continue
		}

		// Root section
		if compLen == 1 {
			targets[comps[0]] = true
			continue
		}

		targets[comps[0]] = false

		nums := strings.Split(comps[1], ".")
		for i, _ := range nums {
			key := fmt.Sprintf("%s/%s", comps[0], strings.Join(nums[:i+1], "."))
			_, ok := targets[key]
			if !ok {
				targets[key] = false
			}
		}

		targets[section] = true
	}

	return targets
}

func version() {
	fmt.Printf("Version: %s\n", VERSION)
	fmt.Printf("Commit:  %s\n", COMMIT)
}

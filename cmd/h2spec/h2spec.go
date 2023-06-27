package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/summerwind/h2spec"
	"github.com/summerwind/h2spec/config"
)

var (
	VERSION string = "2.0.0"
	COMMIT  string = "(Unknown)"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "h2spec [spec...]",
		Short: "Conformance testing tool for HTTP/2 implementation",
		Long:  "Conformance testing tool for HTTP/2 implementation.",
		RunE:  run,
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	flags := cmd.Flags()
	flags.StringP("host", "h", "127.0.0.1", "Target host")
	flags.IntP("port", "p", 0, "Target port")
	flags.StringP("server-name", "n", "", "Server name (SNI)")
	flags.StringP("path", "P", "/", "Target path")
	flags.IntP("timeout", "o", 2, "Time seconds to test timeout")
	flags.Int("max-header-length", 4000, "Maximum length of HTTP header")
	flags.StringP("junit-report", "j", "", "Path for JUnit test report")
	flags.BoolP("strict", "S", false, "Run all test cases including strict test cases")
	flags.Bool("dryrun", false, "Display only the title of test cases")
	flags.BoolP("tls", "t", false, "Connect over TLS")
	flags.StringP("ciphers", "c", "", "List of colon-separated TLS cipher names")
	flags.BoolP("insecure", "k", false, "Don't verify server's certificate")
	flags.StringSliceP("exclude", "x", []string{}, "Disable specific tests")
	flags.StringSliceP("execute-specific-tests", "", []string{}, "Exeute specific tests")
	flags.Bool("exit-on-external-failure", false, "Stop tests execution on an external failure event")
	flags.String("external-failure-source", "", "Path to the file that needs to be tracked for failures")
	flags.String("external-failure-regexp", "", "A regular expression for a falure to be marched with")
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

	path, err := flags.GetString("path")
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

	junitReport, err := flags.GetString("junit-report")
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

	ciphers, err := flags.GetString("ciphers")
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

	exclude, err := flags.GetStringSlice("exclude")
	if err != nil {
		return err
	}

	executeSpecificTests, err  := flags.GetStringSlice("execute-specific-tests")
	if err != nil {
		return err
	}

	exitOnExternalFailure, err := flags.GetBool("exit-on-external-failure")
	if err != nil {
		return err
	}

	externalFailureSource, err := flags.GetString("external-failure-source")
	if err != nil {
		return err
	}

	externalFailureRegexp, err := flags.GetString("external-failure-regexp")
	if err != nil {
		return err
	}

	serverName, err := flags.GetString("server-name")
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
		ServerName:   serverName,
		Path:         path,
		Timeout:      time.Duration(timeout) * time.Second,
		MaxHeaderLen: maxHeaderLen,
		JUnitReport:  junitReport,
		Strict:       strict,
		DryRun:       dryRun,
		TLS:          tls,
		Ciphers:      ciphers,
		Insecure:     insecure,
		Verbose:      verbose,
		Sections:     args,
		Excluded:     exclude,
		ExecuteSpecificTests:  executeSpecificTests,
		ExitOnExternalFailure: exitOnExternalFailure,
		ExternalFailureSource: externalFailureSource,
		ExternalFailureRegexp: externalFailureRegexp,
	}

	success, err := h2spec.Run(c)
	if !success {
		os.Exit(1)
	}

	return err
}

func version() {
	fmt.Printf("Version: %s (%s)\n", VERSION, COMMIT)
}

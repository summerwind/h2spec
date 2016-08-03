package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/summerwind/h2spec"
)

const VERSION = "v1.4.1"

type sections []string

func (s *sections) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *sections) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	port := flag.Int("p", 0, "Target port.")
	host := flag.String("h", "127.0.0.1", "Target host.")
	useTls := flag.Bool("t", false, "Connect over TLS.")
	insecureSkipVerify := flag.Bool("k", false, "Don't verify server's certificate.")
	timeout := flag.Int("o", 2, "Maximum time allowed for test.")
	strict := flag.Bool("S", false, "Strict mode.")
	junit := flag.String("j", "", "Create test report also in JUnit format.")
	version := flag.Bool("version", false, "Display version information and exit.")

	var sectionFlag sections
	flag.Var(&sectionFlag, "s", "Section number on which to run the test")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Println("  -p:        Target port. (Default: 80 or 443)")
		fmt.Println("  -h:        Target host. (Default: 127.0.0.1)")
		fmt.Println("  -t:        Connect over TLS. (Default: false)")
		fmt.Println("  -k:        Don't verify server's certificate. (Default: false)")
		fmt.Println("  -o:        Maximum time allowed for test. (Default: 2)")
		fmt.Println("  -s:        Section number on which to run the test. (Example: -s 6.1 -s 6.2)")
		fmt.Println("  -S:        Run the test cases marked as \"strict\".")
		fmt.Println("  -j:        Creates report also in JUnit format into specified file.")
		fmt.Println("  --version: Display version information and exit.")
		fmt.Println("  --help:    Display this help and exit.")
		os.Exit(1)
	}

	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stderr, "h2spec %s\n", VERSION)
		os.Exit(0)
	}

	if *port == 0 {
		if *useTls {
			*port = 443
		} else {
			*port = 80
		}
	}

	var ctx h2spec.Context
	ctx.Port = *port
	ctx.Host = *host
	ctx.Timeout = time.Duration(*timeout) * time.Second
	ctx.Strict = *strict
	ctx.Junit = *junit
	ctx.Tls = *useTls
	ctx.TlsConfig = &tls.Config{
		InsecureSkipVerify: *insecureSkipVerify,
	}

	if len(sectionFlag) > 0 {
		ctx.Sections = map[string]bool{}
		for _, sec := range sectionFlag {
			splitedSec := strings.Split(sec, ".")
			lastIndex := len(splitedSec) - 1

			num := []string{}
			for i, sec := range splitedSec {
				num = append(num, sec)
				if i != 0 {
					key := strings.Join(num, ".")

					runAll := false
					if i == lastIndex || ctx.Sections[key] {
						runAll = true
					}

					ctx.Sections[key] = runAll
				}
			}
		}
	}

	if !h2spec.Run(&ctx) {
		os.Exit(1)
	}
}

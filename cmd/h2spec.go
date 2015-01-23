package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/summerwind/h2spec"
	"os"
	"time"
)

type sections []string

func (s *sections) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *sections) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	port := flag.Int("p", 80, "Target port")
	host := flag.String("h", "127.0.0.1", "Target host")
	useTls := flag.Bool("t", false, "Connect over TLS")
	insecureSkipVerify := flag.Bool("k", false, "Don't verify server's certificate")
	timeout := flag.Int("o", 3, "Maximum time allowed for test. (Default: 3)")

	var sectionFlag sections
	flag.Var(&sectionFlag, "s", "Section number on which to run the test")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Println("  -p:     Target port. (Default: 80)")
		fmt.Println("  -h:     Target host. (Default: 127.0.0.1)")
		fmt.Println("  -t:     Connect over TLS. (Default: false)")
		fmt.Println("  -k:     Don't verify server's certificate. (Default: false)")
		fmt.Println("  -o:     Maximum time allowed for test. (Default: 3)")
		fmt.Println("  -s:     Section number on which to run the test. (Example: -s 6.1 -s 6.2)")
		fmt.Println("  --help: Display this help and exit.")
		os.Exit(1)
	}

	flag.Parse()

	var ctx h2spec.Context
	ctx.Port = *port
	ctx.Host = *host
	ctx.Timeout = time.Duration(*timeout) * time.Second
	ctx.Tls = *useTls
	ctx.TlsConfig = &tls.Config{
		InsecureSkipVerify: *insecureSkipVerify,
	}

	if len(sectionFlag) > 0 {
		ctx.Sections = map[string]bool{}
		for _, v := range sectionFlag {
			ctx.Sections[v] = true
		}
	}

	h2spec.Run(&ctx)
}

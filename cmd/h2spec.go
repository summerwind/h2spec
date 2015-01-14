package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/summerwind/h2spec"
	"os"
)

func main() {
	port := flag.Int("p", 80, "Target port")
	host := flag.String("h", "127.0.0.1", "Target host")
	useTls := flag.Bool("t", false, "Connect over TLS")
	insecureSkipVerify := flag.Bool("k", false, "Don't verify server's certificate")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Println("  -p:     Target port. (Default: 80)")
		fmt.Println("  -h:     Target host. (Default: 127.0.0.1)")
		fmt.Println("  -t:     Connect over TLS. (Default: false)")
		fmt.Println("  -k:     Don't verify server's certificate. (Default: false)")
		fmt.Println("  --help: Display this help and exit.")
		os.Exit(1)
	}

	flag.Parse()

	var ctx h2spec.Context
	ctx.Port = *port
	ctx.Host = *host
	ctx.Tls = *useTls
	ctx.TlsConfig = &tls.Config{
		InsecureSkipVerify: *insecureSkipVerify,
	}

	h2spec.Run(&ctx)
}

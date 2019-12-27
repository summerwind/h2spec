# h2spec

h2spec is a conformance testing tool for HTTP/2 implementation.  
This tool is compliant with [RFC 7540 (HTTP/2)](http://www.rfc-editor.org/rfc/rfc7540.txt) and [RFC 7541 (HPACK)](http://www.rfc-editor.org/rfc/rfc7541.txt).

## Install

Go to the [releases page](https://github.com/summerwind/h2spec/releases), find the version you want, and download the zip file or tarball file. The docker image is also available in [Docker Hub](https://hub.docker.com/r/summerwind/h2spec/).

## Your server

Your server should respond on `GET /` or `POST /` requests with status 200 response with non-empty data.

## Usage

```
Conformance testing tool for HTTP/2 implementation.

Usage:
  h2spec [spec...] [flags]

Flags:
  -c, --ciphers string          List of colon-separated TLS cipher names
      --dryrun                  Display only the title of test cases
      --help                    Display this help and exit
  -h, --host string             Target host (default "127.0.0.1")
  -k, --insecure                Don't verify server's certificate
  -j, --junit-report string     Path for JUnit test report
      --max-header-length int   Maximum length of HTTP header (default 4000)
  -P, --path string             Target path (default "/")
  -p, --port int                Target port
  -S, --strict                  Run all test cases including strict test cases
  -o, --timeout int             Time seconds to test timeout (default 2)
  -t, --tls                     Connect over TLS
  -v, --verbose                 Output verbose log
      --version                 Display version information and exit
```

### Running a specific test case

You can choose a test case to run by specifying the *Spec ID* as the command argument. For example, if you want to run test cases for HTTP/2, run h2spec as following:

```
$ h2spec http2
```

If you add a section number after the *Spec ID*, test cases related to a specific section will be run. For example, if you want to run test cases related to 6.3 of HTTP/2, run h2spec as following:

```
$ h2spec http2/6.3
```

If you add a test number after the section number, you can run the specific test case individually. For example, to run only the first test case related to 6.3 of HTTP/2 6.3, run h2spec as following:

```
$ h2spec http2/6.3/1
```

The *Spec ID* can be specified multiple times.

```
$ h2spec http2/6.3 generic
```

Currently supported *Spec IDs* are as follows. `generic` is the original spec of h2spec, includes generic test cases for HTTP/2 servers.

Spec ID | Description
--- | ---
http2 | Test cases for RFC 7540 (HTTP/2)
hpack | Test cases for RFC 7541 (HPACK)
generic | Generic test cases for HTTP/2 servers

### Dryrun Mode

To display the list of test cases to be run, use *Dryrun Mode* as follows:

```
$ h2spec --dryrun
```

### Strict Mode

When *Strict Mode* is enabled, h2spec will run the test cases related to the contents requested with the `SHOULD` notation in each specification. It is useful for more rigorous verification of HTTP/2 implementation.

```
$ h2spec --strict
```

## Screenshot

![Sceenshot](https://cloud.githubusercontent.com/assets/230145/22183160/9e9fbb4c-e0fa-11e6-9383-e2cc1ed6750a.png)

## Build

To build from source, you need to install [Go](https://golang.org) and export `GO111MODULE=on` first.

To build:
```
$ make build
```

To test:
```
$ make test
```

## License

h2spec is made available under MIT license.

package spec

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
	"golang.org/x/net/http2"
)

type Server struct {
	listeners []net.Listener
	config    *config.Config
	testCases map[int]*ClientTestCase
	spec      *ClientTestGroup
}

func Listen(c *config.Config, tg *ClientTestGroup) (*Server, error) {
	testCases := make(map[int]*ClientTestCase)
	tg.ClientTestCases(testCases, c.FromPort)

	server := &Server{
		listeners: make([]net.Listener, 0),
		config:    c,
		testCases: testCases,
		spec:      tg,
	}

	for port, tc := range testCases {
		var err error
		var listener net.Listener

		addr := fmt.Sprintf("%s:%d", c.Host, port)

		if c.TLS {
			tlsConfig, err := c.TLSConfig()
			if err != nil {
				return nil, err
			}

			listener, err = tls.Listen("tcp", addr, tlsConfig)
			if err != nil {
				return nil, err
			}
		} else {
			listener, err = net.Listen("tcp", addr)
			if err != nil {
				return nil, err
			}
		}

		server.listeners = append(server.listeners, listener)
		go server.RunListener(listener, tc)
	}

	return server, nil
}

func (server *Server) RunListener(listener net.Listener, tc *ClientTestCase) {
	for {
		baseConn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		conn, err := Accept(server.config, baseConn)
		if err != nil {
			log.Println(err)
			continue
		}

		go server.handleConn(conn, tc)
	}
}

func (server *Server) RunForever() {
	http.HandleFunc("/", server.home)
	http.ListenAndServe(server.config.Addr(), nil)
}

func (server *Server) Close() {
	for _, listener := range server.listeners {
		listener.Close()
	}
}

func (server *Server) handleConn(conn *Conn, tc *ClientTestCase) {
	start := time.Now()

	err := conn.Handshake()
	if err != nil {
		log.Println(red(err))
		return
	}
	request, err := conn.ReadRequest()
	if err != nil {
		log.Println(red(err))
		return
	}

	err = tc.Run(server.config, conn, request)
	end := time.Now()

	// Ensure that connection had been closed
	go closeConn(conn, request.StreamID)

	log.ResetLine()

	tr := NewClientTestResult(tc, err, end.Sub(start))
	tr.Print()

	if tc.Result != nil {
		tc.Parent.IncRecursive(tc.Result.Failed, tc.Result.Skipped, -1)
	}

	tc.Result = tr
	tc.Parent.IncRecursive(tc.Result.Failed, tc.Result.Skipped, 1)
}

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlReport(server.spec, server.config))
}

func htmlReport(tg *ClientTestGroup, c *config.Config) string {
	var buffer bytes.Buffer

	passed := tg.PassedCount
	failed := tg.FailedCount
	skipped := tg.SkippedCount

	total := passed + failed + skipped
	tmp := "<div>%d tests, %d passed, %d skipped, %d failed</div>"
	buffer.WriteString(fmt.Sprintf(tmp, total, passed, skipped, failed))

	buffer.WriteString(htmlReportForTestGroup(tg, c))
	return buffer.String()
}

func htmlReportForTestGroup(tg *ClientTestGroup, c *config.Config) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<div>%s</div>", tg.Title()))

	for _, tc := range tg.Tests {
		buffer.WriteString(htmlReportForTestCase(tc, c))
	}

	for _, g := range tg.Groups {
		buffer.WriteString(htmlReportForTestGroup(g, c))
	}

	buffer.WriteString("<br>")
	return buffer.String()
}

func htmlReportForTestCase(tc *ClientTestCase, c *config.Config) string {
	formatter := "<div>%s<a href=\"%s\" target=\"_blank\">%s</a>%s</div>"

	tr := tc.Result

	if tr == nil {
		resultLabel := "<span style=\"color: red;\">&nbsp;&nbsp;</span>"
		return fmt.Sprintf(formatter, resultLabel, tc.FullPath(c), tc.FullPath(c), tc.Desc)
	}

	if !tr.Failed {
		resultLabel := "<span style=\"color: green;\">✔</span>"
		return fmt.Sprintf(formatter, resultLabel, tc.FullPath(c), tc.FullPath(c), tc.Desc)
	}

	var buffer bytes.Buffer

	resultLabel := "<span style=\"color: red;\">✖</span>"
	buffer.WriteString(fmt.Sprintf(formatter, resultLabel, tc.FullPath(c), tc.FullPath(c), tc.Desc))

	err, ok := tr.Error.(*TestError)
	formatter = "<div style=\"padding-left: %dpx; color: %s\">%s</div>"
	if ok {
		msg := fmt.Sprintf("-> %s", tc.Requirement)
		buffer.WriteString(fmt.Sprintf(formatter, 20, "red", msg))

		label := "Expected:"
		for i, ex := range err.Expected {
			if i != 0 {
				label = strings.Repeat("&nbsp;", len(label))
			}
			msg = fmt.Sprintf("%s&nbsp;%s", label, ex)
			buffer.WriteString(fmt.Sprintf(formatter, 30, "yellow", msg))
		}
		msg = fmt.Sprintf("&nbsp;&nbsp;Actual:&nbsp;%s", err.Actual)
		buffer.WriteString(fmt.Sprintf(formatter, 30, "green", msg))

	} else if err != nil {
		errMsg := fmt.Sprintf("Error: %v", err)
		buffer.WriteString(fmt.Sprintf(formatter, 20, "red", errMsg))
	} else {
		errMsg := fmt.Sprintf("Error: %v", tr.Error.Error())
		buffer.WriteString(fmt.Sprintf(formatter, 20, "red", errMsg))
	}

	return buffer.String()
}

func closeConn(conn *Conn, lastStreamID uint32) {
	if !conn.Closed {
		conn.WriteGoAway(lastStreamID, http2.ErrCodeNo, make([]byte, 0))
		time.Sleep(1 * time.Second)
	}

	conn.Close()
}

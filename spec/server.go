package spec

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

type Server struct {
	net.Listener

	config    *config.ClientSpecConfig
	testCases map[string]*ClientTestCase
	spec      *ClientTestGroup
}

func Listen(c *config.ClientSpecConfig, tg *ClientTestGroup) (*Server, error) {
	var err error
	var listener net.Listener
	if c.TLS {
		tlsConfig, err := c.TLSConfig()
		if err != nil {
			return nil, err
		}

		listener, err = tls.Listen("tcp", c.Addr(), tlsConfig)
		if err != nil {
			return nil, err
		}

		log.Println(fmt.Sprintf("Server is listened at https://%s", c.Addr()))
	} else {
		listener, err = net.Listen("tcp", c.Addr())
		if err != nil {
			return nil, err
		}

		log.Println(fmt.Sprintf("Server is listened at http://%s", c.Addr()))
	}

	testCases := make(map[string]*ClientTestCase)
	tg.ClientTestCases(testCases)

	server := Server{
		Listener:  listener,
		config:    c,
		testCases: testCases,
		spec:      tg,
	}
	return &server, nil
}

func (server *Server) RunForever() {
	for {
		baseConn, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		conn, err := Accept(server.config, baseConn)
		if err != nil {
			log.Println(err)
			continue
		}

		go server.handleConn(conn)
	}
}

func (server *Server) handleConn(conn *Conn) {
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

	var path string
	hasPath := false
	for _, f := range request.Headers {
		if f.Name == ":path" {
			path = f.Value
			hasPath = true
		}
	}

	if !hasPath {
		log.Println(red("No :path found in request"))
		return
	}

	if path == "/" {
		server.home(conn, request)
		return
	}

	tc, ok := server.testCases[path]
	if !ok {
		log.Println(red(fmt.Sprintf("No path match: %s", path)))
		server.notFound(conn, request)
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

func (server *Server) home(conn *Conn, req *Request) {
	hp := http2.HeadersFrameParam{
		StreamID:      req.StreamID,
		EndStream:     false,
		EndHeaders:    true,
		BlockFragment: conn.EncodeHeaders(CommonRespHeaders(server.config)),
	}

	conn.WriteHeaders(hp)

	report := htmlReport(server.spec)
	conn.WriteData(req.StreamID, true, report)

	go closeConn(conn, req.StreamID)
}

func (server *Server) notFound(conn *Conn, req *Request) {
	headers := []hpack.HeaderField{
		HeaderField(":status", "404"),
	}

	hp := http2.HeadersFrameParam{
		StreamID:      req.StreamID,
		EndStream:     true,
		EndHeaders:    true,
		BlockFragment: conn.EncodeHeaders(headers),
	}

	conn.WriteHeaders(hp)
	conn.WriteGoAway(req.StreamID, http2.ErrCodeNo, make([]byte, 0))
	conn.Close()
}

func htmlReport(tg *ClientTestGroup) []byte {
	var buffer bytes.Buffer

	passed := tg.PassedCount
	failed := tg.FailedCount
	skipped := tg.SkippedCount

	total := passed + failed + skipped
	tmp := "<div>%d tests, %d passed, %d skipped, %d failed</div>"
	buffer.WriteString(fmt.Sprintf(tmp, total, passed, skipped, failed))

	buffer.WriteString(htmlReportForGroup(tg))
	return buffer.Bytes()
}

func htmlReportForGroup(tg *ClientTestGroup) string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<div>%s</div>", tg.Title()))

	for _, tc := range tg.Tests {
		tmp := "<div><a href=\"%s\" target=\"_blank\">%s</a>&nbsp;%s</div>"
		buffer.WriteString(fmt.Sprintf(tmp, tc.Path(), tc.Path(), tc.Desc))
	}

	for _, g := range tg.Groups {
		buffer.WriteString(htmlReportForGroup(g))
	}

	buffer.WriteString("<br>")
	return buffer.String()
}

func closeConn(conn *Conn, lastStreamID uint32) {
	if conn.Closed {
		return
	}

	conn.WriteGoAway(lastStreamID, http2.ErrCodeNo, make([]byte, 0))
	time.Sleep(3 * time.Second)
	conn.Close()
}

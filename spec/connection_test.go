package spec

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/summerwind/h2spec/config"
)

// TestConnWaitEvent verifies the ConnectionClosedEvent
func TestConnWaitEvent(t *testing.T) {
	srv, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Errorf("Failed to start server: %v", err)
	}
	srvAddr, err := net.ResolveTCPAddr("tcp", srv.Addr().String())
	if err != nil {
		t.Errorf("Failed to get server port: %v", err)
	}

	// Accept client but disconnect immediately
	go func() {
		conn, err := srv.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		// Force a RST when closing
		tcpConn := conn.(*net.TCPConn)
		tcpConn.SetLinger(0)

		conn.Close()
	}()

	config := config.Config{
		Port:    srvAddr.Port,
		Timeout: 2 * time.Second,
		TLS:     false,
	}

	conn, err := Dial(&config)
	if err != nil {
		t.Errorf("Failed to Dial: %v", err)
	}

	ev := conn.WaitEvent()
	if _, ok := ev.(ConnectionClosedEvent); !ok {
		t.Errorf("Expected Connection closed event but got '%v'", ev)
	}
}

package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = "37"
	timeGap1900 = 2208988800
)

// options set of CLI options
type options struct {
	port string
}

func parseFlags() options {
	var opt options
	flag.StringVar(&opt.port, "p", defaultPort, "port")
	flag.Parse()
	return opt
}

func main() {

	opt := parseFlags()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		<-c
		log.Println("Application shutdown")
		cancel()

	}()

	log.Println("Start application")

	// start server
	s, err := newServer(ctx, defaultHost, opt.port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(s.serve())

}

// server is a tcp time server
type server struct {
	lst     net.Listener
	current func() time.Time
	ctx     context.Context
}

// newServer create a new time server
func newServer(ctx context.Context, host, port string) (*server, error) {
	lst, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}

	s := &server{
		lst:     lst,
		current: time.Now,
		ctx:     ctx,
	}

	go func() {
		// close connection on context close
		<-ctx.Done()
		s.close()
	}()

	return s, nil
}

// close server
func (s *server) close() {
	s.lst.Close()
}

// serve starts incoming request handling
func (s *server) serve() error {
	for {
		// listen for tcp connection
		conn, err := s.lst.Accept()
		if err != nil {
			return err
		}
		s.handleConnection(conn)
	}
}

// handleConnection writes current time in RFC 868 format into connection
func (s *server) handleConnection(conn net.Conn) {
	// read message
	message := make([]byte, 0)
	bufio.NewReader(conn).Read(message)
	if len(message) == 0 {
		message = []byte("<empty>")
	}
	log.Println("Request: " + string(message))

	res := make([]byte, 4)
	// The time is the number of seconds since 00:00 (midnight) 1 January 1900 GMT
	now := getTimeRFC868(s.current)
	binary.BigEndian.PutUint32(res, uint32(now))

	// write response and close connection
	conn.Write(res)
	conn.Close()
}

// getTimeRFC868 returns time in format of RFC868
func getTimeRFC868(t func() time.Time) int64 {
	// The time is the number of seconds since 00:00 (midnight) 1 January 1900 GMT
	return t().Unix() + timeGap1900
}

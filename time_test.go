package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	current := time.Now

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		lst, err := net.Listen("tcp", "localhost:3000")
		if err != nil {
			t.Fatal(err)
		}
		s := server{
			lst:     lst,
			current: current,
			ctx:     ctx,
		}
		s.serve()
	}()

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		t.Fatal(err)
	}

	got := make([]byte, 4)
	bufio.NewReader(conn).Read(got)
	t.Log(string(got))

	conn.Close()

	want := make([]byte, 4)
	binary.BigEndian.PutUint32(want, uint32(getTimeRFC868(current)))

	gotS, wantS := string(got), string(want)
	if gotS != wantS {
		t.Errorf("wrong time format %s, want %s", gotS, wantS)
	}

	cancel()

}

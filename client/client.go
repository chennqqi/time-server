package client

import (
	"bufio"
	"encoding/binary"
	"net"
	"strconv"
	"time"
)

const (
	timeGap1900 = 2208988800
)

// Dial sends request to the time server
func Dial(host, port string) (*string, error) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msg := make([]byte, 4)
	bufio.NewReader(conn).Read(msg)

	// parse time into Unix time
	t := parseTime(msg)

	res := strconv.Itoa(int(t.Unix()))

	return &res, nil

}

// parseTime parses and converts time from RFC868 to Unix time
// note, that wrong input will be converted into wrong time
func parseTime(msg []byte) *time.Time {
	// convert input to integer
	i := binary.BigEndian.Uint32(msg)

	// convert time from 1900 to Unix time
	t := int64(i) - timeGap1900

	res := time.Unix(t, 0)
	return &res
}

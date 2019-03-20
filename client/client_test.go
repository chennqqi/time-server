package client

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	type args struct {
		host string
		port string
		time string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "epoch time",
			args: args{
				host: "localhost",
				port: "3000",
				time: "1970-01-01T00:00:00+00:00",
			},
			want:    "0",
			wantErr: false,
		},

		{
			name: "1900 time",
			args: args{
				host: "localhost",
				port: "3000",
				time: "1900-01-01T00:00:00+00:00",
			},
			want:    "-2208988800",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			lst, err := net.Listen("tcp", tt.args.host+":"+tt.args.port)
			if err != nil {
				t.Fatal(err)
			}

			go func(t *testing.T) {
				conn, err := lst.Accept()
				if err != nil {
					t.Fatal(err)
				}
				tm, _ := time.Parse(time.RFC3339, tt.args.time)
				r := tm.Unix() + timeGap1900
				msg := make([]byte, 4)
				binary.BigEndian.PutUint32(msg, uint32(r))
				conn.Write(msg)
				conn.Close()
			}(t)

			got, err := Dial(tt.args.host, tt.args.port)

			lst.Close()

			if (err != nil) != tt.wantErr {
				t.Errorf("Wrong error behavior %v, but expected %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("Wrong result %v, want %v", *got, tt.want)
			}

		})
	}
}

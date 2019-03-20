package main

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/time-server/client"
)

// main run the client
//
// positional parameters:
// host: server host;
// port: server port;
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("wrong parameters\nUsage: %s host port\n", os.Args[0])
		os.Exit(2)
	}

	msg, err := client.Dial(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(*msg)
}

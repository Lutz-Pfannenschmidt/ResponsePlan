package tty2web

import (
	"math/rand"
	"net"
	"strconv"
)

// getRandomPort returns a random available port in the range 35000-65000.
// It will return -1 if no port is available.
func getRandomPort() int {
	min := 35000
	max := 65000

	var port int

	for {
		port = min + rand.Intn(max-min)

		// Check if the port is in use on localhost
		conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			conn.Close()
			return port
		}
	}
}

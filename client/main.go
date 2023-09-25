package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	i := 0
	for {
		time.Sleep(1 * time.Second)
		conn.Write([]byte(fmt.Sprintf("publish Topic%d", i)))
		i++
	}

}

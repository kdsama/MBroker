package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	listner, e := net.Listen("tcp", "localhost:9999")
	if e != nil {
		fmt.Println(e)
	}

	for {
		conn, err := listner.Accept()
		fmt.Println("accepted Connection", conn.LocalAddr())
		if err != nil {
			log.Fatal(err)
		}
		go HandleConnection(conn)
	}

}
func HandleConnection(conn net.Conn) {
	msg := make([]byte, 1024)
	defer conn.Close()
	b, e := conn.Read(msg)
	fmt.Println(b, e)
	fmt.Println(string(msg))

}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	addr  = "localhost:12345"
	proto = "tcp4"
)

func main() {
	conn, err := net.Dial(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	reader := bufio.NewReader(conn)
	id := 0

	for {
		pvb, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		id++
		msg := strings.Trim(string(pvb), "\n")
		msg = strings.Trim(msg, "\r")
		fmt.Printf("Найдена поговорка № %d: %s\n", id, msg)
	}

}

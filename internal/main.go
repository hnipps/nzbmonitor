package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/hnipps/nzbrefresh/pkg/refresh"
)

var q chan string

func main() {
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	q = make(chan string)
	go consumer(q)

	fmt.Println("Service running on localhost:6666")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		response := processCommand(command)
		fmt.Fprintf(conn, "%s\n", response)
	}
}

func processCommand(nzbpath string) string {
	q <- nzbpath
	return "Queued"
}

func consumer(queue <-chan string) {
	for item := range queue {
		fmt.Printf("Processing: %s\n", item)
		refresh.Prepare(
			refresh.WithNZBFile(item),
			refresh.WithCheckOnly(true),
			refresh.WithDebug(true),
			refresh.WithProvider("../provider.json"),
		)
		refresh.Run()
	}
}

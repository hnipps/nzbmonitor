package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hnipps/nzbrefresh/pkg/refresh"
)

var q chan string

func main() {
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	q = make(chan string, 100)

	refresh.Prepare(
		refresh.WithCheckOnly(true),
		refresh.WithDebug(true),
		refresh.WithProvider("./provider.json"),
		refresh.WithMode("pkg"),
	)

	go consumer(q)

	log.Println("Service running on localhost:6666")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
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
	log.Printf("Queued: %s\n", nzbpath)
	return "Queued"
}

func consumer(queue <-chan string) {
	for item := range queue {
		go func(item string) {
			log.Printf("Processing: %s\n", item)
			if _, err := refresh.Run(item); err != nil {
				log.Printf("Failed: %s\n", item)
			} else {
				log.Printf("Completed: %s\n", item)
			}
		}(item)
	}
}

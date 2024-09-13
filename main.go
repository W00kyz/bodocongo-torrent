package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	PORT       = 10000
	PEERS_FILE = "peers.txt"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("Mensagem recebida: %s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("Erro na leitura:", err)
	}
}

func startServer(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Printf("Servidor ouvindo na porta %d\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Erro na aceitação de conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	go startServer(PORT)

	for {
	}
}

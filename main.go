package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	PORT       = 10000
	PEERS_FILE = "peers.txt"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Conexão recebida de ", conn.RemoteAddr().String())

	fileHash, _ := bufio.NewReader(conn).ReadString('\n')
	fileHash = strings.TrimSpace(fileHash)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("Mensagem recebida: %s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("Erro na leitura:", err)
	}
}

// registerIP
func registerIP(address string) {
	file, err := os.OpenFile(PEERS_FILE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("PEERS_FILE não pôde ser aberto.")
	}
	defer file.Close()

	_, err = file.WriteString(address + "\n")
	if err != nil {
		log.Fatal("erro ao escrever no arquivo: ", err)
	}
}

func startServer(port int) {
	selfAddress := GetLocalIP()
	listenAddress := fmt.Sprintf("%s:%d", selfAddress, port)
	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	registerIP(listenAddress)

	fmt.Printf("Servidor está ouvindo no endereço %s\n", listenAddress)

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

	// pegar os pares e o file_hash
	// se houver pares, executa a busca com file_hash e imprime as máquinas encontradas.
	// senão, não imprime nada pois não há pares.
	//
}

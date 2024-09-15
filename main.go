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
	PEERS_FILE = "peers.txt"
)

type conn struct {
	ip   string
	port string
	conn net.Conn
}

func handleConnection(conn net.Conn, peersCh chan net.Conn) {
	fmt.Println("Conexão recebida de ", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("Mensagem recebida: %s\n", msg)

		response := "Recebido: " + msg
		_, err := conn.Write([]byte(response + "\n"))
		if err != nil {
			log.Println("Erro ao enviar resposta:", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Erro na leitura:", err)
	}

	peersCh <- conn
}

func startServer(ip string, port string, peersCh chan net.Conn) {
	selfAddress := ip + ":" + port
	ln, err := net.Listen("tcp", selfAddress)
	if err != nil {
		connectToPeer(ip, port, peersCh)
	} else {
		defer ln.Close()

		fmt.Printf("Servidor está ouvindo no endereço %s\n", selfAddress)

		conn, err := ln.Accept()
		if err != nil {
			log.Println("Erro na aceitação de conexão:", err)
			return
		}
		go handleConnection(conn, peersCh)
	}
}

func connectToPeer(address string, port string, peersCh chan net.Conn) {
	selfAddress := address + ":" + port
	conn, err := net.Dial("tcp", selfAddress)
	if err != nil {
		log.Printf("Erro ao conectar ao peer %s: %v", address, err)
		return
	}

	fmt.Println("Conectado ao peer:", address)

	peersCh <- conn
}

func startAllConnections(peersCh chan net.Conn) {
	file, err := os.Open(PEERS_FILE)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo de peers: %v", err)
	}
	defer file.Close()

	var peers []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		address := scanner.Text()
		address = strings.TrimSpace(address)
		if address != "" {
			peers = append(peers, address)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Erro ao ler o arquivo de peers: %v", err)
	}

	for _, peer := range peers {
		parts := strings.Split(peer, ":")
		ip := parts[0]
		port := parts[1]
		go startServer(ip, port, peersCh)
	}

}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %v", err)
		return
	}
	fmt.Printf("Mensagem enviada para %s: %s\n", conn.RemoteAddr().String(), message)
}

func main() {
	peersCh := make(chan net.Conn)
	startAllConnections(peersCh)
	a := <-peersCh
	sendMessage(a, "oi")
}

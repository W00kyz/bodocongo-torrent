package main

import (
	"log"
	"net"
	"os"
)

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal("IP não foi obtido!")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Hostname não foi obtido!")
	}

	return hostname
}

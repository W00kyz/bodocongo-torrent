package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	PEERS_FILE         = "peers.txt"
	SHARED_FILES_FILE  = "shared-files.txt"
	HASH_RECALCULATION = 30 * time.Minute
)

var mySharedFilesPathsHashMap map[int][]string
var mutex sync.Mutex

func startHashRecalculation(directories []string, interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			newHashMap := processDirectories(directories)

			mutex.Lock()
			mySharedFilesPathsHashMap = newHashMap
			mutex.Unlock()

			fmt.Println("Hashes recalculated and updated.")
		}
	}()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection received from ", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("Message received: %s\n", msg)

		if strings.HasPrefix(msg, "SEARCH") {
			parts := strings.Split(msg, " ")
			hash, err := strconv.Atoi(parts[1])
			if err == nil {
				mutex.Lock()
				_, exists := mySharedFilesPathsHashMap[hash]
				mutex.Unlock()

				if exists {
					response := "FOUND|" + GetHostname() + "|"
					_, err := conn.Write([]byte(response + "\n"))
					if err != nil {
						log.Println("Error while sending message:", err)
					}
				} else {
					response := "NOT FOUND"
					_, err := conn.Write([]byte(response + "\n"))
					if err != nil {
						log.Println("Error while sending message:", err)
					}
				}
			}
		} else {
			response := "Response: " + msg
			_, err := conn.Write([]byte(response + "\n"))
			if err != nil {
				log.Println("Error while sending response:", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error while writing:", err)
	}
}

func loadDirectories(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var directories []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dir := strings.TrimSpace(scanner.Text())
		if dir != "" {
			directories = append(directories, dir)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return directories, nil
}

func startServer(ip, port string) {
	selfAddress := ip + ":" + port
	ln, err := net.Listen("tcp", selfAddress)
	if err != nil {
		log.Fatalf("Error to load server for address %s: %v", selfAddress, err)
		return
	}
	defer ln.Close()

	fmt.Printf("Server listening from address %s\n", selfAddress)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error to accept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func loadPeers(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var peers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		address := strings.TrimSpace(scanner.Text())
		if address != "" {
			peers = append(peers, address)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return peers, nil
}

func searchHashOnPeers(hash int, peers []string) {
	var wg sync.WaitGroup
	var positiveResponses []string
	var mutex sync.Mutex

	fmt.Printf("Hash asked: %d\n", hash)

	for _, peer := range peers {
		wg.Add(1)
		go func(peer string) {
			defer wg.Done()

			conn, err := net.Dial("tcp", peer)
			if err != nil {
				log.Printf("Error connecting to peer %s: %v", peer, err)
				return
			}
			defer conn.Close()

			message := fmt.Sprintf("SEARCH %d\n", hash)
			_, err = conn.Write([]byte(message))
			if err != nil {
				log.Printf("Error sending message to peer %s: %v", peer, err)
				return
			}

			scanner := bufio.NewScanner(conn)
			if scanner.Scan() {
				response := scanner.Text()
				parts := strings.Split(response, "|")
				if parts[0] == "FOUND" {
					mutex.Lock()
					positiveResponses = append(positiveResponses, parts[1])
					mutex.Unlock()
				}
			}

			if err := scanner.Err(); err != nil {
				log.Printf("Error reading response from peer %s: %v", peer, err)
			}
		}(peer)
	}

	wg.Wait()

	if len(positiveResponses) > 0 {
		fmt.Printf("Positive response from: %s\n", strings.Join(positiveResponses, ", "))
	} else {
		fmt.Println("No peers found with the requested hash.")
	}
}

func processDirectories(directories []string) map[int][]string {
	resultMap := make(map[int][]string)
	for _, dir := range directories {
		dirHashes, err := ProcessDirectory(dir)
		if err != nil {
			log.Printf("Error processing directory %s: %v", dir, err)
			continue
		}

		for hash, paths := range dirHashes {
			resultMap[hash] = append(resultMap[hash], paths...)
		}
	}
	return resultMap
}

func main() {
	mySharedFilesPathsHashMap = make(map[int][]string)

	directories, err := loadDirectories(SHARED_FILES_FILE)
	if err != nil {
		log.Fatalf("Error to load directories: %v", err)
	}

	mySharedFilesPathsHashMap = processDirectories(directories)
	startHashRecalculation(directories, HASH_RECALCULATION)

	peers, err := loadPeers(PEERS_FILE)
	if err != nil {
		log.Fatalf("Error to load peers: %v", err)
	}

	go startServer(GetLocalIP(), "8080")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter the hash you want to search (or 'exit' to finish): ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			fmt.Println("Exiting program.")
			break
		}

		hash, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input, please input a number.")
			continue
		}

		searchHashOnPeers(hash, peers)
	}

	select {}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func port8081() {
	reader := bufio.NewReader(os.Stdin)
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send a message to the server
	message := "Hello, mini server!"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
	fmt.Println("Message sent to server:", message)

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Failed to read from server:", err)
				return
			}
			fmt.Println("Received from server:", string(buffer[:n]))
		}
	}()
	for {
		fmt.Print("> ")
		commands, _ := reader.ReadString('\n')
		commands = strings.TrimSpace(commands)
		_, err = conn.Write([]byte(commands))
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("CONNECT " + username))

	if err != nil {
		panic(err)
	}
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading:", err)
				return
			}
			fmt.Println(string(buffer[:n]))
		}
	}()

	for {
		fmt.Print("> ")
		commands, _ := reader.ReadString('\n')
		commands = strings.TrimSpace(commands)

		_, err = conn.Write([]byte(commands))
		if err != nil {
			panic(err)
		}

		if commands == "DISCONNECT" {
			fmt.Println("Disconnected from server.")
			return
		}
		if commands == "PokeBat" {
			port8081()
		}
	}
}

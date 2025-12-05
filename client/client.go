package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	serverConn, connErr := net.Dial("tcp", "localhost:1234")
	if connErr != nil {
		log.Fatal("Connection error:", connErr)
	}
	defer serverConn.Close()

	var clientID int

	serverScanner := bufio.NewScanner(serverConn)
	if serverScanner.Scan() {
		receivedLine := serverScanner.Text()
		_, parseErr := fmt.Sscanf(receivedLine, "YOUR_ID:%d", &clientID)
		if parseErr != nil {
			log.Fatal("Failed to parse client ID:", parseErr)
		}
	} else {
		log.Fatal("Failed to receive client ID")
	}

	fmt.Printf("Connected to chat server as user %d!\n", clientID)
	fmt.Println("Type 'exit' to quit.")
	fmt.Print("> ")

	go func() {
		for serverScanner.Scan() {
			message := serverScanner.Text()
			fmt.Print("\r\033[K")
			if strings.Contains(message, "joined") || strings.Contains(message, "left") {
				fmt.Println(">>>", message)
			} else {
				fmt.Println(message)
			}
			fmt.Print("> ")
		}
	}()

	userReader := bufio.NewReader(os.Stdin)

	for {
		userInput, readErr := userReader.ReadString('\n')
		if readErr != nil {
			log.Println("Error reading input:", readErr)
			break
		}
		userInput = strings.TrimSpace(userInput)

		if userInput == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if userInput != "" {
			_, sendErr := fmt.Fprintln(serverConn, userInput)
			if sendErr != nil {
				log.Println("Error sending message:", sendErr)
				break
			}
		}

		fmt.Print("> ")
	}
}

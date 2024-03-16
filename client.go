package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Tic Tac Toe Client")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go readServerMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your moves in the format 'row col'. Type 'quit' to exit.")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "quit" {
			sendMessage(conn, input)
			break
		}
		fmt.Println("Input being sent: ", input) // Log input being sent
		sendMessage(conn, input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}

func readServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server connection closed.")
			} else {
				fmt.Println("Error reading from server:", err)
			}
			break
		}
		fmt.Print(message)
	}
}

func sendMessage(conn net.Conn, message string) {
	fmt.Fprintf(conn, "%s\n", message)
}

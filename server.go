package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	players            = make([]net.Conn, 0, 2)
	board              = [3][3]string{{" ", " ", " "}, {" ", " ", " "}, {" ", " ", " "}}
	currentPlayerIndex = 0
	player             = "X"
	gameOver           = false
	winner             = ""
	numFilled          = 0
	mu                 sync.Mutex
	moveMade           = make(chan bool)
)

func main() {
	fmt.Println("Tic Tac Toe Server")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Server started. Waiting for players...")

	for len(players) < 2 {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		players = append(players, conn)
		mu.Unlock()

		fmt.Printf("Player %d connected\n", len(players))

		go handlePlayer(conn)
	}

	fmt.Println("Game started!")
	sendMessageToAll("Game started!\n")
	sendBoardToAll()

	for !gameOver {
		sendMessageToAll(fmt.Sprintf("Player %s's turn\n", player))
		playTurn()
		checkGameOver()
	}

	announceWinner()
	closeConnections()
}

func handlePlayer(conn net.Conn) {
	defer func() {
		removePlayer(conn)
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Println("Recieve these input: ", input)
		if input == "quit" {
			return
		}
		handleInput(input, conn)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from player: %v\n", err)
	}
}

// func waitForMove(conn net.Conn) {
// 	scanner := bufio.NewScanner(conn)
// 	for {
// 		if scanner.Scan() {
// 			input := scanner.Text()
// 			if input == "quit" {
// 				return
// 			}
// 			handleInput(input, conn)
// 			return
// 		}
// 	}
// }

func handleInput(input string, conn net.Conn) {
	fmt.Println("Recieve these input: ", input)
	mu.Lock()
	defer mu.Unlock()

	if conn != players[currentPlayerIndex] {
		sendServerMessage(conn, "It's not your turn.\n")
		return
	}

	coords := strings.Split(input, " ")
	if len(coords) != 2 {
		sendServerMessage(conn, "Invalid input. Please enter row and column numbers separated by space.\n")
		return
	}

	row, col := parseInput(coords[0]), parseInput(coords[1])
	if row == -1 || col == -1 {
		sendServerMessage(conn, "Invalid input. Please enter valid row and column numbers.\n")
		return
	}

	if board[row][col] != " " {
		sendServerMessage(conn, "That cell is already taken. Please choose another.\n")
		return
	}

	board[row][col] = player
	numFilled++
	sendBoardToAll()
	moveMade <- true
}

func playTurn() {
	currentPlayer := players[currentPlayerIndex]
	sendServerMessage(currentPlayer, "Your turn. Enter row and column (e.g., 1 2): ")
	// waitForMove(currentPlayer)
	<-moveMade
}

func sendServerMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		log.Println("Error sending message:", err)
		removePlayer(conn)
	}
}

func sendMessageToAll(message string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Current message to be sent: ") // Log current message
	fmt.Print(message)

	for _, conn := range players {
		sendServerMessage(conn, message)
	}
}

func sendBoardToAll() {
	boardStr := "   0   1   2\n"
	for i := 0; i < 3; i++ {
		boardStr += fmt.Sprintf("%d ", i)
		for j := 0; j < 3; j++ {
			if j != 0 {
				boardStr += "|"
			}
			boardStr += fmt.Sprintf(" %s ", board[i][j])
		}
		boardStr += "\n"
		if i != 2 {
			boardStr += "  ---+---+---\n"
		}
	}
	boardStr += "\n"
	fmt.Println("Current Board: ")
	fmt.Print(boardStr) // Log current board
	sendMessageToAll(boardStr)
}

func parseInput(input string) int {
	switch input {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	default:
		return -1
	}
}

func checkGameOver() {
	mu.Lock()
	defer mu.Unlock()

	if checkWin(player) {
		gameOver = true
		winner = player
		return
	}

	if numFilled == 9 {
		gameOver = true
		return
	}

	switchPlayer()
}

func checkWin(symbol string) bool {
	for i := 0; i < 3; i++ {
		if (board[i][0] == symbol && board[i][1] == symbol && board[i][2] == symbol) ||
			(board[0][i] == symbol && board[1][i] == symbol && board[2][i] == symbol) {
			return true
		}
	}

	if (board[0][0] == symbol && board[1][1] == symbol && board[2][2] == symbol) ||
		(board[0][2] == symbol && board[1][1] == symbol && board[2][0] == symbol) {
		return true
	}

	return false
}

func switchPlayer() {
	currentPlayerIndex = (currentPlayerIndex + 1) % 2
	if player == "X" {
		player = "O"
	} else {
		player = "X"
	}
}

func announceWinner() {
	if winner != "" {
		sendMessageToAll(fmt.Sprintf("Player %s wins!\n", winner))
	} else {
		sendMessageToAll("It's a draw!\n")
	}
}

func closeConnections() {
	for _, conn := range players {
		conn.Close()
	}
}

func removePlayer(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for i, playerConn := range players {
		if playerConn == conn {
			log.Printf("Player %d disconnected\n", i+1)
			players = append(players[:i], players[i+1:]...)
			conn.Close()

			if len(players) < 2 {
				gameOver = true
				log.Println("Not enough players. Game over.")
			}
			break
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	board     = [3][3]string{{" ", " ", " "}, {" ", " ", " "}, {" ", " ", " "}}
	player    = "X"
	gameOver  = false
	winner    = ""
	numFilled = 0
)

func main() {
	fmt.Println("Welcome to Tic Tac Toe!")
	printBoard()

	for !gameOver {
		playTurn()
		checkGameOver()
	}

	if winner != "" {
		fmt.Printf("Player %s wins!\n", winner)
	} else {
		fmt.Println("It's a draw!")
	}
}

func printBoard() {
	fmt.Println("   0   1   2")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d ", i)
		for j := 0; j < 3; j++ {
			if j != 0 {
				fmt.Print("|")
			}
			fmt.Printf(" %s ", board[i][j])
		}
		fmt.Println()
		if i != 2 {
			fmt.Println("  ---+---+---")
		}
	}
}

func playTurn() {
	fmt.Printf("Player %s's turn\n", player)
	fmt.Print("Enter row and column (e.g., 1 2): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	input = strings.TrimSpace(input)
	coords := strings.Split(input, " ")
	if len(coords) != 2 {
		fmt.Println("Invalid input. Please enter row and column numbers separated by space.")
		return
	}

	row, col := parseInput(coords[0]), parseInput(coords[1])
	if row == -1 || col == -1 {
		fmt.Println("Invalid input. Please enter valid row and column numbers.")
		return
	}

	if board[row][col] != " " {
		fmt.Println("That cell is already taken. Please choose another.")
		return
	}

	board[row][col] = player
	numFilled++
	printBoard()
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
	if checkWin(player) {
		gameOver = true
		winner = player
		return
	}

	if numFilled == 9 {
		gameOver = true
		return
	}

	player = switchPlayer(player)
}

func checkWin(symbol string) bool {
	for i := 0; i < 3; i++ {
		// Check rows and columns
		if (board[i][0] == symbol && board[i][1] == symbol && board[i][2] == symbol) ||
			(board[0][i] == symbol && board[1][i] == symbol && board[2][i] == symbol) {
			return true
		}
	}

	// Check diagonals
	if (board[0][0] == symbol && board[1][1] == symbol && board[2][2] == symbol) ||
		(board[0][2] == symbol && board[1][1] == symbol && board[2][0] == symbol) {
		return true
	}

	return false
}

func switchPlayer(currentPlayer string) string {
	if currentPlayer == "X" {
		return "O"
	}
	return "X"
}

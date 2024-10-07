package main

import (
	"math/rand"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

// Choices
const (
	Rock     = "rock"
	Paper    = "paper"
	Scissors = "scissors"
	CW = "Client Wins"
	SW = "Server Wins"
)

type UserChoice struct {
	Choice string `json:"choice"`
}

// Rules route Controller
func getRules(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"game": "Rock, paper, scissors",
		"rules": []gin.H{
			{
				"choice":      "Stone",
				"beats":       "Scissors",
				"losesTo":     "Paper",
				"description": "Rock smashes scissors.",
			},
			{
				"choice":      "Paper",
				"beats":       "Stone",
				"losesTo":     "Scissors",
				"description": "Paper covers stone.",
			},
			{
				"choice":      "Scissors",
				"beats":       "Paper",
				"losesTo":     "Stone",
				"description": "Scissors cut paper.",
			},
		},
		"note": "If both players pick the same, it's a draw.",
	})
}

// Game route controller
func postGame(ctx *gin.Context) {
	var clientChoice UserChoice

	//Bind JSON input to clientChoice and evaluate error
	
	if err := ctx.BindJSON(&clientChoice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}
	
	if clientChoice.Choice == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Choice cannot be empty"})
		return
	}

	//Call function and evaluate results
	clientChoiceString, serverChoice, result := winnerEvaluator(clientChoice)
	if result == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"clientChoice": strings.Title(clientChoiceString),
		"serverChoice": strings.Title(serverChoice),
		"result": strings.Title(result),
	})
}

// Game result
func winnerEvaluator(clientChoice UserChoice) (string, string, string) {
	//Convert user's choice to lowercase
	loweredClientChoice := strings.ToLower(clientChoice.Choice)

	//Returns empty strings if user's choice is not valid
	if loweredClientChoice != Rock && loweredClientChoice != Paper && loweredClientChoice != Scissors {
		return "", "", ""
	}

	//Establishes server's hand
	serverChoice := serverChoiceGenerator()

	//Determins winner
	var result string

	if loweredClientChoice == serverChoice {
		result = "Draw"
	} else {
		switch loweredClientChoice {
		case Rock:
			if serverChoice == Scissors {
				result = CW
			}else{
				result = SW
			}
		case Paper:
			if serverChoice == Rock {
				result = CW
			}else{
				result = SW
			}
		case Scissors:
			if serverChoice == Paper {
				result = CW
			}else{
				result = SW
			}
		}
	}

	//Returns
	return loweredClientChoice, serverChoice, result
}

// Choice generator
func serverChoiceGenerator() string {
	//Stores a number from 0 to 2
	randomChoice := rand.Intn(3)
	//Returns the server's choice given the random number
	if randomChoice == 0 {
		return Rock
	} else if randomChoice == 1 {
		return Paper
	} else {
		return Scissors
	}
}

// Main
func main() {
	router := gin.Default()
	router.GET("/rules", getRules)
	router.POST("/game", postGame)
	router.Run("localhost:8080")
}

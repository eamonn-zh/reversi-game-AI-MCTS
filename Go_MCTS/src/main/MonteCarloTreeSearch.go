package main

import (
	/*
		#include <stdlib.h>
	*/
	"C"
	"encoding/json"
	"time"
	"unsafe"
)

// open concurrent mode
var concurrent = false

// Response body
type ReversiGameResult struct {
	State        [8][8]int
	AvailablePos []PiecePosition
	CurrTurn     int
	GameStatus   int
	BlackCount   int
	WhiteCount   int
}

// An reversi game
var reversiGame ReversiGame

// Status code: Continue Game
var CONTINUE = 3

// Define the weighted board
var blackBoardValues = [8][8]int{
	{6, -2, 4, 3, 3, 4, -2, 6},
	{-2, -4, -3, 1, 1, -3, -4, -2},
	{4, -3, 2, 2, 2, 2, -3, 4},
	{3, 1, 2, 0, 0, 2, 1, 3},
	{3, 1, 2, 0, 0, 2, 1, 3},
	{4, -3, 2, 2, 2, 2, -3, 4},
	{-2, -4, -3, 1, 1, -3, -4, -2},
	{6, -2, 4, 3, 3, 4, -2, 6},
}

var whiteBoardValues = [8][8]int{
	{6, -2, 4, 3, 3, 4, -2, 6},
	{-2, -4, -3, 1, 1, -3, -4, -2},
	{4, -3, 2, 2, 2, 2, -3, 4},
	{3, 1, 2, 0, 0, 2, 1, 3},
	{3, 1, 2, 0, 0, 2, 1, 3},
	{4, -3, 2, 2, 2, 2, -3, 4},
	{-2, -4, -3, 1, 1, -3, -4, -2},
	{6, -2, 4, 3, 3, 4, -2, 6},
}

// A heuristic algorithm based on weighted board
func heuristic(currTurn int, availaPos *[]PiecePosition, reverseCounts *[]int, maxReverse int) int {
	tempBlackBoardValues := [8][8]int{
		{6, -2, 4, 3, 3, 4, -2, 6},
		{-2, -4, -3, 1, 1, -3, -4, -2},
		{4, -3, 2, 2, 2, 2, -3, 4},
		{3, 1, 2, 0, 0, 2, 1, 3},
		{3, 1, 2, 0, 0, 2, 1, 3},
		{4, -3, 2, 2, 2, 2, -3, 4},
		{-2, -4, -3, 1, 1, -3, -4, -2},
		{6, -2, 4, 3, 3, 4, -2, 6},
	}
	tempWhiteBoardValues := [8][8]int{
		{6, -2, 4, 3, 3, 4, -2, 6},
		{-2, -4, -3, 1, 1, -3, -4, -2},
		{4, -3, 2, 2, 2, 2, -3, 4},
		{3, 1, 2, 0, 0, 2, 1, 3},
		{3, 1, 2, 0, 0, 2, 1, 3},
		{4, -3, 2, 2, 2, 2, -3, 4},
		{-2, -4, -3, 1, 1, -3, -4, -2},
		{6, -2, 4, 3, 3, 4, -2, 6},
	}
	maxIndex := 0
	maxValue := -5
	var values *[8][8]int
	if currTurn == -1 {
		values = &tempBlackBoardValues
	} else {
		values = &tempWhiteBoardValues
	}

	// Update weights
	for i := 0; i < len(*reverseCounts); i++ {
		if (*availaPos)[i].Equal(PiecePosition{X: 0, Y: 0}) {
			values[0][1] = 6
			values[1][0] = 6
			return i
		} else if (*availaPos)[i].Equal(PiecePosition{X: 0, Y: 7}) {
			values[6][0] = 6
			values[7][1] = 6
			return i
		} else if (*availaPos)[i].Equal(PiecePosition{X: 7, Y: 0}) {
			values[0][6] = 6
			values[1][7] = 6
			return i
		} else if (*availaPos)[i].Equal(PiecePosition{X: 7, Y: 7}) {
			values[6][7] = 6
			values[7][6] = 6
			return i
		}
		// Pick a position with the highest weight
		if (*reverseCounts)[i] == maxReverse && values[(*availaPos)[i].Y][(*availaPos)[i].X] > maxValue {
			maxValue = values[(*availaPos)[i].Y][(*availaPos)[i].X]
			maxIndex = i
		}
	}
	return maxIndex
}

// run MCTS concurrently
func concurrentMCTS(root *MCTNode, secondLimit float64, currTurn int) PiecePosition {
	results := make(chan PiecePosition, 3)

	// run 3 goroutines, it depends on CPU performance
	for i := 0; i < 3; i++ {
		node := CopyMCTNode(root)
		go monteCarloTreeSearch(&node, secondLimit, results)
	}
	var values *[8][8]int
	if currTurn == -1 {
		values = &blackBoardValues
	} else {
		values = &whiteBoardValues
	}
	maxPosition := <-results

	// AI does not have available position to move, so exit
	if maxPosition.X == -1 && maxPosition.Y == -1 {
		// clear the channel
		for i := 0; i < 2; i++ {
			<-results
		}
		return maxPosition
	}
	// Pick the position that has the highest weight
	maxWeight := values[maxPosition.Y][maxPosition.X]
	for i := 0; i < 2; i++ {
		result := <-results
		if values[result.Y][result.X] > maxWeight {
			maxWeight = values[result.Y][result.X]
			maxPosition = result
		}
	}
	return maxPosition
}

// MCTS algorithm
func monteCarloTreeSearch(root *MCTNode, secondLimit float64, results chan PiecePosition) PiecePosition {
	start := time.Now()
	if root.IsFinalLeafNode() {
		if concurrent {
			results <- PiecePosition{X: -1, Y: -1}
		}
		return PiecePosition{X: -1, Y: -1}
	}
	root.InitUnexpChildren()
	node := new(MCTNode)
	// Loop for 5 sec
	for (time.Now().Sub(start)).Seconds() <= secondLimit {
		// Start searching from the root node
		node = root
		// Selection
		for !node.IsFinalLeafNode() {
			if !node.IsFullyExpanded() {
				node = node.ExpandChild()
				node.InitUnexpChildren()
				break
			}
			node = node.FindBestChild(1.96)
		}
		// Simulation
		gameState := CopyGameState(node.GameState)
		var reverseCounts []int
		var maxReverse int
		var index int
		for !gameState.IsGameOver() {
			reverseCounts = []int{}
			maxReverse = 0
			currTurn := gameState.CurrTurn
			availaPos := gameState.GetAvailablePos(currTurn, &reverseCounts, &maxReverse)
			if len(availaPos) == 0 {
				gameState.CurrTurn = currTurn * -1
				continue
			}
			// use heuristic algorithm instead randomly choose
			index = heuristic(currTurn, &availaPos, &reverseCounts, maxReverse)
			gameState.PlayNextStep(availaPos[index])
		}
		result := gameState.CheckResult()
		// Backpropagation
		node.IncreVistedCount()
		node.UpdateScores(result)
		for node.Parent != nil {
			node = node.Parent
			node.IncreVistedCount()
			node.UpdateScores(result)
		}
	}
	bestMove := root.FindBestChild(0).Diff
	if concurrent {
		results <- bestMove
	}
	return bestMove
}

//export startGame
func startGame(choice int) *C.char {
	reversiGame = NewReversiGame()
	if choice == 2 {
		root := NewMCTNode(&reversiGame.GameState, nil, &PiecePosition{X: 0, Y: 0})
		if concurrent {
			reversiGame.PlayNextStep(concurrentMCTS(&root, 5, reversiGame.GameState.CurrTurn))
		} else {
			reversiGame.PlayNextStep(monteCarloTreeSearch(&root, 5, nil))
		}
	}
	reversiGameResult := new(ReversiGameResult)
	reversiGameResult.State = reversiGame.GameState.State
	reversiGameResult.AvailablePos = reversiGame.GameState.GetAvailablePos(reversiGame.GameState.CurrTurn, nil, nil)
	reversiGameResult.CurrTurn = reversiGame.GameState.CurrTurn
	reversiGameResult.GameStatus = CONTINUE
	reversiGameResult.BlackCount = len(reversiGame.GameState.BlackPos)
	reversiGameResult.WhiteCount = len(reversiGame.GameState.WhitePos)
	jsonData, _ := json.Marshal(reversiGameResult)
	return C.CString(string(jsonData))
}

//export waitForAI
func waitForAI() *C.char {
	reversiGameResult := new(ReversiGameResult)
	root := NewMCTNode(&reversiGame.GameState, nil, &PiecePosition{X: 0, Y: 0})
	var pos PiecePosition
	if concurrent {
		pos = concurrentMCTS(&root, 5, reversiGame.GameState.CurrTurn)
	} else {
		pos = monteCarloTreeSearch(&root, 5, nil)
	}
	if pos.X == -1 && pos.Y == -1 {
		reversiGame.Pass()
	} else {
		reversiGame.PlayNextStep(pos)
		if reversiGame.IsGameOver() {
			reversiGameResult.GameStatus = reversiGame.CheckResult()
		} else {
			reversiGameResult.GameStatus = CONTINUE
		}
	}
	reversiGameResult.State = reversiGame.GameState.State
	reversiGameResult.AvailablePos = reversiGame.GameState.GetAvailablePos(reversiGame.GameState.CurrTurn, nil, nil)
	reversiGameResult.CurrTurn = reversiGame.GameState.CurrTurn
	reversiGameResult.BlackCount = len(reversiGame.GameState.BlackPos)
	reversiGameResult.WhiteCount = len(reversiGame.GameState.WhitePos)
	jsonData, _ := json.Marshal(reversiGameResult)
	return C.CString(string(jsonData))
}

//export nextStep
func nextStep(positionX int, positionY int) *C.char {
	reversiGameResult := new(ReversiGameResult)
	reversiGame.PlayNextStep(PiecePosition{positionX, positionY})
	if reversiGame.IsGameOver() {
		reversiGameResult.GameStatus = reversiGame.CheckResult()
	} else {
		reversiGameResult.GameStatus = CONTINUE
	}
	reversiGameResult.State = reversiGame.GameState.State
	reversiGameResult.CurrTurn = reversiGame.GameState.CurrTurn
	reversiGameResult.BlackCount = len(reversiGame.GameState.BlackPos)
	reversiGameResult.WhiteCount = len(reversiGame.GameState.WhitePos)
	jsonData, _ := json.Marshal(reversiGameResult)
	return C.CString(string(jsonData))
}

//export passCurrentTurn
func passCurrentTurn() *C.char {
	reversiGame.Pass()
	return waitForAI()
}

//export freeMemory
func freeMemory(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

func main() {
}

package main

import (
	"encoding/json"
	/*
		typedef struct PiecePosition{
			int X;
			int Y;
		} PiecePosition;
	*/
	"C"
	"time"
)

type ReversiGameResult struct {
	State [8][8]int
	AvailablePos []PiecePosition
	CurrTurn int
	GameStatus int
}


var reversiGame ReversiGame

var blackBoardValues = [8][8]int{}

var whiteBoardValues = [8][8]int{}


var CONTINUE = 3

var temp int

func initWeightedTable() {
	blackBoardValues = [8][8]int{
		{6, -2, 4, 3, 3, 4, -2, 6},
		{-2, -4, -3, 1, 1, -3, -4, -2},
		{4, -3, 2, 2, 2, 2, -3, 4},
		{3, 1, 2, 0, 0,  2, 1, 3},
		{3, 1, 2, 0, 0,  2, 1, 3},
		{4, -3, 2, 2, 2, 2, -3, 4},
		{-2, -4, -3, 1, 1, -3, -4, -2},
		{6, -2, 4, 3, 3, 4, -2, 6},
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			whiteBoardValues[i][j] = blackBoardValues[i][j]
		}
	}
}

func heuristic(currTurn int, availaPos *[]PiecePosition, reverseCounts *[]int, maxReverse int) int{
	maxIndex := 0
	maxValue := -5
	var values *[8][8]int
	if currTurn == -1 {
		values = &blackBoardValues
	} else {
		values = &whiteBoardValues
	}
	for i := 0; i < len(*reverseCounts); i++ {
		if (*availaPos)[i].Equal(PiecePosition{}) {
			values[0][1] = 6
			values[1][0] = 6
			return i
		} else if (*availaPos)[i].Equal(PiecePosition{Y: 7}) {
			values[6][0] = 6
			values[7][1] = 6
			return i
		} else if (*availaPos)[i].Equal(PiecePosition{X: 7}) {
			values[0][6] = 6
			values[1][7] = 6
			return i
		} else if (*availaPos)[i].Equal(PiecePosition{X: 7, Y: 7}) {
			values[6][7] = 6
			values[7][6] = 6
			return i
		}
		if (*reverseCounts)[i] == maxReverse && values[(*availaPos)[i].Y][(*availaPos)[i].X] > maxValue {
			maxValue = values[(*availaPos)[i].Y][(*availaPos)[i].X]
			maxIndex = i
		}
	}
	return maxIndex
}

func monteCarloTreeSearch(root *MCTNode, secondLimit float64, playoutCount *int) PiecePosition {
	start := time.Now()
	if root.IsFinalLeafNode() {
		return PiecePosition{X: -1, Y: -1}
	}
	root.InitUnexpChildren()
	node := new(MCTNode)
	*playoutCount = 0
	for (time.Now().Sub(start)).Seconds() <= secondLimit {
		initWeightedTable()
		node = root
		for !node.IsFinalLeafNode() {
			if !node.IsFullyExpanded() {
				node = node.ExpandChild()
				node.InitUnexpChildren()
				break
			}
			node = node.FindBestChild(1.96)

		}
		gameState := CopyGameState(node.GameState)
		*playoutCount++
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
			index = heuristic(currTurn, &availaPos, &reverseCounts, maxReverse)
			gameState.PlayNextStep(availaPos[index])
		}
		result := gameState.CheckResult()
		// Backpropagation
		node.IncreVistedCount()
		node.UpdateScores(result)
		for node.Parent != nil{
			node = node.Parent
			node.IncreVistedCount()
			node.UpdateScores(result)
		}
	}
	bestMove := root.FindBestChild(0).Diff
	return bestMove
}

//export startGame
func startGame(choice int) *C.char {
	reversiGame = NewReversiGame()
	var temp int
	if choice == 2 {
		root := NewMCTNode(&reversiGame.GameState, nil, &PiecePosition{X: 0, Y: 0})
		reversiGame.PlayNextStep(monteCarloTreeSearch(&root, 5, &temp))
	}

	reversiGameResult := new(ReversiGameResult)
	reversiGameResult.State = reversiGame.GameState.State
	reversiGameResult.AvailablePos = reversiGame.GameState.GetAvailablePos(reversiGame.GameState.CurrTurn, nil, nil)
	reversiGameResult.CurrTurn = reversiGame.GameState.CurrTurn
	reversiGameResult.GameStatus = CONTINUE
	jsonData, _ := json.Marshal(reversiGameResult)
	return C.CString(string(jsonData))
}

//export waitForAI
func waitForAI() *C.char{
	reversiGameResult := new(ReversiGameResult)
	root := NewMCTNode(&reversiGame.GameState, nil, &PiecePosition{X: 0, Y: 0})
	pos := monteCarloTreeSearch(&root, 5, &temp)
	if pos.X == -1 && pos.Y == -1 {
		reversiGame.Pass()
	} else {
		reversiGame.PlayNextStep(pos)
		if reversiGame.IsGameOver() {
			reversiGameResult.GameStatus = reversiGame.CheckResult()
			jsonData, _ := json.Marshal(reversiGameResult)
			return C.CString(string(jsonData))
		}
	}
	reversiGameResult.State = reversiGame.GameState.State
	reversiGameResult.AvailablePos = reversiGame.GameState.GetAvailablePos(reversiGame.GameState.CurrTurn, nil, nil)
	reversiGameResult.CurrTurn = reversiGame.GameState.CurrTurn
	reversiGameResult.GameStatus = CONTINUE
	jsonData, _ := json.Marshal(reversiGameResult)
	return C.CString(string(jsonData))
}

//export nextStep
func nextStep(positionX int, positionY int) *C.char {
	reversiGameResult := new(ReversiGameResult)
	reversiGame.PlayNextStep(PiecePosition{positionX, positionY})
	if reversiGame.IsGameOver() {
		reversiGameResult.GameStatus = reversiGame.CheckResult()
		jsonData, _ := json.Marshal(reversiGameResult)
		return C.CString(string(jsonData))
	}
	reversiGameResult.State = reversiGame.GameState.State
	reversiGameResult.CurrTurn = reversiGame.GameState.CurrTurn
	reversiGameResult.GameStatus = CONTINUE
	jsonData, _ := json.Marshal(reversiGameResult)
	return C.CString(string(jsonData))
}


func main() {
}

package main

const BLACK int = 1
const WHITE int = -1
const DRAW int = 0
const EMPTY int = 0
var positionDiffX = [8]int{-1, 0, 1, -1, 1, -1, 0, 1}
var positionDiffY = [8]int{-1, -1, -1, 0, 0, 1, 1, 1}

type GameState struct {
	State    [8][8]int
	CurrTurn int
	BlackPos []PiecePosition
	WhitePos []PiecePosition
}

func NewGameState() GameState {
	gameState := new(GameState)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			gameState.State[i][j] = 0
		}
	}
	gameState.State[3][4] = BLACK
	gameState.State[4][3] = BLACK
	gameState.State[3][3] = WHITE
	gameState.State[4][4] = WHITE
	gameState.CurrTurn = WHITE
	gameState.BlackPos = append(gameState.BlackPos, PiecePosition{4, 3}, PiecePosition{3, 4})
	gameState.WhitePos = append(gameState.WhitePos, PiecePosition{3, 3}, PiecePosition{4, 4})
	return *gameState
}

func CopyGameState(gameState GameState) GameState {
	newGameState := new(GameState)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++{
			newGameState.State[i][j] = gameState.State[i][j]
		}
	}
	blackPos := make([]PiecePosition, len(gameState.BlackPos))
	copy(blackPos, gameState.BlackPos)
	whitePos := make([]PiecePosition, len(gameState.WhitePos))
	copy(whitePos, gameState.WhitePos)
	newGameState.CurrTurn = gameState.CurrTurn
	newGameState.BlackPos = blackPos
	newGameState.WhitePos = whitePos
	return *newGameState
}

func (gameState *GameState) PlayNextStep(pos PiecePosition)  {
	gameState.CurrTurn *= -1
	gameState.State[pos.Y][pos.X] = gameState.CurrTurn
	var currPosList *[]PiecePosition
	var oppoPosList *[]PiecePosition
	if gameState.CurrTurn == BLACK {
		currPosList = &gameState.BlackPos
		oppoPosList = &gameState.WhitePos
	} else {
		currPosList = &gameState.WhitePos
		oppoPosList = &gameState.BlackPos
	}
	*currPosList = append(*currPosList, pos)
	var temp []PiecePosition
	for i := 0; i < 8; i++ {
		temp = temp[:0]
		newPos := PiecePosition{pos.X + positionDiffX[i], pos.Y + positionDiffY[i]}
		for newPos.Y >= 0 && newPos.Y <= 7 && newPos.X >= 0 && newPos.X <= 7 && gameState.State[newPos.Y][newPos.X] == gameState.CurrTurn* -1 {
			temp = append(temp, newPos)
			newPos.X = newPos.X + positionDiffX[i]
			newPos.Y = newPos.Y + positionDiffY[i]
		}
		if newPos.Y < 0 || newPos.Y > 7 || newPos.X < 0 || newPos.X > 7 || gameState.State[newPos.Y][newPos.X] == EMPTY {
			continue
		}
		for _, j := range temp {
			gameState.State[j.Y][j.X] = gameState.CurrTurn
			*currPosList = append(*currPosList, j)
			for index, k := range *oppoPosList {
				if k == j {
					*oppoPosList = append((*oppoPosList)[:index], (*oppoPosList)[index + 1:]...)
					break
				}
			}
		}
	}
}

func (gameState *GameState) GetAvailablePos(turn int, reverseCounts *[]int, maxReverse *int) []PiecePosition {
	var availablePos []PiecePosition
	var tempPosList *[]PiecePosition
	if turn == BLACK {
		tempPosList = &gameState.BlackPos
	} else {
		tempPosList = &gameState.WhitePos
	}
	for _, pos := range *tempPosList {
		for i := 0; i < 8; i++ {
			newPos := PiecePosition{pos.X + positionDiffX[i], pos.Y + positionDiffY[i]}
			if newPos.Y >= 0 && newPos.Y < 8 && newPos.X >= 0 && newPos.X < 8 && gameState.State[newPos.Y][newPos.X] == 0 {
				tempY := pos.Y - positionDiffY[i]
				tempX := pos.X - positionDiffX[i]
				reverseCount := 1
				for tempX >= 0 && tempX < 8 && tempY >= 0 && tempY < 8 && gameState.State[tempY][tempX] != EMPTY {
					if gameState.State[tempY][tempX] == -1 * turn {
						availablePos = append(availablePos, newPos)
						if reverseCounts != nil {
							*reverseCounts = append(*reverseCounts, reverseCount)
							if reverseCount > *maxReverse{
								*maxReverse = reverseCount
							}
						}
						break
					}
					reverseCount++
					tempY -= positionDiffY[i]
					tempX -= positionDiffX[i]
				}
			}
		}
	}
	return availablePos
}

func (gameState *GameState) IsGameOver() bool {
	return (len(gameState.GetAvailablePos(WHITE, nil, nil)) == 0 &&
		len(gameState.GetAvailablePos(BLACK, nil, nil)) == 0) ||
		len(gameState.BlackPos) == 0 || len(gameState.WhitePos) == 0
}

func (gameState *GameState) CheckResult() int {
	if len(gameState.BlackPos) == len(gameState.WhitePos) {
		return DRAW
	}
	if len(gameState.BlackPos) > len(gameState.WhitePos) {
		return BLACK
	} else{
		return WHITE
	}
}

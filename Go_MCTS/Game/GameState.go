package Game

const BLACK int32 = 1
const WHITE int32 = -1
const DRAW int32 = 0
const EMPTY int32 = 0
var diffX = [8]int32{-1, 0, 1, -1, 1, -1, 0, 1}
var diffY = [8]int32{-1, -1, -1, 0, 0, 1, 1, 1}

type GameState struct {
	state [8][8]int32
	currTurn int32
	blackPos []Position
	whitePos []Position
}

func NewGameState() GameState{
	gameState := new(GameState)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			gameState.state[i][j] = 0
		}
	}
	gameState.state[3][4] = BLACK
	gameState.state[4][3] = BLACK
	gameState.state[3][3] = WHITE
	gameState.state[4][4] = WHITE
	gameState.currTurn = WHITE
	gameState.blackPos = append(gameState.blackPos, Position{4, 3}, Position{3, 4})
	gameState.whitePos = append(gameState.whitePos, Position{3, 3}, Position{4, 4})
	return *gameState
}

func CopyGameState(gameState GameState) GameState{
	return gameState
}

func (gameState *GameState) playNextStep(pos *Position)  {
	gameState.currTurn *= -1
	gameState.state[pos.y][pos.x] = gameState.currTurn
	var curr_pos_list *[]Position
	var oppo_pos_list *[]Position
	if gameState.currTurn == BLACK {
		curr_pos_list = &gameState.blackPos
		oppo_pos_list = &gameState.whit	} else {
		curr_pos_list = &gameState.whitePos
		oppo_pos_list = &gameState.blackPos
	}
	*curr_pos_list = append(*curr_pos_list, *pos)
	temp := []Position{}
	for i := 0; i < 8; i++ {
		temp = temp[:0]
		newPos := Position{pos.x + diffX[i], pos.y + diffY[i]}
		for newPos.y >= 0 && newPos.y <= 7 && newPos.x >= 0 && newPos.x <= 7 && gameState.state[newPos.y][newPos.x] == gameState.currTurn * -1 {
			temp = append(temp, newPos)
			newPos.x = newPos.x + diffX[i]
			newPos.y = newPos.y + diffY[i]
		}
		if newPos.y < 0 || newPos.y > 7 || newPos.x < 0 || newPos.x > 7 || gameState.state[newPos.y][newPos.x] == EMPTY {
			continue
		}
		for _, j := range temp {
			gameState.state[j.y][j.x] = gameState.currTurn
			*curr_pos_list = append(*curr_pos_list, j)
			for index, k := range *oppo_pos_list {
				if k == j {
					*oppo_pos_list = append((*oppo_pos_list)[:index], (*oppo_pos_list)[index + 1:]...)
					break
				}
			}
		}
	}


}


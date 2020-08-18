package main

type ReversiGame struct {
	GameState GameState
}

// Constructor
func NewReversiGame() ReversiGame {
	reversiGame := new(ReversiGame)
	reversiGame.GameState = NewGameState()
	return *reversiGame
}

// if a player has no available position to play, he should give up the current turn
func (reversiGame *ReversiGame) Pass() {
	reversiGame.GameState.CurrTurn = reversiGame.GameState.CurrTurn * -1
}

// receive the user's input of next position, then update the game state
func (reversiGame *ReversiGame) PlayNextStep(pos PiecePosition) {
	reversiGame.GameState.PlayNextStep(pos)
}

// check if game is over
func (reversiGame *ReversiGame) IsGameOver() bool{
	return reversiGame.GameState.IsGameOver()
}

// check game result
func (reversiGame *ReversiGame) CheckResult() int{
	return reversiGame.GameState.CheckResult()
}
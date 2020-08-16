package main

type ReversiGame struct {
	GameState GameState
}

func NewReversiGame() ReversiGame {
	reversiGame := new(ReversiGame)
	reversiGame.GameState = NewGameState()
	return *reversiGame
}

func (reversiGame *ReversiGame) Pass() {
	reversiGame.GameState.CurrTurn = reversiGame.GameState.CurrTurn * -1
}

func (reversiGame *ReversiGame) PlayNextStep(pos PiecePosition) {
	reversiGame.GameState.PlayNextStep(pos)
}

func (reversiGame *ReversiGame) IsGameOver() bool{
	return reversiGame.GameState.IsGameOver()
}

func (reversiGame *ReversiGame) CheckResult() int{
	return reversiGame.GameState.CheckResult()
}
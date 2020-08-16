package main

type PiecePosition struct {
	X int
	Y int
}

func (piecePosition *PiecePosition) Equal(position PiecePosition) bool {
	return piecePosition.X == position.X && piecePosition.Y == position.Y
}
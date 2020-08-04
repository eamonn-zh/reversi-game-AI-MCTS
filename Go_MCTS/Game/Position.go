package Game

type Position struct {
	x int32
	y int32
}

func (position *Position) equal (position2 Position) bool {
	return position.x == position2.x && position.y == position2.y
}
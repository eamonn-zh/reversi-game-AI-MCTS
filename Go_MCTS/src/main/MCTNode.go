package main

import (
	"math"
)

// Monte Carlo Tree Node
type MCTNode struct {
	GameState          GameState
	Parent             *MCTNode
	Children           []*MCTNode
	UnexpandedChildren []*MCTNode
	Diff               PiecePosition
	ScoreMap           map[int]int
	VisitedCount       int
}

// Constructor
func NewMCTNode(gameState *GameState, parent *MCTNode, diff *PiecePosition) MCTNode {
	scoreMap := map[int]int{}
	scoreMap[1] = 0
	scoreMap[0] = 0
	scoreMap[-1] = 0
	node := MCTNode{GameState: *gameState, Parent: parent, Diff: *diff, VisitedCount: 0, ScoreMap: scoreMap}
	return node
}

// Copy constructor for deep copy
func CopyMCTNode(mctNode *MCTNode) MCTNode {
	newMCTNode := new(MCTNode)
	// Copy score maps
	scoreMap := map[int]int{}
	scoreMap[1] = mctNode.ScoreMap[1]
	scoreMap[0] = mctNode.ScoreMap[0]
	scoreMap[-1] = mctNode.ScoreMap[-1]
	newMCTNode.ScoreMap = scoreMap

	newMCTNode.Parent = mctNode.Parent

	newMCTNode.Children = make([]*MCTNode, len(mctNode.Children))
	copy(newMCTNode.Children, mctNode.Children)
	newMCTNode.UnexpandedChildren = make([]*MCTNode, len(mctNode.UnexpandedChildren))
	copy(newMCTNode.UnexpandedChildren, mctNode.UnexpandedChildren)

	newMCTNode.Diff = PiecePosition{mctNode.Diff.X, mctNode.Diff.Y}
	newMCTNode.VisitedCount = mctNode.VisitedCount
	newMCTNode.GameState = CopyGameState(mctNode.GameState)
	return *newMCTNode
}

// Initialize the unexpanded children nodes
func (mctNode *MCTNode) InitUnexpChildren() {
	turn := mctNode.GameState.CurrTurn
	availaPos := mctNode.GameState.GetAvailablePos(turn, nil, nil)
	if len(availaPos) == 0 {
		turn *= -1
		availaPos = mctNode.GameState.GetAvailablePos(turn, nil, nil)
	}
	for _, position := range availaPos {
		childGameState := CopyGameState(mctNode.GameState)
		childGameState.CurrTurn = turn
		childGameState.PlayNextStep(position)
		childNode := NewMCTNode(&childGameState, mctNode, &position)
		mctNode.UnexpandedChildren = append(mctNode.UnexpandedChildren, &childNode)

	}

}

// Move one child node from unexpanded slice to children slice
func (mctNode *MCTNode) ExpandChild() *MCTNode {
	child := mctNode.UnexpandedChildren[len(mctNode.UnexpandedChildren)-1]
	mctNode.UnexpandedChildren = mctNode.UnexpandedChildren[:len(mctNode.UnexpandedChildren)-1]
	mctNode.Children = append(mctNode.Children, child)
	return child
}

// Check if a node's children are all expanded
func (mctNode *MCTNode) IsFullyExpanded() bool {
	return len(mctNode.UnexpandedChildren) == 0
}

// Check if a node is a leaf node
func (mctNode *MCTNode) IsFinalLeafNode() bool {
	return mctNode.GameState.IsGameOver()
}

// Scores = wins count - loses count
func (mctNode *MCTNode) GetScore() int {
	wins := mctNode.ScoreMap[mctNode.GameState.CurrTurn]
	loses := mctNode.ScoreMap[mctNode.GameState.CurrTurn*-1]
	return wins - loses
}

// Find the best child for next simulation
func (mctNode *MCTNode) FindBestChild(c float64) *MCTNode {
	maxUCTValue := float64(mctNode.Children[0].GetScore())/float64(mctNode.Children[0].VisitedCount) +
		math.Sqrt(c*2*math.Log(float64(mctNode.VisitedCount)/float64(mctNode.Children[0].VisitedCount)))
	bestChild := mctNode.Children[0]

	for _, i := range mctNode.Children {
		UCTValue := float64(i.GetScore())/float64(i.VisitedCount) +
			math.Sqrt(c*math.Log(float64(mctNode.VisitedCount)/float64(i.VisitedCount)))
		if UCTValue > maxUCTValue {
			maxUCTValue = UCTValue
			bestChild = i
		}
	}

	return bestChild
}

// Increase the visited count
func (mctNode *MCTNode) IncreVistedCount() {
	mctNode.VisitedCount++
}

func (mctNode *MCTNode) UpdateScores(result int) {
	mctNode.ScoreMap[result]++
}

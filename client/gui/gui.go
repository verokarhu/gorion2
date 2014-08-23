package gui

import (
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

type ButtonMap struct {
	UpperLeft  sf.Vector2f
	LowerRight sf.Vector2f
}

// CalcPosVec calculates a resolution independent vector of the x,y position
func CalcPosVec(pos sf.Vector2i, res sf.Vector2u) (v sf.Vector2f) {
	v.X = float32(pos.X) / float32(res.X)
	v.Y = float32(pos.Y) / float32(res.Y)

	return
}

func (b *ButtonMap) Contains(pos sf.Vector2f) bool {
	if pos.X >= b.UpperLeft.X && pos.X <= b.LowerRight.X && pos.Y >= b.UpperLeft.Y && pos.Y <= b.LowerRight.Y {
		return true
	}

	return false
}

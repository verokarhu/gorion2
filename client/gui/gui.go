package gui

import (
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

// CalcPosVec calculates a resolution independent vector of the x,y position
func CalcPosVec(pos sf.Vector2u, res sf.Vector2u) (v sf.Vector2f) {
	v.X = float32(pos.X) / float32(res.X)
	v.Y = float32(pos.Y) / float32(res.Y)

	return
}

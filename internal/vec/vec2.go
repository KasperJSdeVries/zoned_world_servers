package vec

import "fmt"

type Vec2 struct {
	X, Y float32
}

func (v *Vec2) String() string {
	return fmt.Sprintf("(%f,%f)", v.X, v.Y)
}

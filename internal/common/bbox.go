package common

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AABB struct {
	min, max Vec2
}

func (a1 *AABB) Contains(a2 AABB) bool {
	return a1.min.x < a2.min.x && a1.min.y < a2.min.y && a1.max.x > a2.max.x && a1.max.y > a2.max.y
}

func (b *AABB) ContainsPoint(p Vec2) bool {
	return b.min.x < p.x && b.min.y < p.y && b.max.x > p.x && b.max.y > p.y
}

func (a1 *AABB) Intersects(a2 AABB) bool {
	return !(a2.max.x > a1.min.x || a2.min.x > a1.max.x || a2.max.y < a1.min.y || a2.min.y > a1.max.y)
}

func (bb *AABB) Dimensions() Vec2 {
	return Vec2{
		bb.max.x - bb.min.x,
		bb.max.y - bb.min.y,
	}
}

func (bb *AABB) Area() float32 {
	dx := bb.max.x - bb.min.x
	dy := bb.max.y - bb.min.y
	return dx * dy
}

func (bb *AABB) Draw(screen *ebiten.Image, color color.Color) {
	vector.DrawFilledRect(
		screen,
		bb.min.x,
		bb.min.y,
		bb.Dimensions().x,
		bb.Dimensions().y,
		color,
		true,
	)
	/*
		vector.DrawFilledCircle(screen, bb.min.x, bb.min.y, 2, color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, true)
		vector.DrawFilledCircle(screen, bb.max.x, bb.min.y, 2, color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, true)
		vector.DrawFilledCircle(screen, bb.min.x, bb.max.y, 2, color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, true)
		vector.DrawFilledCircle(screen, bb.max.x, bb.max.y, 2, color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, true)
	*/
}

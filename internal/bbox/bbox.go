package bbox

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/KasperJSdeVries/zoned_world_servers/internal/vec"
)

type AABB struct {
	Min, Max vec.Vec2
}

func (a1 *AABB) Contains(a2 AABB) bool {
	return a1.Min.X < a2.Min.X && a1.Min.Y < a2.Min.Y && a1.Max.X > a2.Max.X && a1.Max.Y > a2.Max.Y
}

func (b *AABB) ContainsPoint(p vec.Vec2) bool {
	return b.Min.X < p.X && b.Min.Y < p.Y && b.Max.X > p.X && b.Max.Y > p.Y
}

func (a1 *AABB) Intersects(a2 AABB) bool {
	return !(a2.Max.X > a1.Min.X || a2.Min.X > a1.Max.X || a2.Max.Y < a1.Min.Y || a2.Min.Y > a1.Max.Y)
}

func (bb *AABB) Dimensions() vec.Vec2 {
	return vec.Vec2{
		X: bb.Max.X - bb.Min.X,
		Y: bb.Max.Y - bb.Min.Y,
	}
}

func (bb *AABB) Area() float32 {
	dx := bb.Max.X - bb.Min.X
	dy := bb.Max.Y - bb.Min.Y
	return dx * dy
}

func (bb *AABB) Draw(screen *ebiten.Image, color color.Color) {
	vector.DrawFilledRect(
		screen,
		bb.Min.X,
		bb.Min.Y,
		bb.Dimensions().X,
		bb.Dimensions().Y,
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

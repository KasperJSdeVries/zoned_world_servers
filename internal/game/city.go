package game

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type City struct {
	posX, posY float64
	radius     float64
	color      color.Color
}

func (c *City) Init() {
	c.posX = rand.Float64() * ScreenWidth * scale
	c.posY = rand.Float64() * ScreenHeight * scale
	c.radius = (rand.Float64()*(cityMaxRadius-cityMinRadius) + cityMinRadius) * scale
	c.color = color.RGBA{
		R: uint8(rand.Uint32()),
		G: uint8(rand.Uint32()),
		B: uint8(rand.Uint32()),
		A: uint8(math.Floor(cityAlpha * 0xff)),
	}
}

func (c *City) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(c.posX/scale), float32(c.posY/scale), float32(c.radius/scale), c.color, true)
}

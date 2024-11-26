package main

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type City struct {
	posX, posY float32
	radius     float32
	color      color.Color
}

func (c *City) Init() {
	c.posX = rand.Float32() * screenWidth * scale
	c.posY = rand.Float32() * screenHeight * scale
	c.radius = (rand.Float32()*(cityMaxRadius-cityMinRadius) + cityMinRadius) * scale
	c.color = color.RGBA{
		R: uint8(rand.Uint32()),
		G: uint8(rand.Uint32()),
		B: uint8(rand.Uint32()),
		A: uint8(math.Floor(cityAlpha * 0xff)),
	}
}

func (c *City) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, c.posX/scale, c.posY/scale, c.radius/scale, c.color, true)
}

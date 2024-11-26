package main

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TargetType uint

const (
	TargetTypePosition TargetType = iota
	TargetTypeCity
)

type PositionTarget struct {
	x, y float32
}

func (t PositionTarget) Type() TargetType             { return TargetTypePosition }
func (t PositionTarget) Position() (float32, float32) { return t.x, t.y }
func (t PositionTarget) Radius() float32              { return 1 * scale }

type CityTarget struct {
	city *City
}

func (t CityTarget) Type() TargetType             { return TargetTypeCity }
func (t CityTarget) Position() (float32, float32) { return t.city.posX, t.city.posY }
func (t CityTarget) Radius() float32              { return t.city.radius }

type Target interface {
	Type() TargetType
	Position() (float32, float32)
	Radius() float32
}

type Player struct {
	posX, posY float32
	target     Target
	color      color.Color
}

func (p *Player) Init(cities [cityCount]City) {
	p.posX = rand.Float32() * screenWidth * scale
	p.posY = rand.Float32() * screenHeight * scale
	p.NewTarget(cities)
}

func (p *Player) NewTarget(cities [cityCount]City) {
	if rand.Float32() < cityTargetChance {
		city := cities[rand.UintN(cityCount)]
		p.target = CityTarget{city: &city}
		p.color = city.color
	} else {
		p.target = PositionTarget{
			rand.Float32() * screenWidth * scale,
			rand.Float32() * screenHeight * scale,
		}
		p.color = color.RGBA{
			R: uint8(rand.Uint32()),
			G: uint8(rand.Uint32()),
			B: uint8(rand.Uint32()),
			A: 0xff,
		}
	}
}

func (p *Player) Update(cities [cityCount]City) {
	tx, ty := p.target.Position()

	dx := tx - p.posX
	dy := ty - p.posY

	mag2 := dx*dx + dy*dy
	mag := float32(math.Sqrt(float64(mag2)))
	speed := playerSpeed
	if p.target.Type() == TargetTypeCity {
		speed *= 2
	}
	p.posX += (1.0 / 60.0) * (dx / mag) * scale * float32(speed)
	p.posY += (1.0 / 60.0) * (dy / mag) * scale * float32(speed)

	if dx*dx+dy*dy < p.target.Radius()*p.target.Radius() {
		p.NewTarget(cities)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, p.posX/scale, p.posY/scale, playerSize, p.color, true)
}

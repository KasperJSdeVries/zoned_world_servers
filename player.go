package main

import (
	"image/color"
	"log"
	"math"
	"math/rand/v2"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TargetType uint

const (
	TargetTypePosition TargetType = iota
	TargetTypeCity
)

type PositionTarget struct {
	x, y float64
}

func (t PositionTarget) Type() TargetType             { return TargetTypePosition }
func (t PositionTarget) Position() (float64, float64) { return t.x, t.y }
func (t PositionTarget) Radius() float64              { return 1 * scale }

type CityTarget struct {
	city *City
}

func (t CityTarget) Type() TargetType             { return TargetTypeCity }
func (t CityTarget) Position() (float64, float64) { return t.city.posX, t.city.posY }
func (t CityTarget) Radius() float64              { return t.city.radius }

type Target interface {
	Type() TargetType
	Position() (float64, float64)
	Radius() float64
}

type Player struct {
	X, Y   float64
	target Target
	color  color.Color
	ID     int
}

var playerId int

func (p *Player) Init(cities [cityCount]City) {
	p.ID = playerId
	playerId++
	p.X = rand.Float64() * screenWidth * scale
	p.Y = rand.Float64() * screenHeight * scale
	p.NewTarget(cities)
}

func (p *Player) Bounds() rtreego.Rect {
	rect, err := rtreego.NewRectFromPoints(
		rtreego.Point{float64(p.X), float64(p.Y)},
		rtreego.Point{float64(p.X), float64(p.Y)},
	)
	if err != nil {
		log.Fatal("could not create bounds for player")
	}
	return rect
}

func (p *Player) NewTarget(cities [cityCount]City) {
	if rand.Float32() < cityTargetChance {
		city := cities[rand.UintN(cityCount)]
		p.target = CityTarget{city: &city}
		p.color = city.color
	} else {
		p.target = PositionTarget{
			rand.Float64() * screenWidth * scale,
			rand.Float64() * screenHeight * scale,
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

	dx := tx - p.X
	dy := ty - p.Y

	mag2 := dx*dx + dy*dy
	mag := math.Sqrt(mag2)
	speed := playerSpeed
	if p.target.Type() == TargetTypeCity {
		speed *= 2
	}
	p.X += (1.0 / 60.0) * (dx / mag) * scale * speed
	p.Y += (1.0 / 60.0) * (dy / mag) * scale * speed

	if dx*dx+dy*dy < p.target.Radius()*p.target.Radius() {
		p.NewTarget(cities)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(p.X/scale), float32(p.Y/scale), playerSize, p.color, true)
}

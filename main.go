package main

import (
	"image/color"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth      = 640
	screenHeight     = 460
	scale            = 64
	playerCount      = 1024
	playerSpeed      = 30
	cityCount        = 8
	cityMinRadius    = 5
	cityMaxRadius    = 20
	cityAlpha        = 0.5
	cityTargetChance = 0.7
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
	if (p.target.Type() == TargetTypeCity) {
		speed *= 2;
	}
	p.posX += (1.0 / 60.0) * (dx / mag) * scale * float32(speed)
	p.posY += (1.0 / 60.0) * (dy / mag) * scale * float32(speed)

	if dx*dx+dy*dy < p.target.Radius()*p.target.Radius() {
		p.NewTarget(cities)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, p.posX/scale, p.posY/scale, 1, p.color, true)
}

type Game struct {
	cities  [cityCount]City
	players [playerCount]Player
}

func NewGame() *Game {
	g := &Game{}
	for i := 0; i < cityCount; i++ {
		g.cities[i].Init()
	}
	for i := 0; i < playerCount; i++ {
		g.players[i].Init(g.cities)
	}
	return g
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		for i := 0; i < playerCount; i++ {
			g.players[i].target = PositionTarget{
				x: float32(x) * scale,
				y: float32(y) * scale,
			}
		}
	}

	for i := 0; i < playerCount; i++ {
		g.players[i].Update(g.cities)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < cityCount; i++ {
		g.cities[i].Draw(screen)
	}
	for i := 0; i < playerCount; i++ {
		g.players[i].Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Test")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

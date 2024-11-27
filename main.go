package main

import (
	"image/color"
	"log"

	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth      = 1280
	screenHeight     = 920
	scale            = 64
	playerCount      = 1024
	playerSize       = 2
	playerSpeed      = 30.0
	cityCount        = 8
	cityMinRadius    = 5
	cityMaxRadius    = 20
	cityAlpha        = 0.5
	cityTargetChance = 0.7
)

type Game struct {
	cities  [cityCount]City
	players [playerCount]Player
	rt      *rtreego.Rtree
	tree    *Quadtree
}

func NewGame() *Game {
	g := &Game{
		rt:   rtreego.NewTree(2, 8, 16),
		tree: NewQuadtree(0, screenWidth*scale, 0, screenHeight*scale, 8, 4, 16),
	}
	for i := 0; i < cityCount; i++ {
		g.cities[i].Init()
	}
	for i := 0; i < playerCount; i++ {
		g.players[i].Init(g.cities)
		g.rt.Insert(&(g.players[i]))
		g.tree.AddPoint(g.players[i])
	}

	return g
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		for i := 0; i < playerCount; i++ {
			g.players[i].target = PositionTarget{
				x: float64(x) * scale,
				y: float64(y) * scale,
			}
		}
	}

	for i := 0; i < playerCount; i++ {
		g.rt.Delete(&g.players[i])
		g.players[i].Update(g.cities)
		g.tree.MovePoint(g.players[i])
		g.rt.Insert(&g.players[i])
	}

	g.tree.BalanceRegions()

	return nil
}

func DrawRect(screen *ebiten.Image, rect rtreego.Rect) {
	min := Vec2{float32(rect.PointCoord(0)) / scale, float32(rect.PointCoord(1)) / scale}
	max := Vec2{min.x + float32(rect.LengthsCoord(0))/scale, min.y + float32(rect.LengthsCoord(1))/scale}
	(&AABB{min, max}).Draw(screen, color.RGBA{0x0f, 0x0b, 0x0d, 0x09})
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < cityCount; i++ {
		g.cities[i].Draw(screen)
	}
	for i := 0; i < playerCount; i++ {
		g.players[i].Draw(screen)
	}
	for _, bbox := range g.rt.GetAllBoundingBoxes() {
		DrawRect(screen, bbox)
	}
	g.tree.DebugRegions(screen)
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

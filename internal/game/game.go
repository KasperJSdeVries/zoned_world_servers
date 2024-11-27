package game

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth      = 1280
	ScreenHeight     = 920
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
	tree    *Quadtree
}

func NewGame() *Game {
	g := &Game{
		tree: NewQuadtree(0, ScreenWidth*scale, 0, ScreenHeight*scale, 2 * scale, 8, 16),
	}
	for i := 0; i < cityCount; i++ {
		g.cities[i].Init()
	}
	for i := 0; i < playerCount; i++ {
		g.players[i].Init(g.cities)
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
		g.players[i].Update(g.cities)
		g.tree.MovePoint(g.players[i])
	}

	for g.tree.BalanceRegions() {
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
	g.tree.DebugRegions(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

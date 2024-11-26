package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth      = 1280
	screenHeight     = 920
	scale            = 64
	playerCount      = 1024
	playerSize       = 2
	playerSpeed      = 30
	cityCount        = 8
	cityMinRadius    = 5
	cityMaxRadius    = 20
	cityAlpha        = 0.5
	cityTargetChance = 0.7
)

type Vec2 struct {
	x, y float32
}

type AABB struct {
	min, max Vec2
}

func (a1 *AABB) Contains(a2 AABB) bool {
	return a1.min.x < a2.min.x && a1.min.y < a2.min.y && a1.max.x > a2.max.x && a1.max.y > a2.max.y
}

func (a1 *AABB) Intersects(a2 AABB) bool {
	return !(a2.max.x > a1.min.x || a2.min.x > a1.max.x || a2.max.y < a1.min.y || a2.min.y > a1.max.y)
}

func (bb *AABB) Area() float32 {
	dx := bb.max.x - bb.min.x
	dy := bb.max.y - bb.min.y
	return dx * dy
}

const (
	maxItems = 64
	minItems = maxItems / 2
)

type NodeKind uint8

const (
	NodeKindLeaf NodeKind = iota
	NodeKindBranch
)

type Node struct {
	kind     NodeKind
	bounds   AABB
	datas [maxItems]*Player
	children [maxItems]*Node
}

func (n *Node) Search(area AABB) []uint64 {
	switch n.kind {
	case NodeKindBranch:
		var found []uint64
		for _, e := range n.children {
			if area.Contains(e.bounds) {
				found = append(found, e.Search(area)...)
			}
		}
		return found
	case NodeKindLeaf:
		var found []uint64
		for _, e := range n.datas {
			if area.Intersects(e.bounds) {
				found = append(found, e)
			}
		}
		return found
	}

	return []uint64{}
}

func (n *Node) ChooseLeaf(bounds AABB) *Node {
	if n.kind == NodeKindLeaf {
		return n
	}

	var best *Node
	metric := math.Inf(1)
	for _, f := range n.children {
		if f.bounds.Contains(bounds) {
			if best == nil || f.bounds.Area() < best.bounds.Area() {
				metric = 0
				best = f
			}
		} else if metric != 0 {
			fMetric := 0.0
			fMetric += math.Min(0, float64(f.bounds.min.x-bounds.min.x))
			fMetric += math.Min(0, float64(f.bounds.min.y-bounds.min.y))
			fMetric += math.Min(0, float64(bounds.max.x-f.bounds.max.x))
			fMetric += math.Min(0, float64(bounds.max.y-f.bounds.max.y))
			if fMetric < metric {
				best = f
				metric = fMetric
			}
		}
	}
	return best.ChooseLeaf(bounds)
}

type RTree struct {
	bounds AABB
	root   *Node
	items  map[uint64]*Player
	count  uint64
}

func (r *RTree) Search(area AABB) []*Player {
	ids := r.root.Search(area)
	var players []*Player
	for _, id := range ids {
		players = append(players, r.items[id])
	}
	return players
}

func (r *RTree) Insert(player *Player) {
	newItemId := r.count
	r.count++
	r.items[newItemId] = player

	itemBounds := AABB{min: Vec2{player.posX, player.posY}, max: Vec2{player.posX, player.posY}}
	l := r.root.ChooseLeaf(itemBounds)
	if l.
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

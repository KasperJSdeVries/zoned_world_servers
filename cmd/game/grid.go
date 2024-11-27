package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type RegionKind uint8

const (
	BranchRegion RegionKind = iota
	LeafRegion
)

// Region represents a rectangular region with elastic boundaries.
type Region struct {
	MinX, MaxX, MinY, MaxY float64 // Core boundaries
	Buffer                 float64 // Elastic buffer
	Points                 map[int]Player
	Children               [4]int
	Parent                 int
	regionKind             RegionKind
}

// Quadtree represents the overall structure of adaptive regions.
type Quadtree struct {
	Regions    map[int]*Region
	rootRegion int
	count      int
	MinPoints  int
	MaxPoints  int
}

// NewRegion creates a new region.
func NewRegion(minX, maxX, minY, maxY, buffer float64, parent int) *Region {
	return &Region{
		MinX:       minX,
		MaxX:       maxX,
		MinY:       minY,
		MaxY:       maxY,
		Buffer:     buffer,
		Points:     make(map[int]Player),
		Parent:     parent,
		regionKind: LeafRegion,
		Children:   [4]int{-1, -1, -1, -1},
	}
}

// NewQuadtree initializes a quadtree with one large region.
func NewQuadtree(minX, maxX, minY, maxY float64, buffer float64, minPoints, maxPoints int) *Quadtree {
	rootRegion := NewRegion(minX, maxX, minY, maxY, buffer, -1)
	qt := &Quadtree{
		Regions:   make(map[int]*Region),
		count:     0,
		MinPoints: minPoints,
		MaxPoints: maxPoints,
	}
	qt.rootRegion = qt.AddRegion(rootRegion)
	return qt
}

func (qt *Quadtree) AddRegion(r *Region) int {
	ri := qt.count
	qt.count++

	qt.Regions[ri] = r
	return ri
}

// AddPoint assigns a point to the correct region.
func (qt *Quadtree) AddPoint(p Player) {
	for ri, region := range qt.Regions {
		if qt.pointInRegion(p, region) {
			region.Points[p.ID] = p
			p.RegionId = ri
			return
		}
	}

	// If no region is found, the point is lost.
	fmt.Printf("Point %v is out of bounds!\n", p)
}

// MovePoint updates a point's location and adjusts regions if needed.
func (qt *Quadtree) MovePoint(p Player) {
	pr := qt.Regions[p.RegionId]
	if pr != nil {
		if pr.regionKind == LeafRegion && qt.pointInRegion(p, pr) {
			pr.Points[p.ID] = p
			return
		}
		if pr.regionKind == LeafRegion && qt.pointInBuffer(p, pr) {
			pr.Points[p.ID] = p
			return
		}
	}

	qt.reassignPoint(p)
}

// reassignPoint moves a point to the nearest region.
func (qt *Quadtree) reassignPoint(p Player) {
	minDistance := math.MaxFloat64
	nearestRegion := -1

	for ri, region := range qt.Regions {
		if region == nil {
			continue
		}
		if region.regionKind == LeafRegion {
			dist := qt.distanceToRegion(p, region)
			if dist < minDistance {
				minDistance = dist
				nearestRegion = ri
			}
		}
	}

	if nearestRegion != -1 {
		qt.Regions[nearestRegion].Points[p.ID] = p
		p.RegionId = nearestRegion
	}
}

// pointInRegion checks if a point is within the core boundary of a region.
func (qt *Quadtree) pointInRegion(p Player, r *Region) bool {
	return p.X >= r.MinX && p.X < r.MaxX && p.Y >= r.MinY && p.Y < r.MaxY
}

// pointInBuffer checks if a point is within the elastic boundary of a region.
func (qt *Quadtree) pointInBuffer(p Player, r *Region) bool {
	return p.X >= r.MinX-r.Buffer && p.X < r.MaxX+r.Buffer && p.Y >= r.MinY-r.Buffer && p.Y < r.MaxY+r.Buffer
}

// distanceToRegion calculates the Euclidean distance from a point to a region.
func (qt *Quadtree) distanceToRegion(p Player, r *Region) float64 {
	centerX := (r.MinX + r.MaxX) / 2
	centerY := (r.MinY + r.MaxY) / 2
	return math.Sqrt((p.X-centerX)*(p.X-centerX) + (p.Y-centerY)*(p.Y-centerY))
}

// BalanceRegions checks and rebalances regions.
func (qt *Quadtree) BalanceRegions() bool {
	didBalance := false
	for i, region := range qt.Regions {
		if region == nil {
			continue
		}

		switch region.regionKind {
		case BranchRegion:
			if qt.Regions[region.Children[0]].regionKind != LeafRegion {
				continue
			}

			pc := 0
			for _, ci := range region.Children {
				c := qt.Regions[ci]
				pc += len(c.Points)
			}

			if pc < qt.MinPoints {
				qt.mergeRegion(i)
				didBalance = true
				break
			}
		case LeafRegion:
			if len(region.Points) > qt.MaxPoints {
				qt.splitRegion(i)
				didBalance = true
				break
			}
		}
	}
	return didBalance
}

// mergeRegion merges a sparse region with its neighbors.
func (qt *Quadtree) mergeRegion(ri int) {
	fmt.Println("merging region", ri)
	r := qt.Regions[ri]

	r.Points = map[int]Player{}
	for _, ci := range r.Children {
		c := qt.Regions[ci]
		for k, v := range c.Points {
			r.Points[k] = v
		}
		qt.Regions[ci] = nil
	}

	r.Children = [4]int{-1, -1, -1, -1}
	r.regionKind = LeafRegion
}

// splitRegion splits a dense region into sub-regions.
func (qt *Quadtree) splitRegion(ri int) {
	fmt.Println("splitting region", ri)
	r := qt.Regions[ri]

	// Subdivide splits a region into 4 sub-regions.
	midX := (r.MinX + r.MaxX) / 2
	midY := (r.MinY + r.MaxY) / 2

	newRegions := []*Region{
		NewRegion(r.MinX, midX, r.MinY, midY, r.Buffer/4, ri), // Top-left
		NewRegion(midX, r.MaxX, r.MinY, midY, r.Buffer/4, ri), // Top-right
		NewRegion(r.MinX, midX, midY, r.MaxY, r.Buffer/4, ri), // Bottom-left
		NewRegion(midX, r.MaxX, midY, r.MaxY, r.Buffer/4, ri), // Bottom-right
	}

	// Reassign points to sub-regions.
	for _, p := range r.Points {
		for _, sub := range newRegions {
			if qt.pointInRegion(p, sub) {
				sub.Points[p.ID] = p
			}
		}
	}

	// Clear points from parent r.
	r.Points = nil
	r.regionKind = BranchRegion

	// Add new regions to quadtree
	for i, region := range newRegions {
		r.Children[i] = qt.AddRegion(region)
	}
}

func (qt *Quadtree) DrawRegion(screen *ebiten.Image, ri int) {
	r := qt.Regions[ri]

	switch r.regionKind {
	case LeafRegion:
		vector.StrokeRect(screen, float32(r.MinX)/scale, float32(r.MinY)/scale, (float32(r.MaxX)-float32(r.MinX))/scale, (float32(r.MaxY)-float32(r.MinY))/scale, 1, color.RGBA{0x0b, 0x0f, 0xa, 0x09}, true)
		break
	case BranchRegion:
		for _, c := range r.Children {
			qt.DrawRegion(screen, c)
		}
	}
}

// DebugRegions prints the current state of regions.
func (qt *Quadtree) DebugRegions(screen *ebiten.Image) {
	qt.DrawRegion(screen, qt.rootRegion)
}

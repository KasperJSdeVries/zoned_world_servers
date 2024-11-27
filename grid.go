package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Region represents a rectangular region with elastic boundaries.
type Region struct {
	MinX, MaxX, MinY, MaxY float64 // Core boundaries
	Buffer                 float64 // Elastic buffer
	Points                 map[int]Player
	ID                     int
}

func (r *Region) Join(n *Region) {
	r.MinX = min(r.MinX, n.MinX)
	r.MaxX = min(r.MaxX, n.MaxX)
	r.MinY = min(r.MinY, n.MinY)
	r.MaxY = min(r.MaxY, n.MaxY)
	for k, v := range n.Points {
		r.Points[k] = v
	}
}

// Quadtree represents the overall structure of adaptive regions.
type Quadtree struct {
	Regions   []*Region
	MinPoints int
	MaxPoints int
}

// NewRegion creates a new region.
func NewRegion(minX, maxX, minY, maxY, buffer float64, id int) *Region {
	return &Region{
		MinX:   minX,
		MaxX:   maxX,
		MinY:   minY,
		MaxY:   maxY,
		Buffer: buffer,
		Points: make(map[int]Player),
		ID:     id,
	}
}

// NewQuadtree initializes a quadtree with one large region.
func NewQuadtree(minX, maxX, minY, maxY float64, buffer float64, minPoints, maxPoints int) *Quadtree {
	rootRegion := NewRegion(minX, maxX, minY, maxY, buffer, 0)
	return &Quadtree{
		Regions:   []*Region{rootRegion},
		MinPoints: minPoints,
		MaxPoints: maxPoints,
	}
}

// AddPoint assigns a point to the correct region.
func (qt *Quadtree) AddPoint(p Player) {
	for _, region := range qt.Regions {
		if qt.pointInRegion(p, region) {
			region.Points[p.ID] = p
			return
		}
	}

	// If no region is found, the point is lost.
	fmt.Printf("Point %v is out of bounds!\n", p)
}

// MovePoint updates a point's location and adjusts regions if needed.
func (qt *Quadtree) MovePoint(p Player) {
	for _, region := range qt.Regions {
		if qt.pointInRegion(p, region) {
			// Update point in the region if it hasn't left the buffer.
			region.Points[p.ID] = p
			return
		}
	}

	// If the point has left its region, reassign it.
	for _, region := range qt.Regions {
		if qt.pointInBuffer(p, region) {
			region.Points[p.ID] = p
			return
		}
	}

	// If no valid region is found, reassign to the closest region.
	qt.reassignPoint(p)
}

// reassignPoint moves a point to the nearest region.
func (qt *Quadtree) reassignPoint(p Player) {
	minDistance := math.MaxFloat64
	var nearestRegion *Region

	for _, region := range qt.Regions {
		dist := qt.distanceToRegion(p, region)
		if dist < minDistance {
			minDistance = dist
			nearestRegion = region
		}
	}

	if nearestRegion != nil {
		nearestRegion.Points[p.ID] = p
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
func (qt *Quadtree) BalanceRegions() {
	for i, region := range qt.Regions {
		if len(region.Points) < qt.MinPoints {
			qt.mergeRegion(region)
		} else if len(region.Points) > qt.MaxPoints {
			qt.splitRegion(i)
		}
	}
}

// mergeRegion merges a sparse region with its neighbors.
func (qt *Quadtree) mergeRegion(r *Region) {
	neighbor := -1
	for i, region := range qt.Regions {
		if region.MaxX == r.MinX && r.MinY == region.MinY && r.MaxY == region.MaxY && len(region.Points)+len(r.Points) > qt.MinPoints {
			neighbor = i
			break
		}
		if region.MinX == r.MaxX && r.MinY == region.MinY && r.MaxY == region.MaxY && len(region.Points)+len(r.Points) > qt.MinPoints {
			neighbor = i
			break
		}
		if region.MaxY == r.MinY && r.MinX == region.MinX && r.MaxX == region.MaxX && len(region.Points)+len(r.Points) > qt.MinPoints {
			neighbor = i
			break
		}
		if region.MinY == r.MaxY && r.MinX == region.MinX && r.MaxX == region.MaxX && len(region.Points)+len(r.Points) > qt.MinPoints {
			neighbor = i
			break
		}
	}
	if neighbor != -1 {
		n := qt.Regions[neighbor]
		qt.Regions = append(qt.Regions[:neighbor], qt.Regions[neighbor+1:]...)
		r.Join(n)
	}
}

// splitRegion splits a dense region into sub-regions.
func (qt *Quadtree) splitRegion(regionIndex int) {
	r := qt.Regions[regionIndex]

	// Subdivide splits a region into 4 sub-regions.
	midX := (r.MinX + r.MaxX) / 2
	midY := (r.MinY + r.MaxY) / 2

	newRegions := []*Region{
		{MinX: r.MinX, MaxX: midX, MinY: r.MinY, MaxY: midY, Points: make(map[int]Player)}, // Top-left
		{MinX: midX, MaxX: r.MaxX, MinY: r.MinY, MaxY: midY, Points: make(map[int]Player)}, // Top-right
		{MinX: r.MinX, MaxX: midX, MinY: midY, MaxY: r.MaxY, Points: make(map[int]Player)}, // Bottom-left
		{MinX: midX, MaxX: r.MaxX, MinY: midY, MaxY: r.MaxY, Points: make(map[int]Player)}, // Bottom-right
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
	qt.Regions = append(qt.Regions[:regionIndex], qt.Regions[regionIndex+1:]...)

	qt.Regions = append(qt.Regions, newRegions...)
}

// DebugRegions prints the current state of regions.
func (qt *Quadtree) DebugRegions(screen *ebiten.Image) {
	for _, region := range qt.Regions {
		vector.StrokeRect(screen, float32(region.MinX)/scale, float32(region.MinY)/scale, (float32(region.MaxX)-float32(region.MinX))/scale, (float32(region.MaxY)-float32(region.MinY))/scale, 1, color.RGBA{0x0b, 0x0f, 0xa, 0x09}, true)
	}
}

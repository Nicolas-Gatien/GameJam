package main

type Rect struct {
	Position
	width  float64
	height float64
}

func NewRect(x, y, width, height float64) Rect {
	return Rect{
		Position: Position{positionX: x, positionY: y},
		width:    width,
		height:   height,
	}
}

func (r Rect) MaxX() float64 {
	return r.positionX + r.width
}

func (r Rect) MaxY() float64 {
	return r.positionY + r.height
}

func (r Rect) Intersects(other Rect) bool {
	return r.positionX <= other.MaxX() && other.positionX <= r.MaxX() && r.positionY <= other.MaxY() && other.positionY <= r.MaxY()
}

func (p *Gecko) Collider() Rect {
	_ = p.SPRITE.Bounds()

	// return NewRect(
	// 	p.positionX,
	// 	p.positionY,
	// )
	return NewRect(0, 0, 0, 0)
}

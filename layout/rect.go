package layout

// Rect represents a rectangular region with z-ordering
type Rect struct {
	X      int
	Y      int
	Z      int
	Width  int
	Height int
}

// NewRect creates a new rectangle with z-ordering
func NewRect(x, y, z, width, height int) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Z:      z,
		Width:  width,
		Height: height,
	}
}

// NewRectXY creates a new rectangle with default z=0
func NewRectXY(x, y, width, height int) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Z:      0,
		Width:  width,
		Height: height,
	}
}

// Zero returns a zero-sized rectangle at origin
func Zero() Rect {
	return Rect{}
}

// WithZ returns a copy of the rectangle with a different z value
func (r Rect) WithZ(z int) Rect {
	r.Z = z
	return r
}

// Right returns the x coordinate of the right edge
func (r Rect) Right() int {
	return r.X + r.Width
}

// Bottom returns the y coordinate of the bottom edge
func (r Rect) Bottom() int {
	return r.Y + r.Height
}

// Area returns the area of the rectangle
func (r Rect) Area() int {
	return r.Width * r.Height
}

// IsEmpty returns true if the rectangle has zero area
func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Contains checks if a point is inside the rectangle
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.Right() && y >= r.Y && y < r.Bottom()
}

// Intersects checks if two rectangles overlap
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.Right() && r.Right() > other.X &&
		r.Y < other.Bottom() && r.Bottom() > other.Y
}

// Intersection returns the intersection of two rectangles
// Preserves Z from the receiver
func (r Rect) Intersection(other Rect) Rect {
	x := max(r.X, other.X)
	y := max(r.Y, other.Y)
	right := min(r.Right(), other.Right())
	bottom := min(r.Bottom(), other.Bottom())

	if right <= x || bottom <= y {
		return Zero()
	}

	return NewRect(x, y, r.Z, right-x, bottom-y)
}

// Union returns the smallest rectangle containing both rectangles
// Preserves Z from the receiver
func (r Rect) Union(other Rect) Rect {
	if r.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return r
	}

	x := min(r.X, other.X)
	y := min(r.Y, other.Y)
	right := max(r.Right(), other.Right())
	bottom := max(r.Bottom(), other.Bottom())

	return NewRect(x, y, r.Z, right-x, bottom-y)
}

// Inset returns a new rectangle inset by the given amounts
// Preserves Z from the receiver
func (r Rect) Inset(top, right, bottom, left int) Rect {
	return NewRect(
		r.X+left,
		r.Y+top,
		r.Z,
		max(0, r.Width-left-right),
		max(0, r.Height-top-bottom),
	)
}

// InsetAll returns a new rectangle inset by the same amount on all sides
func (r Rect) InsetAll(amount int) Rect {
	return r.Inset(amount, amount, amount, amount)
}

// Offset returns a new rectangle offset by the given amounts
// Preserves Z from the receiver
func (r Rect) Offset(dx, dy int) Rect {
	return NewRect(r.X+dx, r.Y+dy, r.Z, r.Width, r.Height)
}

// Center returns the center point of the rectangle
func (r Rect) Center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

// Equals checks if two rectangles are equal (including Z)
func (r Rect) Equals(other Rect) bool {
	return r.X == other.X && r.Y == other.Y && r.Z == other.Z &&
		r.Width == other.Width && r.Height == other.Height
}

// SplitHorizontal splits the rectangle horizontally at the given offset
// Both resulting rectangles preserve Z from the receiver
func (r Rect) SplitHorizontal(offset int) (Rect, Rect) {
	if offset < 0 {
		offset = 0
	}
	if offset > r.Height {
		offset = r.Height
	}

	top := NewRect(r.X, r.Y, r.Z, r.Width, offset)
	bottom := NewRect(r.X, r.Y+offset, r.Z, r.Width, r.Height-offset)
	return top, bottom
}

// SplitVertical splits the rectangle vertically at the given offset
// Both resulting rectangles preserve Z from the receiver
func (r Rect) SplitVertical(offset int) (Rect, Rect) {
	if offset < 0 {
		offset = 0
	}
	if offset > r.Width {
		offset = r.Width
	}

	left := NewRect(r.X, r.Y, r.Z, offset, r.Height)
	right := NewRect(r.X+offset, r.Y, r.Z, r.Width-offset, r.Height)
	return left, right
}

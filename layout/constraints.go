package layout

// Size represents a width and height
type Size struct {
	Width  int
	Height int
}

// NewSize creates a new size
func NewSize(width, height int) Size {
	return Size{Width: width, Height: height}
}

// Constraint represents size constraints for layout
type Constraint struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

// Unconstrained returns a constraint with no limits
func Unconstrained() Constraint {
	return Constraint{
		MinWidth:  0,
		MaxWidth:  MaxInt,
		MinHeight: 0,
		MaxHeight: MaxInt,
	}
}

// Exact returns a constraint for an exact size
func Exact(width, height int) Constraint {
	return Constraint{
		MinWidth:  width,
		MaxWidth:  width,
		MinHeight: height,
		MaxHeight: height,
	}
}

// AtMost returns a constraint with maximum bounds
func AtMost(width, height int) Constraint {
	return Constraint{
		MinWidth:  0,
		MaxWidth:  width,
		MinHeight: 0,
		MaxHeight: height,
	}
}

// AtLeast returns a constraint with minimum bounds
func AtLeast(width, height int) Constraint {
	return Constraint{
		MinWidth:  width,
		MaxWidth:  MaxInt,
		MinHeight: height,
		MaxHeight: MaxInt,
	}
}

// Between returns a constraint with both min and max bounds
func Between(minW, maxW, minH, maxH int) Constraint {
	return Constraint{
		MinWidth:  minW,
		MaxWidth:  maxW,
		MinHeight: minH,
		MaxHeight: maxH,
	}
}

// Constrain applies the constraints to a size
func (c Constraint) Constrain(size Size) Size {
	return Size{
		Width:  clamp(size.Width, c.MinWidth, c.MaxWidth),
		Height: clamp(size.Height, c.MinHeight, c.MaxHeight),
	}
}

// ConstrainWidth applies width constraints
func (c Constraint) ConstrainWidth(width int) int {
	return clamp(width, c.MinWidth, c.MaxWidth)
}

// ConstrainHeight applies height constraints
func (c Constraint) ConstrainHeight(height int) int {
	return clamp(height, c.MinHeight, c.MaxHeight)
}

// ShrinkWidth returns a new constraint with reduced width
func (c Constraint) ShrinkWidth(amount int) Constraint {
	return Constraint{
		MinWidth:  max(0, c.MinWidth-amount),
		MaxWidth:  max(0, c.MaxWidth-amount),
		MinHeight: c.MinHeight,
		MaxHeight: c.MaxHeight,
	}
}

// ShrinkHeight returns a new constraint with reduced height
func (c Constraint) ShrinkHeight(amount int) Constraint {
	return Constraint{
		MinWidth:  c.MinWidth,
		MaxWidth:  c.MaxWidth,
		MinHeight: max(0, c.MinHeight-amount),
		MaxHeight: max(0, c.MaxHeight-amount),
	}
}

// Shrink returns a new constraint reduced by the given amounts
func (c Constraint) Shrink(width, height int) Constraint {
	return c.ShrinkWidth(width).ShrinkHeight(height)
}

// HasBoundedWidth returns true if width has an upper bound
func (c Constraint) HasBoundedWidth() bool {
	return c.MaxWidth < MaxInt
}

// HasBoundedHeight returns true if height has an upper bound
func (c Constraint) HasBoundedHeight() bool {
	return c.MaxHeight < MaxInt
}

// IsTight returns true if min and max are equal for both dimensions
func (c Constraint) IsTight() bool {
	return c.MinWidth == c.MaxWidth && c.MinHeight == c.MaxHeight
}

// MaxInt is a large integer for representing unbounded constraints
const MaxInt = 1<<31 - 1

// clamp restricts a value to a range
func clamp(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// Alignment represents alignment within a container
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// Align calculates the offset for aligning content within a container
func Align(contentSize, containerSize int, alignment Alignment) int {
	switch alignment {
	case AlignCenter:
		return (containerSize - contentSize) / 2
	case AlignEnd:
		return containerSize - contentSize
	default:
		return 0
	}
}


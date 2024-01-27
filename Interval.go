package interval

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

type IInterval[T constraints.Integer | constraints.Float] interface {
	Lower() T
	SetLower(lower T) error
	Upper() T
	SetUpper(upper T) error
	LowerUnbounded() bool
	SetLowerUnbounded(lowerUnbounded bool) error
	UpperUnbounded() bool
	SetUpperUnbounded(upperUnbounded bool) error
	LowerIncluded() bool
	SetLowerIncluded(lowerIncluded bool) error
	UpperIncluded() bool
	SetUpperIncluded(upperIncluded bool) error
	String() string
	Equal(x IInterval[T]) bool
	IsEmpty() bool
	LtBeginOf(x IInterval[T]) bool
	LeEndOf(x IInterval[T]) bool
	Contains(x IInterval[T]) bool
	Has(value T) bool
	Intersect(x IInterval[T]) IInterval[T]
	Move(x T) IInterval[T]
	Subtract(x IInterval[T]) (IInterval[T], IInterval[T])
	Adjoin(x IInterval[T]) IInterval[T]
	Encompass(x IInterval[T]) IInterval[T]
}

type Interval[T constraints.Integer | constraints.Float] struct {
	// begin of this interval.
	lower T
	// if lowerIncluded is true, this interval is inclusive of the lower point.
	lowerIncluded bool
	// if lowerUnbounded is true, the value of lower is ignored and the interval is considered endles on the lowerside
	lowerUnbounded bool
	// end of this interval.
	upper T
	// if upperIncluded is true, this interval is inclusive of the upper point.
	upperIncluded bool
	// if upperUnbounded is true, the value of upper is ignored and the interval is considered endles on the upperside
	upperUnbounded bool
}

func NewInterval[T constraints.Integer | constraints.Float](lower, upper T, lowerIncluded, lowerUnbounded, upperIncluded, upperUnbounded bool) *Interval[T] {
	interval := new(Interval[T])
	interval.lower = lower
	interval.upper = upper
	interval.lowerIncluded = lowerIncluded
	interval.upperIncluded = upperIncluded
	interval.lowerUnbounded = lowerUnbounded
	interval.upperUnbounded = upperUnbounded
	return interval
}

func (i *Interval[T]) Lower() T {
	return i.lower
}

func (i *Interval[T]) SetLower(lower T) error {
	i.lower = lower
	return nil
}

func (i *Interval[T]) Upper() T {
	return i.upper
}

func (i *Interval[T]) SetUpper(upper T) error {
	i.upper = upper
	return nil
}

func (i *Interval[T]) LowerUnbounded() bool {
	return i.lowerUnbounded
}

func (i *Interval[T]) SetLowerUnbounded(lowerUnbounded bool) error {
	i.lowerUnbounded = lowerUnbounded
	return nil
}

func (i *Interval[T]) UpperUnbounded() bool {
	return i.upperUnbounded
}

func (i *Interval[T]) SetUpperUnbounded(upperUnbounded bool) error {
	i.upperUnbounded = upperUnbounded
	return nil
}

func (i *Interval[T]) LowerIncluded() bool {
	return i.lowerIncluded
}

func (i *Interval[T]) SetLowerIncluded(lowerIncluded bool) error {
	i.lowerIncluded = lowerIncluded
	return nil
}

func (i *Interval[T]) UpperIncluded() bool {
	return i.upperIncluded
}

func (i *Interval[T]) SetUpperIncluded(upperIncluded bool) error {
	i.upperIncluded = upperIncluded
	return nil
}

func (i *Interval[T]) String() string {
	var b strings.Builder
	if i.lowerUnbounded {
		b.WriteByte('<')
	}
	if i.lowerIncluded {
		b.WriteByte('[')
	} else {
		b.WriteByte('(')
	}
	fmt.Fprintf(&b, "%v", i.lower)
	b.WriteString(", ")
	fmt.Fprintf(&b, "%v", i.upper)
	if i.upperIncluded {
		b.WriteByte(']')
	} else {
		b.WriteByte(')')
	}
	if i.upperUnbounded {
		b.WriteByte('>')
	}
	return b.String()
}

// Equal returns true if receiver interval is equals x_interval_string interval.
func (i *Interval[T]) Equal(x IInterval[T]) bool {
	if x == nil {
		return false
	}
	if x.UpperUnbounded() && i.upperUnbounded {
		return (i.lower == x.Lower() &&
			i.lowerIncluded == x.LowerIncluded() &&
			i.lowerUnbounded == x.LowerUnbounded()) || (x.IsEmpty() && i.IsEmpty())
	}
	if x.LowerUnbounded() && i.lowerUnbounded {
		return (i.upper == x.Upper() &&
			i.upperIncluded == x.UpperIncluded() &&
			i.upperUnbounded == x.UpperUnbounded()) || (x.IsEmpty() && i.IsEmpty())
	}
	return (i.lower == x.Lower() &&
		i.upper == x.Upper() &&
		i.lowerIncluded == x.LowerIncluded() &&
		i.upperIncluded == x.UpperIncluded() &&
		i.lowerUnbounded == x.LowerUnbounded() &&
		i.upperUnbounded == x.UpperUnbounded()) || (x.IsEmpty() && i.IsEmpty())
}

// IsEmpty returns true if receiver interval has no value.
func (i *Interval[T]) IsEmpty() bool {
	if i.upperUnbounded || i.lowerUnbounded {
		return false
	}
	if i.lower < i.upper {
		return false
	} else if i.lower == i.upper {
		return !i.lowerIncluded || !i.upperIncluded
	}
	return true
}

// LtBeginOf returns true if receiver interval is less than begin of x_interval_string interval.
func (i *Interval[T]) LtBeginOf(x IInterval[T]) bool {
	if x == nil || x.IsEmpty() {
		return false
	}
	if i.IsEmpty() {
		return false
	}
	if i.upperUnbounded {
		return false
	}
	if x.LowerUnbounded() {
		return false
	}
	if i.upper < x.Lower() {
		return true
	} else if i.upper == x.Lower() {
		return !i.upperIncluded || !x.LowerIncluded()
	}
	return false
}

// LeEndOf returns true if receiver interval is less than or equal to end of x_interval_string interval.
func (i *Interval[T]) LeEndOf(x IInterval[T]) bool {
	if x == nil || x.IsEmpty() {
		return false
	}
	if i.IsEmpty() {
		return false
	}
	if i.upperUnbounded {
		return false
	} else if x.UpperUnbounded() {
		return true
	}
	if i.upper < x.Upper() {
		return true
	}
	if i.upper == x.Upper() {
		return (i.upperIncluded && x.UpperIncluded()) || !i.upperIncluded
	}
	return false
}

// Contains returns true if x_interval_string interval is completely covered by receiver interval.
func (i *Interval[T]) Contains(x IInterval[T]) bool {
	if x == nil || x.IsEmpty() {
		return true
	}
	if i.IsEmpty() {
		return false
	}
	lowerSide := false
	upperSide := false
	if x.LowerUnbounded() && !i.lowerUnbounded {
		lowerSide = false
	} else {
		if i.lowerUnbounded {
			lowerSide = true
		}
		if i.lower < x.Lower() {
			lowerSide = true
		}
		if i.lower == x.Lower() && !x.LowerIncluded() {
			lowerSide = true
		}
		if i.lower == x.Lower() && i.lowerIncluded {
			lowerSide = true
		}
	}
	if x.UpperUnbounded() && !i.upperUnbounded {
		upperSide = false
	} else {
		if i.upperUnbounded {
			upperSide = true
		}
		if i.upper > x.Upper() {
			upperSide = true
		}
		if i.upper == x.Upper() && !x.UpperIncluded() {
			upperSide = true
		}
		if i.upper == x.Upper() && i.upperIncluded {
			upperSide = true
		}
	}
	return lowerSide && upperSide
}

func (i *Interval[T]) Has(value T) bool {
	if i.lowerUnbounded && i.upperUnbounded {
		return true
	}
	returnValue := true
	if !i.lowerUnbounded && i.upperUnbounded {
		if i.lowerIncluded {
			returnValue = value >= i.lower
		} else {
			returnValue = value > i.lower
		}
	} else if !i.upperUnbounded && i.lowerUnbounded {
		if i.upperIncluded {
			returnValue = value <= i.upper
		} else {
			returnValue = value < i.upper
		}
	} else if !i.lowerUnbounded && !i.upperUnbounded {
		if i.lowerIncluded && i.upperIncluded {
			returnValue = value >= i.lower && value <= i.upper
		} else if !i.lowerIncluded && i.upperIncluded {
			returnValue = value > i.lower && value <= i.upper
		} else if i.lowerIncluded && !i.upperIncluded {
			returnValue = value >= i.lower && value < i.upper
		} else if !i.lowerIncluded && !i.upperIncluded {
			returnValue = value > i.lower && value < i.upper
		}
	}
	return returnValue
}

// Intersect returns the intersection of receiver interval with x_interval_string interval.
func (i *Interval[T]) Intersect(x IInterval[T]) IInterval[T] {
	if x == nil || x.IsEmpty() || i.IsEmpty() {
		return nil
	}
	r := NewInterval[T](
		x.Lower(),
		x.Upper(),
		x.LowerIncluded(),
		x.LowerUnbounded(),
		x.UpperIncluded(),
		x.UpperUnbounded(),
	)
	if !i.lowerUnbounded && !x.LowerUnbounded() {
		if i.lower > x.Lower() {
			r.SetLower(i.lower)
			r.SetLowerIncluded(i.lowerIncluded)
		} else if i.lower == x.Lower() && !i.lowerIncluded {
			r.SetLowerIncluded(false)
		}
	} else if x.LowerUnbounded() && !i.lowerUnbounded {
		r.SetLower(i.lower)
		r.SetLowerIncluded(i.lowerIncluded)
		r.SetLowerUnbounded(false)
	}
	if !i.upperUnbounded && !x.UpperUnbounded() {
		if i.upper < x.Upper() {
			r.SetUpper(i.upper)
			r.SetUpperIncluded(i.upperIncluded)
		} else if i.upper == x.Upper() && !i.upperIncluded {
			r.SetUpperIncluded(false)
		}
	} else if x.UpperUnbounded() && !i.upperUnbounded {
		r.SetUpper(i.upper)
		r.SetUpperIncluded(i.upperIncluded)
		r.SetUpperUnbounded(false)
	}
	return maybeEmpty(r)
}

func maybeEmpty[T constraints.Integer | constraints.Float](x IInterval[T]) IInterval[T] {
	if x == nil || x.IsEmpty() {
		return nil
	}
	return x
}

// Move returns an interval that adds number x_interval_string to begin and end of receiver interval.
func (i *Interval[T]) Move(x T) IInterval[T] {
	if i.IsEmpty() {
		return nil
	}
	return NewInterval[T](
		i.lower+x,
		i.upper+x,
		i.lowerIncluded,
		i.lowerUnbounded,
		i.upperIncluded,
		i.upperUnbounded,
	)
}

// Subtract returns two intervals, one on the before of x_interval_string and one on the
// after of x_interval_string, corresponding to the subtraction of x_interval_string from the receiver
// interval. The returned intervals are always within the range of the
// receiver interval.
func (i *Interval[T]) Subtract(x IInterval[T]) (IInterval[T], IInterval[T]) {
	in := i.Intersect(x)
	if in == nil || in.IsEmpty() {
		if i.LtBeginOf(x) {
			return i, nil
		}
		return nil, i
	}
	var r1, r2 IInterval[T]
	if x.LowerUnbounded() {
		r1 = nil
	} else {
		lower := i.lower
		if i.lower > in.Lower() {
			lower = in.Lower()
		}
		r1 = maybeEmpty(NewInterval[T](lower, in.Lower(), i.lowerIncluded, i.lowerUnbounded, !in.LowerIncluded(), false))
	}
	if x.UpperUnbounded() {
		r2 = nil
	} else {
		r2 = maybeEmpty(NewInterval[T](in.Upper(), i.upper, !in.UpperIncluded(), false, i.upperIncluded, i.upperUnbounded))
		if r2 != nil {
			if r2.Lower() > r2.Upper() && i.upperUnbounded {
				r2.SetUpper(r2.Lower())
			}
		}
	}
	return r1, r2
}

// Adjoin returns the union of two intervals, if the intervals are exactly
// adjacent, or the zero interval if they are not.
// Monica says: The Adjoin keyword in Swift is used to combine two ranges, specifically closed ranges, into a single range. It operates on ranges of any type that conforms to the Comparable protocol, allowing for flexibility in combining different types of ranges.
//
// The Adjoin keyword takes two closed ranges as input and returns a new closed range that represents the union of the two input ranges.
// If the input ranges are exactly adjacent, meaning they share a common boundary point,
// the Adjoin operation successfully combines them into a single range.
// However, if the input ranges are not adjacent, the Adjoin operation results in an empty range,
// represented by an empty closed range with the same lower and upper bounds.
func (i *Interval[T]) Adjoin(x IInterval[T]) IInterval[T] {
	if x == nil || x.IsEmpty() || i.IsEmpty() {
		return nil
	}
	if i.lowerUnbounded || i.upperUnbounded || x.UpperUnbounded() || x.LowerUnbounded() {
		return nil
	}
	if i.lower == x.Upper() && (i.lowerIncluded || x.UpperIncluded()) {
		x.SetUpper(i.upper)
		x.SetUpperIncluded(i.upperIncluded)
		return x
	}
	if i.upper == x.Lower() && (i.upperIncluded || x.LowerIncluded()) {
		x.SetLower(i.lower)
		x.SetLowerIncluded(i.lowerIncluded)
		return x
	}
	return nil
}

// Encompass returns an interval that covers the exact extents of two intervals.
func (i *Interval[T]) Encompass(x IInterval[T]) IInterval[T] {
	if x == nil || x.IsEmpty() {
		return i
	}
	if i.IsEmpty() {
		return x
	}
	if i.lower < x.Lower() {
		x.SetLower(i.lower)
		x.SetLowerIncluded(i.lowerIncluded)
	} else if i.lower == x.Lower() && i.lowerIncluded {
		x.SetLowerIncluded(true)
	}
	if i.upper > x.Upper() {
		x.SetUpper(i.upper)
		x.SetUpperIncluded(i.upperIncluded)
	} else if i.upper == x.Upper() && i.upperIncluded {
		x.SetUpperIncluded(true)
	}
	if i.lowerUnbounded {
		x.SetLowerUnbounded(true)
	}
	if i.upperUnbounded {
		x.SetUpperUnbounded(true)
	}
	return x
}

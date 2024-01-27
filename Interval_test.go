package interval

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
	"testing"
)

/*
Reading the test-sets
The '========' stands for interval-range
Upper- and lower-Included is true
The '<' is for lowerUnbound, no lower-end
The '>' is for upperUnbound, no upper-end
The '*' is for Included is false (depending on the side of the string for upper en lower included)

The '$' in the middlepart makes lower and upper equal

The interval-string is divided in three parts by |,(or is empty)
First part indicates the lower-unbounded/included, the third part for the right side
It is allowed to have as many characters as is convenient on the first and the last part,
only the indicated characters are meaningful

It looks hard to understand and is maybe confusing, but read on,
So:
*|=================|  lowerIncluded is false, and lower  is 0
|=================|   lowerIncluded is true, and lower is 0
<|================| lowerUnbounded is true and lower is 0
<*|==============| lowerUnbounded is true, and lowerIncluded is false and lower is 0
*|---=================|  lowerIncluded is false, and lower  is 4
|---=================|   lowerIncluded is true, and lower is 4
|---<================| lowerUnbounded is true and lower is 4
<*|---==============| lowerUnbounded is true, and lowerIncluded is false and lower is 4

|---&---|   lowerIncluded and upperIncluded are true, and lower and upper are 4

# On the upperside are mutatis mutandis the same rules

So, these two intervals below, although visible overlapping, do NOT overlap
"<*|====--------------|" 0,4
"*|----=======|*"       5,12
So it can be written like this, which is much more readable
"<*|====--------------| 0,4"
"* |----=======|*"       5,12 or
" *|----=======|*"

Remember, this is only for creating tests-sets, it has nothing to do with a notation language of meaning outside these test-sets
*/
func parseInterval[T constraints.Integer | constraints.Float](s string) (*Interval[T], error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, "|")
	if len(parts) != 3 {
		if len(parts) == 1 {
			begin := strings.IndexAny(s, "*=")
			end := strings.LastIndexAny(s, "*=")
			return NewInterval[T](T(begin), T(end), s[begin] == '=', false, s[end] == '=', false), nil
		}
		return nil, errors.New(fmt.Sprintf("The interval string '%s' is not wellformed, it must have 2 '|' (pipes).", s))
	}
	leftside := parts[0]
	rightside := parts[2]
	interval := parts[1]
	if strings.ContainsAny(interval, "*<>") {
		return nil, errors.New(fmt.Sprintf("The interval string '%s' is not wellformed, it has not allowed characters in the middlepart", s))
	}
	begin := strings.Index(interval, "=")
	end := strings.LastIndex(interval, "=") + 1
	beginendequal := strings.Index(interval, "&")
	if begin == -1 && end == 0 {
		begin = beginendequal
		end = beginendequal
	}
	lowerunbounded, lowerincluded, upperunbounded, upperincluded := false, true, false, true
	if len(leftside) > 0 {
		lowerunbounded = strings.Contains(leftside, "<")
		lowerincluded = !strings.Contains(leftside, "*")
	}
	if len(rightside) > 0 {
		upperunbounded = strings.Contains(rightside, ">")
		upperincluded = !strings.Contains(rightside, "*")
	}
	r := NewInterval(T(begin), T(end), lowerincluded, lowerunbounded, upperincluded, upperunbounded)
	return r, nil
}

func TestIntervalHas(t *testing.T) {
	testIntervalHas[int](t)
	testIntervalHas[float64](t)
}

func TestParseInterval(t *testing.T) {
	testParseInterval[int](t)
	testParseInterval[float64](t)
}

func TestIntervalLtBeginOf(t *testing.T) {
	testIntervalLtBeginOf[int](t)
	testIntervalLtBeginOf[float64](t)
}

func TestIntervalLeEndOf(t *testing.T) {
	testIntervalLeEndOf[int](t)
	testIntervalLeEndOf[float64](t)
}

func TestIntervalContains(t *testing.T) {
	testIntervalContains[int](t)
	testIntervalContains[float64](t)
}

func TestIntervalIntersect(t *testing.T) {
	testIntervalIntersect[int](t)
	testIntervalIntersect[float64](t)
}

func TestIntervalSubtract(t *testing.T) {
	testIntervalSubtract[int](t)
	testIntervalSubtract[float64](t)
}

func TestIntervalAdjoin(t *testing.T) {
	testIntervalAdjoin[int](t)
	testIntervalAdjoin[float64](t)
}

func TestIntervalEncompass(t *testing.T) {
	testIntervalEncompass[int](t)
	testIntervalEncompass[float64](t)
}

func testParseInterval[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsParseInterval {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.s)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			if i.Lower() != T(tc.begin) ||
				i.Upper() != T(tc.end) ||
				i.LowerIncluded() != tc.lowerIncluded ||
				i.LowerUnbounded() != tc.lowerUnbounded ||
				i.UpperIncluded() != tc.upperIncluded ||
				i.UpperUnbounded() != tc.upperUnbounded {
				t.Errorf("String s: %s want tc = %v but get %v", tc.s, tc, i)
			}
		})
	}
}

func testIntervalLtBeginOf[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalLtBeginOf {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
				return
			}
			a, b := i.LtBeginOf(x), x.LtBeginOf(i)
			if a != tc.i_Before_x {
				t.Errorf("want %s.LtBeginOf(%s) = %v (in test) but is %v, counter: %v\n%s\n%s",
					i, x, tc.i_Before_x, a, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if b != tc.x_Before_i {
				t.Errorf("want %s.LtBeginOf(%s) = %v (in test) but is %v, counter: %v\n%s\n%s",
					x, i, tc.x_Before_i, b, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
		})
	}
}

func testIntervalLeEndOf[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalLeEndOf {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			a, b := i.LeEndOf(x), x.LeEndOf(i)
			if a != tc.i_LeEnd_x {
				t.Errorf("want %s.LeEndOf(%s) = %v (in test) but is %v, counter: %v\n%s\n%s",
					i, x, tc.i_LeEnd_x, a, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if b != tc.x_LeEnd_i {
				t.Errorf("want %s.LeEndOf(%s) = %v (in test) but is %v, counter: %v\n%s\n%s",
					x, i, tc.x_LeEnd_i, b, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
		})
	}
}

func testIntervalContains[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalContains {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			c, d := i.Contains(x), x.Contains(i)
			if c != tc.i_Cover_x {
				t.Errorf("want %s.Contains(%s) = %v (in test) but is %v, counter: %v\n%s\n%s",
					i, x, tc.i_Cover_x, c, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if d != tc.x_Cover_i {
				t.Errorf("want %s.Contains(%s) = %v (in test) but is %v, counter: %v\n%s\n%s",
					x, i, tc.x_Cover_i, d, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
		})
	}
}

func testIntervalIntersect[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalIntersect {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}

			e := i.Intersect(x)
			we, er := parseInterval[T](tc.i_intersect_x)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			if e == nil && we == nil {
				return
			}
			if e == nil && we != nil {
				t.Errorf("want %s.Intersect(%s) = %s, (%s) (result conform test) but is actually %s, counter: %v-a\n%s\n%s",
					i, x, we, tc.i_intersect_x, "nil", tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if e != nil && we == nil {
				t.Errorf("want %s.Intersect(%s) = %s, (%s) (result conform test) but is actually %s, counter: %v-a\n%s\n%s",
					i, x, "nil", tc.i_intersect_x, e, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if !e.Equal(we) {
				t.Errorf("want %s.Intersect(%s) = %s, (%s) (result conform test) but is actually %s, counter: %v-a\n%s\n%s",
					i, x, we, tc.i_intersect_x, e, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
		})
	}
}

func testIntervalAdjoin[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalAdjoin {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}

			l := i.Adjoin(x)
			wl, er := parseInterval[T](tc.i_Adjoin_x)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			if l == nil && wl == nil {
				return
			}
			if l == nil && wl != nil {
				t.Errorf("\nwant %s.Adjoin(%s) = %s (result conform test)\n but is actually %s, counter: %v\n%s\n%s",
					i, x, wl, "nil", tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
			if l != nil && wl == nil {
				t.Errorf("\nwant %s.Adjoin(%s) = %s (result conform test)\n but is actually %s, counter: %v\n%s\n%s",
					i, x, "nil", l, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
			if !l.Equal(wl) {
				t.Errorf("want %s.Adjoin(%s) = %s but get %s", i, x, wl, l)
				t.Errorf("\nwant %s.Adjoin(%s) = %s (result conform test)\n but is actually %s, counter: %v\n%s\n%s",
					i, x, wl, l, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
		})
	}
}

func testIntervalSubtract[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalSubtract {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			g, h := i.Subtract(x)
			wg, er := parseInterval[T](tc.i_Subtract_x_before)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			wh, er := parseInterval[T](tc.i_Subtract_x_after)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			if g == nil && wg == nil || h == nil && wh == nil {
				return
			}
			if g == nil && wg != nil || h == nil && wh != nil {
				if g == nil && h != nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, wg, wh, tc.i_Subtract_x_before, tc.i_Subtract_x_after, "nil", h, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				} else if g != nil && h == nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, wg, wh, tc.i_Subtract_x_before, tc.i_Subtract_x_after, g, "nil", tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				} else if g == nil && h == nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, wg, wh, tc.i_Subtract_x_before, tc.i_Subtract_x_after, "nil", "nil", tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				}
			}
			if g != nil && wg == nil || h != nil && wh == nil {
				if wg == nil && wh != nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, "nil", wh, tc.i_Subtract_x_before, tc.i_Subtract_x_after, g, h, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				} else if wg != nil && wh == nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, wg, "nil", tc.i_Subtract_x_before, tc.i_Subtract_x_after, g, h, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				} else if wg == nil && wh == nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, "nil", "nil", tc.i_Subtract_x_before, tc.i_Subtract_x_after, g, h, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				}
				t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
					i, x, "nil", "nil", tc.i_Subtract_x_before, tc.i_Subtract_x_after, g, h, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if !g.Equal(wg) || !h.Equal(wh) {
				t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
					i, x, wg, wh, tc.i_Subtract_x_before, tc.i_Subtract_x_after, g, h, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			j, k := x.Subtract(i)
			wj, er := parseInterval[T](tc.x_Subtract_i_before)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			wk, er := parseInterval[T](tc.x_Subtract_i_after)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			if j == nil && k == nil && wj == nil && wk == nil {
				return
			}
			if j != nil && wj == nil || k != nil && wk == nil {
				if wj == nil && wk != nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, "nil", wk, tc.i_Subtract_x_before, tc.i_Subtract_x_after, j, k, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				} else if wj != nil && wk == nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, wj, "nil", tc.i_Subtract_x_before, tc.i_Subtract_x_after, j, k, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				} else if wj == nil && wk == nil {
					t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
						i, x, "nil", "nil", tc.i_Subtract_x_before, tc.i_Subtract_x_after, j, k, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				}
				t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
					i, x, "nil", "nil", tc.i_Subtract_x_before, tc.i_Subtract_x_after, j, k, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if !j.Equal(wj) || !k.Equal(wk) {
				t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-a\n%s\n%s",
					i, x, wj, wk, tc.i_Subtract_x_before, tc.i_Subtract_x_after, j, k, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
				return
			}
			if !j.Equal(wj) || !k.Equal(wk) {
				t.Errorf("\nwant %s.Subtract(%s) = %s, %s\n%s,%s (result conform test)\n but is actually %s, %s, counter: %v-b\n%s\n%s",
					x, i, wj, wk, tc.x_Subtract_i_before, tc.x_Subtract_i_after, j, k, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
		})
	}
}

func testIntervalEncompass[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsIntervalEncompass {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.test.i_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			x, er := parseInterval[T](tc.test.x_interval_string)
			if er != nil {
				t.Errorf(er.Error())
				return
			}

			o := i.Encompass(x)
			wo, er := parseInterval[T](tc.i_Encompass_x)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			if o == nil && wo == nil {
				return
			}
			if o == nil && wo == nil {
				return
			}
			if o == nil && wo != nil {
				t.Errorf("\nwant %s.Encompass(%s) = %s (result conform test)\n but is actually %s, counter: %v-a\n%s\n%s",
					i, x, wo, "nil", tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
			if o != nil && wo == nil {
				t.Errorf("\nwant %s.Encompass(%s) = %s (result conform test)\n but is actually %s, counter: %v-a\n%s\n%s",
					i, x, "nil", o, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
			if !o.Equal(wo) {
				t.Errorf("\nwant %s.Encompass(%s) = %s (result conform test)\n but is actually %s, counter: %v-a\n%s\n%s",
					i, x, wo, o, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
				return
			}
		})
	}
}

func testIntervalHas[T constraints.Integer | constraints.Float](t *testing.T) {
	for n, tc := range testsHAS {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i, er := parseInterval[T](tc.s)
			if er != nil {
				t.Errorf(er.Error())
				return
			}
			o := i.Has(T(tc.value))
			if o != tc.result {
				t.Errorf("\nwant %s.Has(%v) = %v (result conform test)\n but is actually %v, counter: %v",
					i, tc.value, tc.result, o, tc.counter)
				return
			}
		})
	}
}

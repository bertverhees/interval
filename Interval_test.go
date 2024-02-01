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
			if we == nil {
				if ee == nil {
					return
				}else{
					t.Errorf("want %s.Intersect(%s) = %s, (%s) (result conform test) but is actually %s, counter: %v-a\n%s\n%s",
						i, x, "nil", tc.i_intersect_x, e, tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				}
			}else{
				if ee == nil {
					t.Errorf("want %s.Intersect(%s) = %s, (%s) (result conform test) but is actually %s, counter: %v-a\n%s\n%s",
						i, x, we, tc.i_intersect_x, "nil", tc.test.counter, tc.test.i_interval_string, tc.test.x_interval_string)
					return
				}
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
			if wl == nil {
				if l == nil {
					return
				}else{
					t.Errorf("\nwant %s.Adjoin(%s) = %s (result conform test)\n but is actually %s, counter: %v\n%s\n%s",
						i, x, "nil", l, tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
					return
				}
			}else {
				if l == nil{
					t.Errorf("\nwant %s.Adjoin(%s) = %s (result conform test)\n but is actually %s, counter: %v\n%s\n%s",
						i, x, wl, "nil", tc.test.counter, tc.test.x_interval_string, tc.test.i_interval_string)
					return
				}
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

/*
Reading the test-sets
The '========' stands for interval-range
Upper- en lower-Included is true
The '<' is for lowerUnbound, no lower-end
The '>' is for upperUnbound, no upper-end
The '*' is for Included is false (depending on the side of the string for upper en lower included)

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

# On the upperside are mutatis mutandis the same rules

So, these two intervals below, although visible overlapping, do NOT overlap
<*|====--------------| 0,4
*|----=======|*       5,12
So it can be written like this, which is much more readable
So it can be written like this, which is much more readable
"<*|====--------------| 0,4"
"* |----=======|*"       5,12 or
" *|----=======|*"
Remember, this is only for creating tests-sets, it has nothing to do with a notation language of meaning outside these test-sets
*/

type testGeneral struct {
	i_interval_string string
	x_interval_string string
	counter           string
}

var testsGeneralSets = []testGeneral{
	{ //0
		i_interval_string: "  |=====---------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "0",
	},
	{ //1
		i_interval_string: "  |======--------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "1",
	},
	{ //2
		i_interval_string: "  |---=====------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "2",
	},
	{ //3
		i_interval_string: "  |-------=====--------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "3",
	},
	{ //4
		i_interval_string: "  |------========------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "4",
	},
	{ //5
		i_interval_string: "  |-----------=====----|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "5",
	},
	{ //6
		i_interval_string: "  |--------------======|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "6",
	},
	{ //7
		i_interval_string: "  |---------------=====|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "7",
	},
	{ //8
		i_interval_string: " *|======--------------|  ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "8",
	},
	{ //9
		i_interval_string: " *|------========------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "9",
	},
	{ //10
		i_interval_string: " *|--------------======|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "10",
	},
	{ //11
		i_interval_string: "  |======--------------|* ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "11",
	},
	{ //12
		i_interval_string: "  |------========------|* ",
		x_interval_string: "  |------========------|  ",
		counter:           "12",
	},
	{ //13
		i_interval_string: "  |--------------======|* ",
		x_interval_string: "  |------========------|  ",
		counter:           "13",
	},
	{ //14
		i_interval_string: " <|-------=====--------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "14",
	},
	{ //15
		i_interval_string: " <|------========------|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "15",
	},
	{ //16
		i_interval_string: " <|-----------=====----|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "16",
	},
	{ //17
		i_interval_string: " <|--------------======|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "17",
	},
	{ //18
		i_interval_string: " <|---------------=====|  ",
		x_interval_string: "  |------========------|  ",
		counter:           "18",
	},
	{ //19
		i_interval_string: "  |=====---------------|> ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "19",
	},
	{ //20
		i_interval_string: "  |======--------------|> ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "20",
	},
	{ //21
		i_interval_string: "  |---=====------------|> ",
		x_interval_string: "  |------=======-------|  ",
		counter:           "21",
	},
	{ //22
		i_interval_string: "  |-------=====--------|> ",
		x_interval_string: "  |------========------|  ",
		counter:           "22",
	},
	{ //23
		i_interval_string: "  |------========------|> ",
		x_interval_string: "  |------========------|  ",
		counter:           "23",
	},
	{ //24
		i_interval_string: "  |======--------------|  ",
		x_interval_string: " *|------=======-------|  ",
		counter:           "24",
	},
	{ //25
		i_interval_string: "  |------========------|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "25",
	},
	{ //26
		i_interval_string: "  |--------------======|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "26",
	},
	{ //27
		i_interval_string: "  |======--------------|  ",
		x_interval_string: "  |------=======-------|*  ",
		counter:           "27",
	},
	{ //28
		i_interval_string: "  |------========------|  ",
		x_interval_string: "  |------========------|* ",
		counter:           "28",
	},
	{ //29
		i_interval_string: "  |--------------======|  ",
		x_interval_string: "  |------========------|*  ",
		counter:           "29",
	},
	{ //30
		i_interval_string: " *|======--------------|  ",
		x_interval_string: " *|------=======-------|  ",
		counter:           "30",
	},
	{ //31
		i_interval_string: " *|------========------|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "31",
	},
	{ //32
		i_interval_string: " *|--------------======|  ",
		x_interval_string: " *|------========------|  ",
		counter:           "32",
	},
	{ //33
		i_interval_string: "  |======--------------|* ",
		x_interval_string: "  |------=======-------|* ",
		counter:           "33",
	},
	{ //34
		i_interval_string: "  |------========------|* ",
		x_interval_string: "  |------========------|* ",
		counter:           "34",
	},
	{ //35
		i_interval_string: "  |--------------======|* ",
		x_interval_string: "  |------========------|* ",
		counter:           "35",
	},
	{ //36
		i_interval_string: "  |=====---------------|  ",
		x_interval_string: " <|------=======-------|  ",
		counter:           "36",
	},
	{ //37
		i_interval_string: "  |======--------------|  ",
		x_interval_string: " <|------=======-------|  ",
		counter:           "37",
	},
	{ //38
		i_interval_string: "  |---=====------------|  ",
		x_interval_string: " <|------=======-------|  ",
		counter:           "38",
	},
	{ //39
		i_interval_string: "  |-------=====--------|  ",
		x_interval_string: " <|------========------|  ",
		counter:           "39",
	},
	{ //40
		i_interval_string: "  |------========------|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "40",
	},
	{ //41
		i_interval_string: "  |-----------=====----|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "41",
	},
	{ //42
		i_interval_string: "  |--------------======|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "42",
	},
	{ //43
		i_interval_string: "  |---------------=====|  ",
		x_interval_string: "  |------========------|> ",
		counter:           "43",
	},
	{ //44
		i_interval_string: " <|-------=====--------|  ",
		x_interval_string: " <|------========------|  ",
		counter:           "44",
	},
	{ //45
		i_interval_string: " <|------========------|  ",
		x_interval_string: " <|------========------|  ",
		counter:           "45",
	},
	{ //46
		i_interval_string: "  |-------=====--------|> ",
		x_interval_string: "  |------========------|> ",
		counter:           "46",
	},
	{ //47
		i_interval_string: " <|------========------|> ",
		x_interval_string: " <|------========------|> ",
		counter:           "47",
	},
	{ //48
		i_interval_string: " <|--------====--------|> ",
		x_interval_string: " <|------========------|> ",
		counter:           "48",
	},
}

var testsParseInterval = []struct {
	s              string
	begin          int
	lowerIncluded  bool
	lowerUnbounded bool
	end            int
	upperIncluded  bool
	upperUnbounded bool
}{
	{
		s:              "  |=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " |=====|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	//-------------
	{
		s:              "  |--=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	//-------------
	{
		s:              "  |--=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              " *|--=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====--|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              " <|--=====--|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "  |--=====--|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          2,
		end:            7,
	},
	{
		s:              "<*|--=====--|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          2,
		end:            7,
	},
	//-------------
	{
		s:              "  |=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====--|",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " *|=====--|*",
		lowerUnbounded: false,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====--|>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              " <|=====--|>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "  |=====--|*>",
		lowerUnbounded: false,
		upperUnbounded: true,
		lowerIncluded:  true,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====--|",
		lowerUnbounded: true,
		upperUnbounded: false,
		lowerIncluded:  false,
		upperIncluded:  true,
		begin:          0,
		end:            5,
	},
	{
		s:              "<*|=====--|*>",
		lowerUnbounded: true,
		upperUnbounded: true,
		lowerIncluded:  false,
		upperIncluded:  false,
		begin:          0,
		end:            5,
	},
}

var testsIntervalSubtract = []struct {
	//i
	test testGeneral
	//g,h
	// i_interval_string.Subtract(x_interval_string)
	i_Subtract_x_before, i_Subtract_x_after string
	//j,k
	// x_interval_string.Subtract(i_interval_string)
	x_Subtract_i_before, x_Subtract_i_after string
}{
	{
		test:                testsGeneralSets[0],
		i_Subtract_x_before: "  |=====|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "|------=======-------|",
	},
	{
		test:                testsGeneralSets[1],
		i_Subtract_x_before: "  |======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|",
	},
	{
		test:                testsGeneralSets[2],
		i_Subtract_x_before: "  |---===|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|--------=====-----|",
	},
	{
		test:                testsGeneralSets[3],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "  |------=|*",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[4],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[5],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------==|",
		x_Subtract_i_before: "  |------=====|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[6],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "  |------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[7],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "  |---------------=====|",
		x_Subtract_i_before: "  |------========------|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[8],
		i_Subtract_x_before: " *|======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|",
	},
	{
		test:                testsGeneralSets[9],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "  |------&-|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[10],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "  |------========------|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[11],
		i_Subtract_x_before: "  |======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "  |------=======---|",
	},
	{
		test:                testsGeneralSets[12],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "  |--------------&|",
	},
	{
		test:                testsGeneralSets[13],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  " *|--------------======|*",
		x_Subtract_i_before: "  |------========------|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[14],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[15],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[16],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  " *|--------------==----|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[17],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[18],
		i_Subtract_x_before: " <|------&|*",
		i_Subtract_x_after:  " *|--------------======|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[19],
		i_Subtract_x_before: " |======|*",
		i_Subtract_x_after:  "*|-------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[20],
		i_Subtract_x_before: " |======|*",
		i_Subtract_x_after:  "*|-------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[21],
		i_Subtract_x_before: " |---===|*",
		i_Subtract_x_after:  "*|-------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[22],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------&-------|>",
		x_Subtract_i_before: " |------=|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[23],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------&-------|>",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[24],
		i_Subtract_x_before: "  |======|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|",
	},
	{
		test:                testsGeneralSets[25],
		i_Subtract_x_before: "  |------&|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[26],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------======|",
		x_Subtract_i_before: "*|------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[27],
		i_Subtract_x_before: "|======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------=======|*",
	},
	{
		test:                testsGeneralSets[28],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "|--------------&|",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[29],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "|--------------======|",
		x_Subtract_i_before: " |------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[30],
		i_Subtract_x_before: "*|======|",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "*|------=======|",
	},
	{
		test:                testsGeneralSets[31],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[32],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "*|--------------======|",
		x_Subtract_i_before: "*|------========|",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[33],
		i_Subtract_x_before: "  |======|*",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "  |------=======|*",
	},
	{
		test:                testsGeneralSets[34],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[35],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "|--------------======|*",
		x_Subtract_i_before: "|------========|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[36],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "<|&|*",
		x_Subtract_i_after:  " *|-----========|",
	},
	{
		test:                testsGeneralSets[37],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "<|&|*",
		x_Subtract_i_after:  "*|------=======|",
	},
	{
		test:                testsGeneralSets[38],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "<|---&|*",
		x_Subtract_i_after:  "*|--------=====|",
	},
	{
		test:                testsGeneralSets[39],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: " <|------=|*",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[40],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "*|--------------&-------|>",
	},
	{
		test:                testsGeneralSets[41],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: " |------=====|*",
		x_Subtract_i_after:  "*|----------------&|>",
	},
	{
		test:                testsGeneralSets[42],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: " |------========|*",
		x_Subtract_i_after:  "*|--------------------&|>",
	},
	{
		test:                testsGeneralSets[43],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "|------=========------|*",
		x_Subtract_i_after:  "*|--------------------&|>",
	},
	{
		test:                testsGeneralSets[44],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  " *|------------==|",
	},
	{
		test:                testsGeneralSets[45],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[46],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "  |------=|*",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[47],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
	{
		test:                testsGeneralSets[48],
		i_Subtract_x_before: "",
		i_Subtract_x_after:  "",
		x_Subtract_i_before: "",
		x_Subtract_i_after:  "",
	},
}

var testsIntervalAdjoin = []struct {
	test testGeneral
	// i_interval_string.Adjoin(x_interval_string)
	i_Adjoin_x string
}{
	{
		test:       testsGeneralSets[0],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[1],
		i_Adjoin_x: "|=============|",
	},
	{
		test:       testsGeneralSets[2],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[3],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[4],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[5],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[6],
		i_Adjoin_x: "  |------==============|",
	},
	{
		test:       testsGeneralSets[7],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[8],
		i_Adjoin_x: "*|=============|",
	},
	{
		test:       testsGeneralSets[9],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[10],
		i_Adjoin_x: "|------==============|",
	},
	{
		test:       testsGeneralSets[11],
		i_Adjoin_x: "|=============|",
	},
	{
		test:       testsGeneralSets[12],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[13],
		i_Adjoin_x: "|------==============|*",
	},
	{
		test:       testsGeneralSets[14],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[15],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[16],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[17],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[18],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[19],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[20],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[21],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[22],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[23],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[24],
		i_Adjoin_x: "|=============|",
	},
	{
		test:       testsGeneralSets[25],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[26],
		i_Adjoin_x: "*|------==============|",
	},
	{
		test:       testsGeneralSets[27],
		i_Adjoin_x: "|=============|*",
	},
	{
		test:       testsGeneralSets[28],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[29],
		i_Adjoin_x: "|------==============|",
	},
	{
		test:       testsGeneralSets[30],
		i_Adjoin_x: "*|=============|",
	},
	{
		test:       testsGeneralSets[31],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[32],
		i_Adjoin_x: "*|------==============|",
	},
	{
		test:       testsGeneralSets[33],
		i_Adjoin_x: "|=============|*",
	},
	{
		test:       testsGeneralSets[34],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[35],
		i_Adjoin_x: "|------==============|*",
	},
	{
		test:       testsGeneralSets[36],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[37],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[38],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[39],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[40],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[41],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[42],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[43],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[44],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[45],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[46],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[47],
		i_Adjoin_x: "",
	},
	{
		test:       testsGeneralSets[48],
		i_Adjoin_x: "",
	},
}

var testsIntervalContains = []struct {
	test      testGeneral
	i_Cover_x bool
	x_Cover_i bool
}{
	{
		test:      testsGeneralSets[0],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[1],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[2],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[3],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[4],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[5],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[6],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[7],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[8],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[9],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[10],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[11],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[12],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[13],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[14],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[15],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[16],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[17],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[18],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[19],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[20],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[21],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[22],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[23],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[24],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[25],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[26],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[27],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[28],
		i_Cover_x: true,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[29],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[30],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[31],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[32],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[33],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[34],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[35],
		i_Cover_x: false,
		x_Cover_i: false,
	},
	{
		test:      testsGeneralSets[36],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[37],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[38],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[39],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[40],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[41],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[42],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[43],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[44],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[45],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[46],
		i_Cover_x: false,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[47],
		i_Cover_x: true,
		x_Cover_i: true,
	},
	{
		test:      testsGeneralSets[48],
		i_Cover_x: true,
		x_Cover_i: true,
	},
}

var testsIntervalLeEndOf = []struct {
	test      testGeneral
	i_LeEnd_x bool
	x_LeEnd_i bool
}{
	{
		test:      testsGeneralSets[0],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[1],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[2],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[3],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[4],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[5],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[6],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[7],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[8],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[9],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[10],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[11],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[12],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[13],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[14],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[15],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[16],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[17],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[18],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[19],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[20],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[21],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[22],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[23],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[24],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[25],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[26],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[27],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[28],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[29],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[30],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[31],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[32],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[33],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[34],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[35],
		i_LeEnd_x: false,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[36],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[37],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[38],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[39],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[40],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[41],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[42],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[43],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[44],
		i_LeEnd_x: true,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[45],
		i_LeEnd_x: true,
		x_LeEnd_i: true,
	},
	{
		test:      testsGeneralSets[46],
		i_LeEnd_x: false,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[47],
		i_LeEnd_x: false,
		x_LeEnd_i: false,
	},
	{
		test:      testsGeneralSets[48],
		i_LeEnd_x: false,
		x_LeEnd_i: false,
	},
}

var testsIntervalLtBeginOf = []struct {
	test       testGeneral
	i_Before_x bool
	x_Before_i bool
}{
	{
		test:       testsGeneralSets[0],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[1],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[2],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[3],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[4],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[5],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[6],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[7],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[8],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[9],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[10],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[11],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[12],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[13],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[14],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[15],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[16],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[17],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[18],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[19],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[20],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[21],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[22],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[23],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[24],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[25],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[26],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[27],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[28],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[29],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[30],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[31],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[32],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[33],
		i_Before_x: true,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[34],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[35],
		i_Before_x: false,
		x_Before_i: true,
	},
	{
		test:       testsGeneralSets[36],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[37],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[38],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[39],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[40],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[41],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[42],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[43],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[44],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[45],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[46],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[47],
		i_Before_x: false,
		x_Before_i: false,
	},
	{
		test:       testsGeneralSets[48],
		i_Before_x: false,
		x_Before_i: false,
	},
}

var testsIntervalIntersect = []struct {
	test          testGeneral
	i_intersect_x string
}{
	{
		test:          testsGeneralSets[0],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[1],
		i_intersect_x: "|------&------|",
	},
	{
		test:          testsGeneralSets[2],
		i_intersect_x: "|------==------|",
	},
	{
		test:          testsGeneralSets[3],
		i_intersect_x: "|-------=====------|",
	},
	{
		test:          testsGeneralSets[4],
		i_intersect_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[5],
		i_intersect_x: "|-----------===------|",
	},
	{
		test:          testsGeneralSets[6],
		i_intersect_x: "|--------------&------|",
	},
	{
		test:          testsGeneralSets[7],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[8],
		i_intersect_x: "|------&------|",
	},
	{
		test:          testsGeneralSets[9],
		i_intersect_x: "*|------========------|",
	},
	{
		test:          testsGeneralSets[10],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[11],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[12],
		i_intersect_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[13],
		i_intersect_x: "|--------------&|",
	},
	{
		test:          testsGeneralSets[14],
		i_intersect_x: "|------======|",
	},
	{
		test:          testsGeneralSets[15],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[16],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[17],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[18],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[19],
		i_intersect_x: "|------=======|",
	},
	{
		test:          testsGeneralSets[20],
		i_intersect_x: "|------=======|",
	},
	{
		test:          testsGeneralSets[21],
		i_intersect_x: "|------=======|",
	},
	{
		test:          testsGeneralSets[22],
		i_intersect_x: "|-------=======|",
	},
	{
		test:          testsGeneralSets[23],
		i_intersect_x: "|------========|",
	},
	{
		test:          testsGeneralSets[24],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[25],
		i_intersect_x: "*|------========|",
	},
	{
		test:          testsGeneralSets[26],
		i_intersect_x: "|--------------&|",
	},
	{
		test:          testsGeneralSets[27],
		i_intersect_x: "|------&|",
	},
	{
		test:          testsGeneralSets[28],
		i_intersect_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[29],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[30],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[31],
		i_intersect_x: "*|------========|",
	},
	{
		test:          testsGeneralSets[32],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[33],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[34],
		i_intersect_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[35],
		i_intersect_x: "",
	},
	{
		test:          testsGeneralSets[36],
		i_intersect_x: "|=====----|",
	},
	{
		test:          testsGeneralSets[37],
		i_intersect_x: "|======----|",
	},
	{
		test:          testsGeneralSets[38],
		i_intersect_x: "|---=====------------|",
	},
	{
		test:          testsGeneralSets[39],
		i_intersect_x: "|-------=====--------|",
	},
	{
		test:          testsGeneralSets[40],
		i_intersect_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[41],
		i_intersect_x: "|-----------=====----|",
	},
	{
		test:          testsGeneralSets[42],
		i_intersect_x: "|--------------======|",
	},
	{
		test:          testsGeneralSets[43],
		i_intersect_x: "|---------------=====|",
	},
	{
		test:          testsGeneralSets[44],
		i_intersect_x: "<|------======------|",
	},
	{
		test:          testsGeneralSets[45],
		i_intersect_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[46],
		i_intersect_x: "|-------=====--------|>",
	},
	{
		test:          testsGeneralSets[47],
		i_intersect_x: "<|------========------|>",
	},
	{
		test:          testsGeneralSets[48],
		i_intersect_x: "<|------========------|>",
	},
}

var testsIntervalEncompass = []struct {
	// i_interval_string.Encompass(x_interval_string)
	i_Encompass_x string
	test          testGeneral
}{
	{
		test:          testsGeneralSets[0],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[1],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[2],
		i_Encompass_x: "|---==========|",
	},
	{
		test:          testsGeneralSets[3],
		i_Encompass_x: "|------========|",
	},
	{
		test:          testsGeneralSets[4],
		i_Encompass_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[5],
		i_Encompass_x: "|------==========------|",
	},
	{
		test:          testsGeneralSets[6],
		i_Encompass_x: "|------==============------|",
	},
	{
		test:          testsGeneralSets[7],
		i_Encompass_x: "|------==============------|",
	},
	{
		test:          testsGeneralSets[8],
		i_Encompass_x: "*|=============|",
	},
	{
		test:          testsGeneralSets[9],
		i_Encompass_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[10],
		i_Encompass_x: "|------==============------|",
	},
	{
		test:          testsGeneralSets[11],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[12],
		i_Encompass_x: "|------========------|",
	},
	{
		test:          testsGeneralSets[13],
		i_Encompass_x: "|------==============------|*",
	},
	{
		test:          testsGeneralSets[14],
		i_Encompass_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[15],
		i_Encompass_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[16],
		i_Encompass_x: "<|------==========------|",
	},
	{
		test:          testsGeneralSets[17],
		i_Encompass_x: "<|------==============------|",
	},
	{
		test:          testsGeneralSets[18],
		i_Encompass_x: "<|------==============------|",
	},
	{
		test:          testsGeneralSets[19],
		i_Encompass_x: "|=============------|>",
	},
	{
		test:          testsGeneralSets[20],
		i_Encompass_x: "|=============|>", //0 13
	},
	{
		test:          testsGeneralSets[21],
		i_Encompass_x: "|---==========|>",
	},
	{
		test:          testsGeneralSets[22],
		i_Encompass_x: "|------========|>", //6 14
	},
	{
		test:          testsGeneralSets[23],
		i_Encompass_x: "|------========|>",
	},
	{
		test:          testsGeneralSets[24],
		i_Encompass_x: "|=============|",
	},
	{
		test:          testsGeneralSets[25],
		i_Encompass_x: "|------========|",
	},
	{
		test:          testsGeneralSets[26],
		i_Encompass_x: "*|------==============--|", //6 20
	},
	{
		test:          testsGeneralSets[27],
		i_Encompass_x: "|=============|*",
	},
	{
		test:          testsGeneralSets[28],
		i_Encompass_x: "|------========|",
	},
	{
		test:          testsGeneralSets[29],
		i_Encompass_x: "|------==============--|",
	},
	{
		test:          testsGeneralSets[30],
		i_Encompass_x: "*|=============|",
	},
	{
		test:          testsGeneralSets[31],
		i_Encompass_x: "*|------========|",
	},
	{
		test:          testsGeneralSets[32],
		i_Encompass_x: "*|------==============--|",
	},
	{
		test:          testsGeneralSets[33],
		i_Encompass_x: "|=============|*",
	},
	{
		test:          testsGeneralSets[34],
		i_Encompass_x: "|------========|*",
	},
	{
		test:          testsGeneralSets[35],
		i_Encompass_x: "|------==============--|*",
	},
	{
		test:          testsGeneralSets[36],
		i_Encompass_x: "<|=============|",
	},
	{
		test:          testsGeneralSets[37],
		i_Encompass_x: "<|=============|",
	},
	{
		test:          testsGeneralSets[38],
		i_Encompass_x: "<|---==========|",
	},
	{
		test:          testsGeneralSets[39],
		i_Encompass_x: "<|------========|",
	},
	{
		test:          testsGeneralSets[40],
		i_Encompass_x: "|------========|>",
	},
	{
		test:          testsGeneralSets[41],
		i_Encompass_x: "|------==========|>",
	},
	{
		test:          testsGeneralSets[42],
		i_Encompass_x: "|------==============|>",
	},
	{
		test:          testsGeneralSets[43],
		i_Encompass_x: "|------==============|>",
	},
	{
		test:          testsGeneralSets[44],
		i_Encompass_x: "<|------========|",
	},
	{
		test:          testsGeneralSets[45],
		i_Encompass_x: "<|------========------|",
	},
	{
		test:          testsGeneralSets[46],
		i_Encompass_x: "|------========|>",
	},
	{
		test:          testsGeneralSets[47],
		i_Encompass_x: "<|------========------|>",
	},
	{
		test:          testsGeneralSets[48],
		i_Encompass_x: "<|------========|>",
	},
}

var testsHAS = []struct {
	s       string
	value   int
	result  bool
	counter string
}{
	{
		s:       "  |===|",
		value:   3,
		result:  true,
		counter: "0",
	},
	{
		s:       "  |===|",
		value:   0,
		result:  true,
		counter: "1",
	},
	{
		s:       "  |===|",
		value:   -1,
		result:  false,
		counter: "2",
	},
	{
		s:       "  |===|",
		value:   6,
		result:  false,
		counter: "3",
	},
	//----------------------------
	{
		s:       "*|===|",
		value:   3,
		result:  true,
		counter: "4",
	},
	{
		s:       "*|===|",
		value:   0,
		result:  false,
		counter: "5",
	},
	{
		s:       "*|===|",
		value:   -1,
		result:  false,
		counter: "6",
	},
	{
		s:       "*|===|",
		value:   6,
		result:  false,
		counter: "7",
	},
	//------------------------
	{
		s:       "|=====|*",
		value:   5,
		result:  false,
		counter: "8",
	},
	{
		s:       "|=====|*",
		value:   0,
		result:  true,
		counter: "9",
	},
	{
		s:       "|=====|*",
		value:   -1,
		result:  false,
		counter: "10",
	},
	{
		s:       "|=====|*",
		value:   6,
		result:  false,
		counter: "11",
	},
	//------------------------
	{
		s:       "*|=====|*",
		value:   5,
		result:  false,
		counter: "12",
	},
	{
		s:       "*|=====|*",
		value:   0,
		result:  false,
		counter: "13",
	},
	{
		s:       "*|=====|*",
		value:   -1,
		result:  false,
		counter: "14",
	},
	{
		s:       "*|=====|*",
		value:   6,
		result:  false,
		counter: "15",
	},
	//----------------------------
	{
		s:       " <|=====|",
		value:   5,
		result:  true,
		counter: "16",
	},
	{
		s:       " <|=====|",
		value:   0,
		result:  true,
		counter: "17",
	},
	{
		s:       " <|=====|",
		value:   -1,
		result:  true,
		counter: "18",
	},
	{
		s:       " <|=====|",
		value:   6,
		result:  false,
		counter: "19",
	},
	//----------------------------
	{
		s:       " |=====|>",
		value:   5,
		result:  true,
		counter: "20",
	},
	{
		s:       " |=====|>",
		value:   0,
		result:  true,
		counter: "21",
	},
	{
		s:       " |=====|>",
		value:   -1,
		result:  false,
		counter: "22",
	},
	{
		s:       " |=====|>",
		value:   6,
		result:  true,
		counter: "23",
	},
	//----------------------------
	{
		s:       " <|=====|>",
		value:   5,
		result:  true,
		counter: "24",
	},
	{
		s:       " <|=====|>",
		value:   0,
		result:  true,
		counter: "25",
	},
	{
		s:       " <|=====|>",
		value:   -1,
		result:  true,
		counter: "26",
	},
	{
		s:       " <|=====|>",
		value:   6,
		result:  true,
		counter: "27",
	},
	//----------------------------
	{
		s:       "  |=====|*>",
		value:   5,
		result:  true,
		counter: "28",
	},
	{
		s:       "  |=====|*>",
		value:   0,
		result:  true,
		counter: "29",
	},
	{
		s:       "  |=====|*>",
		value:   -1,
		result:  false,
		counter: "30",
	},
	{
		s:       "  |=====|*>",
		value:   6,
		result:  true,
		counter: "31",
	},
	//----------------------------
	{
		s:       "<*|=====|",
		value:   5,
		result:  true,
		counter: "32",
	},
	{
		s:       "<*|=====|",
		value:   0,
		result:  true,
		counter: "33",
	},
	{
		s:       "<*|=====|",
		value:   -1,
		result:  true,
		counter: "34",
	},
	{
		s:       "<*|=====|",
		value:   6,
		result:  false,
		counter: "35",
	},
	//----------------------------
	{
		s:       "<*|=====|*>",
		value:   5,
		result:  true,
		counter: "36",
	},
	{
		s:       "<*|=====|*>",
		value:   0,
		result:  true,
		counter: "37",
	},
	{
		s:       "<*|=====|*>",
		value:   -1,
		result:  true,
		counter: "38",
	},
	{
		s:       "<*|=====|*>",
		value:   6,
		result:  true,
		counter: "39",
	},
	//----------------------------
}

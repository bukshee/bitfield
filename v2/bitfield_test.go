package bitfield_test

import (
	"fmt"
	"testing"
	. "github.com/bukshee/bitfield/v2"
)

func assert(t *testing.T, a, b interface{}) {
	if a == b {
		return
	}
	t.Helper()
	t.Errorf("%s != %s", a, b)
}

// test that function call do panic
func doesPanic(f func()) bool {
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		f()
	}()
	return didPanic
}

func Test1(t *testing.T) {
	if !doesPanic(func() { NewBitField(-3) }) {
		t.Error("should panic")
	}

	if !doesPanic(func() { New(4).Resize(-1) }) {
		t.Error("should panic")
	}

	a := NewBitField(0)
	assert(t, a.Len(), 0)

	a = New(65).SetAll()
	assert(t, a.OnesCount(), 65)

	assert(t, New(3).Equal(New(4)), false)

	New(3).Set()
	a = New(3).Set(0, -1).Not()
	assert(t, a.String(), "010")

	assert(t, a.Equal(New(3).Set(1)), true)

	a = New(129).Set(0, -1).Clear(123, -3).Not().Not()
	assert(t, a.OnesCount(), 2)

	if !a.Get(0) || !a.Get(-a.Len()) || !a.Get(a.Len()) {
		t.Error("should be true")
	}

	b := New(129).Set(4, -1, 23, 11)
	c := New(129).Set(0, -1)
	assert(t, b.And(c).OnesCount(), 1)

	if !doesPanic(func() { New(5).And(New(121)) }) {
		t.Error("should panic")
	}
	if !doesPanic(func() { New(5).Or(New(121)) }) {
		t.Error("should panic")
	}
	if !doesPanic(func() { New(5).Xor(New(121)) }) {
		t.Error("should panic")
	}

	c = New(129).Set(73, -2).ClearAll().Set(-1)
	d := New(129).Set(128)
	assert(t, c.Equal(d), true)

	d = New(129).Not()
	assert(t, c.Equal(d), false)
	assert(t, c.Get(127), false)
	assert(t, c.Get(-1), true)

	a = New(4).Flip(-1)
	assert(t, a.String(), "0001")

	a = New(4).Not().Flip(-1)
	assert(t, a.String(), "1110")

	a = New(4).SetAll().Flip(0, -1).Mid(1, 2)
	assert(t, a.String(), "11")

	d = New(65).Set(-1)
	e := d.Clone()
	assert(t, d.Equal(e), true)

	assert(t, d.Xor(d).OnesCount(), 0)

	assert(t, d.Set(11).Or(e).Get(11), true)
}

func TestResize(t *testing.T) {
	if New(4).Resize(0).Len() != 0 {
		t.Error("should be 0")
	}
	a := New(3)
	a.Resize(3).Set(1)
	assert(t, a.String(), "000")

	a.Mut().Resize(3).Set(1)
	assert(t, a.String(), "010")

	a = New(1).SetAll().Resize(4)
	assert(t, a.String(), "1000")

	a = New(10).SetAll().Resize(6)
	assert(t, a.String(), "111111")

	a = New(27).Set(0, -1).Resize(65)
	assert(t, a.Get(26), true)

	a = New(65).Set(-1).Right(5)
	assert(t, a.String(), "00001")

	a = New(65).SetAll().Resize(40)
	assert(t, a.OnesCount(), 40)

	a = New(3).SetAll()
	a.Resize(4)
	assert(t, a.String(), "111")

	a.Mut().Resize(4)
	assert(t, a.String(), "1110")
}

func TestShift(t *testing.T) {
	tests := []struct {
		size     int
		shift    int
		expected int
	}{
		{100, 54, 46},
		{129, -3, 126},
		{6, 2, 4},
		{6, 0, 6},
		{0, 5, 0},
		{193, 192, 1},
		{6, 6, 0},
		{6, -6, 0},
		{193, -192, 1},
		{193, -1, 192},
		{191, -100, 91},
	}

	for _, tt := range tests {
		c := New(tt.size).SetAll().Shift(tt.shift).OnesCount()
		if c != tt.expected {
			t.Errorf("New(%d).SetAll().Shift(%d) had %d bits set, however %d was expected",
				tt.size, tt.shift, c, tt.expected)
		}
	}

	a := New(3)
	a.SetAll().Shift(-1)
	assert(t, a.String(), "000")

	a.Mut().SetAll().Shift(-1)
	assert(t, a.String(), "110")
}

func TestMid(t *testing.T) {
	a := New(5).SetAll().Mid(1, 1)
	assert(t, a.String(), "1")

	a = New(121).SetAll().Mid(-3, 3)
	assert(t, a.String(), "111")

	if !doesPanic(func() {
		New(5).Left(-1)
	}) {
		t.Error("should panic")
	}

	if !doesPanic(func() {
		New(5).Right(-1)
	}) {
		t.Error("should panic")
	}

	if !doesPanic(func() {
		New(10).Mid(3, -1)
	}) {
		t.Error("should panic")
	}

	a = New(65).Set(3).Left(3)
	assert(t, a.String(), "000")

	a = New(60).SetAll().Right(10)
	b := New(10).SetAll()

	assert(t, a.Equal(b), true)

	a = New(10).SetAll().Right(11)
	b = New(10).SetAll().Left(11)
	assert(t, a.Equal(b), true)

	assert(t, New(4).Left(0).Len(), 0)
}

func TestAppend(t *testing.T) {
	// trivial cases
	a := New(0).Append(New(0))
	assert(t, a.Equal(New(0)), true)

	a = New(0).Append(New(1))
	assert(t, a.Equal(New(1)), true)

	// real cases
	a = New(3).SetAll().Append(New(3))
	assert(t, a.String(), "111000")
}

func TestRotate(t *testing.T) {

	assert(t, New(0).Rotate(4).Len(), 0)

	a := New(65).Set(63).Rotate(1)
	assert(t, a.Equal(New(65).Set(64)), true)

	a = New(65).Set(0).Rotate(-1)
	assert(t, a.Equal(New(65).Set(64)), true)

	const len = 163
	a = New(len).Set(0)
	for i := -len * 2; i < len*2; i++ {
		r := a.Clone().Rotate(i)
		if !r.Equal(New(len).Set(i)) {
			t.Errorf("@%d rotate failed", i)
		}
	}

	a = New(3).Set(0)
	a.Rotate(-1) // discarded
	assert(t, a.String(), "100")

	a.Mut().Rotate(-1)
	assert(t, a.String(), "001")
}

func TestMut(t *testing.T) {
	a := New(65).Mut()
	a.SetAll().Clear(0, -1).Flip(3, 4)
	assert(t, a.OnesCount(), 61)

	a = New(65).Mut().SetAll()
	a.Xor(a)
	assert(t, a.OnesCount(), 0)

	a = New(65)
	assert(t, New(5).Copy(a), false)

	New(65).Set(0, -1).Copy(a)
	assert(t, a.OnesCount(), 2)
	assert(t, a.Get(-1), true)
}

func Benchmark1(b *testing.B) {
	a := New(365)
	for n := 0; n < b.N; n++ {
		a.Set(0)
	}
}

func ExampleBitField_Shift_e1() {
	bf := NewBitField(3).Set(0).Shift(1)
	fmt.Println(bf)
	// Output: 010
}

func ExampleBitField_Shift_e2() {
	bf := NewBitField(3).Set(0).Shift(3)
	fmt.Println(bf)
	// Output: 000
}

func ExampleBitField_Clear() {
	bf := NewBitField(4).SetAll().Clear(0, -1)
	fmt.Println(bf)
	// Output: 0110
}

func ExampleBitField_Mut() {
	bf := NewBitField(4).Set(0)
	bf.Set(1) // this is set then discarded!
	fmt.Println("without Mut():", bf)

	bf.Mut().Set(1)
	fmt.Println("with Mut():", bf)
	// Output: without Mut(): 1000
	// with Mut(): 1100

}

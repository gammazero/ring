package ring

import (
	"fmt"
	"testing"
	"unicode"
)

func TestRingEmpty(t *testing.T) {
	const ringSize = 10
	r := New[string](ringSize)
	if r.Len() != 0 {
		t.Error("r.Len() =", r.Len(), "expect 0")
	}
	if r.Cap() != ringSize {
		t.Error("expected r.Cap() ==", ringSize, "got", r.Cap())
	}
	idx := r.Index(func(item string) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil derue")
	}
	idx = r.RIndex(func(item string) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil derue")
	}
}

func TestRingNil(t *testing.T) {
	var r *Ring[int]
	if r.Len() != 0 {
		t.Error("expected r.Len() == 0")
	}
	if r.Cap() != 0 {
		t.Error("expected r.Cap() == 0")
	}
	r.Rotate(5)
	idx := r.Index(func(item int) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil derue")
	}
	idx = r.RIndex(func(item int) bool {
		return true
	})
	if idx != -1 {
		t.Error("should return -1 index for nil derue")
	}
}

func TestCycleForward(t *testing.T) {
	r := New[rune](5)
	for _, char := range "hello" {
		r.PushBack(char)
	}
	if r.Front() != 'h' {
		t.Fatal("expected first character to be 'h'")
	}
	if r.Back() != 'o' {
		t.Fatal("expected last character to be 'o'")
	}

	r.PushBack('w')
	if r.Front() != 'e' {
		t.Fatal("expected first character to be 'h'")
	}
	if r.Back() != 'w' {
		t.Fatal("expected last character to be 'w'")
	}
	r.PushBack('o')
	if r.Front() != 'l' {
		t.Fatal("expected first character to be 'l'")
	}
	if r.Back() != 'o' {
		t.Fatal("expected last character to be 'o'")
	}

	for _, char := range "rld" {
		r.PushBack(char)
	}
	if r.Front() != 'w' {
		t.Fatal("expected first character to be 'w'")
	}
	if r.Back() != 'd' {
		t.Fatal("expected last character to be 'd'")
	}

	for i, char := range "world" {
		if r.PopFront() != char {
			t.Fatal("Unexpected character", char, "at position", i)
		}
	}
}

func TestCycleBackward(t *testing.T) {
	r := New[rune](5)
	for _, char := range "hello" {
		r.PushFront(char)
	}
	if r.Front() != 'o' {
		t.Fatal("expected first character to be 'o'")
	}
	if r.Back() != 'h' {
		t.Fatal("expected last character to be 'h'")
	}

	r.PushFront('w')
	if r.Front() != 'w' {
		t.Fatal("expected first character to be 'w'")
	}
	if r.Back() != 'e' {
		t.Fatal("expected last character to be 'e'")
	}
	r.PushFront('o')
	if r.Front() != 'o' {
		t.Fatal("expected first character to be 'o'")
	}
	if r.Back() != 'l' {
		t.Fatal("expected last character to be 'l'")
	}

	for _, char := range "rld" {
		r.PushFront(char)
	}
	if r.Front() != 'd' {
		t.Fatal("expected first character to be 'd'")
	}
	if r.Back() != 'w' {
		t.Fatal("expected last character to be 'w'")
	}
}

func TestFrontBack(t *testing.T) {
	r := New[string](4)
	r.PushBack("foo")
	r.PushBack("bar")
	r.PushBack("baz")
	if r.Front() != "foo" {
		t.Error("wrong value at front of ring")
	}
	if r.Back() != "baz" {
		t.Error("wrong value at back of ring")
	}

	if r.PopFront() != "foo" {
		t.Error("wrong value removed from front of ring")
	}
	if r.Front() != "bar" {
		t.Error("wrong value remaining at front of ring")
	}
	if r.Back() != "baz" {
		t.Error("wrong value remaining at back of ring")
	}

	if r.PopBack() != "baz" {
		t.Error("wrong value removed from back of ring")
	}
	if r.Front() != "bar" {
		t.Error("wrong value remaining at front of ring")
	}
	if r.Back() != "bar" {
		t.Error("wrong value remaining at back of ring")
	}
}

func TestSimple(t *testing.T) {
	const ringSize = 10
	r := New[int](ringSize)

	for i := 0; i < ringSize; i++ {
		r.PushBack(i)
	}
	if r.Front() != 0 {
		t.Fatalf("expected 0 at front, got %d", r.Front())
	}
	if r.Back() != ringSize-1 {
		t.Fatalf("expected %d at back, got %d", ringSize-1, r.Back())
	}

	for i := 0; i < ringSize; i++ {
		if r.Front() != i {
			t.Error("peek", i, "had value", r.Front())
		}
		x := r.PopFront()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}

	for i := 0; i < ringSize; i++ {
		r.PushFront(i)
	}
	for i := ringSize - 1; i >= 0; i-- {
		x := r.PopFront()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestWrap(t *testing.T) {
	const ringSize = 10
	r := New[int](ringSize)

	for i := 0; i < ringSize; i++ {
		r.PushBack(i)
	}

	for i := 0; i < 3; i++ {
		r.PopFront()
		r.PushBack(ringSize + i)
	}

	for i := 0; i < ringSize; i++ {
		if r.Front() != i+3 {
			t.Error("peek", i, "had value", r.Front())
		}
		r.PopFront()
	}
}

func TestWrapReverse(t *testing.T) {
	const ringSize = 10
	r := New[int](ringSize)

	for i := 0; i < ringSize; i++ {
		r.PushFront(i)
	}
	for i := 0; i < 3; i++ {
		r.PopBack()
		r.PushFront(ringSize + i)
	}

	for i := 0; i < ringSize; i++ {
		if r.Back() != i+3 {
			t.Error("peek", i, "had value", r.Front())
		}
		r.PopBack()
	}
}

func TestLen(t *testing.T) {
	const ringSize = 10
	r := New[int](ringSize)

	if r.Len() != 0 {
		t.Error("empty ring length not 0")
	}

	for i := 0; i < ringSize; i++ {
		r.PushBack(i)
		if r.Len() != i+1 {
			t.Error("adding: ring with", i, "elements has length", r.Len())
		}
	}

	r.PushBack(10)
	if r.Len() != ringSize {
		t.Fatal("wrong size for full ring")
	}
	r.PushBack(11)
	if r.Len() != ringSize {
		t.Fatal("wrong size for full ring")
	}

	for i := 0; i < ringSize; i++ {
		r.PopFront()
		if r.Len() != ringSize-i-1 {
			t.Error("removing: ring with", ringSize-i-i, "elements has length", r.Len())
		}
	}
}

func TestBack(t *testing.T) {
	const ringSize = 10
	r := New[int](ringSize)

	for i := 0; i < ringSize+5; i++ {
		r.PushBack(i)
		if r.Back() != i {
			t.Errorf("Back returned wrong value")
		}
	}
}

func checkRotate(t *testing.T, size int) {
	r := New[int](size)

	for i := 0; i < size; i++ {
		r.PushBack(i)
	}

	for i := 0; i < r.Len(); i++ {
		x := i
		for n := 0; n < r.Len(); n++ {
			if r.At(n) != x {
				t.Fatalf("a[%d] != %d after rotate and copy", n, x)
			}
			x++
			if x == r.Len() {
				x = 0
			}
		}
		r.Rotate(1)
		if r.Back() != i {
			t.Fatal("wrong value during rotation")
		}
	}
	for i := r.Len() - 1; i >= 0; i-- {
		r.Rotate(-1)
		if r.Front() != i {
			t.Fatal("wrong value during reverse rotation")
		}
	}
}

func TestRotate(t *testing.T) {
	checkRotate(t, 10)

	r := New[int](10)

	for i := 0; i < r.Cap(); i++ {
		r.PushBack(i)
	}
	r.Rotate(11)
	if r.Front() != 1 {
		t.Error("rotating 11 places should have been same as one")
	}
	r.Rotate(-21)
	if r.Front() != 0 {
		t.Error("rotating -21 places should have been same as one -1")
	}
	r.Rotate(r.Len())
	if r.Front() != 0 {
		t.Error("should not have rotated")
	}
	r.Reset()
	r.PushBack(0)
	r.Rotate(13)
	if r.Front() != 0 {
		t.Error("should not have rotated")
	}
}

func TestAt(t *testing.T) {
	r := New[int](10)

	for i := 0; i < r.Cap(); i++ {
		r.PushBack(i)
	}

	// Front to back.
	for j := 0; j < r.Len(); j++ {
		if r.At(j) != j {
			t.Errorf("fwd: index %d doesn't contain %d", j, j)
		}
	}

	// Back to front
	for j := 1; j <= r.Len(); j++ {
		if r.At(r.Len()-j) != r.Len()-j {
			t.Errorf("index %d doesn't contain %d", r.Len()-j, r.Len()-j)
		}
	}
}

func TestSet(t *testing.T) {
	r := New[int](1000)

	for i := 0; i < 1000; i++ {
		r.PushBack(i)
		r.Set(i, i+50)
	}

	// Front to back.
	for j := 0; j < r.Len(); j++ {
		if r.At(j) != j+50 {
			t.Errorf("index %d doesn't contain %d", j, j+50)
		}
	}
}

func TestReset(t *testing.T) {
	r := New[int](100)

	for i := 0; i < 100; i++ {
		r.PushBack(i)
	}
	if r.Len() != 100 {
		t.Error("push: ring with 100 elements has length", r.Len())
	}
	cap := len(r.buf)
	r.Reset()
	if r.Len() != 0 {
		t.Error("empty ring length not 0 after clear")
	}
	if len(r.buf) != cap {
		t.Error("ring capacity changed after clear")
	}

	// Check that there are no remaining references after Reset()
	for i := 0; i < len(r.buf); i++ {
		if r.buf[i] != 0 {
			t.Logf("r[%d] = %d", i, r.buf[i])
			t.Error("ring has non-nil deleted elements after Reset()")
			break
		}
	}
}

func TestIndex(t *testing.T) {
	r := New[rune](16)
	for _, x := range "Hello, 世界" {
		r.PushBack(x)
	}
	idx := r.Index(func(item rune) bool {
		c := item
		return unicode.Is(unicode.Han, c)
	})
	if idx != 7 {
		t.Fatal("Expected index 7, got", idx)
	}
	idx = r.Index(func(item rune) bool {
		c := item
		return c == 'H'
	})
	if idx != 0 {
		t.Fatal("Expected index 0, got", idx)
	}
	idx = r.Index(func(item rune) bool {
		return false
	})
	if idx != -1 {
		t.Fatal("Expected index -1, got", idx)
	}
}

func TestRIndex(t *testing.T) {
	r := New[rune](16)
	for _, x := range "Hello, 世界" {
		r.PushBack(x)
	}
	idx := r.RIndex(func(item rune) bool {
		c := item
		return unicode.Is(unicode.Han, c)
	})
	if idx != 8 {
		t.Fatal("Expected index 8, got", idx)
	}
	idx = r.RIndex(func(item rune) bool {
		c := item
		return c == 'H'
	})
	if idx != 0 {
		t.Fatal("Expected index 0, got", idx)
	}
	idx = r.RIndex(func(item rune) bool {
		return false
	})
	if idx != -1 {
		t.Fatal("Expected index -1, got", idx)
	}
}

func TestInsert(t *testing.T) {
	r := New[rune](16)
	for _, x := range "ABCDEFG" {
		r.PushBack(x)
	}
	r.Insert(4, 'x') // ABCDxEFG
	if r.At(4) != 'x' {
		t.Error("expected x at position 4, got", r.At(4))
	}

	r.Insert(2, 'y') // AByCDxEFG
	if r.At(2) != 'y' {
		t.Error("expected y at position 2")
	}
	if r.At(5) != 'x' {
		t.Error("expected x at position 5")
	}

	r.Insert(0, 'b') // bAByCDxEFG
	if r.Front() != 'b' {
		t.Error("expected b inserted at front, got", r.Front())
	}

	r.Insert(r.Len(), 'e') // bAByCDxEFGe

	for i, x := range "bAByCDxEFGe" {
		if r.PopFront() != x {
			t.Error("expected", x, "at position", i)
		}
	}

	rs := New[string](16)

	for i := 0; i < rs.Cap(); i++ {
		rs.PushBack(fmt.Sprint(i))
	}
	// derue: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
	// buffer: [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
	for i := 0; i < rs.Cap()/2; i++ {
		rs.PopFront()
	}
	// derue: 8 9 10 11 12 13 14 15
	// buffer: [_,_,_,_,_,_,_,_,8,9,10,11,12,13,14,15]
	for i := 0; i < rs.Cap()/4; i++ {
		rs.PushBack(fmt.Sprint(rs.Cap() + i))
	}
	// derue: 8 9 10 11 12 13 14 15 16 17 18 19
	// buffer: [16,17,18,19,_,_,_,_,8,9,10,11,12,13,14,15]

	at := rs.Len() - 2
	rs.Insert(at, "x")
	// derue: 8 9 10 11 12 13 14 15 16 17 x 18 19
	// buffer: [16,17,x,18,19,_,_,_,8,9,10,11,12,13,14,15]
	if rs.At(at) != "x" {
		t.Error("expected x at position", at)
	}
	if rs.At(at) != "x" {
		t.Error("expected x at position", at)
	}

	rs.Insert(2, "y")
	// derue: 8 9 y 10 11 12 13 14 15 16 17 x 18 19
	// buffer: [16,17,x,18,19,_,_,8,9,y,10,11,12,13,14,15]
	if rs.At(2) != "y" {
		t.Error("expected y at position 2")
	}
	if rs.At(at+1) != "x" {
		t.Error("expected x at position 5")
	}

	rs.Insert(0, "b")
	// derue: b 8 9 y 10 11 12 13 14 15 16 17 x 18 19
	// buffer: [16,17,x,18,19,_,b,8,9,y,10,11,12,13,14,15]
	if rs.Front() != "b" {
		t.Error("expected b inserted at front, got", rs.Front())
	}

	rs.Insert(rs.Len(), "e")
	if rs.Cap() != rs.Len() {
		t.Fatal("Expected full buffer")
	}
	// derue: b 8 9 y 10 11 12 13 14 15 16 17 x 18 19 e
	// buffer: [16,17,x,18,19,e,b,8,9,y,10,11,12,13,14,15]
	for i, x := range []string{"16", "17", "x", "18", "19", "e", "b", "8", "9", "y", "10", "11", "12", "13", "14", "15"} {
		if rs.buf[i] != x {
			t.Error("expected", x, "at buffer position", i)
		}
	}
	for i, x := range []string{"b", "8", "9", "y", "10", "11", "12", "13", "14", "15", "16", "17", "x", "18", "19", "e"} {
		if rs.Front() != x {
			t.Error("expected", x, "at position", i, "got", rs.Front())
		}
		rs.PopFront()
	}
}

func TestRemove(t *testing.T) {
	r := New[rune](16)
	for _, x := range "ABCDEFG" {
		r.PushBack(x)
	}

	if r.Remove(4) != 'E' { // ABCDFG
		t.Error("expected E from position 4")
	}

	if r.Remove(2) != 'C' { // ABDFG
		t.Error("expected C at position 2")
	}
	if r.Back() != 'G' {
		t.Error("expected G at back")
	}

	if r.Remove(0) != 'A' { // BDFG
		t.Error("expected to remove A from front")
	}
	if r.Front() != 'B' {
		t.Error("expected G at back")
	}

	if r.Remove(r.Len()-1) != 'G' { // BDF
		t.Error("expected to remove G from back")
	}
	if r.Back() != 'F' {
		t.Error("expected F at back")
	}

	if r.Len() != 3 {
		t.Error("wrong length")
	}
}

func TestFrontBackOutOfRangePanics(t *testing.T) {
	const msg = "should panic when peeking empty ring"
	r := New[rune](16)

	assertPanics(t, msg, func() {
		r.Front()
	})
	assertPanics(t, msg, func() {
		r.Back()
	})

	r.PushBack(1)
	r.PopFront()

	assertPanics(t, msg, func() {
		r.Front()
	})
	assertPanics(t, msg, func() {
		r.Back()
	})
}

func TestPopFrontOutOfRangePanics(t *testing.T) {
	r := New[rune](16)

	assertPanics(t, "should panic when removing empty ring", func() {
		r.PopFront()
	})

	r.PushBack(1)
	r.PopFront()

	assertPanics(t, "should panic when removing emptied ring", func() {
		r.PopFront()
	})
}

func TestPopBackOutOfRangePanics(t *testing.T) {
	r := New[rune](16)

	assertPanics(t, "should panic when removing empty ring", func() {
		r.PopBack()
	})

	r.PushBack(1)
	r.PopBack()

	assertPanics(t, "should panic when removing emptied ring", func() {
		r.PopBack()
	})
}

func TestAtOutOfRangePanics(t *testing.T) {
	r := New[rune](16)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	assertPanics(t, "should panic when negative index", func() {
		r.At(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		r.At(4)
	})
}

func TestSetOutOfRangePanics(t *testing.T) {
	r := New[rune](16)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	assertPanics(t, "should panic when negative index", func() {
		r.Set(-4, 1)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		r.Set(4, 1)
	})
}

func TestInsertOutOfRangePanics(t *testing.T) {
	r := New[string](16)

	assertPanics(t, "should panic when inserting out of range", func() {
		r.Insert(1, "X")
	})

	r.PushBack("A")

	assertPanics(t, "should panic when inserting at negative index", func() {
		r.Insert(-1, "Y")
	})

	assertPanics(t, "should panic when inserting out of range", func() {
		r.Insert(2, "B")
	})
}

func TestRemoveOutOfRangePanics(t *testing.T) {
	r := New[string](16)

	assertPanics(t, "should panic when removing from empty ring", func() {
		r.Remove(0)
	})

	r.PushBack("A")

	assertPanics(t, "should panic when removing at negative index", func() {
		r.Remove(-1)
	})

	assertPanics(t, "should panic when removing out of range", func() {
		r.Remove(1)
	})
}

func assertPanics(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: didn't panic as expected", name)
		}
	}()

	f()
}

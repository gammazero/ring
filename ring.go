package ring

import "fmt"

// Ring is a fixed-size circular buffer of items of the type sepcified by the
// type argument. Pushing an item onto a full Ring overwrites the item at the
// other end of the ring.
type Ring[T any] struct {
	buf   []T
	head  int
	tail  int
	count int
}

func New[T any](capacity int) *Ring[T] {
	return &Ring[T]{
		buf: make([]T, capacity),
	}
}

// Cap returns the current capacity of the Ring. If r is nil, r.Cap() is zero.
func (r *Ring[T]) Cap() int {
	if r == nil {
		return 0
	}
	return len(r.buf)
}

// Len returns the number of elements currently stored in the Ring. If r is
// nil, r.Len() is zero.
func (r *Ring[T]) Len() int {
	if r == nil {
		return 0
	}
	return r.count
}

func (r *Ring[T]) Full() bool {
	return r.count == len(r.buf)
}

// PushBack appends an element to the back of the Ring. Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack. Wraps by overwriting front when Ring is full.
func (r *Ring[T]) PushBack(elem T) {
	r.buf[r.tail] = elem
	r.tail = r.next(r.tail)

	// If full, move head. Otherwise, increment count.
	if r.count == len(r.buf) {
		r.head = r.next(r.head)
	} else {
		r.count++
	}
}

// PushFront prepends an element to the front of the Ring. Implements FIFO when
// elements are removed with PopBack(), and LIFO when elements are removed with
// PopFront. Wraps by overwriting back when Ring is full.
func (r *Ring[T]) PushFront(elem T) {
	// Calculate new head position.
	r.head = r.prev(r.head)
	r.buf[r.head] = elem

	// If full, move tail. Otherwise, increment count.
	if r.count == len(r.buf) {
		r.tail = r.prev(r.tail)
	} else {
		r.count++
	}
}

// PopFront removes and returns the element from the front of the Ring.
// Implements FIFO when used with PushBack(). If the Ring is empty, the call
// panics.
func (r *Ring[T]) PopFront() T {
	if r.count <= 0 {
		panic("PopFront called when empty")
	}
	ret := r.buf[r.head]
	var zero T
	r.buf[r.head] = zero
	// Calculate new head position.
	r.head = r.next(r.head)
	r.count--
	return ret
}

// PopBack removes and returns the element from the back of the Ring.
// Implements LIFO when used with PushBack(). If the Ring is empty, the call
// panics.
func (r *Ring[T]) PopBack() T {
	if r.count <= 0 {
		panic("PopBack called when empty")
	}
	// Calculate new tail position
	r.tail = r.prev(r.tail)

	// Remove value at tail.
	ret := r.buf[r.tail]
	var zero T
	r.buf[r.tail] = zero
	r.count--
	return ret
}

// Front returns the element at the front of the Ring. This is the element that
// would be returned by PopFront(). This call panics if the Ring is empty.
func (r *Ring[T]) Front() T {
	if r.count <= 0 {
		panic("Front called when empty")
	}
	return r.buf[r.head]
}

// Back returns the element at the back of the Ring. This is the element that
// would be returned by PopBack(). This call panics if the Ring is empty.
func (r *Ring[T]) Back() T {
	if r.count <= 0 {
		panic("Back called when empty")
	}

	return r.buf[r.prev(r.tail)]
}

// At returns the element at index i in the Ring without removing the element
// from the Ring. This method accepts only non-negative index values. At(0)
// refers to the first element and is the same as Front(). At(Len()-1) refers
// to the last element and is the same as Back(). If the index is invalid, the
// call panics.
//
// The purpose of At is to allow Ring to serve as a more general purpose
// circular buffer, where items are only added to and removed from the ends of
// the Ring, but may be read from any place within the Ring. Consider the case
// of a fixed-size circular log buffer: A new entry is pushed onto one end and
// when full the oldest is popped from the other end. All the log entries in
// the buffer must be readable without altering the buffer contents.
func (r *Ring[T]) At(i int) T {
	if i < 0 || i >= r.Len() {
		panic(outOfRangeText(i, r.Len()))
	}
	return r.buf[(r.head+i)%len(r.buf)]
}

// Set assigns the item to index i in the Ring. Set indexes the Ring the same
// as At but perform the opposite operation. If the index is invalid, the call
// panics.
func (r *Ring[T]) Set(i int, item T) {
	if i < 0 || i >= r.Len() {
		panic(outOfRangeText(i, r.Len()))
	}
	r.buf[(r.head+i)%len(r.buf)] = item
}

// Rotate rotates the Ring n steps front-to-back. If n is negative, rotates
// back-to-front. Having Ring provide Rotate() allows a more efficient
// implementation, than only Pop and Push methods. that operates by only moving
// the head and tail of the Ring. If Len() is one or less, or Ring is nil, then
// Rotate does nothing.
func (r *Ring[T]) Rotate(n int) {
	if r.Len() <= 1 {
		return
	}
	// Rotating a multiple of count is same as no rotation.
	n %= r.count
	if n == 0 {
		return
	}

	l := len(r.buf)

	// If no empty space in buffer, only move head and tail indexes.
	if r.head == r.tail {
		// Calculate new head and tail.
		r.head = (r.head + n + l) % l
		r.tail = r.head
		return
	}

	var zero T

	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new head and tail.
			r.head = (r.head - 1 + l) % l
			r.tail = (r.tail - 1 + l) % l
			// Put tail value at head and remove value at tail.
			r.buf[r.head] = r.buf[r.tail]
			r.buf[r.tail] = zero
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put head value at tail and remove value at head.
		r.buf[r.tail] = r.buf[r.head]
		r.buf[r.head] = zero
		// Calculate new head and tail.
		r.head = (r.head + 1) % l
		r.tail = (r.tail + 1) % l
	}
}

// Index returns the index into the Ring of the first item satisfying f(item),
// or -1 if none do. If Ring is nil, then -1 is always returned. Search is
// linear starting with index 0.
func (r *Ring[T]) Index(f func(T) bool) int {
	if r.Len() > 0 {
		for i := 0; i < r.count; i++ {
			if f(r.buf[(r.head+i)%len(r.buf)]) {
				return i
			}
		}
	}
	return -1
}

// RIndex is the same as Index, but searches from Back to Front. The index
// returned is from Front to Back, where index 0 is the index of the item
// returned by Front().
func (r *Ring[T]) RIndex(f func(T) bool) int {
	if r.Len() > 0 {
		l := r.Len()
		for i := r.count - 1; i >= 0; i-- {
			if f(r.buf[(r.head+i)%l]) {
				return i
			}
		}
	}
	return -1
}

// Insert is used to insert an element into the middle of the Ring, before the
// element at the specified index. Insert(0,e) is the same as PushFront(e) and
// Insert(Len(),e) is the same as PushBack(e). Accepts only non-negative index
// values, and panics if index is out of range.
//
// Important: Ring is optimized for O(1) operations at the ends of the Ring,
// not for operations in the the middle. Complexity of this function is
// constant plus linear in the lesser of the distances between the index and
// either of the ends of the Ring.
func (r *Ring[T]) Insert(at int, item T) {
	if at < 0 || at > r.count {
		panic(outOfRangeText(at, r.Len()))
	}
	if r.Full() {
		panic("cannot insert into full ring")
	}
	if at*2 < r.count {
		r.PushFront(item)
		front := r.head
		for i := 0; i < at; i++ {
			next := r.next(front)
			r.buf[front], r.buf[next] = r.buf[next], r.buf[front]
			front = next
		}
		return
	}
	swaps := r.count - at
	r.PushBack(item)
	back := r.prev(r.tail)
	for i := 0; i < swaps; i++ {
		prev := r.prev(back)
		r.buf[back], r.buf[prev] = r.buf[prev], r.buf[back]
		back = prev
	}
}

// Remove removes and returns an element from the middle of the Ring, at the
// specified index. Remove(0) is the same as PopFront() and Remove(Len()-1) is
// the same as PopBack(). Accepts only non-negative index values, and panics if
// index is out of range.
//
// Important: Ring is optimized for O(1) operations at the ends of the Ring,
// not for operations in the the middle. Complexity of this function is
// constant plus linear in the lesser of the distances between the index and
// either of the ends of the Ring.
func (r *Ring[T]) Remove(at int) T {
	if at < 0 || at >= r.Len() {
		panic(outOfRangeText(at, r.Len()))
	}

	rm := (r.head + at) % len(r.buf)
	if at*2 < r.count {
		for i := 0; i < at; i++ {
			prev := r.prev(rm)
			r.buf[prev], r.buf[rm] = r.buf[rm], r.buf[prev]
			rm = prev
		}
		return r.PopFront()
	}
	swaps := r.count - at - 1
	for i := 0; i < swaps; i++ {
		next := r.next(rm)
		r.buf[rm], r.buf[next] = r.buf[next], r.buf[rm]
		rm = next
	}
	return r.PopBack()
}

// Reset resets the Ring to be empty, but it retains the underlying storage for
// use by future writes.
func (r *Ring[T]) Reset() {
	var zero T
	l := len(r.buf)
	h := r.head
	for i := 0; i < r.Len(); i++ {
		r.buf[(h+i)%l] = zero
	}
	r.head = 0
	r.tail = 0
	r.count = 0
}

// prev returns the previous buffer position wrapping around buffer.
func (r *Ring[T]) prev(i int) int {
	l := len(r.buf)
	return (i - 1 + l) % l
}

// next returns the next buffer position wrapping around buffer.
func (r *Ring[T]) next(i int) int {
	return (i + 1) % len(r.buf)
}

// Resize resizes the Ring to have the specified capacity. Any items present in
// the Ring are copied into the resized ring.
func (r *Ring[T]) Resize(newSize int) {
	if r.count == newSize {
		return
	}

	newBuf := make([]T, newSize)
	var n int
	if r.tail > r.head {
		n = copy(newBuf, r.buf[r.head:r.tail])
	} else {
		n = copy(newBuf, r.buf[r.head:])
		if n < len(r.buf) {
			n += copy(newBuf[n:], r.buf[:r.tail])
		}
	}

	r.count = n
	r.head = 0
	r.tail = n
	r.buf = newBuf
}

func outOfRangeText(i, len int) string {
	return fmt.Sprintf("ring: index out of range %d with length %d", i, len)
}

package main

import (
	"errors"
	"fmt"
)

// ErrIsEmpty 缓冲区为空
var ErrIsEmpty = errors.New("ring buffer is empty")

// RingSlice 自动扩容循环缓冲区
type RingSlice struct {
	buf      []int64
	initSize int
	size     int
	vr       int
	r        int // next position to read
	w        int // next position to write
	isEmpty  bool
	callback func(diff int)
}

// New 返回一个初始大小为 size 的 RingSlice
func NewRingSlice(size int) *RingSlice {
	return &RingSlice{
		buf:      make([]int64, size),
		initSize: size,
		size:     size,
		isEmpty:  true,
	}
}

// NewWithData 特殊场景使用，RingSlice 会持有data，不会自己申请内存去拷贝
func NewRingSliceWithData(data []int64) *RingSlice {
	return &RingSlice{
		buf:      data,
		size:     len(data),
		initSize: len(data),
	}
}

func (r *RingSlice) WithData(data []int64) {
	r.r = 0
	r.w = 0
	r.vr = 0
	r.isEmpty = false
	r.size = len(data)
	r.initSize = len(data)
	r.buf = data
}

// VirtualFlush 刷新虚读指针
// VirtualXXX 系列配合使用
func (r *RingSlice) VirtualFlush() {
	r.r = r.vr
	if r.r == r.w {
		r.isEmpty = true
	}
}

// VirtualRevert 还原虚读指针
// VirtualXXX 系列配合使用
func (r *RingSlice) VirtualRevert() {
	r.vr = r.r
}

// VirtualRead 虚读，不移动 read 指针，需要配合 VirtualFlush 和 VirtualRevert 使用
// VirtualXXX 系列配合使用
func (r *RingSlice) VirtualRead(p []int64) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.isEmpty {
		return 0, ErrIsEmpty
	}
	n = len(p)
	if r.w > r.vr {
		if n > r.w-r.vr {
			n = r.w - r.vr
		}
		copy(p, r.buf[r.vr:r.vr+n])
		// move vr
		r.vr = (r.vr + n) % r.size
		if r.vr == r.w {
			r.isEmpty = true
		}
		return
	}
	if n > r.size-r.vr+r.w {
		n = r.size - r.vr + r.w
	}
	if r.vr+n <= r.size {
		copy(p, r.buf[r.vr:r.vr+n])
	} else {
		// head
		copy(p, r.buf[r.vr:r.size])
		// tail
		copy(p[r.size-r.vr:], r.buf[0:n-r.size+r.vr])
	}

	// move vr
	r.vr = (r.vr + n) % r.size
	return
}

// VirtualLength 虚拟长度，虚读后剩余可读数据长度
// VirtualXXX 系列配合使用
func (r *RingSlice) VirtualLength() int {
	if r.w == r.vr {
		if r.isEmpty {
			return 0
		}
		return r.size
	}

	if r.w > r.vr {
		return r.w - r.vr
	}

	return r.size - r.vr + r.w
}

func (r *RingSlice) RetrieveAll() {
	r.r = 0
	r.w = 0
	r.vr = 0
	r.isEmpty = true
}

func (r *RingSlice) Retrieve(len int) {
	if r.isEmpty || len <= 0 {
		return
	}

	if len < r.Length() {
		r.r = (r.r + len) % r.size
		r.vr = r.r

		if r.w == r.r {
			r.isEmpty = true
		}
	} else {
		r.RetrieveAll()
	}
}

func (r *RingSlice) Peek(len int) (first []int64, end []int64) {
	if r.isEmpty || len <= 0 {
		return
	}

	if r.w > r.r {
		if len > r.w-r.r {
			len = r.w - r.r
		}

		first = r.buf[r.r : r.r+len]
		return
	}

	if len > r.size-r.r+r.w {
		len = r.size - r.r + r.w
	}
	if r.r+len <= r.size {
		first = r.buf[r.r : r.r+len]
	} else {
		// head
		first = r.buf[r.r:r.size]
		// tail
		end = r.buf[0 : len-r.size+r.r]
	}
	return
}

func (r *RingSlice) PeekAll() (first []int64, end []int64) {
	if r.isEmpty {
		return
	}

	if r.w > r.r {
		first = r.buf[r.r:r.w]
		return
	}

	first = r.buf[r.r:r.size]
	end = r.buf[0:r.w]
	return
}

func (r *RingSlice) Read(p []int64) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.isEmpty {
		return 0, ErrIsEmpty
	}
	n = len(p)
	if r.w > r.r {
		if n > r.w-r.r {
			n = r.w - r.r
		}
		copy(p, r.buf[r.r:r.r+n])
		// move readPtr
		r.r = (r.r + n) % r.size
		if r.r == r.w {
			r.isEmpty = true
		}
		r.vr = r.r
		return
	}
	if n > r.size-r.r+r.w {
		n = r.size - r.r + r.w
	}
	if r.r+n <= r.size {
		copy(p, r.buf[r.r:r.r+n])
		for i := r.r; i < r.r+n; i++ {
			r.buf[i] = 0
		}
	} else {
		// head
		copy(p, r.buf[r.r:r.size])
		for i := r.r; i < r.size; i++ {
			r.buf[i] = 0
		}
		// tail
		copy(p[r.size-r.r:], r.buf[0:n-r.size+r.r])
		for i := 0; i < n-r.size+r.r; i++ {
			r.buf[i] = 0
		}
	}

	// move readPtr
	r.r = (r.r + n) % r.size
	if r.r == r.w {
		r.isEmpty = true
	}
	r.vr = r.r
	return
}
func (r *RingSlice) ReadOneDiff(n int) (result int64, err error) {
	t := r.r + n
	if t >= r.size {
		t -= r.size
	}
	if t < 0 {
		t += r.size
	} else {
		t = t % r.size
	}

	return r.buf[t], nil
}
func (r *RingSlice) ReadOneAbs(t int) (result int64, err error) {
	if t < 0 {
		t += r.size
	} else {
		t = t % r.size
	}

	return r.buf[t], nil
}

func (r *RingSlice) Jump(n int) {
	// move readPtr
	r.r = (r.r + n) % r.size
	if r.r == r.w {
		r.isEmpty = true
	}
	r.vr = r.r
	return
}
func (r *RingSlice) VirtualTurnTrue() {
	r.r = r.vr
	return
}

func (r *RingSlice) Write(p []int64) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	n = len(p)
	free := r.free()
	if free < n {
		r.makeSpace(n - free)
	}
	if r.w >= r.r {
		if r.size-r.w >= n {
			copy(r.buf[r.w:], p)
			r.w += n
		} else {
			copy(r.buf[r.w:], p[:r.size-r.w])
			copy(r.buf[0:], p[r.size-r.w:])
			r.w += n - r.size
		}
	} else {
		copy(r.buf[r.w:], p)
		r.w += n
	}

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return
}
func (r *RingSlice) WriteSingleDiff(diff int, p int64) (n int, err error) {
	if diff == 0 {
		err := r.WriteSingle(p)
		return 1, err
	}
	if diff < 0 {
		t := r.w + diff
		if t < 0 {
			t += r.size
		}
		r.buf[t] = p
		return 1, nil
	}
	n = diff + 1
	free := r.free()
	if free < n {
		r.makeSpace(n - free)
	}
	if r.w >= r.r {
		if r.size-r.w >= n {
			//copy(r.buf[r.w:], p)
			r.buf[r.w+n-1] = p
			r.w += n
		} else {
			//copy(r.buf[r.w:], p[:r.size-r.w])
			//copy(r.buf[0:], p[r.size-r.w:])
			r.buf[0+n-1-(r.size-r.w)] = p
			r.w += n - r.size
		}
	} else {
		//copy(r.buf[r.w:], p)
		r.buf[r.w+n-1] = p
		r.w += n
	}

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return
}
func (r *RingSlice) WriteSingleAbs(index int, p int64) (n int, err error) {
	if index < r.w {
		r.buf[index] = p
		return 1, nil
	}
	//if index == r.w {
	//	r.buf[index] = p
	//	r.w++
	//	r.isEmpty = false
	//	return 1, nil
	//}
	if index > r.size-2 {
		t := r.r
		r.makeSpaceN(index + 2) //参数是增加的个数
		index -= t
	}

	r.buf[index] = p
	if r.w < index {
		r.w = index
	}
	r.w += 1

	if r.w == r.size {
		r.w = 0
	}

	r.isEmpty = false

	return
}
func (r *RingSlice) WriteSingle(data int64) (err error) {
	// 检查剩余空间，如有必要则扩展
	if r.free() < 1 {
		r.makeSpace(1)
	}

	// 写入数据
	r.buf[r.w] = data
	r.w = (r.w + 1) % r.size

	r.isEmpty = false

	return nil
}
func (r *RingSlice) Length() int {
	if r.w == r.r {
		if r.isEmpty {
			return 0
		}
		return r.size
	}

	if r.w > r.r {
		return r.w - r.r
	}

	return r.size - r.r + r.w
}

func (r *RingSlice) Capacity() int {
	return r.size
}

// int64s 返回所有可读数据，此操作不会移动读指针，仅仅是拷贝全部数据
func (r *RingSlice) CopyAll() (buf []int64) {
	if r.isEmpty {
		return
	}

	if r.w > r.r {
		buf = make([]int64, r.w-r.r)
		copy(buf, r.buf[r.r:r.w])
		return
	}

	buf = make([]int64, r.size-r.r+r.w)
	copy(buf, r.buf[r.r:r.size])
	copy(buf[r.size-r.r:], r.buf[0:r.w])
	return
}

func (r *RingSlice) IsFull() bool {
	return !r.isEmpty && r.w == r.r
}

func (r *RingSlice) IsEmpty() bool {
	return r.isEmpty
}

func (r *RingSlice) Reset() {
	r.r = 0
	r.vr = 0
	r.w = 0
	r.isEmpty = true
	if r.size > r.initSize {
		r.buf = make([]int64, r.initSize)
		r.size = r.initSize
	}
}

func (r *RingSlice) String() string {
	return fmt.Sprintf("Ring Buffer: \n\tCap: %d\n\tReadable int64s: %d\n\tWriteable int64s: %d\n\tBuffer: %v\n", r.size, r.Length(), r.free(), r.buf)
}

func (r *RingSlice) makeSpace(len int) {
	vlen := r.VirtualLength()
	newSize := r.grow(r.size + len)
	newBuf := make([]int64, newSize)
	oldLen := r.Length()
	_, _ = r.Read(newBuf)

	r.w = oldLen
	r.callback(-r.r)
	r.r = 0
	r.vr = oldLen - vlen
	r.size = newSize
	r.buf = newBuf
}
func (r *RingSlice) makeSpaceN(newSize int) {
	vlen := r.VirtualLength()
	newSize = r.grow(newSize)
	newBuf := make([]int64, newSize)
	oldLen := r.Length()
	r.callback(-r.r)
	_, _ = r.Read(newBuf)

	r.w = oldLen
	r.r = 0
	r.vr = oldLen - vlen
	r.size = newSize
	r.buf = newBuf
}

func (r *RingSlice) free() int {
	if r.w == r.r {
		if r.isEmpty {
			return r.size
		}
		return 0
	}
	if r.w < r.r {
		return r.r - r.w
	}

	return r.size - r.w + r.r
}

func copyByte(f, e []int64) []int64 {
	buf := make([]int64, len(f)+len(e))
	_ = copy(buf, f)
	_ = copy(buf[len(f):], e)
	return buf
}

func (r *RingSlice) grow(cap int) int {
	newcap := r.size
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if r.size < 1024 {
			newcap = doublecap
		} else {
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
	return newcap
}

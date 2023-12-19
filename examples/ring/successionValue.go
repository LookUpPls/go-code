package main

import (
	"fmt"
)

type succession struct {
	diff   int
	bp     int
	len    int
	right  int
	body   *RingSlice
	inited bool
}

func (s *succession) init(len int) {
	s.inited = false
	s.body = NewRingSlice(len)
	s.body.callback = func(d int) {
		s.diff += -d
	}
}
func (s *succession) put(id int, time int64) {
	if id < s.diff {
		return
	}
	if s.inited == false {
		s.inited = true
		s.diff = id // 相对于写指针的第几个元素
		s.body.WriteSingleAbs(0, time)
		s.bp = 0 //连续串的最大值 指针
		s.len = 1
		s.right = 1
	} else {
		index := id - s.diff
		s.body.WriteSingleAbs(index, time)
		if index == s.bp+1 {
			s.bp++
			s.len++
			last := time
			for {
				v, _ := s.body.ReadOneAbs(s.bp + 1)
				if v == 0 || s.bp == s.body.size-1 {
					break
				}
				last = v
				s.bp++
				s.len++
			}

			//if s.len > 10 {
			// 取出连续字串
			read := make([]int64, s.len)
			s.body.Read(read)
			//s.bp = 0
			fmt.Printf("\n--wirte last %d,len %d  %v\n\n", last, s.len, read)
			s.len = 0
			//}
		}
	}
}

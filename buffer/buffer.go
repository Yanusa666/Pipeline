package buffer

import (
	"sync"
)

type RingIntBuffer struct {
	array []int64
	pos   int64
	size  int64
	m     sync.Mutex
}

func New(size int64) *RingIntBuffer {
	return &RingIntBuffer{
		array: make([]int64, size),
		pos:   -1,
		size:  size,
		m:     sync.Mutex{},
	}
}

func (r *RingIntBuffer) Push(el int64) {
	r.m.Lock()
	defer r.m.Unlock()

	if r.pos == r.size-1 {
		for i := int64(1); i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else {
		r.pos++
		r.array[r.pos] = el
	}
}

func (r *RingIntBuffer) Get() []int64 {
	if r.pos < 0 {
		return nil
	}

	r.m.Lock()
	defer r.m.Unlock()

	output := r.array[:r.pos+1]
	r.pos = -1

	return output
}

type StageInt func(<-chan bool, <-chan int) <-chan int

type PipeLineInt struct {
	stages []StageInt
	done   <-chan bool
}

package usecase

import "sync/atomic"

type Generator struct {
	current int32
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (gen *Generator) Next(postsCount int) []int {
	ids := make([]int, 0, postsCount)

	for i := 0; i < postsCount; i++ {
		ids = append(ids, int(atomic.AddInt32(&gen.current, 1)))
	}

	return ids
}
package utils

import (
	"errors"

	"gitlab.wildberries.ru/logisticcloud/logistic-route/lml-calculator/internal/rand"
)

type node struct {
	key   float64
	val   int
	tower []*node
}

type randomizable interface {
	Float64() float64
	Intn(n int) int
}

type Option func(*Skiplist)

func WithRnd(rnd randomizable) Option {
	return func(sl *Skiplist) {
		sl.rnd = rnd
	}
}

func WithMax(m int) Option {
	return func(sl *Skiplist) {
		sl.max = m
	}
}

func WithP(p float64) Option {
	return func(sl *Skiplist) {
		sl.p = p
	}
}

type Skiplist struct {
	head   *node
	height int
	max    int
	p      float64
	rnd    randomizable
}

func NewSkipList(opts ...Option) *Skiplist {
	sl := &Skiplist{
		max:    16,
		p:      0.5,
		height: 1,
	}
	for _, o := range opts {
		o(sl)
	}

	sl.head = &node{
		tower: make([]*node, sl.max),
	}

	if sl.rnd == nil {
		sl.rnd = rand.New(1 << 23)
	}

	return sl
}

func (sl *Skiplist) randLevel() int {
	dice := sl.rnd.Float64()
	height := 1
	for dice < sl.p && height < sl.max {
		height++
		dice = sl.rnd.Float64()
	}

	return height
}

// PopMin returns minimal ( leftmost ) element - it's key (float64) and val (int)
//
// returned node deleted from skiplist
func (sl *Skiplist) PopMin() (float64, int, error) {
	if sl.head.tower[0] == nil {
		return 0, 0, errors.New("empty list")
	}

	res := sl.head.tower[0]

	for h := sl.height - 1; h >= 0; h-- {
		sl.head.tower[h] = sl.head.tower[h].tower[h]
		if sl.head.tower[h] == nil {
			sl.height--
		}
	}

	return res.key, res.val, nil
}

func (sl *Skiplist) Insert(key float64, val int) {
	path := make([]*node, sl.max)
	curr := sl.head

	for h := sl.height - 1; h >= 0; h-- {
		for curr.tower[h] != nil && curr.tower[h].key < key {
			curr = curr.tower[h]
		}
		path[h] = curr
	}

	curr = curr.tower[0]

	if curr != nil && curr.key == key {
		curr.val = val

		return
	}

	newHeight := sl.randLevel()
	newNode := &node{key: key, val: val, tower: make([]*node, sl.max)}

	for h := 0; h < newHeight; h++ {
		curr = path[h]

		if curr == nil { // если превысили высоту текущей башни, т.е. на предложенном уровне еще нет списка
			curr = sl.head
		}
		newNode.tower[h] = curr.tower[h]
		curr.tower[h] = newNode
	}

	if newHeight > sl.height {
		sl.height = newHeight
	}
}

func (sl *Skiplist) IsEmpty() bool {
	return sl.head.tower[0] == nil
}

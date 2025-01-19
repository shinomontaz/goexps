package utils

type Qlist struct {
	list []item
}

func NewQlist() *Qlist {
	return &Qlist{
		list: []item{},
	}
}

func (q *Qlist) Clear() {
	q.list = q.list[:0]
}

func (q *Qlist) IsEmpty() bool {
	return len(q.list) == 0
}

func (q *Qlist) Len() int {
	return len(q.list)
}

func (q *Qlist) Insert(w float64, i int) {
	q.list = append(q.list, item{i: i, w: w})
}

func (q *Qlist) PopMin() int {
	min := q.list[0].w
	minIdx := q.list[0].i
	minQidx := 0

	for idx, it := range q.list {
		if it.w < min {
			min = it.w
			minIdx = it.i
			minQidx = idx
		}
	}

	q.list = append(q.list[:minQidx], q.list[minQidx+1:]...)

	return minIdx
}

func (q *Qlist) GetMin() int {
	min := q.list[0].w
	minIdx := q.list[0].i

	for _, it := range q.list {
		if it.w < min {
			min = it.w
			minIdx = it.i
		}
	}

	return minIdx
}

func (q *Qlist) Discard(i int) {
	for idx, it := range q.list {
		if it.i == i {
			q.list = append(q.list[:idx], q.list[idx+1:]...)
			return
		}
	}
}

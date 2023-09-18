package syncx

import "container/list"

type List[k any] struct {
	list.List
}

func NewList[k any]() *List[k] {
	return &List[k]{}
}

func (l *List[k]) PushBack(v k) {
	l.List.PushBack(v)
}

func (l *List[k]) PushFront(v k) {
	l.List.PushFront(v)
}

func (l *List[k]) PopBack() k {
	return l.List.Remove(l.List.Back()).(k)
}

func (l *List[k]) PopFront() k {
	return l.List.Remove(l.List.Front()).(k)
}

func (l *List[k]) Front() k {
	return l.List.Front().Value.(k)
}

func (l *List[k]) Back() k {
	return l.List.Back().Value.(k)
}

func (l *List[k]) Len() int {
	return l.List.Len()
}

func (l *List[k]) Remove(e *list.Element) k {
	return l.List.Remove(e).(k)
}

func (l *List[k]) Range(f func(k)) {
	for e := l.List.Front(); e != nil; e = e.Next() {
		f(e.Value.(k))
	}
}

func (l *List[k]) Clear() {
	l.List.Init()
}

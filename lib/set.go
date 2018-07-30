package lib

import "sync"

type Set struct {
	m sync.Map
	bid sync.Map
}

func NewSet() *Set{
	set := &Set{
		m: sync.Map{},
		bid: sync.Map{},
	}
	return set
}

// 判断是否存在，如果不存在，存储并返回true，否则返回false
func (s *Set) Add(item int) bool {
	_, ok := s.m.LoadOrStore(item, true)
	if ok {
		return false
	} else {
		return true
	}
}

// 判断是否存在，如果不存在，存储并返回true，否则返回false
func (s *Set) AddBid(item int) bool {
	_, ok := s.bid.LoadOrStore(item, true)
	if ok {
		return false
	} else {
		return true
	}
}


// 判断是否存在
func (s *Set) Exist(item int) bool {
	_, exist := s.m.Load(item)
	return exist
}


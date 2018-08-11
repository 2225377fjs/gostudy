package lib

import "sync"

type Set struct {
	m map[int]string
	bid map[int]bool
	lock sync.Mutex
}

func NewSet() *Set{
	set := &Set{
		m: map[int]string{},
		bid: map[int]bool{},
	}
	return set
}

// 判断是否存在，如果不存在，存储并返回true，否则返回false
func (s *Set) Add(item int, reason string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, exist := s.m[item]; if exist {
		return false
	} else {
		s.m[item] = reason
		return true
	}
}

// 判断是否存在，如果不存在，存储并返回true，否则返回false
func (s *Set) AddBid(item int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, exist := s.bid[item]; if exist {
		return false
	} else {
		s.bid[item] = true
		return true
	}
}


// 判断是否存在
func (s *Set) Exist(item int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, exist := s.m[item]
	return exist
}

func (s *Set) GetReason(item int) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	reason, _ := s.m[item]
	return reason
}


// 判断是否存在
func (s *Set) BidExist(item int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, exist := s.bid[item]
	return exist
}
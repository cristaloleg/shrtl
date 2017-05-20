package main

import (
	"strconv"
	"sync"
	"time"
)

type storage struct {
	sync.RWMutex
	data map[string]string
}

func newStorage() storage {
	return storage{
		data: make(map[string]string),
	}
}

func (s *storage) Add(url string) (name string) {
	name = generateName()
	s.Lock()
	s.data[name] = url
	s.Unlock()
	return name
}

func (s *storage) Get(name string) (url string, ok bool) {
	s.RLock()
	url, ok = s.data[name]
	s.RUnlock()
	return url, ok
}

func generateName() string {
	cputime := time.Now().Unix()
	return strconv.FormatInt(cputime, 36)
}

package main

type stats struct {
	data map[string]uint64
	ch   chan string
}

func newStats() stats {
	s := stats{
		data: make(map[string]uint64),
		ch:   make(chan string, 1024),
	}

	go s.listen()

	return s
}

func (s *stats) Add(name string) {
	s.data[name] = 0
}

func (s *stats) Inc(name string) {
	s.ch <- name
}

func (s *stats) listen() {
	for {
		select {
		case name := <-s.ch:
			s.data[name]++
		}
	}
}

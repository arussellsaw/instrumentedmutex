package instrumentedmutex

import (
	"sync"
	"time"
)

type Mutex struct {
	sync.Mutex

	Record  func(time.Duration)
	Sampler func() bool
}

var timeNow = time.Now

func (m *Mutex) Lock() {
	if !m.Sampler() {
		m.Mutex.Lock()
		return
	}
	c := make(chan struct{})
	go func() {
		start := timeNow()
		<-c
		m.Record(time.Since(start))
	}()
	m.Mutex.Lock()
	c <- struct{}{}
	return
}

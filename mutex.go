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
	if m.Sampler == nil || !m.Sampler() {
		m.Mutex.Lock()
		return
	}
	c := make(chan struct{}, 1)
	go func() {
		start := timeNow()
		<-c
		if m.Record != nil {
			m.Record(timeNow().Sub(start))
		}
	}()
	m.Mutex.Lock()
	c <- struct{}{}
	return
}

type RWMutex struct {
	sync.RWMutex

	RecordRead func(time.Duration)
	Record     func(time.Duration)
	Sampler    func() bool
}

func (m *RWMutex) Lock() {
	if m.Sampler == nil || !m.Sampler() {
		m.RWMutex.Lock()
		return
	}
	c := make(chan struct{}, 1)
	go func() {
		start := timeNow()
		<-c
		if m.Record != nil {
			m.Record(timeNow().Sub(start))
		}
	}()
	m.RWMutex.Lock()
	c <- struct{}{}
	return
}

func (m *RWMutex) RLock() {
	if m.Sampler == nil || !m.Sampler() {
		m.RWMutex.RLock()
		return
	}
	c := make(chan struct{}, 1)
	go func() {
		start := timeNow()
		<-c
		if m.RecordRead != nil {
			m.RecordRead(timeNow().Sub(start))
		}
	}()
	m.RWMutex.RLock()
	c <- struct{}{}
	return
}

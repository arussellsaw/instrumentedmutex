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
	sync.Mutex

	RecordRead func(time.Duration)
	Record     func(time.Duration)
	Sampler    func() bool
}

func (m *RMutex) Lock() {
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

func (m *RMutex) RLock() {
	if m.Sampler == nil || !m.Sampler() {
		m.Mutex.RLock()
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
	m.Mutex.RLock()
	c <- struct{}{}
	return
}

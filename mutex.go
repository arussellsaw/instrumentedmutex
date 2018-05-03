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
	start := timeNow()
	m.Mutex.Lock()
	if m.Record != nil {
		m.Record(timeNow().Sub(start))
	}
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
	start := timeNow()
	m.RWMutex.Lock()
	if m.Record != nil {
		m.Record(timeNow().Sub(start))
	}
}

func (m *RWMutex) RLock() {
	if m.Sampler == nil || !m.Sampler() {
		m.RWMutex.RLock()
		return
	}
	start := timeNow()
	m.RWMutex.RLock()
	if m.RecordRead != nil {
		m.RecordRead(timeNow().Sub(start))
	}
}

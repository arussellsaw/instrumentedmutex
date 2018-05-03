package instrumentedmutex

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var timeNow = time.Now

// Mutex satisfies the same interface as sync.Mutex with the ability to record time taken
// to acquire the lock
type Mutex struct {
	sync.Mutex

	Record  func(time.Duration)
	Sampler func() bool
}

// Lock the mutex, recording the time taken if Sampler returns true
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

// RWMutex satisfies the same interface as sync.RWMutex with the ability to record time taken
// to acquire the lock
type RWMutex struct {
	sync.RWMutex

	RecordRead func(time.Duration)
	Record     func(time.Duration)
	Sampler    func() bool
}

// Lock the mutex, recording the time taken if Sampler returns true
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

// RLock the mutex, recording the time taken if Sampler returns true
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

// NewSampler returns a Sampler func that will return true `n` times in `in` calls
func NewSampler(n, in int) func() bool {
	s := make([]bool, in)
	for i := range s {
		if i < n {
			s[i] = true
		}
	}

	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	x := uint64(0)
	return func() bool {
		m := atomic.AddUint64(&x, uint64(1))
		i := m % uint64(len(s))
		v := s[i]
		return v
	}
}

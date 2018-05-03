package instrumentedmutex

import (
	"sync"
	"testing"
	"time"
)

func TestNoSamplerDoesntRecord(t *testing.T) {
	m := Mutex{
		Record: func(d time.Duration) {
			t.Errorf("did not expect recordfn to be called")
		},
	}
	m.Lock()
	m.Unlock()
}

func TestSamplerRecords(t *testing.T) {
	var (
		called bool
	)
	m := Mutex{
		Sampler: func() bool { return true },
		Record: func(d time.Duration) {
			called = true
			if d != 10*time.Second {
				t.Errorf("expected 10s, got %s", d)
			}
		},
	}
	now := time.Now()
	i := 0
	timeNow = func() time.Time {
		t := now.Add(time.Duration(i) * 10 * time.Second)
		i++
		return t
	}
	m.Mutex.Lock() // lock underlying mutex to avoid recording this time
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.Lock()
		m.Unlock()
		if !called {
			t.Errorf("expected Record func to be called")
		}
		wg.Done()
	}()
	m.Unlock()
	wg.Wait()
}

func TestRWNoSamplerDoesntRecord(t *testing.T) {
	m := RWMutex{
		Record: func(d time.Duration) {
			t.Errorf("did not expect recordfn to be called")
		},
		RecordRead: func(d time.Duration) {
			t.Errorf("did not expect recordfn to be called")
		},
	}
	m.Lock()
	m.Unlock()
	m.RLock()
	m.RUnlock()
}

func TestRWSamplerRecords(t *testing.T) {
	var (
		called, rCalled bool
	)

	m := RWMutex{
		Sampler: func() bool { return true },
		Record: func(d time.Duration) {
			called = true
			if d != 10*time.Second {
				t.Errorf("expected 10s, got %s", d)
			}
		},
		RecordRead: func(d time.Duration) {
			rCalled = true
			if d != 10*time.Second {
				t.Errorf("expected 10s, got %s", d)
			}
		},
	}
	now := time.Now()
	i := 0
	timeNow = func() time.Time {
		t := now.Add(time.Duration(i) * 10 * time.Second)
		i++
		return t
	}
	m.RWMutex.Lock() // lock underlying mutex to avoid recording this time
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.Lock()
		m.Unlock()
		if !called {
			t.Errorf("expected Record func to be called")
		}
		if rCalled {
			t.Errorf("expected RecordRead func not to be called")
		}
		wg.Done()
	}()
	m.Unlock()
	wg.Wait()
	called = false
	m.RWMutex.Lock() // lock underlying mutex to avoid recording this time
	wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.RLock()
		m.RUnlock()
		if !rCalled {
			t.Errorf("expected RecordRead func to be called")
		}
		if called {
			t.Errorf("expected Record func not to be called")
		}
		wg.Done()
	}()
	m.Unlock()
	wg.Wait()
}

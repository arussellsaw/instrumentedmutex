# Instrumented Mutex

[![GoDoc](https://godoc.org/github.com/arussellsaw/instrumentedmutex?status.svg)](http://godoc.org/github.com/arussellsaw/instrumentedmutex)

this package provides a way to record time taken to acquire a mutex lock, whilst still satisfying the same interface as sync.Mutex/RWMutex. The Sampler func is there so that you can configure a rate at which you wish to record timings, so that you can choose how much perf you're willing to sacrifice. Calls that aren't sampled should be just as fast as a regular mutex lock, but ones that are sampled do incur overhead.

example implementation, samples mutex wait timing on 1% of calls to Lock() or RLock():

```go
mutexTimer := instrumentation.GetTimer("targetter.mutex.write")
mutexReadTimer := instrumentation.GetTimer("targetter.mutex.read")

t.mutex = instrumentedmutex.RWMutex{
	Record: func(d time.Duration) {
		mutexTimer.RecordTime(d)
	},
	RecordRead: func(d time.Duration) {
		mutexReadTimer.RecordTime(d)
	},
	Sampler: instrumentedmutex.NewSampler(1, 100),// create a random sampler, sampler.IsTrue() will return true once in 100 calls
}
```

# Instrumented Mutex

this package provides a way to record time taken to acquire a mutex lock, the Sampler func is there so that you can configure a rate at which you wish to record timings, so that you can choose how much perf you're willing to sacrifice. Calls that aren't sampled should be just as fast as a regular mutex lock, but ones that are sampled do incur overhead.

i haven't provided a sampler, but they are easy enough to implement and i figured i'd err on the side of making as lightweight package as possible.

example implementation:

```go
mutexTimer := instrumentation.GetTimer("targetter.mutex.write")
mutexReadTimer := instrumentation.GetTimer("targetter.mutex.read")

sampler := sampler.New(1, 100) // create a random sampler, sampler.IsTrue() will return true once in 100 calls
t.mutex = instrumentedmutex.RWMutex{
	Record: func(d time.Duration) {
		mutexTimer.RecordTime(d)
	},
	RecordRead: func(d time.Duration) {
		mutexReadTimer.RecordTime(d)
	},
	Sampler: sampler.IsTrue,
}
```

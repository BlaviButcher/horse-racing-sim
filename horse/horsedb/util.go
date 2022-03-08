package horsedb

import "sync"

// Badger has some operations that will cause all other writes to fail. This provides a mechanism for
// all writing routines to wait until all other writes have completed and track the number of current write
// happening.

type WriteMu struct {
	mu sync.RWMutex
}

// StartWriting should be called when a normal update is happening to the database.
// doneWriting() must be deferred to say the update is done.
func (m *WriteMu) StartWriting() { m.mu.RLock() }

// DoneWriting tells the mutex the routine is done its writes. It must be used in conjunction with
// startWriting.
func (m *WriteMu) DoneWriting() { m.mu.RUnlock() }

// StartStopWorldWrites locks all writes and blocks any new writes from occurring until
// doneStopWorldWrites is called.
func (m *WriteMu) StartStopWorldWrites() { m.mu.Lock() }

// DoneStopWorldWrites unlocks a stopWorld and allows new writes to occur.
func (m *WriteMu) DoneStopWorldWrites() { m.mu.Unlock() }

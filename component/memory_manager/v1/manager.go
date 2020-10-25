package v1

import "sync"

type MemoryManager struct {
	mu          *sync.Mutex
	totalMemory int64
	usedMemory  int64
}

func NewMemoryManager(totalMemory int64) *MemoryManager {
	return &MemoryManager{
		mu:          new(sync.Mutex),
		totalMemory: totalMemory,
		usedMemory:  0,
	}
}

func (mm *MemoryManager) Malloc(n int64) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if n > mm.totalMemory-mm.usedMemory {
		return false
	}
	mm.usedMemory += n
	return true
}

func (mm *MemoryManager) Free(n int64) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	if mm.usedMemory < n {
		return false
	}
	mm.usedMemory -= n
	return true
}

func (mm *MemoryManager) Left() int64 {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return mm.totalMemory - mm.usedMemory
}

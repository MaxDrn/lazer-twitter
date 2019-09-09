package logging

import (
	"sync"
)

type EntryPool interface {
	Allocate() *Entry
	Release(*Entry)
}

func NewEntryPool() *entryPool {
	return &entryPool{
		mu: new(sync.Mutex),
	}
}

type entryPool struct {
	mu    *sync.Mutex
	slotA *Entry
	slotB *Entry
}

func (b *entryPool) Allocate() *Entry {
	b.mu.Lock()
	c := b.slotA
	if c != nil {
		b.slotA = b.slotB
		b.slotB = nil
	}
	b.mu.Unlock()
	if c == nil {
		c = new(Entry)
	}

	return c
}

func (b *entryPool) Release(c *Entry) {
	b.mu.Lock()
	b.slotB = b.slotA
	b.slotA = c
	b.mu.Unlock()
}

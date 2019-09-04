package registry

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/fid-dev/go-pflog/pkg/logging"
)

var ErrNotRegistered = errors.New("ErrNotRegistered")

func NewBucket() *Bucket {
	b := new(Bucket)
	b.mu = new(sync.RWMutex)
	b.reg = make(map[string]reflect.Type)

	return b
}

type Bucket struct {
	mu  *sync.RWMutex
	reg map[string]reflect.Type
}

func (b *Bucket) Register(containers ...logging.Container) {
	b.mu.Lock()

	for _, c := range containers {
		if _, ok := b.reg[string(c.Kind())]; ok {
			panic(fmt.Sprintf("duplicate registration of plfog container kind: %s", c.Kind()))
		}
		b.reg[string(c.Kind())] = reflect.TypeOf(c)
	}

	b.mu.Unlock()
}

func (b *Bucket) Lookup(kind string) (logging.Container, error) {
	b.mu.RLock()

	if t, ok := b.reg[kind]; ok {
		b.mu.RUnlock()
		return reflect.New(t).Interface().(logging.Container), nil
	}

	b.mu.RUnlock()
	return nil, ErrNotRegistered
}

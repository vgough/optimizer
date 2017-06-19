package optimizer

// #include <stdbool.h>
import "C"

import "sync"

type callbackFN func(p *C.double, r *C.double) C.bool

// registry holds a mapping of integer to a runtime-registered
// callback function.
type registry struct {
	mu    sync.Mutex
	index int
	fns   map[int]callbackFN
}

var callbackRegistry = &registry{
	fns: make(map[int]callbackFN),
}

func (r *registry) register(f callbackFN) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.index++
	r.fns[r.index] = f
	return r.index
}

func (r *registry) deregister(i int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.fns, i)
}

func (r *registry) lookup(index C.int) callbackFN {
	r.mu.Lock()
	f := r.fns[int(index)]
	r.mu.Unlock()

	return f
}

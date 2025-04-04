/*
 * Copyright The Alfred.go Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package function

import (
	"alfred/internal/log"
	"context"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

// VMPool manages a pool of Goja VMs
type VMPool struct {
	pool     chan *goja.Runtime
	minSize  int
	maxSize  int
	mutex    sync.Mutex
	current  int
	stopChan chan struct{} // Channel to stop cleanup goroutine
}

var (
	globalPool *VMPool
	once       sync.Once
)

// initializePool creates a new VM pool with the specified size
func initializePool(minSize, maxSize int) *VMPool {
	pool := &VMPool{
		pool:     make(chan *goja.Runtime, maxSize),
		minSize:  minSize,
		maxSize:  maxSize,
		current:  minSize,
		stopChan: make(chan struct{}),
	}

	// Initialize the pool with minimum number of VMs
	for i := 0; i < minSize; i++ {
		vm := createVM()
		pool.pool <- vm
	}

	// Start cleanup routine
	go pool.cleanup()

	return pool
}

// acquireVM gets a VM from the pool or creates a new one if needed
func (p *VMPool) acquireVM() *goja.Runtime {
	select {
	case vm := <-p.pool:
		return vm
	default:
		// No VM available in pool, try to create new one
		p.mutex.Lock()
		if p.current < p.maxSize {
			p.current++
			p.mutex.Unlock()
			return createVM()
		}
		p.mutex.Unlock()
		// If we've reached maxSize, wait for an available VM
		return <-p.pool
	}
}

// GetPool returns the global VM pool instance
func GetPool() *VMPool {
	once.Do(func() {
		globalPool = initializePool(1, 1000) // Min 1, Max 1000 VMs
	})
	return globalPool
}

// releaseVM returns a VM to the pool or discards it if pool is full
func (p *VMPool) releaseVM(vm *goja.Runtime) {
	select {
	case p.pool <- vm:
		// VM successfully returned to pool
	default:
		// Pool is full, discard the VM and decrease counter
		p.mutex.Lock()
		p.current--
		p.mutex.Unlock()
	}
}

// cleanup periodically removes excess VMs
func (p *VMPool) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.mutex.Lock()
			excess := p.current - p.minSize
			if excess > 0 {
				// Try to remove excess VMs
				for i := 0; i < excess; i++ {
					select {
					case <-p.pool:
						p.current--
					default:
						// No more VMs to remove
						break
					}
				}
			}
			p.mutex.Unlock()
		case <-p.stopChan:
			return
		}
	}
}

// createVM creates a new Goja VM instance
func createVM() *goja.Runtime {
	vm := goja.New()
	new(require.Registry).Enable(vm)
	console.Enable(vm)
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	//Add global vars
	globalVars := getDbWrapper()
	err := addVarsToGogaVmGlobalConext(vm, globalVars...)
	if err != nil {
		log.Fatal(context.Background(), "fail to add global vars to Goja vm global context", err)
	}

	return vm
}

// Shutdown gracefully stops the pool and cleanup routine
func (p *VMPool) Shutdown() {
	close(p.stopChan)

	// Clear the pool
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Drain the pool
	for len(p.pool) > 0 {
		<-p.pool
	}
	p.current = 0
}

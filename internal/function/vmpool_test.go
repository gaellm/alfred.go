package function

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitializePool(t *testing.T) {
	// Initialize a pool with minSize = 2 and maxSize = 5
	pool := initializePool(2, 5)

	// Ensure the pool is initialized with the correct size
	assert.Equal(t, 2, len(pool.pool))
	assert.Equal(t, 2, pool.current)
	assert.Equal(t, 2, pool.minSize)
	assert.Equal(t, 5, pool.maxSize)

	// Shutdown the pool to clean up
	pool.Shutdown()
}

func TestAcquireVM(t *testing.T) {
	// Initialize a pool with minSize = 1 and maxSize = 3
	pool := initializePool(1, 3)

	// Acquire a VM from the pool
	vm := pool.acquireVM()
	assert.NotNil(t, vm)

	// Ensure the pool size decreases after acquiring a VM
	assert.Equal(t, 0, len(pool.pool))
	assert.Equal(t, 1, pool.current)

	// Acquire another VM, which should create a new one
	vm2 := pool.acquireVM()
	assert.NotNil(t, vm2)
	assert.Equal(t, 2, pool.current)

	// Shutdown the pool to clean up
	pool.Shutdown()
}

func TestReleaseVM(t *testing.T) {
	// Initialize a pool with minSize = 1 and maxSize = 3
	pool := initializePool(1, 3)

	// Acquire a VM and then release it back to the pool
	vm := pool.acquireVM()
	pool.releaseVM(vm)

	// Ensure the pool size not increases after releasing the VM
	assert.Equal(t, 1, len(pool.pool))
	assert.Equal(t, 1, pool.current)

	// Acquire two more VMs to reach maxSize
	vm1 := pool.acquireVM()
	vm2 := pool.acquireVM()
	vm3 := pool.acquireVM()

	assert.Equal(t, 3, pool.current)

	// Release a VM when the pool is full
	pool.releaseVM(vm1)
	assert.Equal(t, 1, len(pool.pool)) // only vm1 in the available vms
	assert.Equal(t, 3, pool.current)   // Pool should not exceed maxSize

	pool.releaseVM(vm3)
	pool.releaseVM(vm2)
	assert.Equal(t, 3, len(pool.pool)) // all vm available in the pool

	// Shutdown the pool to clean up
	pool.Shutdown()
}

func TestPoolRespectsMaxSize(t *testing.T) {
	// Initialize a pool with minSize = 1 and maxSize = 2
	pool := initializePool(1, 2)

	// Acquire two VMs, reaching maxSize
	vm1 := pool.acquireVM()
	vm2 := pool.acquireVM()

	assert.Equal(t, 2, pool.current)

	// Try to acquire a third VM, which should block until one is released
	go func() {
		time.Sleep(100 * time.Millisecond)
		pool.releaseVM(vm1)
	}()

	start := time.Now()
	vm3 := pool.acquireVM()
	dur := time.Since(start)

	assert.NotNil(t, vm3)
	assert.GreaterOrEqual(t, dur.Milliseconds(), time.Duration.Milliseconds(80)) // Ensure it waited
	assert.Equal(t, 2, pool.current)

	pool.releaseVM(vm2)

	// Shutdown the pool to clean up
	pool.Shutdown()
}

func TestCleanupRoutine(t *testing.T) {
	// Initialize a pool with minSize = 1 and maxSize = 3
	pool := initializePool(1, 3)

	pool.cleanupFreq = 3 * time.Second

	// Add two extra VMs to the pool
	vm1 := pool.acquireVM()
	vm2 := pool.acquireVM()
	pool.releaseVM(vm1)
	pool.releaseVM(vm2)

	// Ensure the pool has 2 VMs
	assert.Equal(t, 2, pool.current)

	// Wait for the cleanup routine to run
	time.Sleep(5 * time.Second) // Cleanup runs every 5 minutes

	// Ensure the pool size is reduced to minSize
	assert.Equal(t, 1, len(pool.pool))
	assert.Equal(t, 1, pool.current)

	// Shutdown the pool to clean up
	pool.Shutdown()
}

func TestShutdown(t *testing.T) {
	// Initialize a pool with minSize = 2 and maxSize = 5
	pool := initializePool(2, 5)

	// Shutdown the pool
	pool.Shutdown()

	// Ensure the pool is empty and the current count is 0
	assert.Equal(t, 0, len(pool.pool))
	assert.Equal(t, 0, pool.current)
}

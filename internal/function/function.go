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
	"alfred/internal/db"
	"alfred/internal/helper"
	"alfred/internal/log"
	"alfred/internal/mock"
	"alfred/pkg/request"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/google/uuid"
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

const FUNC_UPDATE_HELPERS = "updateHelpers"
const FUNC_ALFRED = "alfred"
const FUNC_SETUP = "setup"

type Function struct {
	FileName             string
	FileContent          string
	HasFuncUpdateHelpers bool
	HasFuncAlfred        bool
	HasFuncSetup         bool
}

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

// GetPool returns the global VM pool instance
func GetPool() *VMPool {
	once.Do(func() {
		globalPool = initializePool(1, 1000) // Min 1, Max 1000 VMs
	})
	return globalPool
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

	// Generate a random UUID-based path for the database
	randomUUID := uuid.New().String()
	dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("alfred_badger_%s", randomUUID))

	// Initialize the database manager
	dbManager := db.GetDBManager()
	err := dbManager.Init(dbPath)
	if err != nil {
		panic("Failed to initialize Badger DB: " + err.Error())
	}

	// Expose database functions to JavaScript
	if err := vm.Set("dbSet", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.ToValue("dbSet requires a key and a value as arguments"))
		}
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		err := dbManager.Set(key, value)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}
		return goja.Undefined()
	}); err != nil {
		log.Warn(context.Background(), "failed to set dbSet function in vm:", err)
	}

	if err := vm.Set("dbGet", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("dbGet requires a key as an argument"))
		}
		key := call.Argument(0).String()
		value, err := dbManager.Get(key)
		if err != nil {
			// Return undefined if the key is not found
			return goja.Undefined()
		}
		return vm.ToValue(value)
	}); err != nil {
		log.Warn(context.Background(), "failed to set dbGet function in vm:", err)
	}

	if err := vm.Set("dbDelete", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("dbDelete requires a key as an argument"))
		}
		key := call.Argument(0).String()
		err := dbManager.Delete(key)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}
		return goja.Undefined()
	}); err != nil {
		log.Warn(context.Background(), "failed to set dbDelete function in vm:", err)
	}

	// Add the dbLoadFile function with batch writes
	if err := vm.Set("dbLoadFile", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("dbLoadFile requires a file path as an argument"))
		}
		filePath := call.Argument(0).String()

		// Read the file contents using os.ReadFile
		content, err := os.ReadFile(filePath)
		if err != nil {
			panic(vm.ToValue("Failed to read file: " + err.Error()))
		}

		// Parse the JSON into a map
		var data map[string]string
		err = json.Unmarshal(content, &data)
		if err != nil {
			panic(vm.ToValue("Failed to parse JSON file: " + err.Error()))
		}

		// Use a batch write for faster imports
		writeBatch := dbManager.GetDB().NewWriteBatch() // Get the Badger DB instance
		defer writeBatch.Cancel()                       // Ensure the batch is canceled if something goes wrong

		for key, value := range data {
			err = writeBatch.Set([]byte(key), []byte(value))
			if err != nil {
				panic(vm.ToValue("Failed to add key-value pair to batch: " + err.Error()))
			}
		}

		// Commit the batch
		err = writeBatch.Flush()
		if err != nil {
			panic(vm.ToValue("Failed to commit batch: " + err.Error()))
		}

		return vm.ToValue("File loaded successfully using batch writes")
	}); err != nil {
		panic("Failed to set dbLoadFile function in VM: " + err.Error())
	}

	return vm
}

func CreateFunction(fileName string, fileContent []byte) (Function, error) {

	var err error
	f := Function{fileName, string(fileContent), false, false, false}

	//CheckIfSetupFuncExists
	f.HasFuncSetup, err = f.CheckIfFuncExists(FUNC_SETUP)
	if err != nil {
		return f, err
	}

	f.HasFuncAlfred, err = f.CheckIfFuncExists(FUNC_ALFRED)
	if err != nil {
		return f, err
	}

	f.HasFuncUpdateHelpers, err = f.CheckIfFuncExists(FUNC_UPDATE_HELPERS)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (f *Function) SetupFunc() error {

	if !f.HasFuncSetup {
		return errors.New("function file " + f.FileName + " not contains " + FUNC_SETUP + " function")
	}

	var setup func() error
	pool := GetPool()
	vm := pool.acquireVM()
	defer pool.releaseVM(vm)

	//load js functions in vm
	_, err := vm.RunString(f.FileContent)
	if err != nil {

		err = errors.New(f.FileName + ": " + err.Error())
		return err
	}

	err = vm.ExportTo(vm.Get(FUNC_SETUP), &setup)
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return err
	}

	err = setup()
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return err
	}

	return nil
}

func (f *Function) UpdateHelpersListener(helpers []helper.Helper) ([]helper.Helper, error) {

	if !f.HasFuncUpdateHelpers {
		return helpers, errors.New("function file " + f.FileName + " not contains " + FUNC_UPDATE_HELPERS + " function")
	}

	var updateHelpers func([]helper.Helper) ([]helper.Helper, error)

	pool := GetPool()
	vm := pool.acquireVM()
	defer pool.releaseVM(vm)

	//load js functions in vm
	_, err := vm.RunString(f.FileContent)
	if err != nil {

		err = errors.New(f.FileName + ": " + err.Error())
		return helpers, err
	}

	err = vm.ExportTo(vm.Get(FUNC_UPDATE_HELPERS), &updateHelpers)
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return helpers, err
	}

	updatedHelpers, err := updateHelpers(helpers)
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return helpers, err
	}

	return updatedHelpers, nil
}

func (f *Function) AlfredFunc(m mock.Mock, helpers []helper.Helper, req request.Req, res request.Res) (request.Res, error) {

	if !f.HasFuncAlfred {
		return res, errors.New("function file " + f.FileName + " not contains " + FUNC_ALFRED + " function")
	}

	var alfred func(mock.Mock, []helper.Helper, request.Req, request.Res) (request.Res, error)
	pool := GetPool()
	vm := pool.acquireVM()
	defer pool.releaseVM(vm)

	//load js functions in vm
	_, err := vm.RunString(f.FileContent)
	if err != nil {

		err = errors.New(f.FileName + ": " + err.Error())
		return res, err
	}

	err = vm.ExportTo(vm.Get(FUNC_ALFRED), &alfred)
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return res, err
	}

	resUpdated, err := alfred(m, helpers, req, res)
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return res, err
	}

	return resUpdated, nil
}

func (f *Function) CheckIfFuncExists(funcName string) (bool, error) {

	pool := GetPool()
	vm := pool.acquireVM()
	defer pool.releaseVM(vm)

	//load js functions in vm
	_, err := vm.RunString(f.FileContent)
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return false, err
	}

	v, err := vm.RunString("typeof " + funcName + " === 'function'")
	if err != nil {
		err = errors.New(f.FileName + ": " + err.Error())
		return false, err
	}

	return v.Export().(bool), nil

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

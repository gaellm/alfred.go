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
	"alfred/internal/helper"
	"alfred/internal/mock"
	"alfred/pkg/request"
	"errors"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

// VMPool manages a pool of Goja VMs
type VMPool struct {
	pool chan *goja.Runtime
	size int
}

var (
	globalPool *VMPool
	once       sync.Once
)

const FUNC_UPDATE_HELPERS = "updateHelpers"
const FUNC_ALFRED = "alfred"

type Function struct {
	FileName             string
	FileContent          string
	HasFuncUpdateHelpers bool
	HasFuncAlfred        bool
}

// initializePool creates a new VM pool with the specified size
func initializePool(size int) *VMPool {
	pool := &VMPool{
		pool: make(chan *goja.Runtime, size),
		size: size,
	}

	// Initialize the pool with VMs
	for i := 0; i < size; i++ {
		vm := createVM()
		pool.pool <- vm
	}

	return pool
}

// GetPool returns the global VM pool instance
func GetPool() *VMPool {
	once.Do(func() {
		globalPool = initializePool(10) // Adjust pool size as needed
	})
	return globalPool
}

// acquireVM gets a VM from the pool
func (p *VMPool) acquireVM() *goja.Runtime {
	return <-p.pool
}

// releaseVM returns a VM to the pool
func (p *VMPool) releaseVM(vm *goja.Runtime) {
	p.pool <- vm
}

// createVM creates a new Goja VM instance
func createVM() *goja.Runtime {
	vm := goja.New()
	new(require.Registry).Enable(vm)
	console.Enable(vm)
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	/*
		time.AfterFunc(timeout, func() {
			vm.Interrupt("halt")
		})*/

	return vm
}

func CreateFunction(fileName string, fileContent []byte) (Function, error) {

	var err error
	f := Function{fileName, string(fileContent), false, false}

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

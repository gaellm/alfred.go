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
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestAddVarsToGogaVmGlobalConext(t *testing.T) {
	// Create a new goja.Runtime instance
	vm := goja.New()

	// Define mock wrappers
	wrappers := []wrapper{
		{
			name: "mockFunction1",
			value: func(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
				return func(call goja.FunctionCall) goja.Value {
					return vm.ToValue("mockFunction1 called")
				}
			},
		},
		{
			name: "mockFunction2",
			value: func(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
				return func(call goja.FunctionCall) goja.Value {
					return vm.ToValue("mockFunction2 called")
				}
			},
		},
	}

	// Add the wrappers to the global context
	err := addVarsToGogaVmGlobalConext(vm, wrappers...)
	assert.NoError(t, err)

	// Test if the functions are correctly added to the global context
	result1, err := vm.RunString("mockFunction1()")
	assert.NoError(t, err)
	assert.Equal(t, "mockFunction1 called", result1.String())

	result2, err := vm.RunString("mockFunction2()")
	assert.NoError(t, err)
	assert.Equal(t, "mockFunction2 called", result2.String())

	// Test calling a non-existent function
	_, err = vm.RunString("nonExistentFunction()")
	assert.Error(t, err)
}

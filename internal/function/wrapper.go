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
	"github.com/dop251/goja"
)

type wrapper struct {
	name  string
	value func(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value
}

// Set the specified variable in the global context.
func addVarsToGogaVmGlobalConext(vm *goja.Runtime, wrappers ...wrapper) error {

	var err error

	for _, w := range wrappers {

		if err = vm.Set(w.name, w.value(vm)); err != nil {
			return err
		}

	}

	return err

}

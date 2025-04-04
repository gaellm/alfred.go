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

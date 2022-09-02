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
	"regexp"
	"testing"
)

func TestCreateFunctionCollectionFromFolder(t *testing.T) {

	collection, err := CreateFunctionCollectionFromFolder("../../user-files/functions/")
	if err != nil {

		t.Errorf("create collection from folder failed with error: " + err.Error())
	}

	if !collection.IsEmpty() {

		var testHelper helper.Helper
		testHelper.Name = "my-city"
		testHelper.Value = "Rennes"
		testHelper.Regex = regexp.MustCompile(".*date[(]([^)]*)[)].*")
		testHelper.AddPrivateParam("toto", "titi")

		testHelperArray := []helper.Helper{testHelper}

		for _, f := range collection {

			if f.HasFuncUpdateHelpers {

				updatedHelpers, err := f.UpdateHelpersListener(testHelperArray)
				if err != nil {
					t.Errorf(err.Error())
				}

				if updatedHelpers[0].Value != "Gotham" {
					t.Errorf("failed to update helpers from js : '" + updatedHelpers[0].Value + "' should be Gotham")
				}

				if updatedHelpers[0].Regex == nil {
					t.Errorf("regex lost")
				}

				if updatedHelpers[0].GetPrivateParam("toto") != "titi" {
					t.Errorf("private params lost")
				}
			}
		}

	} else {
		t.Logf("collection is empty")
	}

}

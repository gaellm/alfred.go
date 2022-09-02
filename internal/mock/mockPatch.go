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

package mock

import (
	"encoding/json"
)

func MockPatch(m *Mock, jsonStr []byte) error {

	err := json.Unmarshal(jsonStr, &m)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	patchedMock, err := BuildMockFromJson(jsonData)
	if err != nil {
		return err
	}

	MergeMock(m, patchedMock)

	return nil

}

func MergeMock(dest *Mock, src Mock) {

	dest.Request = src.Request
	dest.Response = src.Response
	dest.jsonBytes = src.jsonBytes
	dest.requestHelpers = src.requestHelpers
	dest.dateHelpers = src.dateHelpers
	dest.randomHelpers = src.randomHelpers
	dest.FunctionFile = src.FunctionFile
	dest.Actions = src.Actions

}

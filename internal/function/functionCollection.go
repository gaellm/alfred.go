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

import "errors"

type FunctionCollection []Function

func (c FunctionCollection) GetJs(fileName string) string {

	for _, f := range c {

		if f.FileName == fileName {
			return f.FileContent
		}
	}

	return ""
}

func (c FunctionCollection) IsEmpty() bool {
	return len(c) == 0
}

func (c FunctionCollection) GetFunction(fileName string) (Function, error) {

	for _, f := range c {

		if f.FileName == fileName {
			return f, nil
		}
	}

	return Function{}, errors.New("no fonction files named " + fileName)
}

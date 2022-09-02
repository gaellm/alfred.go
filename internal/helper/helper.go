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

package helper

import (
	"encoding/json"
	"regexp"
)

//helper types
const REQUEST = "req"
const DATE = "time"
const RANDOM = "random"

//helper params
const PARAM_NAME = "name"
const PARAM_REGEX = "regex"

//json tag used with users js functions
type Helper struct {
	Type          string
	String        string `json:"str"`
	Value         string `json:"value"`
	Target        string
	Name          string `json:"name"`
	Regex         *regexp.Regexp
	privateParams map[string]string
}

func (h *Helper) AddPrivateParam(paramName string, paramValue string) {

	if h.privateParams == nil {
		h.privateParams = map[string]string{}
	}

	h.privateParams[paramName] = paramValue
}

func (h Helper) GetPrivateParam(paramName string) string {

	return h.privateParams[paramName]
}

func (h Helper) GetPrivateParams() map[string]string {

	if h.privateParams == nil {
		h.privateParams = map[string]string{}
	}

	return h.privateParams
}

func (h Helper) HasValue() bool {

	return h.Value != ""
}

func (h Helper) GetJsonMarshal() string {

	jsonH, _ := json.Marshal(h)
	return string(jsonH)
}

func (h Helper) Clone() Helper {

	var clone Helper
	{
		clone.Type = h.Type
		clone.String = h.String
		clone.Value = h.Value
		clone.Target = h.Target
		clone.Name = h.Name
		clone.Regex = h.Regex
		clone.privateParams = h.privateParams
	}

	return clone
}

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
	"alfred/internal/helper"
)

type MockRequest struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type MockResponse struct {
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

type Mock struct {
	Name      string       `json:"name"`
	Request   MockRequest  `json:"request"`
	Response  MockResponse `json:"response"`
	jsonBytes []byte
	helpers   map[string][]helper.Helper
}

func (m *Mock) AddHelper(h helper.Helper) {

	if !m.HasHelper() {
		m.helpers = make(map[string][]helper.Helper)
	}

	if m.HasHelperType(h.Type) {
		m.helpers[h.Type] = append(m.helpers[h.Type], h)
		return
	}

	var helpers []helper.Helper
	helpers = append(helpers, h)

	m.helpers[h.Type] = helpers
}

func (m Mock) HasHelper() bool {

	return len(m.helpers) > 0
}

func (m *Mock) UpdateHelpers(h map[string][]helper.Helper) {

	m.helpers = h
}

func (m Mock) GetHelpers() map[string][]helper.Helper {

	return m.helpers
}

func (m Mock) GetJsonHelpers() string {

	var jsonHelpers string
	helpers := m.GetHelpers()

	for i := range helpers {
		for y := range helpers[i] {
			jsonHelpers += helpers[i][y].GetJsonMarshal()

		}
	}
	return jsonHelpers
}

func (m Mock) HasHelperType(helperType string) bool {

	return len(m.helpers[helperType]) > 0
}

func (m Mock) GetRequestMethod() string {

	if m.Request.Method != "" {
		return m.Request.Method
	}

	return "GET"
}

func (m *Mock) SaveJsonBytes(jsonData []byte) {

	m.jsonBytes = jsonData
}

func (m Mock) GetJsonBytes() []byte {

	return m.jsonBytes
}

func (m Mock) GetRequestUrl() string {
	return m.Request.Url
}

func (m Mock) GetName() string {

	if m.Name != "" {
		return m.Name
	}

	return m.GetRequestMethod() + "-" + m.GetRequestUrl()
}

func (m Mock) GetResponseStatus() int {
	return m.Response.Status
}

func (m Mock) GetResponseBody() string {
	return m.Response.Body
}

func (m Mock) GetResponseHeaders() map[string]string {
	return m.Response.Headers
}

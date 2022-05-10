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
	Name           string       `json:"name"`
	Request        MockRequest  `json:"request"`
	Response       MockResponse `json:"response"`
	jsonBytes      []byte
	requestHelpers []helper.Helper
}

func (m *Mock) AddRequestHelper(h helper.Helper) {

	m.requestHelpers = append(m.requestHelpers, h)
}

func (m Mock) HasRequestHelper() bool {

	return len(m.requestHelpers) > 0
}

func (m Mock) UpdateRequestHelpers(h []helper.Helper) Mock {

	m.requestHelpers = h
	return m
}

// An array of truct keep references to struct, so we have to clone
// the array before returning it.
func (m Mock) GetRequestHelpers() []helper.Helper {

	var clones []helper.Helper

	for _, h := range m.requestHelpers {

		var clone helper.Helper
		{
			clone.Type = h.Type
			clone.String = h.String
			clone.Value = h.Value
			clone.Target = h.Target
		}

		clones = append(clones, clone)
	}

	return clones
}

func (m Mock) GetJsonHelpers() string {

	var jsonHelpers string
	helpers := m.GetRequestHelpers()

	for i := range helpers {

		jsonHelpers += helpers[i].GetJsonMarshal()
	}
	return jsonHelpers
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

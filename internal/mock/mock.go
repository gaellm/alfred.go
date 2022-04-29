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
	Name     string       `json:"name"`
	Request  MockRequest  `json:"request"`
	Response MockResponse `json:"response"`
}

func (m Mock) GetRequestMethod() string {

	if m.Request.Method != "" {
		return m.Request.Method
	}

	return "GET"
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

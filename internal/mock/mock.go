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
	"encoding/json"
	"math/rand"
	"net/url"
	"regexp"
	"time"
)

type MockRequest struct {
	Method         string `json:"method"`
	Url            string `json:"url"`
	UrlRegexStr    string `json:"urlRegex"`
	UrlTransformed string

	//use to manage url helpers
	RegexUrl *regexp.Regexp
}

type MockResponse struct {
	Status          int               `json:"status"`
	Body            json.RawMessage   `json:"body"`
	BodyFile        string            `json:"body-file"`
	Headers         map[string]string `json:"headers"`
	MinResponseTime int               `json:"minResponseTime"`
	MaxResponseTime int               `json:"maxResponseTime"`
}

type MockAction struct {
	Type             string            `json:"type"`
	MinScheduledTime int               `json:"minScheduledTime"`
	MaxScheduledTime int               `json:"maxScheduledTime"`
	Method           string            `json:"method"`
	Url              string            `json:"url"`
	Timeout          string            `json:"timeout"`
	Body             string            `json:"body"`
	Headers          map[string]string `json:"headers"`
}

type Mock struct {
	Name             string       `json:"name"`
	Request          MockRequest  `json:"request"`
	Response         MockResponse `json:"response"`
	jsonBytes        []byte
	requestHelpers   []helper.Helper
	dateHelpers      []helper.Helper
	randomHelpers    []helper.Helper
	pathRegexHelpers []helper.Helper
	FunctionFile     string       `json:"function-file"`
	Actions          []MockAction `json:"actions"`
}

func (m *Mock) AddRequestHelper(h helper.Helper) {

	m.requestHelpers = append(m.requestHelpers, h)
}

func (m *Mock) SetRegexUrl() {
	r := regexp.MustCompile(m.Request.UrlRegexStr)
	m.Request.RegexUrl = r
	m.Request.UrlTransformed = "/" + url.QueryEscape(m.Request.UrlRegexStr)
}

func (m *Mock) AddPathRegexHelper(h helper.Helper) {

	m.pathRegexHelpers = append(m.pathRegexHelpers, h)
}

func (m *Mock) AddDatetHelper(h helper.Helper) {

	m.dateHelpers = append(m.dateHelpers, h)
}

func (m *Mock) AddRandomHelper(h helper.Helper) {

	m.randomHelpers = append(m.randomHelpers, h)
}

func (m Mock) HasRequestHelper() bool {

	return len(m.requestHelpers) > 0
}

func (m Mock) HasDatetHelper() bool {

	return len(m.dateHelpers) > 0
}

func (m Mock) HasPathRegexHelper() bool {

	return len(m.pathRegexHelpers) > 0
}

func (m Mock) HasRandomHelper() bool {

	return len(m.randomHelpers) > 0
}

func (m Mock) HasRegexUrl() bool {

	return len(m.Request.UrlRegexStr) > 0
}

func (m *Mock) HasFunctionFile() bool {

	return m.FunctionFile != ""
}

func (m *Mock) GetFunctionFile() string {

	return m.FunctionFile
}

func (m Mock) HasHelper() bool {

	return m.HasDatetHelper() || m.HasRequestHelper() || m.HasRandomHelper() || m.HasPathRegexHelper()
}

func (m Mock) UpdateRequestHelpers(h []helper.Helper) Mock {

	m.requestHelpers = h
	return m
}

func (m Mock) UpdatePathRegexHelpers(h []helper.Helper) Mock {

	m.pathRegexHelpers = h
	return m
}

func (m Mock) UpdateDateHelpers(h []helper.Helper) Mock {

	m.dateHelpers = h
	return m
}

func (m Mock) UpdateRandomHelpers(h []helper.Helper) Mock {

	m.randomHelpers = h
	return m
}

// An array of truct keep references to struct, so we have to clone
// the array before returning it.
func (m Mock) GetRequestHelpers() []helper.Helper {

	var clones []helper.Helper

	for _, h := range m.requestHelpers {
		clones = append(clones, h.Clone())
	}

	return clones
}

// An array of truct keep references to struct, so we have to clone
// the array before returning it.
func (m Mock) GetDateHelpers() []helper.Helper {

	var clones []helper.Helper

	for _, h := range m.dateHelpers {
		clones = append(clones, h.Clone())
	}

	return clones
}

func (m Mock) GetPathRegexHelpers() []helper.Helper {

	var clones []helper.Helper

	for _, h := range m.pathRegexHelpers {
		clones = append(clones, h.Clone())
	}

	return clones
}

// An array of truct keep references to struct, so we have to clone
// the array before returning it.
func (m Mock) GetRandomHelpers() []helper.Helper {

	var clones []helper.Helper

	for _, h := range m.randomHelpers {
		clones = append(clones, h.Clone())
	}

	return clones
}

func (m Mock) GetJsonHelpers() string {

	var jsonHelpers string
	helpers := append(
		append(
			append([]helper.Helper{}, m.GetRequestHelpers()...),
			m.GetDateHelpers()...,
		), m.GetRandomHelpers()...)

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

	if m.Request.UrlRegexStr != "" {
		return m.Request.UrlTransformed
	}

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

	// Declare a variable to hold the unmarshaled text
	var body string

	// Unmarshal the raw message into the string variable
	err := json.Unmarshal(m.Response.Body, &body)
	if err != nil {
		// Handle error
		return string(m.Response.Body)
	}
	return body
}

func (m Mock) GetResponseHeaders() map[string]string {
	return m.Response.Headers
}

func (m *Mock) GetDelay() time.Duration {

	if m.Response.MaxResponseTime <= m.Response.MinResponseTime {
		return time.Duration(m.Response.MinResponseTime) * time.Millisecond
	}

	return time.Duration(rand.Intn(m.Response.MaxResponseTime-m.Response.MinResponseTime)+m.Response.MinResponseTime) * time.Millisecond
}

func (m Mock) GetActions() []MockAction {

	return m.Actions
}

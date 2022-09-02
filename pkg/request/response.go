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

package request

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string
	Body    string
	Headers map[string]string
}

type Res struct {
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

func (r *Res) SetHeader(key string, value string) {

	if r.Headers == nil {
		r.Headers = map[string]string{}
	}

	r.Headers[key] = value

}

func (r *Res) Stringify() string {

	jsonStr, _ := json.Marshal(r)

	return string(jsonStr)
}

func (resp *Response) SetHeaders(headers http.Header) {

	resp.Headers = make(map[string]string)

	for k, v := range headers {

		resp.Headers[k] = v[0]

	}
}

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
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestSendRequest(t *testing.T) {

	var r Request

	err := r.SetUrl("https://httpbin.org/post?test=toto")
	if err != nil {
		t.Errorf("TestSendRequest return an error: %v", err)
	}

	err = r.SetTimeout("6s")
	if err != nil {
		t.Errorf("TestSendRequest return an error: %v", err)
	}

	values := map[string]string{"foo": "baz"}
	jsonData, _ := json.Marshal(values)
	r.Body = jsonData

	err = r.SetMethod("POST")
	if err != nil {
		t.Errorf("TestSendRequest return an error: %v", err)
	}

	queryArgs := make(map[string]string)
	queryArgs["titi"] = "tutu"
	r.AddQueryArgs(queryArgs)

	response, err := r.Send(context.Background())
	if err != nil {
		t.Errorf("TestSendRequest return an error: %v", err)
	}

	if response.Body == "" {
		t.Errorf("TestSendRequest return an empty response body")
	}
	t.Logf("%v", response.Body)

	if len(response.Headers) < 1 {
		t.Errorf("TestSendRequest return a bad response header")
	}
	t.Logf("%v", fmt.Sprint(response.Headers))

	if response.Status == "" {
		t.Errorf("TestSendRequest return a bad response status")
	}
	t.Logf("%v", response.Status)

	err = r.SetUrl("https://httpbin.org/get?test=toto")
	if err != nil {
		t.Errorf("TestSendRequest return an error: %v", err)
	}

	r.method = "GET"
	getResponse, err := r.Send(context.Background())
	if err != nil {
		t.Errorf("TestSendRequest return an error: %v", err)
	}
	t.Logf("%v", getResponse.Body)
}

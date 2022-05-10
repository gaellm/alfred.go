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
	"errors"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/tidwall/gjson"
)

func RequestHelperWatcher(data []byte, contentType string, h []Helper) ([]Helper, error) {

	println("RequestHelperWatcher h ref:")
	println(&h)

	//JSON
	if strings.Contains(contentType, "json") {

		helpers, err := jsonWatcher(data, h)
		if err != nil {
			return h, err
		}

		return helpers, nil

		//XML
	} else if strings.Contains(contentType, "xml") {

		xml := strings.NewReader(string(data))

		//https://github.com/basgys/goxml2json
		json, err := xj.Convert(xml)
		if err != nil {
			return h, err
		}

		helpers, err := jsonWatcher(json.Bytes(), h)
		if err != nil {
			return h, err
		}

		println("RequestHelperWatcher helpers ref:")
		println(&helpers)

		return helpers, nil
	}

	return h, errors.New("content type unknown")
}

func jsonWatcher(d []byte, h []Helper) ([]Helper, error) {

	println("jsonWatcher h ref:")
	println(&h)

	for i, helper := range h {

		if helper.Value != "" {
			continue
		}

		h[i].Value = gjson.Get(string(d), helper.Target).String()
	}

	return h, nil
}

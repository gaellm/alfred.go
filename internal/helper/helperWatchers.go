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
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

func RequestHelperWatcher(data []byte, contentType string, h map[string][]Helper) (map[string][]Helper, error) {

	if strings.Contains(contentType, "json") {

		helpers, err := jsonWatcher(data, h)
		if err != nil {
			return h, err
		}
		return helpers, nil
	}

	return h, errors.New("content type unknown")

}

func jsonWatcher(d []byte, h map[string][]Helper) (map[string][]Helper, error) {

	for i, helper := range h[REQUEST] {

		if helper.Value != "" {
			continue
		}

		// to improve
		r := regexp.MustCompile(`alfred\.req\.(.*)`)
		helperTarget := r.FindStringSubmatch(helper.Target)[1]

		h[REQUEST][i].Value = gjson.Get(string(d), helperTarget).String()

	}

	return h, nil
}

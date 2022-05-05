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
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func RequestHelperWatcher(c *gin.Context, h map[string][]Helper) (map[string][]Helper, error) {

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return h, err
	}

	switch c.ContentType() {
	case "application/json":
		helpers, err := jsonWatcher(data)
		if err != nil {
			return h, err
		}
		return helpers, nil
	}

	return h, errors.New("content type unknown")

}

func jsonWatcher(d []byte) (map[string][]Helper, error) {

	return make(map[string][]Helper), errors.New("TODO")

}

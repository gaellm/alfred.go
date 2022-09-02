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
	"net/http"
	"net/url"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func RequestHelperWatcher(data []byte, c *gin.Context, h []Helper) ([]Helper, error) {

	var err error

	//Watch HTTP query
	h = queryWatcher(c.Request.URL.Query(), h)

	//Watch URL params
	h = paramWatcher(c.Params, h)

	//Watch HTTP body
	h, err = bodyWatcher(data, c, h)
	if err != nil {
		return h, err
	}

	//Watch HTTP headers
	h = headersWatcher(c.Request.Header, h)

	return h, err
}

func jsonWatcher(d []byte, h []Helper) ([]Helper, error) {

	for i, helper := range h {

		if helper.Value != "" {
			continue
		}

		h[i].Value = gjson.Get(string(d), helper.Target).String()
	}

	return h, nil
}

func textWatcher(d []byte, h []Helper) ([]Helper, error) {

	for i, helper := range h {

		if helper.Value != "" || helper.Regex == nil {
			continue
		}

		subMatch := helper.Regex.FindSubmatch(d)
		if len(subMatch) > 1 {
			h[i].Value = string(subMatch[1])
		}

	}

	return h, nil
}

func queryWatcher(query url.Values, h []Helper) []Helper {

	if len(query) > 0 {

		for i, helper := range h {

			if helper.Value != "" {
				continue
			}

			h[i].Value = query.Get(helper.Target)
		}
	}

	return h
}

func paramWatcher(params gin.Params, h []Helper) []Helper {

	if len(params) > 0 {

		for i, helper := range h {

			if helper.Value != "" {
				continue
			}

			h[i].Value = params.ByName(helper.Target)
		}
	}

	return h
}

func bodyWatcher(data []byte, c *gin.Context, h []Helper) ([]Helper, error) {

	//JSON
	if strings.Contains(c.ContentType(), "json") {

		helpers, err := jsonWatcher(data, h)
		if err != nil {
			return h, err
		}

		return helpers, nil

		//XML
	} else if strings.Contains(c.ContentType(), "xml") {

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

		return helpers, nil
	} else if strings.Contains(c.ContentType(), "text") {

		helpers, err := textWatcher(data, h)
		if err != nil {
			return h, err
		}

		return helpers, nil
	}

	return h, errors.New("content type unknown")
}

func headersWatcher(headers http.Header, h []Helper) []Helper {

	if len(headers) > 0 {

		for i, helper := range h {

			if helper.Value != "" {
				continue
			}

			//iterate headers and break if found
			for header, headerArray := range headers {

				if strings.EqualFold(header, helper.Target) {

					//get one value if multi header case
					for _, headerValue := range headerArray {

						h[i].Value = headerValue
						break
					}

					break
				}
			}
		}
	}

	return h
}

func DateWatcher(h []Helper) ([]Helper, error) {

	for i, helper := range h {

		if helper.HasValue() {
			continue
		}

		//add the date command result
		h[i].Value, _ = GetTargetDateStringValue(helper)

	}

	return h, nil
}

func RandomWatcher(h []Helper) ([]Helper, error) {

	for i, helper := range h {

		if helper.HasValue() {
			continue
		}

		methodParams := []string{}
		privateParams := helper.GetPrivateParams()

		if len(privateParams) > 0 {

			for _, v := range privateParams {

				methodParams = append(methodParams, v)
			}
		}

		h[i].Value = GetFakerValueStr(helper.Target, methodParams)

	}

	return h, nil
}

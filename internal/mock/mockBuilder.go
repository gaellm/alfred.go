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
	"alfred/internal/conf"
	"alfred/internal/helper"
	"alfred/internal/log"
	"bytes"
	"context"
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

func BuildMockFromJson(jsonData []byte) (Mock, error) {

	var mock Mock
	err := json.Unmarshal(jsonData, &mock)
	if err != nil {
		return mock, err
	}

	if mock.Response.BodyFile != "" {

		mock.Response.Body, err = getBodyFileContent(mock.Response.BodyFile)
		if err != nil {
			return mock, err
		}

		jsonData, err = json.Marshal(mock)
		if err != nil {
			return mock, err
		}

		err = json.Unmarshal(jsonData, &mock)
		if err != nil {
			return mock, err
		}
	}

	//clean json before save
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, jsonData)
	if err != nil {
		return mock, err
	}
	mock.SaveJsonBytes(buffer.Bytes())

	//find and create helpers
	helpers, err := helper.HelpersBuilder(buffer.Bytes())
	if err != nil {
		return mock, err
	}

	for _, h := range helpers {

		if h.Type == helper.REQUEST {

			mock.AddRequestHelper(h)

		} else if h.Type == helper.DATE {
			mock.AddDatetHelper(h)
		} else if h.Type == helper.RANDOM {

			mock.AddRandomHelper(h)
		}

		log.Debug(context.Background(), "helper "+h.Name+" found :'"+h.Target+"'"+" of type : '"+h.Type+"'", zap.String("mock-name", mock.GetName()))
	}

	return mock, nil
}

func getBodyFileContent(bodyFileName string) (string, error) {

	config, _ := conf.GetConfiguration()

	filePath := config.Alfred.Core.BodiesDir + bodyFileName

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(fileBytes), nil

}

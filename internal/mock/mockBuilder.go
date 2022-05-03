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
	"alfred/internal/log"
	"bytes"
	"context"
	"encoding/json"
	"regexp"

	"go.uber.org/zap"
)

func BuildMockFromJson(jsonData []byte) (Mock, error) {

	var mock Mock
	err := json.Unmarshal(jsonData, &mock)
	if err != nil {
		return mock, err
	}

	//clean json before save
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, jsonData)
	if err != nil {
		return mock, err
	}
	mock.SaveJsonBytes(buffer.Bytes())

	//find and create helpers
	for _, helperStrings := range findHelpersStrings(buffer.Bytes()) {

		h, err := CreateHelper(helperStrings[0], helperStrings[1])
		if err != nil {
			log.Error(context.Background(), "error creating helper", err, zap.String("mock-name", mock.GetName()))
		}

		log.Debug(context.Background(), "helper found :'"+h.Target+"'"+" of type : '"+h.Type+"'", zap.String("mock-name", mock.GetName()))
		mock.AddHelper(h)
	}

	return mock, nil
}

func findHelpersStrings(jsonData []byte) [][]string {

	r := regexp.MustCompile(`{{[ ]?([^{^}]*?)[ ]?}}`)
	return r.FindAllStringSubmatch(string(jsonData), -1)
}

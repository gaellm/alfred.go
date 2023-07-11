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

package server

import (
	"alfred/internal/log"
	"alfred/internal/mock"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func PatchMock(w http.ResponseWriter, r *http.Request, mockCollection mock.MockCollection) {

	requestRecover(w, r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), "failed to read request body", err)
	}

	//json Unmarshal to get name
	var rm map[string]interface{}
	err = json.Unmarshal(data, &rm)
	if err != nil {
		log.Error(r.Context(), "received json unmarschal fail", err)
	}

	var body []byte

	for _, m := range mockCollection.Mocks {

		if m.Name == rm["name"] {

			oldMock := m.GetJsonBytes()

			err = mock.MockPatch(m, data)
			if err != nil {
				log.Error(r.Context(), "mock patch error", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("mock patch error: " + err.Error()))
				return
			}

			log.Info(r.Context(), "mock has been patched",
				zap.String("mock-name", m.GetName()),
				zap.String("mock-conf-before", string(oldMock)),
				zap.String("mock-conf-after", string(m.GetJsonBytes())))
			body, _ = json.Marshal(m)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}
	}

	mockList, _ := json.MarshalIndent(mockCollection.GetMockInfoList(), "", "   ")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Sir ! This mock does not exists, however, I've found this mock list, if it can help: \n" + string(mockList) + "\n (Alfred)"))
}

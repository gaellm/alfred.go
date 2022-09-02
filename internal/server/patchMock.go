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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PatchMock(c *gin.Context, mocks mock.MockCollection) {

	requestRecover(c)

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(c.Request.Context(), "failed to read request body", err)
	}

	//json Unmarshal to get name
	var rm map[string]interface{}
	err = json.Unmarshal(data, &rm)
	if err != nil {
		log.Error(c.Request.Context(), "received json unmarschal fail", err)
	}

	var body []byte

	for _, m := range mocks {

		if m.Name == rm["name"] {

			oldMock := m.GetJsonBytes()

			err = mock.MockPatch(m, data)
			if err != nil {
				log.Error(c.Request.Context(), "mock patch error", err)
				c.String(http.StatusInternalServerError, "mock patch error: "+err.Error())
				return
			}

			log.Info(c.Request.Context(), "mock has been patched",
				zap.String("mock-name", m.GetName()),
				zap.String("mock-conf-before", string(oldMock)),
				zap.String("mock-conf-after", string(m.GetJsonBytes())))
			body, _ = json.Marshal(m)
			c.String(http.StatusOK, string(body))
			return
		}
	}

	mockList, _ := json.MarshalIndent(mocks.GetMockInfoList(), "", "   ")

	c.String(http.StatusNotFound, string("Hello Sir ! This mock does not exists, however, I've found this mock list, if it can help: \n"+string(mockList)+"\n (Alfred)"))

}

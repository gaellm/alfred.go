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
	"alfred/internal/helper"
	"alfred/internal/log"
	"alfred/internal/mock"
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddMocksRoutes(c *gin.Engine, mocks mock.MockCollection) {

	ctx := context.Background()

	for _, m := range mocks {

		m := m

		log.Debug(ctx, "Creating route for mock '"+m.GetName()+"'", zap.String("mock-conf", string(m.GetJsonBytes())))
		c.Handle(m.GetRequestMethod(), m.GetRequestUrl(), func(c *gin.Context) {

			requestRecover(c)

			data, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.Error(c.Request.Context(), "failed to read request body", err,
					zap.String("mock-name", m.GetName()),
					zap.String("request-path", c.Request.RequestURI),
				)
			}

			log.Debug(c.Request.Context(), "received a mock request, gona use mock '"+m.GetName()+"'",
				zap.String("request-path", c.Request.RequestURI),
				zap.String("request-body", string(data)),
				zap.String("mock-conf", string(m.GetJsonBytes())),
			)

			body := m.GetResponseBody()

			if m.HasRequestHelper() {
				log.Debug(c.Request.Context(), "start to populate request helper(s)",
					zap.String("mock-name", m.GetName()),
					zap.String("request-path", c.Request.RequestURI),
					zap.String("request-body", string(data)),
					zap.String("mock-conf", string(m.GetJsonBytes())),
				)

				// Populate mock request helpers
				helpersPopulated, err := helper.RequestHelperWatcher(data, c, m.GetRequestHelpers())
				if err != nil {
					log.Warn(c.Request.Context(), "helpers request watcher in error", err,
						zap.String("mock-name", m.GetName()),
						zap.String("request-path", c.Request.RequestURI),
						zap.String("request-body", string(data)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
					)
				}

				log.Debug(c.Request.Context(), "helper(s) populated",
					zap.String("mock-name", m.GetName()),
					zap.String("request-path", c.Request.RequestURI),
					zap.String("request-body", string(data)),
					zap.String("mock-conf", string(m.GetJsonBytes())),
					zap.String("helpers", m.UpdateRequestHelpers(helpersPopulated).GetJsonHelpers()),
				)

				// Replace helpers inside mock response body
				body, err = helper.HelperReplacement(m.GetResponseBody(), helpersPopulated)

				// Replace helpers inside mock response headers
				for k, v := range m.GetResponseHeaders() {

					v, err = helper.HelperReplacement(v, helpersPopulated)
					c.Header(k, v)
				}

				if err != nil {
					log.Warn(c.Request.Context(), "error during helpers replacement", err,
						zap.String("request-path", c.Request.RequestURI),
						zap.String("request-body", string(data)),
						zap.String("mock-conf", string(m.GetJsonBytes())),
						zap.String("response-body", m.GetResponseBody()),
						zap.String("helpers", m.UpdateRequestHelpers(helpersPopulated).GetJsonHelpers()))
				}

			} else {

				//set headers
				for k, v := range m.GetResponseHeaders() {
					c.Header(k, v)
				}
			}

			//time.Sleep(5 * time.Second)

			//set status and body to end response
			c.String(m.GetResponseStatus(), body)

		})

	}
}

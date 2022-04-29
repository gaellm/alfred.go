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
	"context"

	"github.com/gin-gonic/gin"
)

func AddMocksRoutes(c *gin.Engine, mocks mock.MockCollection) {

	ctx := context.Background()

	for _, m := range mocks {

		m := m

		log.Debug(ctx, "Creating route for mock '"+m.GetName()+"' "+"with url '"+m.GetRequestUrl()+"' "+"body '"+m.GetResponseBody()+"'")

		c.Handle(m.GetRequestMethod(), m.GetRequestUrl(), func(c *gin.Context) {

			for k, v := range m.GetResponseHeaders() {
				c.Header(k, v)
			}
			c.String(m.GetResponseStatus(), m.GetResponseBody())

		})

	}
}

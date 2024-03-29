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
	"testing"
	"time"
)

func TestGetDelay(t *testing.T) {

	var testMock Mock
	var delay time.Duration

	testMock.Response.MinResponseTime = 45

	delay = testMock.GetDelay()
	if delay != 45*time.Millisecond {
		t.Errorf("mock delay failed, it should be 45ms")
	}

	testMock.Response.MinResponseTime = 3
	testMock.Response.MaxResponseTime = 30

	delay = testMock.GetDelay()

	if delay < 3*time.Millisecond && delay > 30*time.Millisecond {

		t.Errorf("mock delay failed, it should be between 3 and 30ms")
	}

}

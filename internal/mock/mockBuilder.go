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
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func BuildMockFromJson(jsonData []byte) (Mock, error) {

	var mock Mock
	err := json.Unmarshal(jsonData, &mock)
	if err != nil {
		return mock, err
	}

	return mock, nil
}

func TestMock() {
	mockTest := []byte(`
	{
		"name": "Test-mock",
		"request": {
			"method": "GET",
			"url": "/some/thing"
		}
	}`)

	//var mymock Mock
	mymock, err := BuildMockFromJson(mockTest)
	if err != nil {
		log.Error(context.Background(), "Error during mock build from json", err, zap.String("text-provided", string(mockTest)))
		panic(fmt.Errorf("fatal error, mock build from json: %w", err))
	}

	fmt.Println(mymock.GetRequestMethod())
}

/*
type Student struct {
    Name string
    Roll int
}

func (s Student) Print() {
    fmt.Println(s)
}

func main() {

    jack := Student{"Jack", 123}

    jack.Print() // {Jack 123}

}

*/

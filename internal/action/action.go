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

package action

import (
	"alfred/internal/helper"
	"alfred/internal/mock"
	"alfred/pkg/request"
	"errors"
	"math/rand"
	"time"
)

const SEND_REQUEST_TYPE = "send-request"
const SEND_REQUEST_DEFAULT_METHOD = "GET"

func GetDelayDuration(act mock.MockAction) time.Duration {

	if act.MaxScheduledTime <= act.MinScheduledTime {
		return time.Duration(act.MinScheduledTime) * time.Millisecond
	}

	return time.Duration(rand.Intn(act.MaxScheduledTime-act.MinScheduledTime)+act.MinScheduledTime) * time.Millisecond
}

func CreateRequestFromMockAction(act mock.MockAction, helpers []helper.Helper) (request.Request, error) {

	if act.Type == SEND_REQUEST_TYPE {

		act, _ = replaceHelperRequestAction(act, helpers)

		var req request.Request

		var err error
		act, err = sanitizeRequestAction(act)
		if err != nil {
			return request.Request{}, err
		}

		err = req.SetUrl(act.Url)
		if err != nil {
			return req, err
		}

		err = req.SetMethod(act.Method)
		if err != nil {
			return req, err
		}

		err = req.SetTimeout(act.Timeout)
		if err != nil {
			return req, err
		}

		req.Body = []byte(act.Body)
		req.Headers = act.Headers

		return req, nil
	}

	return request.Request{}, errors.New(act.Type + " can't generate send request action, type have to be: " + SEND_REQUEST_TYPE)
}

func replaceHelperRequestAction(act mock.MockAction, helpers []helper.Helper) (mock.MockAction, error) {

	if act.Type == SEND_REQUEST_TYPE {

		if len(helpers) > 0 {

			var err error

			act.Url, err = helper.HelperReplacement(act.Url, helpers)
			if err != nil {
				return act, err
			}

			act.Body, err = helper.HelperReplacement(act.Body, helpers)
			if err != nil {
				return act, err
			}

			//replace headers helpers
			if len(act.Headers) > 0 {

				headers := make(map[string]string)
				for k, v := range act.Headers {

					key, err := helper.HelperReplacement(k, helpers)
					if err != nil {
						return act, err
					}
					value, err := helper.HelperReplacement(v, helpers)
					if err != nil {
						return act, err
					}
					headers[key] = value
				}
			}
		}

		return act, nil
	}

	return act, errors.New(act.Type + " type can't replace helper in send request action, type have to be: " + SEND_REQUEST_TYPE)

}

func sanitizeRequestAction(act mock.MockAction) (mock.MockAction, error) {

	if act.Url == "" {
		return act, errors.New("empty action request url")
	}

	if act.Method == "" {
		act.Method = SEND_REQUEST_DEFAULT_METHOD
	}

	return act, nil

}

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
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GlobalDelayPlayload struct {
	ResponseTime int `json:"minResponseTime"`
	Duration     int `json:"duration"`
}

func DelayMocks(alfredGlobalDelay *time.Duration, w http.ResponseWriter, r *http.Request) {

	requestRecover(w, r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), "failed to read request body", err)
	}

	var playload GlobalDelayPlayload
	err = json.Unmarshal(data, &playload)
	if err != nil {
		log.Error(r.Context(), "received json unmarschal fail", err)
	}

	responseTimeDuration := time.Duration(playload.ResponseTime) * time.Millisecond
	delayDuration := time.Duration(playload.Duration) * time.Millisecond

	//set response time offset
	*alfredGlobalDelay = responseTimeDuration
	log.Info(r.Context(), "a global response time offset of : "+strconv.Itoa(playload.ResponseTime)+"ms has been set for "+strconv.Itoa(playload.Duration)+"ms")

	//goroutine for duration and responseTime reset
	go func(during time.Duration, alfredGlobalDelay *time.Duration) {

		time.Sleep(during)
		*alfredGlobalDelay = 0
		log.Info(r.Context(), "global response time offset reset")

	}(delayDuration, alfredGlobalDelay)

}

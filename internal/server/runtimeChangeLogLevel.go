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

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOG_LEVEL_DEBUG = "debug"
const LOG_LEVEL_INFO = "info"
const LOG_LEVEL_ERROR = "error"

var LOG_LEVEL_LIST = []string{LOG_LEVEL_DEBUG, LOG_LEVEL_ERROR, LOG_LEVEL_INFO}
var ZAP_LEVEL_MAP = map[string]zapcore.Level{
	LOG_LEVEL_INFO:  zap.InfoLevel,
	LOG_LEVEL_DEBUG: zap.DebugLevel,
	LOG_LEVEL_ERROR: zap.ErrorLevel,
}

type ChangeLogLevelReqPlayload struct {
	ConfiguredLevel string `json:"configuredLevel"`
}

type ChangeLogLevelResPlayload struct {
	PreviousLevel   string `json:"previousLevel"`
	ConfiguredLevel string `json:"configuredLevel"`
	EffectiveLevel  string `json:"effectiveLevel"`
}

func ChangingLoggingLevelRuntime(w http.ResponseWriter, r *http.Request) {

	//requestRecover(c)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), "failed to read request body", err)
	}

	var playload ChangeLogLevelReqPlayload

	err = json.Unmarshal(data, &playload)
	if err != nil {
		log.Error(r.Context(), "received json unmarschal fail", err)
	}

	var res ChangeLogLevelResPlayload

	res.ConfiguredLevel = playload.ConfiguredLevel
	res.PreviousLevel = log.GetLevel()

	err = log.SetLevel(playload.ConfiguredLevel)
	res.EffectiveLevel = log.GetLevel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	resJsonStr, _ := json.MarshalIndent(res, "", "   ")
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonStr)
}

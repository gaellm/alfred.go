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

package helper

import (
	"errors"
	"regexp"
	"strings"
)

var (
	TYPES = [...]string{REQUEST, DATE, RANDOM, PATH_REGEX}
)

func createHelper(helperString string, helperTarget string) (Helper, error) {

	var h Helper
	params := getHelperStringParams(helperTarget)

	helperName, exists := params[PARAM_NAME]
	if exists {
		h.Name = helperName
	}

	h.String = helperString

	r := regexp.MustCompile(`alfred\.\w*\.([^@]*).*`)
	h.Target = strings.TrimSpace(r.FindStringSubmatch(helperTarget)[1])

	var err error
	h.Type, err = detectHelperType(helperTarget)
	if err != nil {
		return h, err
	}

	helperRegex, helperRegexExists := params[PARAM_REGEX]
	if helperRegexExists {
		h.Regex, _ = regexp.Compile(helperRegex)

	}

	h, error := sanitizeHelper(h)
	if error != nil {
		return h, error
	}

	return h, err
}

func detectHelperType(helperTarget string) (string, error) {

	s := strings.Split(helperTarget, ".")

	for _, t := range TYPES {
		if s[1] == t {
			return s[1], nil
		}
	}

	return "", errors.New("helper type '" + s[1] + "' is not handled by Alfred")
}

func findHelpersStrings(jsonData []byte) [][]string {

	var helpersStrings [][]string

	r := regexp.MustCompile(`{{[ ]?([^{}]*?[)]?)[ ]?}}`)

	allStringSubmatch := r.FindAllStringSubmatch(string(jsonData), -1)

	// keep unique
	for _, stringSubmatch := range allStringSubmatch {

		founded := false
		for _, submatch := range helpersStrings {

			if stringSubmatch[1] == submatch[1] {
				founded = true
				continue
			}
		}

		if !founded {
			helpersStrings = append(helpersStrings, stringSubmatch)
		}
	}
	return helpersStrings
}

func isKnownParam(param string) bool {

	//knownParams := [...]string{PARAM_NAME, PARAM_TYPE, PARAM_DESC}
	knownParams := [...]string{PARAM_NAME, PARAM_REGEX}

	for _, knownParam := range knownParams {

		if knownParam == param {
			return true
		}
	}

	return false
}

func getHelperStringParams(helperStr string) map[string]string {

	params := make(map[string]string)

	//Get helpers params
	paramsRegexp := regexp.MustCompile(`@(\w*):'([^']*)'`)
	paramsStringsSubmatch := paramsRegexp.FindAllStringSubmatch(string(helperStr), -1)

	for _, paramStringsSubmatch := range paramsStringsSubmatch {

		if len(paramStringsSubmatch) > 1 && isKnownParam(paramStringsSubmatch[1]) {
			params[paramStringsSubmatch[1]] = paramStringsSubmatch[2]
		}
	}

	return params
}

func HelpersBuilder(buffer []byte) ([]Helper, error) {

	var helpers []Helper
	helpersStrings := findHelpersStrings(buffer)

	//find and create helpers
	for _, helperStrings := range helpersStrings {

		h, err := createHelper(helperStrings[0], helperStrings[1])
		if err != nil {
			return helpers, err
		}

		helpers = append(helpers, h)
	}

	return helpers, nil
}

func sanitizeHelper(h Helper) (Helper, error) {

	//check date targets
	if h.Type == DATE {

		return sanitizeDateHelper(h)
	} else if h.Type == RANDOM {

		return sanitizeRandomHelper(h)
	}

	return h, nil
}

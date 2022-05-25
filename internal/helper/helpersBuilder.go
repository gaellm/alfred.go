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
	TYPES = [...]string{REQUEST}
)

func createHelper(helperString string, helperTarget string) (Helper, error) {

	var h Helper

	h.String = helperString

	r := regexp.MustCompile(`alfred\.req\.(.*)`)
	h.Target = r.FindStringSubmatch(helperTarget)[1]

	var err error
	h.Type, err = detectHelperType(helperTarget)

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

	r := regexp.MustCompile(`{{[ ]?([^{^}]*?)[ ]?}}`)
	return r.FindAllStringSubmatch(string(jsonData), -1)
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

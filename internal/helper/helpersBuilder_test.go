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

import "testing"

func TestCreateHelper(t *testing.T) {

	helper, err := createHelper("{{ alfred.req.test.titi }}", "alfred.req.test.titi")
	if err != nil {
		t.Errorf("Create helper fail with error " + err.Error())
	}

	if helper.HasValue() {
		t.Errorf("Helper HasValue is: true, want: false.")
	}

	if helper.String != "{{ alfred.req.test.titi }}" {
		t.Errorf("Helper string is: %s, want: %s.", helper.String, "{{ alfred.req.test.titi }}")
	}

	if helper.Target != "test.titi" {
		t.Errorf("Helper target is: %s, want: %s.", helper.Target, "test.titi")
	}

	if helper.Type != "req" {
		t.Errorf("Helper type is: %s, want: %s.", helper.Type, "req")
	}

	helperParam, errParam := createHelper("{{ alfred.req.test.titi @name:'toto' }}", "alfred.req.test.titi @name:'toto'")
	if errParam != nil {
		t.Errorf("Create helper fail with error " + err.Error())
	}

	if helperParam.Name != "toto" {
		t.Errorf("Helper Name is not toto.")
	}

	if helperParam.Target != "test.titi" {
		t.Errorf("Helper target is: %s, want: %s.", helperParam.Target, "test.titi")
	}

}

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func TestFindHelpersStrings(t *testing.T) {

	jsonData := []byte("{\"name\":\"postmock\",\"request\":{\"method\":\"POST\",\"url\":\"/some/thing/:test.tyty\"},\"response\":{\"status\":200,\"body\":\"Hello world! {{ alfred.req.test.titi }}\",\"headers\":{\"Content-Type\":\"text/plain\",\"Test\":\"{{ alfred.req.test.tyty }}\"}}}")

	s0 := []string{"{{ alfred.req.test.titi }}", "{{ alfred.req.test.tyty }}"}
	s1 := []string{"alfred.req.test.titi", "alfred.req.test.tyty"}

	for _, helperStrings := range findHelpersStrings(jsonData) {

		if !isElementExist(s0, helperStrings[0]) {
			t.Errorf("helperStrings[0] failed with '%s' value", helperStrings[0])
		}

		if !isElementExist(s1, helperStrings[1]) {
			t.Errorf("helperStrings[1] failed with '%s' value", helperStrings[1])
		}

	}
}

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
	"testing"
)

//https://golangexample.com/faker-for-golang-random-data-generator-compatible-with-postman-dynamic-variables/

func TestFakerNoParamReturnString(t *testing.T) {

	alfredFaker := AlfredFaker{}

	res := alfredFaker.FakerNoParamReturnString("RandomLoremSentences")

	if res == "" {
		t.Errorf("random return empty string")
	}

	t.Logf("%v", res)

}

func TestFakerNoParamReturnInt(t *testing.T) {

	alfredFaker := AlfredFaker{}

	res := alfredFaker.FakerNoParamReturnInt("RandomInt")

	if res == "" {
		t.Error("random return empty string")
	}

}

func TestFakerNoParamReturnBool(t *testing.T) {

	alfredFaker := AlfredFaker{}

	res := alfredFaker.FakerNoParamReturnBool("RandomBoolean")

	if res == "" {
		t.Error("random return empty string")
	}

	t.Logf("%v", res)
}

func TestFakerTwoParamReturnInt(t *testing.T) {

	alfredFaker := AlfredFaker{}

	res := alfredFaker.FakerTwoParamReturnInt("RandomIntBetween", 0, 15)

	if res == "" {
		t.Error("random return empty string")
	}

	t.Logf("%v", res)
}

func TestFakerNoParamReturnUuid(t *testing.T) {

	alfredFaker := AlfredFaker{}

	res := alfredFaker.FakerNoParamReturnUuid("RandomUUID")

	if res == "" {
		t.Error("random return empty string")
	}

	t.Logf("%v", res)
}

func TestFakerNoParamReturnFloat64(t *testing.T) {

	alfredFaker := AlfredFaker{}

	res := alfredFaker.FakerNoParamReturnFloat64("RandomAddressLatitude")

	if res == "" {
		t.Error("random return empty string")
	}

	t.Logf("%v", res)

}

func TestCheckFakerMethodName(t *testing.T) {

	_, err := CheckFakerMethodName("RandomAnimalsImage")
	if err != nil {
		t.Error("check faker method name failed with error: " + err.Error())
	}
}

func TestGetFakerReflectRef(t *testing.T) {

	reflectMethod := getFakerReflectRef("RandomBoolean")
	if reflectMethod != "FakerNoParamReturnBool" {
		t.Error("get faker reflection name failed : " + reflectMethod + " instead of FakerNoParamReturnBool")
	}

	reflectMethod = getFakerReflectRef("RandomDomainWord")
	if reflectMethod != "FakerNoParamReturnString" {
		t.Error("get faker reflection name failed : " + reflectMethod + " instead of FakerNoParamReturnString")
	}

}

func TestGetFakerValueStr(t *testing.T) {

	allMethodsToTest := fakerMethodNamesAllowed

	for _, methodName := range allMethodsToTest {

		str := GetFakerValueStr(methodName, []string{"1", "100"})

		if str == "" {
			t.Error("get faker value for : " + methodName + " failed")
		}

		t.Logf("%v", str)
	}
}

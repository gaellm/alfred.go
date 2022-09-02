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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ddosify/go-faker/faker"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AlfredFaker struct {
}

var fakerMethodNamesAllowed = [116]string{
	"RandomGuid",
	"RandomUUID",
	"RandomAlphanumeric",
	"RandomBoolean",
	"RandomInt",
	"RandomSafeColorName",
	"RandomSafeColorHex",
	"RandomAbbreviation",
	"RandomIP",
	"RandomIpv6",
	"RandomIntBetween",
	"RandomMACAddress",
	"RandomPassword",
	"RandomLocale",
	"RandomUserAgent",
	"RandomProtocol",
	"RandomSemver",
	"RandomPersonFirstName",
	"RandomPersonLastName",
	"RandomPersonFullName",
	"RandomPersonNamePrefix",
	"RandomPersonNameSuffix",
	"RandomJobArea",
	"RandomJobDescriptor",
	"RandomJobTitle",
	"RandomJobType",
	"RandomPhoneNumberExt",
	"RandomAddressCity",
	"RandomAddresStreetName",
	"RandomAddressStreetAddress",
	"RandomAddressCountry",
	"RandomCountryCode",
	"RandomAddressLatitude",
	"RandomAddressLongitude",
	"RandomAvatarImage",
	"RandomImageURL",
	"RandomAbstractImage",
	"RandomAnimalsImage",
	"RandomBusinessImage",
	"RandomCatsImage",
	"RandomCityImage",
	"RandomFoodImage",
	"RandomNightlifeImage",
	"RandomFashionImage",
	"RandomPeopleImage",
	"RandomNatureImage",
	"RandomSportsImage",
	"RandomTransportImage",
	"RandomDataImageUri",
	"RandomBankAccount",
	"RandomBankAccountName",
	"RandomCreditCardMask",
	"RandomBankAccountBic",
	"RandomBankAccountIban",
	"RandomTransactionType",
	"RandomCurrencyCode",
	"RandomCurrencyName",
	"RandomCurrencySymbol",
	"RandomBitcoin",
	"RandomCompanyName",
	"RandomCompanySuffix",
	"RandomBs",
	"RandomBsAdjective",
	"RandomBsBuzzWord",
	"RandomBsNoun",
	"RandomCatchPhrase",
	"RandomCatchPhraseAdjective",
	"RandomCatchPhraseDescriptor",
	"RandomCatchPhraseNoun",
	"RandomDatabaseColumn",
	"RandomDatabaseType",
	"RandomDatabaseCollation",
	"RandomDatabaseEngine",
	"RandomDateFuture",
	"RandomDatePast",
	"RandomDateRecent",
	"RandomWeekday",
	"RandomMonth",
	"RandomDomainName",
	"RandomDomainSuffix",
	"RandomDomainWord",
	"RandomEmail",
	"RandomExampleEmail",
	"RandomUsername",
	"RandomUrl",
	"RandomFileName",
	"RandomFileType",
	"RandomFileExtension",
	"RandomCommonFileName",
	"RandomCommonFileType",
	"RandomCommonFileExtension",
	"RandomFilePath",
	"RandomDirectoryPath",
	"RandomMimeType",
	"RandomPrice",
	"RandomProduct",
	"RandomProductAdjective",
	"RandomProductMaterial",
	"RandomProductName",
	"RandomDepartment",
	"RandomNoun",
	"RandomVerb",
	"RandomIngVerb",
	"RandomAdjective",
	"RandomWord",
	"RandomWords",
	"RandomPhrase",
	"RandomLoremWord",
	"RandomLoremWords",
	"RandomLoremSentence",
	"RandomLoremSentences",
	"RandomLoremParagraph",
	"RandomLoremParagraphs",
	"RandomLoremText",
	"RandomLoremSlug",
	"RandomLoremLines",
}

func sanitizeRandomHelper(h Helper) (Helper, error) {

	//First letter go to upper case
	h.Target = cases.Title(language.Und, cases.NoLower).String(h.Target)

	methodName, err := CheckFakerMethodName(h.Target)
	if err != nil {
		return h, err
	}

	//Get method params
	randomParams, err := getParamsFromFakerMethodName(h.Target)
	if err != nil {
		return h, err
	}
	for i, v := range randomParams {

		h.AddPrivateParam("param-"+strconv.Itoa(i), v)
	}

	// Target is now method name only
	h.Target = methodName

	return h, nil
}

func getParamsFromFakerMethodName(fakerMethodName string) ([]string, error) {

	if strings.Contains(fakerMethodName, "(") {

		methodParams := regexp.MustCompile(`\w+[(](.*)[)].*`).FindStringSubmatch(fakerMethodName)
		if len(methodParams) <= 1 {
			return []string{}, errors.New("bad random format '" + fakerMethodName + "' need somthing like RandomIntBetween(1,1OO)")
		}

		return strings.Split(methodParams[1], ","), nil
	}

	return []string{}, nil
}

func getFakerReflectRef(fakerMethodName string) string {

	var fakerReflectRef = make(map[string]string)

	//UUID
	fakerReflectRef["RandomUUID"] = "FakerNoParamReturnUuid"
	fakerReflectRef["RandomGuid"] = "FakerNoParamReturnUuid"
	//bool
	fakerReflectRef["RandomBoolean"] = "FakerNoParamReturnBool"
	//int
	fakerReflectRef["RandomInt"] = "FakerNoParamReturnInt"
	//float64
	fakerReflectRef["RandomAddressLatitude"] = "FakerNoParamReturnFloat64"
	fakerReflectRef["RandomAddressLongitude"] = "FakerNoParamReturnFloat64"
	//2 params
	fakerReflectRef["RandomIntBetween"] = "FakerTwoParamReturnInt"

	if fakerReflectRef[fakerMethodName] != "" {
		return fakerReflectRef[fakerMethodName]
	}

	return "FakerNoParamReturnString"
}

func GetFakerValueStr(fakerMethodName string, fakerMethodParams []string) string {

	alfredFaker := AlfredFaker{}
	methodName := getFakerReflectRef(fakerMethodName)

	//reflection
	methodVal := reflect.ValueOf(&alfredFaker).MethodByName(methodName)
	methodIface := methodVal.Interface()

	if methodName == "FakerTwoParamReturnInt" {

		method := methodIface.(func(string, int, int) string)

		param1, _ := strconv.Atoi(fakerMethodParams[0])
		param2, _ := strconv.Atoi(fakerMethodParams[1])

		return method(fakerMethodName, param1, param2)

	}

	method := methodIface.(func(string) string)
	return method(fakerMethodName)
}

func CheckFakerMethodName(fakerMethodName string) (string, error) {

	//extract method if params
	if strings.Contains(fakerMethodName, "(") {
		methodWithTwoParams := regexp.MustCompile(`(\w+)[(].*`).FindStringSubmatch(fakerMethodName)
		if len(methodWithTwoParams) <= 1 {
			return "", errors.New("bad random format '" + fakerMethodName + "' need somthing like RandomIntBetween(1,1OO)")
		}
		fakerMethodName = methodWithTwoParams[1]
	}

	for _, v := range fakerMethodNamesAllowed {
		if v == fakerMethodName {
			return fakerMethodName, nil
		}
	}

	return "", errors.New("the random method name '" + fakerMethodName + "' not exists, allowed methods are: " + fmt.Sprint(fakerMethodNamesAllowed))
}

func (a *AlfredFaker) FakerNoParamReturnString(fakerMethodName string) string {
	faker := faker.NewFaker()

	methodVal := reflect.ValueOf(&faker).MethodByName(fakerMethodName)
	methodIface := methodVal.Interface()

	//no param return string
	method := methodIface.(func() string)

	return method()
}

func (a *AlfredFaker) FakerNoParamReturnUuid(fakerMethodName string) string {
	faker := faker.NewFaker()

	methodVal := reflect.ValueOf(&faker).MethodByName(fakerMethodName)
	methodIface := methodVal.Interface()

	//no param return string
	method := methodIface.(func() uuid.UUID)

	return fmt.Sprint(method())
}

func (a *AlfredFaker) FakerNoParamReturnBool(fakerMethodName string) string {
	faker := faker.NewFaker()

	methodVal := reflect.ValueOf(&faker).MethodByName(fakerMethodName)
	methodIface := methodVal.Interface()

	//no param return string
	method := methodIface.(func() bool)

	return fmt.Sprint(method())
}

func (a *AlfredFaker) FakerNoParamReturnInt(fakerMethodName string) string {
	faker := faker.NewFaker()

	methodVal := reflect.ValueOf(&faker).MethodByName(fakerMethodName)
	methodIface := methodVal.Interface()

	//no param return string
	method := methodIface.(func() int)

	return fmt.Sprint(method())

}

func (a *AlfredFaker) FakerNoParamReturnFloat64(fakerMethodName string) string {
	faker := faker.NewFaker()

	methodVal := reflect.ValueOf(&faker).MethodByName(fakerMethodName)
	methodIface := methodVal.Interface()

	//no param return string
	method := methodIface.(func() float64)

	return fmt.Sprint(method())

}

func (a *AlfredFaker) FakerTwoParamReturnInt(fakerMethodName string, param1 int, param2 int) string {
	faker := faker.NewFaker()

	methodVal := reflect.ValueOf(&faker).MethodByName(fakerMethodName)
	methodIface := methodVal.Interface()

	//no param return string
	method := methodIface.(func(int, int) int)

	return fmt.Sprint(method(param1, param2))

}

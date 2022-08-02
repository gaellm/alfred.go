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

	"github.com/ddosify/go-faker/faker"
	"github.com/google/uuid"
)

type AlfredFaker struct {
}

var fakerMethodNamesAllowed = [115]string{
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

	if fakerReflectRef[fakerMethodName] != "" {
		return fakerReflectRef[fakerMethodName]
	}

	return "FakerNoParamReturnString"
}

func GetFakerValueStr(fakerMethodName string) string {

	alfredFaker := AlfredFaker{}
	methodName := getFakerReflectRef(fakerMethodName)

	//reflection
	methodVal := reflect.ValueOf(&alfredFaker).MethodByName(methodName)
	methodIface := methodVal.Interface()
	method := methodIface.(func(string) string)

	return method(fakerMethodName)
}

func CheckFakerMethodName(fakerMethodName string) error {

	for _, v := range fakerMethodNamesAllowed {
		if v == fakerMethodName {
			return nil
		}
	}

	return errors.New("the random method name '" + fakerMethodName + "' not exists, allowed methods are: " + fmt.Sprint(fakerMethodNamesAllowed))
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

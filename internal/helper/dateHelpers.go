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
	"strconv"
	"strings"
	"time"
)

const DATE_REF_DATE = "date"
const DATE_REF_NOW = "now"
const DATE_UTC_STR = ".utc"

const DATE_PRIVATE_PARAMS_REF_NAME = "dateRef"
const DATE_PRIVATE_PARAMS_ISUTC_NAME = "isUTC"
const DATE_PRIVATE_PARAMS_FORMAT_NAME = "format"
const DATE_PRIVATE_PARAMS_ADD_NAME = "add"

const DATE_PRIVATE_PARAMS_ADD_VALUE_NAME = "addValue"

var DATE_REFS = [4]string{DATE_REF_DATE, DATE_REF_NOW}

func sanitizeDateHelper(h Helper) (Helper, error) {

	ref, err := checkDateHelperRef(h)
	if err != nil {
		return h, err
	}

	h.AddPrivateParam(DATE_PRIVATE_PARAMS_REF_NAME, ref)

	if isDateHelperUtc(h) {
		h.AddPrivateParam(DATE_PRIVATE_PARAMS_ISUTC_NAME, "true")
	}

	dateFormat, formatErr := getTimeFormatFromStr(h.Target)
	if formatErr != nil {
		return h, formatErr
	}
	h.AddPrivateParam(DATE_PRIVATE_PARAMS_FORMAT_NAME, dateFormat)

	addValue, addError := getTimeAddFromStr(h.Target)
	if addError != nil {
		return h, addError
	}
	h.AddPrivateParam(DATE_PRIVATE_PARAMS_ADD_VALUE_NAME, addValue)

	if ref == DATE_REF_DATE {

		h.Value, err = GetTargetDateStringValue(h)
		if err != nil {
			return h, err
		}
	}

	return h, nil
}

func buildDateTimeFromStr(dateStr string) (time.Time, error) {

	paramsStr := regexp.MustCompile(".*date[(]([^)]*)[)].*").FindStringSubmatch(dateStr)
	if len(paramsStr) <= 1 {
		return time.Time{}, errors.New("bad date format '" + dateStr + "' need somthing like date(2009,01,03,4,2,0,0)")
	}

	params := strings.Split(paramsStr[1], ",")
	if len(params) < 7 {
		return time.Time{}, errors.New("bad date arguments numbers : '" + dateStr + "' need somthing like date(2009,01,03,4,2,0,0)")
	}

	var dateParams [7]int
	for i, v := range params {

		var err error

		dateParams[i], err = strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return time.Time{}, errors.New("bad date string argument: " + dateStr + " - " + err.Error())
		}
	}

	return time.Date(dateParams[0], time.Month(dateParams[1]), dateParams[2], dateParams[3], dateParams[4], dateParams[5], dateParams[6], time.Local), nil

}

func getTimeFormatFromStr(dateStr string) (string, error) {

	if !strings.Contains(dateStr, "format") {
		return "", nil
	}

	formatStr := regexp.MustCompile(".*format[(]'([^)]*)'[)].*").FindStringSubmatch(dateStr)

	if len(formatStr) <= 1 {
		return "", errors.New("bad format '" + dateStr + "' need somthing like .format('unix') or .format('2006-01-02T15:04:05.000Z')")
	}

	return formatStr[1], nil
}

func getTimeAddFromStr(dateStr string) (string, error) {

	if !strings.Contains(dateStr, "add") {
		return "", nil
	}

	addStr := regexp.MustCompile(".*add[(]'([^)]*)'[)].*").FindStringSubmatch(dateStr)
	if len(addStr) <= 1 {
		return "", errors.New("bad add '" + dateStr + "' need somthing like .add(10ms) or .add(-1h)")
	}

	durationStr := strings.TrimSpace(addStr[1])

	_, err := time.ParseDuration(durationStr)
	if err != nil {
		return "", err
	}

	return durationStr, nil
}

func checkDateHelperRef(h Helper) (string, error) {

	s := regexp.MustCompile("([^ .(]*)[ .(]?.*").FindAllStringSubmatch(h.Target, -1)[0][1]

	for _, ref := range DATE_REFS {
		if s == ref {

			return s, nil
		}
	}

	return "", errors.New("alfred " + DATE + " helper reference " + h.Target + " unknown")
}

func isDateHelperUtc(h Helper) bool {

	return strings.Contains(h.Target, DATE_UTC_STR)

}

// called on each requests (target is checked during startup)
func GetTargetDateStringValue(h Helper) (string, error) {

	ref := h.GetPrivateParam(DATE_PRIVATE_PARAMS_REF_NAME)
	isUtc := h.GetPrivateParam(DATE_PRIVATE_PARAMS_ISUTC_NAME)
	format := h.GetPrivateParam(DATE_PRIVATE_PARAMS_FORMAT_NAME)
	addValue := h.GetPrivateParam(DATE_PRIVATE_PARAMS_ADD_VALUE_NAME)

	var theDate time.Time

	if ref == DATE_REF_NOW {

		theDate = time.Now()
	} else {

		theDate, _ = buildDateTimeFromStr(h.Target)
	}

	//if Not UTC
	if isUtc != "" {

		theDate = theDate.UTC()
	}

	if addValue != "" {
		duration, _ := time.ParseDuration(addValue)
		theDate = theDate.Add(duration)
	}

	if format != "" {

		if format == "unix" {
			unixDateStr := strconv.Itoa(int(theDate.Unix()))
			return unixDateStr, nil
		}
		return theDate.Format(format), nil
	}

	return theDate.String(), nil
}

/*


alfred.time.now.format('')
alfred.time.now.format('unix')
alfred.time.now.utc



		func main() {
	    t := time.Now()

	    //Add 1 hours
	    newT := t.Add(time.Hour * 1)
	    fmt.Printf("Adding 1 hour\n: %s\n", newT)

	    //Add 15 min
	    newT = t.Add(time.Minute * 15)
	    fmt.Printf("Adding 15 minute\n: %s\n", newT)

	    //Add 10 sec
	    newT = t.Add(time.Second * 10)
	    fmt.Printf("Adding 10 sec\n: %s\n", newT)

	    //Add 100 millisecond
	    newT = t.Add(time.Millisecond * 10)
	    fmt.Printf("Adding 100 millisecond\n: %s\n", newT)

	    //Add 1000 microsecond
	    newT = t.Add(time.Millisecond * 10)
	    fmt.Printf("Adding 1000 microsecond\n: %s\n", newT)

	    //Add 10000 nanosecond
	    newT = t.Add(time.Nanosecond * 10000)
	    fmt.Printf("Adding 1000 nanosecond\n: %s\n", newT)

	    //Add 1 year 2 month 4 day
	    newT = t.AddDate(1, 2, 4)
	    fmt.Printf("Adding 1 year 2 month 4 day\n: %s\n", newT)
	}

			.Add(time.Hour * 3)

			.AddDate(0,0,-1)

			now

*/

//time.Now().Add(time.Millisecond * 5)
//time.Date(2009,01,03,4,2,0,0,time.UTC)

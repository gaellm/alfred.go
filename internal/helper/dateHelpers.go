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
	"strings"
	"time"
)

const DATE_REF_DATE = "date"
const DATE_REF_NOW = "now"

var DATE_REFS = [4]string{DATE_REF_DATE, DATE_REF_NOW}

func sanitizeDateHelper(h Helper) (Helper, error) {

	ref, err := checkDateHelperRef(h)
	if err != nil {
		return h, err
	}

	return h.AddPrivateParam("dateRef", ref), nil
}

func checkDateHelperRef(h Helper) (string, error) {

	s := strings.Split(h.Target, ".")

	for _, ref := range DATE_REFS {
		if s[0] == ref {

			return s[0], nil
		}
	}

	return "", errors.New("alfred " + DATE + " helper reference " + h.Target + " unknown")
}

// called on each requests (target is checked during startup)
func GetTargetDateStringValue(h Helper) (string, error) {

	ref := h.GetPrivateParam("dateRef") //now date ?

	if ref == DATE_REF_NOW {
		return time.Now().String(), nil
	}

	return time.Now().UTC().String(), nil

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

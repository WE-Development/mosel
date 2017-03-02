/*
 * Copyright 2017 Robin Engel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package commons

import (
	"testing"
)

func TestContainsStr(t *testing.T) {
	arr := [...]string{"1", "2", "3", "4", "5"}
	ex := [...]string{"4", "2"}

	res := ExcludeStr(arr[:], ex[:])

	expected := [...]string{"1", "3", "5"}

	for i := range res {
		if res[i] != expected[i] {
			t.Fail()
		}
	}
}

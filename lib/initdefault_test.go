/*
 * Copyright 2018 InfAI (CC SES)
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

package lib

import (
	"encoding/json"
	"fmt"
)

func defaultUseStr(msg string, defaultStr string) (result []map[string]interface{}, err error) {
	features := map[string]interface{}{}
	initDefault := InitDefault{}
	err = json.Unmarshal([]byte(msg), &features)
	if err != nil {
		return result, err
	}
	json.Unmarshal([]byte(defaultStr), &initDefault)
	if err != nil {
		return result, err
	}
	result = initDefault.Use(features)
	return
}

func ExampleInitDefault_Use1() {
	features := `{"a": "foo", "b": 42, "c":[1,2], "d": {"foo": "bar"}}`
	defaultStr := `[
              {
                "a":{"feature":"a"}, "b":{"value":"b"}, "c":{"feature":"c"}, "d":{"value": [3,4]}
              }
            ]`

	result, err := defaultUseStr(features, defaultStr)
	fmt.Println(len(result), result[0]["a"], result[0]["b"], result[0]["c"], result[0]["d"], err)

	//output:
	//1 foo b [1 2] [3 4] <nil>

}

func ExampleInitDefault_Use2() {
	features := `{"a": "foo", "b": 42, "c":[1,2], "d": {"foo": "bar"}}`
	defaultStr := `[
              {
                "a":{"feature":"a"}
              },
              {
                "foo":{"value":"bar"}
              },
              {
                "bar":{"feature":"bar"}
              }
            ]`

	fmt.Println(defaultUseStr(features, defaultStr))

	//output:
	//[map[a:foo] map[foo:bar] map[bar:<nil>]] <nil>
}

func ExampleInitDefault_Use3() {
	features := `{"a": "foo", "b": 42, "c":[1,2], "d": {"foo": "bar"}}`
	defaultStr := `[
              {
                "a":{"feature":"a"}
              },
              {
                "b":{"value":"bstr"}
              },
			  {
                "c":{"feature":"c"}
              },
			  {
                "d":{"value": [3,4]}
              }
            ]`

	fmt.Println(defaultUseStr(features, defaultStr))

	//output:
	//[map[a:foo] map[b:bstr] map[c:[1 2]] map[d:[3 4]]] <nil>
}

func ExampleInitDefault_Use4() {
	features := `{"a": "foo", "b": 42, "c":[1,2], "d": {"foo": "bar"}}`
	defaultStr := `[
              {
                "a":{"feature":"a"}
              },
              {
                "b":{"value":"bstr"}
              },
			  {
                "c":{"feature":"c"}
              },
			  {
                "d":{"feature": "d.foo"}
              }
            ]`

	fmt.Println(defaultUseStr(features, defaultStr))

	//output:
	//[map[a:foo] map[b:bstr] map[c:[1 2]] map[d:bar]] <nil>
}

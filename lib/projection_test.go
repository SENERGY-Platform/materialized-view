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

func checkProjectionUseStr(msg string, projectionsStr string) (result map[string]interface{}, err error) {
	features := map[string]interface{}{}
	projections := Projection{}
	err = json.Unmarshal([]byte(msg), &features)
	if err != nil {
		return result, err
	}
	json.Unmarshal([]byte(projectionsStr), &projections)
	if err != nil {
		return result, err
	}
	result = projections.Use(features)
	return
}

func ExampleProjection_Use() {
	features := `{"a": "foo", "b": 42, "c":[1,2], "d": {"foo": "bar"}}`
	fmt.Println(checkProjectionUseStr(features, ``))
	fmt.Println(checkProjectionUseStr(features, `[]`))
	fmt.Println(checkProjectionUseStr(features, `[""]`))
	fmt.Println(checkProjectionUseStr(features, `["x"]`))
	fmt.Println(checkProjectionUseStr(features, `["a"]`))
	fmt.Println(checkProjectionUseStr(features, `["b"]`))
	fmt.Println(checkProjectionUseStr(features, `["c"]`))
	fmt.Println(checkProjectionUseStr(features, `["d"]`))

	ab, err1 := checkProjectionUseStr(features, `["a", "b"]`)
	all, err2 := checkProjectionUseStr(features, `["*"]`)
	all_plus, err3 := checkProjectionUseStr(features, `["*", "a"]`)

	fmt.Println(len(ab), ab["a"], ab["b"], ab["c"], ab["foo"], err1)
	fmt.Println(len(all), all["a"], all["b"], all["c"], all["d"], all["foo"], err2)
	fmt.Println(len(all_plus), all_plus["a"], all_plus["b"], all_plus["c"], all_plus["d"], all_plus["foo"], err3)

	//output:
	//map[] <nil>
	//map[] <nil>
	//map[:<nil>] <nil>
	//map[x:<nil>] <nil>
	//map[a:foo] <nil>
	//map[b:42] <nil>
	//map[c:[1 2]] <nil>
	//map[d:map[foo:bar]] <nil>
	//2 foo 42 <nil> <nil> <nil>
	//4 foo 42 [1 2] map[foo:bar] <nil> <nil>
	//4 foo 42 [1 2] map[foo:bar] <nil> <nil>
}

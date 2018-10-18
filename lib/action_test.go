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
	"sort"
)

func ExampleActions_many() {
	target := Target{Id: "id", Features: map[string]interface{}{}, Name: "test"}

	actionsStr := `[
          {
            "type": "remove",
            "fields": ["list"],
            "scale": "many"
          },
          {
            "type": "insert",
            "if": [{"feature": "condition", "operation": "==", "value": "ok"}],
            "fields": ["list"],
            "scale": "many"
          }
        ]`

	actions := Actions{}
	err := json.Unmarshal([]byte(actionsStr), &actions)
	fmt.Println(err)

	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"not-ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"ok"}, map[string]interface{}{"element": "b"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "c", "condition":"ok"}, map[string]interface{}{"element": "c"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "d", "condition":"remove"}, map[string]interface{}{"element": "d"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"remove"}, map[string]interface{}{"element": "b"}))

	//output:
	//<nil>
	//{map[list:[]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:b]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:b] map[element:c]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:b] map[element:c]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:c]]] <nil> id true false false test} <nil>
}

func test_helper_print(m map[string]interface{}){
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Print(k, ":", m[k], ",")
	}
	fmt.Print("\n")
}

func ExampleActions_one() {
	target := Target{Id: "id", Features: map[string]interface{}{}, Name: "test"}

	actionsStr := `[
          {
            "type": "remove",
            "fields": ["", "field1"],
            "scale": "one"
          },
          {
            "type": "insert",
            "if": [{"feature": "condition", "operation": "==", "value": "ok"}],
            "fields": ["field1", "field2", ""],
            "scale": "one"
          },
          {
            "type": "insert",
            "fields": ["field3", "field4"],
            "scale": "one"
          }
        ]`

	actions := Actions{}
	err := json.Unmarshal([]byte(actionsStr), &actions)
	fmt.Println(err)

	result, err := actions.Do(target, map[string]interface{}{"element": "a", "condition":"not-ok"}, map[string]interface{}{"element": "a"})
	fmt.Println(err)
	test_helper_print(result.Features)

	result, err = actions.Do(target, map[string]interface{}{"element": "a", "condition":"ok"}, map[string]interface{}{"element": "a"})
	fmt.Println(err)
	test_helper_print(result.Features)

	result, err = actions.Do(target, map[string]interface{}{"element": "b", "condition":"ok"}, map[string]interface{}{"element": "b"})
	fmt.Println(err)
	test_helper_print(result.Features)

	result, err = actions.Do(target, map[string]interface{}{"element": "c", "condition":"ok"}, map[string]interface{}{"element": "c"})
	fmt.Println(err)
	test_helper_print(result.Features)

	result, err = actions.Do(target, map[string]interface{}{"element": "d", "condition":"remove"}, map[string]interface{}{"element": "d"})
	fmt.Println(err)
	test_helper_print(result.Features)

	result, err = actions.Do(target, map[string]interface{}{"element": "b", "condition":"remove"}, map[string]interface{}{"element": "b"})
	fmt.Println(err)
	test_helper_print(result.Features)


	//output:
	//<nil>
	//<nil>
	//field3:map[element:a],field4:map[element:a],
	//<nil>
	//element:a,field1:map[element:a],field2:map[element:a],field3:map[element:a],field4:map[element:a],
	//<nil>
	//element:b,field1:map[element:b],field2:map[element:b],field3:map[element:b],field4:map[element:b],
	//<nil>
	//element:c,field1:map[element:c],field2:map[element:c],field3:map[element:c],field4:map[element:c],
	//<nil>
	//field2:map[element:c],field3:map[element:d],field4:map[element:d],
	//<nil>
	//field2:map[element:c],field3:map[element:b],field4:map[element:b],
}

func ExampleActions_unknown_action() {
	target := Target{Id: "id", Features: map[string]interface{}{}, Name: "test"}

	actionsStr := `[
          {
            "type": "foo",
            "fields": ["list"],
            "scale": "many"
          },
          {
            "type": "foo",
            "if": [{"feature": "condition", "operation": "==", "value": "ok"}],
            "fields": ["list"],
            "scale": "many"
          }
        ]`

	actions := Actions{}
	err := json.Unmarshal([]byte(actionsStr), &actions)
	fmt.Println(err)

	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"not-ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"ok"}, map[string]interface{}{"element": "b"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "c", "condition":"ok"}, map[string]interface{}{"element": "c"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "d", "condition":"remove"}, map[string]interface{}{"element": "d"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"remove"}, map[string]interface{}{"element": "b"}))

	//output:
	//<nil>
	//{map[] <nil> id true false false test} unknown action type foo
	//{map[] <nil> id true false false test} unknown action type foo
	//{map[] <nil> id true false false test} unknown action type foo
	//{map[] <nil> id true false false test} unknown action type foo
	//{map[] <nil> id true false false test} unknown action type foo
	//{map[] <nil> id true false false test} unknown action type foo
}

func ExampleActions_unknown_scale() {
	target := Target{Id: "id", Features: map[string]interface{}{}, Name: "test"}

	actionsStr := `[
          {
            "type": "remove",
            "fields": ["list"],
            "scale": "foo"
          },
          {
            "type": "insert",
            "if": [{"feature": "condition", "operation": "==", "value": "ok"}],
            "fields": ["list"],
            "scale": "bar"
          }
        ]`

	actions := Actions{}
	err := json.Unmarshal([]byte(actionsStr), &actions)
	fmt.Println(err)

	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"not-ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"ok"}, map[string]interface{}{"element": "b"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "c", "condition":"ok"}, map[string]interface{}{"element": "c"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "d", "condition":"remove"}, map[string]interface{}{"element": "d"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"remove"}, map[string]interface{}{"element": "b"}))

	//output:
	//<nil>
	//{map[] <nil> id true false false test} unknown scale type foo
	//{map[] <nil> id true false false test} unknown scale type foo
	//{map[] <nil> id true false false test} unknown scale type foo
	//{map[] <nil> id true false false test} unknown scale type foo
	//{map[] <nil> id true false false test} unknown scale type foo
	//{map[] <nil> id true false false test} unknown scale type foo
}

func ExampleActions_remove_target() {
	target := Target{Id: "id", Features: map[string]interface{}{}, Name: "test"}

	actionsStr := `[
          {
            "type": "remove",
            "fields": ["list"],
            "scale": "many"
          },
          {
            "type": "insert",
            "if": [{"feature": "condition", "operation": "==", "value": "ok"}],
            "fields": ["list"],
            "scale": "many"
          },
          {
            "type": "remove_target",
            "if": [{"feature": "condition", "operation": "==", "value": "remove_target"}]
          }
        ]`

	actions := Actions{}
	err := json.Unmarshal([]byte(actionsStr), &actions)
	fmt.Println(err)

	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"not-ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "a", "condition":"ok"}, map[string]interface{}{"element": "a"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"ok"}, map[string]interface{}{"element": "b"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "c", "condition":"ok"}, map[string]interface{}{"element": "c"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "d", "condition":"remove"}, map[string]interface{}{"element": "d"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "b", "condition":"remove"}, map[string]interface{}{"element": "b"}))
	fmt.Println(actions.Do(target, map[string]interface{}{"element": "foo", "condition":"remove_target"}, map[string]interface{}{"element": "foo"}))

	//output:
	//<nil>
	//{map[list:[]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:b]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:b] map[element:c]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:b] map[element:c]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:c]]] <nil> id true false false test} <nil>
	//{map[list:[map[element:a] map[element:c]]] <nil> id true true false test} <nil>
}
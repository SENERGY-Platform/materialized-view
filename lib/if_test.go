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

func ifConditionsCheckStr(msg string, cond string) (result bool, err error) {
	features := map[string]interface{}{}
	conditions := IfConditions{}
	err = json.Unmarshal([]byte(msg), &features)
	if err != nil {
		return result, err
	}
	json.Unmarshal([]byte(cond), &conditions)
	if err != nil {
		return result, err
	}
	result = conditions.CheckFeatures(features)
	return
}

func ExampleIfConditions_Equal() {
	msg_1 := `{"a": "a", "b": 42, "foo":"bar"}`
	if_1 := `[{"feature": "a", "operation": "==", "value": "a"}]`
	if_2 := `[{"feature": "b", "operation": "==", "value": 42}]`
	if_3 := `[{"feature": "b", "operation": "==", "value": "a"}]`
	if_4 := `[{"feature": "a", "operation": "==", "value": 42}]`
	if_5 := `[{"feature": "c", "operation": "==", "value": 42}]`
	if_6 := `[{"feature": "a", "operation": "==", "value": "b"}]`
	if_7 := `[{"feature": "b", "operation": "==", "value": 13}]`
	if_8 := `[{"feature": "a", "operation": "==", "value": "a"},{"feature": "b", "operation": "==", "value": 42}]`
	if_9 := `[{"feature": "a", "operation": "==", "value": "a"},{"feature": "b", "operation": "==", "value": 13}]`
	if_10 := `[{"feature": "a", "operation": "==", "value": "a"},{"feature": "b", "operation": "==", "value": "nope"}]`
	fmt.Println(ifConditionsCheckStr(msg_1, if_1))
	fmt.Println(ifConditionsCheckStr(msg_1, if_2))
	fmt.Println(ifConditionsCheckStr(msg_1, if_3))
	fmt.Println(ifConditionsCheckStr(msg_1, if_4))
	fmt.Println(ifConditionsCheckStr(msg_1, if_5))
	fmt.Println(ifConditionsCheckStr(msg_1, if_6))
	fmt.Println(ifConditionsCheckStr(msg_1, if_7))
	fmt.Println(ifConditionsCheckStr(msg_1, if_8))
	fmt.Println(ifConditionsCheckStr(msg_1, if_9))
	fmt.Println(ifConditionsCheckStr(msg_1, if_10))

	//output:
	//true <nil>
	//true <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//true <nil>
	//false <nil>
	//false <nil>
}

func ExampleIfConditions_UnEqual() {
	msg_1 := `{"a": "a", "b": 42, "foo":"bar", "bool": true}`
	if_1 := `[{"feature": "a", "operation": "!=", "value": "a"}]`
	if_2 := `[{"feature": "b", "operation": "!=", "value": 42}]`
	if_3 := `[{"feature": "b", "operation": "!=", "value": "a"}]`
	if_4 := `[{"feature": "a", "operation": "!=", "value": 42}]`
	if_5 := `[{"feature": "c", "operation": "!=", "value": 42}]`
	if_6 := `[{"feature": "a", "operation": "!=", "value": "b"}]`
	if_7 := `[{"feature": "b", "operation": "!=", "value": 13}]`
	if_8 := `[{"feature": "a", "operation": "!=", "value": "a"},{"feature": "b", "operation": "!=", "value": 42}]`
	if_9 := `[{"feature": "a", "operation": "!=", "value": "a"},{"feature": "b", "operation": "!=", "value": 13}]`
	if_10 := `[{"feature": "a", "operation": "!=", "value": "a"},{"feature": "b", "operation": "!=", "value": "nope"}]`
	if_11 := `[{"feature": "a", "operation": "!=", "value": "b"},{"feature": "b", "operation": "!=", "value": "nope"}]`
	if_12 := `[{"feature": "c", "operation": "!=", "value": undefined}]`
	if_13 := `[{"feature": "c", "operation": "!=", "value": null}]`

	if_14 := `[{"feature": "bool", "operation": "!=", "value": null}]`
	if_15 := `[{"feature": "bool", "operation": "!=", "value": true}]`
	if_16 := `[{"feature": "bool", "operation": "!=", "value": false}]`
	if_17 := `[{"feature": "bool", "operation": "==", "value": true}]`
	if_18 := `[{"feature": "bool", "operation": "==", "value": false}]`

	if_19 := `[{"feature": "nobool", "operation": "!=", "value": true}]`
	if_20 := `[{"feature": "nobool", "operation": "!=", "value": false}]`
	if_21 := `[{"feature": "nobool", "operation": "==", "value": true}]`
	if_22 := `[{"feature": "nobool", "operation": "==", "value": false}]`

	fmt.Println(ifConditionsCheckStr(msg_1, if_1))
	fmt.Println(ifConditionsCheckStr(msg_1, if_2))
	fmt.Println(ifConditionsCheckStr(msg_1, if_3))
	fmt.Println(ifConditionsCheckStr(msg_1, if_4))
	fmt.Println(ifConditionsCheckStr(msg_1, if_5))
	fmt.Println(ifConditionsCheckStr(msg_1, if_6))
	fmt.Println(ifConditionsCheckStr(msg_1, if_7))
	fmt.Println(ifConditionsCheckStr(msg_1, if_8))
	fmt.Println(ifConditionsCheckStr(msg_1, if_9))
	fmt.Println(ifConditionsCheckStr(msg_1, if_10))
	fmt.Println(ifConditionsCheckStr(msg_1, if_11))
	fmt.Println(ifConditionsCheckStr(msg_1, if_12))
	fmt.Println(ifConditionsCheckStr(msg_1, if_13))

	fmt.Println(ifConditionsCheckStr(msg_1, if_14))
	fmt.Println(ifConditionsCheckStr(msg_1, if_15))
	fmt.Println(ifConditionsCheckStr(msg_1, if_16))
	fmt.Println(ifConditionsCheckStr(msg_1, if_17))
	fmt.Println(ifConditionsCheckStr(msg_1, if_18))
	fmt.Println(ifConditionsCheckStr(msg_1, if_19))
	fmt.Println(ifConditionsCheckStr(msg_1, if_20))
	fmt.Println(ifConditionsCheckStr(msg_1, if_21))
	fmt.Println(ifConditionsCheckStr(msg_1, if_22))

	//output:
	//false <nil>
	//false <nil>
	//true <nil>
	//true <nil>
	//true <nil>
	//true <nil>
	//true <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//true <nil>
	//true <nil>
	//false <nil>
	//true <nil>
	//false <nil>
	//true <nil>
	//true <nil>
	//false <nil>
	//true <nil>
	//true <nil>
	//false <nil>
	//false <nil>
}

func ExampleIfConditions_StringContains() {
	msg_1 := `{"a": "a", "b": 42, "foo":"bar", "right": "rw"}`
	if_1 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "r"}]`
	if_2 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "rw"}]`
	if_3 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "w"}]`
	if_4 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "r"}, {"feature": "right", "operation": "feature_str_contains_value", "value": "w"}]`
	if_5 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "rx"}]`
	if_6 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "x"}]`
	if_7 := `[{"feature": "b", "operation": "feature_str_contains_value", "value": "x"}]`
	if_8 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": 13}]`
	if_9 := `[{"feature": "c", "operation": "feature_str_contains_value", "value": "x"}]`
	if_10 := `[{"feature": "right", "operation": "feature_str_contains_value", "value": "r"}, {"feature": "right", "operation": "feature_str_contains_value", "value": 2}]`
	fmt.Println(ifConditionsCheckStr(msg_1, if_1))
	fmt.Println(ifConditionsCheckStr(msg_1, if_2))
	fmt.Println(ifConditionsCheckStr(msg_1, if_3))
	fmt.Println(ifConditionsCheckStr(msg_1, if_4))
	fmt.Println(ifConditionsCheckStr(msg_1, if_5))
	fmt.Println(ifConditionsCheckStr(msg_1, if_6))
	fmt.Println(ifConditionsCheckStr(msg_1, if_7))
	fmt.Println(ifConditionsCheckStr(msg_1, if_8))
	fmt.Println(ifConditionsCheckStr(msg_1, if_9))
	fmt.Println(ifConditionsCheckStr(msg_1, if_10))

	//output:
	//true <nil>
	//true <nil>
	//true <nil>
	//true <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//false <nil>
	//false <nil>
}

func ExampleIfConditions_unknown() {
	msg_1 := `{"a": "a", "b": 42, "foo":"bar", "right": "rw"}`
	if_1 := `[{"feature": "right", "operation": "foo", "value": "r"}]`
	if_2 := `[{"feature": "right", "value": "rw"}]`
	fmt.Println(ifConditionsCheckStr(msg_1, if_1))
	fmt.Println(ifConditionsCheckStr(msg_1, if_2))

	//output:
	//false <nil>
	//false <nil>
}

func ExampleIfConditions_Sub() {
	msg_1 := `{"a": {"b":"c"}}`
	if_1 := `[{"feature": "a.b", "operation": "==", "value": "c"}]`
	if_2 := `[{"feature": "a.b", "operation": "==", "value": "d"}]`
	if_3 := `[{"feature": "a.c", "operation": "==", "value": "a"}]`
	if_4 := `[{"feature": "a.c", "operation": "!=", "value": "a"}]`
	fmt.Println(ifConditionsCheckStr(msg_1, if_1))
	fmt.Println(ifConditionsCheckStr(msg_1, if_2))
	fmt.Println(ifConditionsCheckStr(msg_1, if_3))
	fmt.Println(ifConditionsCheckStr(msg_1, if_4))

	//output:
	//true <nil>
	//false <nil>
	//false <nil>
	//true <nil>
}

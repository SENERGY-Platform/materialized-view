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
	"fmt"
	"time"

	"github.com/bouk/monkey"
)

func ExampleMsgToFeatures() {
	msg1 := `{"a": "abc"}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", Temp: true}}, []byte(msg1)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", Temp: true}}, []byte(msg1)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", Temp: true, Omitempty: true}}, []byte(msg1)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", Temp: true, Omitempty: true}}, []byte(msg1)))

	msg2 := `{"b": ["abc"]}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.b+"}}, []byte(msg2)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.b+", Temp: true}}, []byte(msg2)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.c+"}}, []byte(msg2)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.c+", Temp: true}}, []byte(msg2)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.b+", Omitempty: true}}, []byte(msg2)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.b+", Temp: true, Omitempty: true}}, []byte(msg2)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.c+", Omitempty: true}}, []byte(msg2)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.c+", Temp: true, Omitempty: true}}, []byte(msg2)))

	msg3 := `{"c": []}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.c+"}}, []byte(msg3)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.c+", Temp: true}}, []byte(msg3)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.d+"}}, []byte(msg3)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.d+", Temp: true}}, []byte(msg3)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.c+", Omitempty: true}}, []byte(msg3)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.c+", Temp: true, Omitempty: true}}, []byte(msg3)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.d+", Omitempty: true}}, []byte(msg3)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "c", Path: "$.d+", Temp: true, Omitempty: true}}, []byte(msg3)))

	msg4 := `{"d": ""}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.d+"}}, []byte(msg4)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.d+", Temp: true}}, []byte(msg4)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.e+"}}, []byte(msg4)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.e+", Temp: true}}, []byte(msg4)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.d+", Omitempty: true}}, []byte(msg4)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.d+", Temp: true, Omitempty: true}}, []byte(msg4)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.e+", Omitempty: true}}, []byte(msg4)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.e+", Temp: true, Omitempty: true}}, []byte(msg4)))

	msg5 := `{"e": 0}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.e+"}}, []byte(msg5)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.e+", Temp: true}}, []byte(msg5)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.f+"}}, []byte(msg5)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.f+", Temp: true}}, []byte(msg5)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.e+", Omitempty: true}}, []byte(msg5)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.e+", Temp: true, Omitempty: true}}, []byte(msg5)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.f+", Omitempty: true}}, []byte(msg5)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.f+", Temp: true, Omitempty: true}}, []byte(msg5)))

	msg6 := `{"f":[{"g": 1}, {"g": 2}]}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "f", Path: "$.f[*].g+"}}, []byte(msg6)))

	//output:
	//map[a:abc] map[a:abc] <nil>
	//map[a:abc] map[] <nil>
	//map[a:<nil>] map[a:<nil>] <nil>
	//map[a:<nil>] map[] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:abc] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[b:[abc]] map[b:[abc]] <nil>
	//map[b:[abc]] map[] <nil>
	//map[b:<nil>] map[b:<nil>] <nil>
	//map[b:<nil>] map[] <nil>
	//map[b:[abc]] map[b:[abc]] <nil>
	//map[b:[abc]] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[c:[]] map[c:[]] <nil>
	//map[c:[]] map[] <nil>
	//map[c:<nil>] map[c:<nil>] <nil>
	//map[c:<nil>] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[d:] map[d:] <nil>
	//map[d:] map[] <nil>
	//map[d:<nil>] map[d:<nil>] <nil>
	//map[d:<nil>] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[e:0] map[e:0] <nil>
	//map[e:0] map[] <nil>
	//map[e:<nil>] map[e:<nil>] <nil>
	//map[e:<nil>] map[] <nil>
	//map[e:0] map[e:0] <nil>
	//map[e:0] map[] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[f:[1 2]] map[f:[1 2]] <nil>
}

func ExampleMsgBooleanToFeatures() {
	msg1 := `{"a": true, "b": false, "c": {"d": true, "e": false}}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "b", Path: "$.b+"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "d", Path: "$.c.d+"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "e", Path: "$.c.e+"}}, []byte(msg1)))

	msg2 := `{
		"command": "PUT",
		"id": "process_id_2",
		"processmodel":{
			"process":{"definitions":{"process":{"_id":"process_name_2"}}},
			"publish": true
		}
	}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "publish", Path: "$.processmodel.publish+"}}, []byte(msg2)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "name", Path: "$.processmodel.process.definitions.process._id+"}}, []byte(msg2)))

	//output:
	//map[a:true] map[a:true] <nil>
	//map[b:false] map[b:false] <nil>
	//map[d:true] map[d:true] <nil>
	//map[e:false] map[e:false] <nil>
	//map[publish:true] map[publish:true] <nil>
	//map[name:process_name_2] map[name:process_name_2] <nil>
}

func ExampleMsgToFeaturesDefaults() {
	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	msg1 := `{"a": "abc"}`
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", Default: "default"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", Default: "default"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Default: "default"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", Default: "default", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", Default: "default", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Default: "default", Omitempty: true}}, []byte(msg1)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", DefaultRef: "default"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", DefaultRef: "default"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", DefaultRef: "default"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", DefaultRef: "default", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", DefaultRef: "default", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", DefaultRef: "default", Omitempty: true}}, []byte(msg1)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", DefaultRef: "time.epoch_millis"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", DefaultRef: "time.epoch_millis"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", DefaultRef: "time.epoch_millis"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", DefaultRef: "time.epoch_millis", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", DefaultRef: "time.epoch_millis", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", DefaultRef: "time.epoch_millis", Omitempty: true}}, []byte(msg1)))

	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", DefaultRef: "time.epoch_second"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", DefaultRef: "time.epoch_second"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", DefaultRef: "time.epoch_second"}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.a+", DefaultRef: "time.epoch_second", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", Path: "$.b+", DefaultRef: "time.epoch_second", Omitempty: true}}, []byte(msg1)))
	fmt.Println(MsgToFeatures([]Feature{{Name: "a", DefaultRef: "time.epoch_second", Omitempty: true}}, []byte(msg1)))

	//output:
	//map[a:abc] map[a:abc] <nil>
	//map[a:default] map[a:default] <nil>
	//map[a:default] map[a:default] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:default] map[a:default] <nil>
	//map[a:default] map[a:default] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:<nil>] map[a:<nil>] <nil>
	//map[a:<nil>] map[a:<nil>] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[] map[] <nil>
	//map[] map[] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:138157323000] map[a:138157323000] <nil>
	//map[a:138157323000] map[a:138157323000] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:138157323000] map[a:138157323000] <nil>
	//map[a:138157323000] map[a:138157323000] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:138157323] map[a:138157323] <nil>
	//map[a:138157323] map[a:138157323] <nil>
	//map[a:abc] map[a:abc] <nil>
	//map[a:138157323] map[a:138157323] <nil>
	//map[a:138157323] map[a:138157323] <nil>
}

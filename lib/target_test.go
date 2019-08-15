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
	"time"

	"testing"
)

var testConfStrTargetWhereSorted = `{
  "server_port":		          "8080",
  "log_level":		              "CALL",

    "zookeeper_url": "zk",
  "consumer_group": "matview_1",
  "debug": false,

  "force_user": "false",
  "force_auth": "false",

  "elastic_url": "http://elastic:9200",
  "elastic_retry": 3,

  "events": {
    "test":[
      {
        "type": "root",
        "target": "test",
        "features": [
          {"name": "id", "path": "$.id+"},
          {"name": "s", "default_ref":"time.epoch_second"},
          {"name": "ms", "default_ref": "time.epoch_millis"}
        ],
        "actions": [
          {
            "type": "insert",
            "fields": [""],
            "scale": "one"
          }
        ]
      }
    ]
  },

  "elastic_mapping": {
    "test":{
		"id":         {"type": "keyword"},
		"s": {"type":   "date","format": "epoch_second"},
		"ms": {"type":   "date","format": "epoch_millis"}
    }
  }
}
`

func initTargetWherSortedTestContainer() (purge func(), err error) {
	return initTestContainer(testConfStrTargetWhereSorted)
}

func testHelperGetTargetsWhereSorted(t *testing.T, sortingConf string, expectedIds []string) {
	t.Helper()
	sorting := Sorting{}
	err := json.Unmarshal([]byte(sortingConf), &sorting)
	if err != nil {
		t.Fatal(err)
	}
	result, err := GetTargetsWhereSorted("test", WhereConditions{}, Features{}, sorting)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(expectedIds) {
		t.Fatal("len unequal", len(result), len(expectedIds))
	}
	for index, target := range result {
		if _, ok := target.Features["ms"]; !ok {
			t.Fatal("ms not set")
		}
		if _, ok := target.Features["ms"]; !ok {
			t.Fatal("s not set")
		}
		if target.Features["id"] != expectedIds[index] {
			t.Fatal("id unequal", target.Features["id"], expectedIds[index], result, expectedIds)
		}
	}
}

func TestGetTargetsWhereSorted(t *testing.T) {
	purge, err := initTargetWherSortedTestContainer()
	defer purge()
	if err != nil {
		t.Fatal(err)
	}
	err = test_sendEvent("test", map[string]interface{}{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	err = test_sendEvent("test", map[string]interface{}{
		"id": "2",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	err = test_sendEvent("test", map[string]interface{}{
		"id": "3",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	testHelperGetTargetsWhereSorted(t, `{"by": "s", "asc":true}`, []string{"1", "2", "3"})
	testHelperGetTargetsWhereSorted(t, `{"by": "s", "asc":false}`, []string{"3", "2", "1"})

	testHelperGetTargetsWhereSorted(t, `{"by": "s", "asc":true, "limit":2}`, []string{"1", "2"})
	testHelperGetTargetsWhereSorted(t, `{"by": "s", "asc":false, "limit": 2}`, []string{"3", "2"})

	testHelperGetTargetsWhereSorted(t, `{"by": "ms", "asc":true}`, []string{"1", "2", "3"})
	testHelperGetTargetsWhereSorted(t, `{"by": "ms", "asc":false}`, []string{"3", "2", "1"})

	testHelperGetTargetsWhereSorted(t, `{"by": "ms", "asc":true, "limit":2}`, []string{"1", "2"})
	testHelperGetTargetsWhereSorted(t, `{"by": "ms", "asc":false, "limit": 2}`, []string{"3", "2"})
}

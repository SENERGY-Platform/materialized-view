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
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/SmartEnergyPlatform/util/http/request"
)

var testConfStr = `{
  "server_port":		          "8080",
  "log_level":		              "CALL",

  "amqp_url": "amqp://guest:guest@rabbitmq:5672/",
  "amqp_consumer_name": "matview_1",
  "AmqpReconnectTimeout": 10,

  "force_user": "false",
  "force_auth": "false",

  "elastic_url": "http://elastic:9200",
  "elastic_retry": 3,

  "events": {
    "permission":[
      {
        "type":"child",
        "target": "deviceinstance",
        "where": [{"target_feature": "device.id", "operation":"==", "event_feature":"resource"}],
        "if": [{"feature": "kind", "operation": "==", "value": "deviceinstance"}],
        "features": [
          {"name": "command", "path": "$.command+", "temp":true},
          {"name": "user", "path": "$.User+", "omitempty": true},
          {"name": "right", "path": "$.Right+", "temp":true},
          {"name": "kind", "path": "$.Kind+", "temp":true},
          {"name": "resource", "path": "$.Resource+", "temp":true}
        ],
        "actions": [
          {
            "type": "remove",
            "fields": ["admin", "read", "write", "execute"],
            "scale": "many"
          },
          {
            "type": "insert",
            "if": [
              {"feature": "command", "operation": "==", "value": "PUT"},
              {"feature": "right", "operation": "feature_str_contains_value", "value": "a"}
            ],
            "fields": ["admin"],
            "scale": "many"
          },
          {
            "type": "insert",
            "if": [
              {"feature": "command", "operation": "==", "value": "PUT"},
              {"feature": "right", "operation": "feature_str_contains_value", "value": "r"}
            ],
            "fields": ["read"],
            "scale": "many"
          },
          {
            "type": "insert",
            "if": [
              {"feature": "command", "operation": "==", "value": "PUT"},
              {"feature": "right", "operation": "feature_str_contains_value", "value": "w"}
            ],
            "fields": ["write"],
            "scale": "many"
          },
          {
            "type": "insert",
            "if": [
              {"feature": "command", "operation": "==", "value": "PUT"},
              {"feature": "right", "operation": "feature_str_contains_value", "value": "x"}
            ],
            "fields": ["execute"],
            "scale": "many"
          }
        ]
      },
      {
        "type": "root",
        "target": "permission_event",
        "features": [
          {"name": "command", "path": "$.command+"},
          {"name": "user", "path": "$.User+", "omitempty": true},
          {"name": "right", "path": "$.Right+"},
          {"name": "kind", "path": "$.Kind+"},
          {"name": "resource", "path": "$.Resource+"}
        ],
        "actions": [
          {
            "type": "insert",
            "fields": [""],
            "scale": "one"
          }
        ]
      }
    ],

    "deviceinstance":[
      {
        "type": "root",
        "target": "deviceinstance",
        "id_feature":"id",
        "features": [
          {"name": "command", "path": "$.command+", "temp":true},
          {"name": "owner", "path": "$.owner+", "temp":true},
          {"name": "id", "path": "$.id+"},
          {"name": "name", "path": "$.device_instance.name+"}
        ],
        "actions": [
          {
            "type": "remove_target",
            "if": [{"feature": "command", "operation": "==", "value": "DELETE"}]
          },
          {
            "type": "insert",
            "if": [{"feature": "command", "operation": "==", "value": "PUT"}],
            "fields": ["device"],
            "scale": "one"
          }
        ],
        "init":[
          {
            "target":"permission_event",
            "where":[
              {"target_feature": "resource", "operation":"==", "event_feature":"id"},
              {"target_feature": "kind", "operation":"==", "value":"deviceinstance"}
            ],
            "default": [{"user":{"feature":"owner"}, "right":{"value":"rwxa"}, "command":{"value":"PUT"}}],
            "transform": [
              {"name": "right", "path": "$.right+", "temp":true},
			  {"name": "command", "path": "$.command+", "temp":true},
              {"Name": "user", "Path": "$.user+"}
            ],
            "actions": [
              {
                "type": "remove",
                "fields": ["admin", "read", "write", "execute"],
                "scale": "many"
              },
              {
                "type": "insert",
                "if": [
                  {"feature": "command", "operation": "==", "value": "PUT"},
                  {"feature": "right", "operation": "feature_str_contains_value", "value": "a"}
                ],
                "fields": ["admin"],
                "scale": "many"
              },
              {
                "type": "insert",
                "if": [
                  {"feature": "command", "operation": "==", "value": "PUT"},
                  {"feature": "right", "operation": "feature_str_contains_value", "value": "r"}
                ],
                "fields": ["read"],
                "scale": "many"
              },
              {
                "type": "insert",
                "if": [
                  {"feature": "command", "operation": "==", "value": "PUT"},
                  {"feature": "right", "operation": "feature_str_contains_value", "value": "w"}
                ],
                "fields": ["write"],
                "scale": "many"
              },
              {
                "type": "insert",
                "if": [
                  {"feature": "command", "operation": "==", "value": "PUT"},
                  {"feature": "right", "operation": "feature_str_contains_value", "value": "x"}
                ],
                "fields": ["execute"],
                "scale": "many"
              }
            ]
          }
        ]
      }
    ]
  },

  "queries": {
    "deviceinstance":{
      "r": {
        "selection":{
          "or": [
            {"condition": {"feature": "read.user", "operation": "==", "ref": "user"}}
          ]
        },
        "projection": ["device"]
      },
      "w": {
        "selection":{
          "or": [
            {"condition": {"feature": "write.user", "operation": "==", "ref": "user"}}
          ]
        },
        "projection": ["device"]
      },
      "x": {
        "selection":{
          "or": [
            {"condition": {"feature": "execute.user", "operation": "==", "ref": "user"}}
          ]
        },
        "projection": ["device"]
      },
      "a": {
        "selection":{
          "or": [
            {"condition": {"feature": "admin.user", "operation": "==", "ref": "user"}}
          ]
        },
        "projection": ["device"]
      }
    }
  },

  "elastic_mapping": {
    "deviceinstance":{
      "device":{
        "properties": {
          "id":         {"type": "keyword"},
          "name":         {"type": "keyword", "copy_to": "feature_search"},
          "description":  {"type": "text",    "copy_to": "feature_search"},
          "usertag":      {"type": "keyword", "copy_to": "feature_search"},
          "tag":          {"type": "keyword", "copy_to": "feature_search"},
          "devicetype":   {"type": "keyword"},
          "uri":          {"type": "keyword"},
          "img":          {"type": "keyword"}
        }
      },
      "read":{
        "properties": {
          "user":         {"type": "keyword"}
        }
      },
      "write":{
        "properties": {
          "user":         {"type": "keyword"}
        }
      },
      "admin":{
        "properties": {
          "user":         {"type": "keyword"}
        }
      },
      "execute":{
        "properties": {
          "user":         {"type": "keyword"}
        }
      }
    },
    "permission_event": {
      "resource": {"type": "keyword"},
      "user":     {"type": "keyword"},
      "kind":     {"type": "keyword"},
      "command":  {"type": "keyword"}
    }
  }
}
`

func testGetFrePort() string {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	parts := strings.Split(l.Addr().String(), ":")
	return parts[len(parts)-1]
}

func initIntegrationTestContainer() (purge func(), err error) {
	config := ConfigStruct{}
	err = json.Unmarshal([]byte(testConfStr), &config)
	if err != nil {
		log.Fatalf("Could not unmarshal config: %s", err)
	}
	Config = &config
	Config.ServerPort = testGetFrePort()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dockerrabbitmq, err := pool.Run("rabbitmq", "3-management", []string{})
	if err != nil {
		log.Fatalf("Could not start dockerrabbitmq: %s", err)
	}
	dockeresultelastic, err := pool.Run("elasticsearch", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start dockeresultelastic: %s", err)
	}
	purge = func() {
		conn.Close()
		pool.Purge(dockeresultelastic)
		pool.Purge(dockerrabbitmq)
	}

	time.Sleep(2 * time.Second)

	Config.ElasticUrl = "http://localhost:" + dockeresultelastic.GetPort("9200/tcp")
	Config.AmqpUrl = "amqp://guest:guest@localhost:" + dockerrabbitmq.GetPort("5672/tcp") + "/"

	if err := pool.Retry(func() error {
		localclient, err := elastic.NewClient(elastic.SetURL(Config.ElasticUrl), elastic.SetRetrier(newRetrier()))
		if err != nil {
			return err
		}
		ping, _, err := elastic.NewPingService(localclient).Do(context.Background())
		if err != nil {
			return err
		}
		if ping.Version.Number == "" {
			return errors.New("empty ping result")
		}
		GetClient()
		client = createClient()
		log.Println(Config.ElasticUrl, client)
		return nil
	}); err != nil {
		purge()
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err := pool.Retry(func() error {
		return InitEventHandling()
	}); err != nil {
		purge()
		log.Fatalf("Could not connect to docker: %s", err)
	}
	go StartApi()
	return
}

func testHelperGetAll(index string) (all []map[string]interface{}, err error) {
	allQuery := elastic.NewMatchAllQuery()
	resp, err := GetClient().Search().Index(index).Type(ElasticResourceType).Query(allQuery).Do(context.Background())
	if err != nil {
		return all, err
	}
	for _, hit := range resp.Hits.Hits {
		value := map[string]interface{}{}
		json.Unmarshal(*hit.Source, &value)
		all = append(all, value)
	}
	return
}

func testHelperCompare(a interface{}, b interface{}) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	json.Unmarshal(ab, &a)
	json.Unmarshal(bb, &b)
	return reflect.DeepEqual(a, b)
}

func testHelperCheckHttpGet(t *testing.T, path string, expected interface{}) {
	t.Helper()
	var result interface{}
	err := request.Get("http://localhost:"+Config.ServerPort+path, &result)
	if err != nil {
		t.Fatal(err)
	}
	if !testHelperCompare(result, expected) {
		t.Fatal("\ngot:\n", result, "\nwant:\n", expected)
	}
}

func testHelperCheckHttpPost(t *testing.T, path string, body interface{}, expected interface{}) {
	t.Helper()
	var result interface{}
	err, _, _ := request.Post("http://localhost:"+Config.ServerPort+path, body, &result)
	if err != nil {
		t.Fatal(err)
	}
	if !testHelperCompare(result, expected) {
		t.Fatal("\ngot:\n", result, "\nwant:\n", expected)
	}
}

func TestEventsAndRest(t *testing.T) {
	purge, err := initIntegrationTestContainer()
	defer purge()
	if err != nil {
		t.Fatal(err)
	}
	err = sendEvent("permission", map[string]interface{}{
		"command":  "PUT",
		"User":     "user1",
		"Right":    "rw",
		"Kind":     "deviceinstance",
		"Resource": "device1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	all, err := testHelperGetAll("permission_event")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{
		"command":  "PUT",
		"user":     "user1",
		"right":    "rw",
		"kind":     "deviceinstance",
		"resource": "device1",
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], expected) {
		t.Fatal(all)
	}
	all, err = testHelperGetAll("deviceinstance")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Fatal()
	}

	err = sendEvent("deviceinstance", map[string]interface{}{
		"command":         "PUT",
		"id":              "device1",
		"device_instance": map[string]string{"name": "device_name"},
		"owner":           "user1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	err = sendEvent("deviceinstance", map[string]interface{}{
		"command": "DELETE",
		"id":      "UNKNOWN",
	})

	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	all, err = testHelperGetAll("deviceinstance")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) == 0 {
		t.Fatal()
	}
	expected = map[string]interface{}{
		"device": map[string]interface{}{
			"id":   "device1",
			"name": "device_name",
		},
		"read":    []map[string]interface{}{{"user": "user1"}},
		"write":   []map[string]interface{}{{"user": "user1"}},
		"admin":   []interface{}{},
		"execute": []interface{}{},
	}
	if len(all) != 1 || !testHelperCompare(all[0], expected) {
		t.Fatal("\ngot:\n", all[0], "\nwant:\n", expected)
	}

	err = sendEvent("permission", map[string]interface{}{
		"command":  "PUT",
		"User":     "user2",
		"Right":    "rx",
		"Kind":     "deviceinstance",
		"Resource": "device1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	all, err = testHelperGetAll("deviceinstance")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) == 0 {
		t.Fatal()
	}
	expected = map[string]interface{}{
		"device": map[string]interface{}{
			"id":   "device1",
			"name": "device_name",
		},
		"read":    []map[string]interface{}{{"user": "user1"}, {"user": "user2"}},
		"write":   []map[string]interface{}{{"user": "user1"}},
		"admin":   []map[string]interface{}{},
		"execute": []map[string]interface{}{{"user": "user2"}},
	}
	if len(all) != 1 || !testHelperCompare(all[0], expected) {
		t.Fatal("\ngot:\n", all[0], "\nwant:\n", expected)
	}

	err = sendEvent("permission", map[string]interface{}{
		"command":  "DELETE",
		"User":     "user2",
		"Kind":     "deviceinstance",
		"Resource": "device1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	all, err = testHelperGetAll("deviceinstance")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) == 0 {
		t.Fatal()
	}
	expected = map[string]interface{}{
		"device": map[string]interface{}{
			"id":   "device1",
			"name": "device_name",
		},
		"read":    []map[string]interface{}{{"user": "user1"}},
		"write":   []map[string]interface{}{{"user": "user1"}},
		"admin":   []map[string]interface{}{},
		"execute": []map[string]interface{}{},
	}
	if len(all) != 1 || !testHelperCompare(all[0], expected) {
		t.Fatal("\ngot:\n", all[0], "\nwant:\n", expected)
	}

	err = sendEvent("deviceinstance", map[string]interface{}{
		"command":         "PUT",
		"id":              "device2",
		"device_instance": map[string]string{"name": "device_name_2"},
		"owner":           "user1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	all, err = testHelperGetAll("deviceinstance")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) == 0 {
		t.Fatal()
	}
	expected2 := map[string]interface{}{
		"device": map[string]interface{}{
			"id":   "device2",
			"name": "device_name_2",
		},
		"read":    []map[string]interface{}{{"user": "user1"}},
		"write":   []map[string]interface{}{{"user": "user1"}},
		"admin":   []map[string]interface{}{{"user": "user1"}},
		"execute": []map[string]interface{}{{"user": "user1"}},
	}
	if len(all) != 2 {
		t.Fatal("\ngot:\n", all, "\nwant:\n", expected, expected2)
	}
	if !testHelperCompare(all[1], expected) {
		t.Fatal("\ngot:\n", all[1], "\nwant:\n", expected)
	}
	if !testHelperCompare(all[0], expected2) {
		t.Fatal("\ngot:\n", all[0], "\nwant:\n", expected2)
	}

	err = sendEvent("deviceinstance", map[string]interface{}{
		"command": "DELETE",
		"id":      "device2",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	all, err = testHelperGetAll("deviceinstance")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatal("\ngot:\n", all, "\nwant:\n", expected)
	}
	if !testHelperCompare(all[0], expected) {
		t.Fatal("\ngot:\n", all[0], "\nwant:\n", expected)
	}

	err = sendEvent("deviceinstance", map[string]interface{}{
		"command":         "PUT",
		"id":              "device3",
		"device_instance": map[string]string{"name": "device_name_3"},
		"owner":           "user3",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sendEvent("deviceinstance", map[string]interface{}{
		"command":         "PUT",
		"id":              "device4",
		"device_instance": map[string]string{"name": "device_name_4"},
		"owner":           "user3",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r", []interface{}{})
	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r?user=user1", []interface{}{map[string]interface{}{
		"device": map[string]interface{}{
			"id":   "device1",
			"name": "device_name",
		},
	}})
	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		}})
	testHelperCheckHttpGet(t, "/search/deviceinstance/devx/r?user=user1", []interface{}{})

	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r/1/0?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
	})

	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r/1/1?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r/1/0/device.name/asc?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
	})

	testHelperCheckHttpGet(t, "/search/deviceinstance/dev/r/1/0/device.name/desc?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpGet(t, "/get/deviceinstance/r?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpGet(t, "/get/deviceinstance/r/1/1?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpGet(t, "/get/deviceinstance/r/1/0/device.name/asc?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
	})

	testHelperCheckHttpGet(t, "/get/deviceinstance/r/1/0/device.name/desc?user=user3", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpGet(t, "/get/deviceinstance/r/1/0/device.name/desc?user=foo", []interface{}{})

	err = sendEvent("deviceinstance", map[string]interface{}{
		"command":         "PUT",
		"id":              "device5",
		"device_instance": map[string]string{"name": "somename"},
		"owner":           "user5",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sendEvent("deviceinstance", map[string]interface{}{
		"command":         "PUT",
		"id":              "device6",
		"device_instance": map[string]string{"name": "somename"},
		"owner":           "user5",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/nope?user=user5", []interface{}{})

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/somename?user=user5", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device5",
				"name": "somename",
			},
		},
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device6",
				"name": "somename",
			},
		},
	})

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/somename/5/0?user=user5", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device5",
				"name": "somename",
			},
		},
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device6",
				"name": "somename",
			},
		},
	})

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/somename/1/0?user=user5", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device5",
				"name": "somename",
			},
		},
	})

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/somename/1/1?user=user5", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device6",
				"name": "somename",
			},
		},
	})

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/somename/1/0/device.id/asc?user=user5", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device5",
				"name": "somename",
			},
		},
	})

	testHelperCheckHttpGet(t, "/select/field/deviceinstance/r/device.name/somename/1/0/device.id/desc?user=user5", []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device6",
				"name": "somename",
			},
		},
	})

	testHelperCheckHttpPost(t, "/select/field/deviceinstance/r/device.id?user=user3", []string{"device3", "device4"}, []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpPost(t, "/select/field/deviceinstance/r/device.id/1/0?user=user3", []string{"device3", "device4"}, []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
	})

	testHelperCheckHttpPost(t, "/select/field/deviceinstance/r/device.id/1/1?user=user3", []string{"device3", "device4"}, []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})

	testHelperCheckHttpPost(t, "/select/field/deviceinstance/r/device.id/1/0/device.name/asc?user=user3", []string{"device3", "device4"}, []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device3",
				"name": "device_name_3",
			},
		},
	})

	testHelperCheckHttpPost(t, "/select/field/deviceinstance/r/device.id/1/0/device.name/desc?user=user3", []string{"device3", "device4"}, []interface{}{
		map[string]interface{}{
			"device": map[string]interface{}{
				"id":   "device4",
				"name": "device_name_4",
			},
		},
	})
}

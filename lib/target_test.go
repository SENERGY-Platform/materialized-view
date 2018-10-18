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
	"time"

	"testing"

	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
)

var testConfStrTargetWhereSorted = `{
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
	config := ConfigStruct{}
	err = json.Unmarshal([]byte(testConfStrTargetWhereSorted), &config)
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
	return
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
	err = sendEvent("test", map[string]interface{}{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	err = sendEvent("test", map[string]interface{}{
		"id": "2",
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	err = sendEvent("test", map[string]interface{}{
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

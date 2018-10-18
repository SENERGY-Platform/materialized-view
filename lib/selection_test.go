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
	"errors"
	"log"

	"testing"

	"time"

	"encoding/json"

	"net/url"

	"reflect"

	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

var testDefaultConfig = `{
	"server_port":		          "8080",
	"log_level":		              "CALL",

	"amqp_url": "amqp://guest:guest@rabbitmq:5672/",
	"amqp_consumer_name": "matview_1",
	"AmqpReconnectTimeout": 10,

	"force_user": "true",
	"force_auth": "true",

    "elastic_url": "http://elastic:9200",
    "elastic_retry": 3,

    "elastic_mapping": {
        "test": {
            "searchable":         {"type": "keyword", "copy_to": "feature_search"},
            "selectable":         {"type": "keyword"},
            "number":             {"type": "double"},
            "list": {
                "properties": {
                    "a": {"type": "keyword"},
                    "b": {"type": "keyword"}
                }
            }
        }
    }
}
`

func initElasticSelectionTest() (elasticClient *elastic.Client, purge func(), err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dockeresultelastic, err := pool.Run("elasticsearch", "latest", []string{})
	purge = func() {
		pool.Purge(dockeresultelastic)
	}
	if err != nil {
		log.Fatalf("Could not start dockeresultelastic: %s", err)
	}
	config := ConfigStruct{}
	err = json.Unmarshal([]byte(testDefaultConfig), &config)
	if err != nil {
		log.Fatalf("Could not unmarshal config: %s", err)
	}
	Config = &config
	Config.ElasticUrl = "http://localhost:" + dockeresultelastic.GetPort("9200/tcp")

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
		elasticClient = createClient()
		client = elasticClient
		log.Println(Config.ElasticUrl, client)
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return
}

func selectionTestResult(endpoint string, jwt jwt_http_router.Jwt, values url.Values) (all []string, hits []string, err error) {
	query := elastic.NewBoolQuery()
	err = UseSelection(query, "test", endpoint, jwt, values)
	if err != nil {
		return all, hits, err
	}
	resp, err := GetClient().Search().Index("test").Type(ElasticResourceType).Query(query).Do(context.Background())
	if err != nil {
		return all, hits, err
	}
	for _, hit := range resp.Hits.Hits {
		hits = append(hits, string(*hit.Source))
	}

	allQuery := elastic.NewMatchAllQuery()
	resp, err = GetClient().Search().Index("test").Type(ElasticResourceType).Query(allQuery).Do(context.Background())
	if err != nil {
		return all, hits, err
	}
	for _, hit := range resp.Hits.Hits {
		all = append(all, string(*hit.Source))
	}

	return
}

func testSelectionCheck(t *testing.T, result_all []string, result_hits []string, result_err error, expected_all []string, expected_hits []string) {
	t.Helper()
	if result_err != nil {
		t.Fatal("result err:", result_err)
	}
	if len(result_all) != len(expected_all) || (len(result_all) > 0 && !reflect.DeepEqual(result_all, expected_all)) {
		t.Fatal("result_all != expected_all:", result_all, expected_all)
	}
	if len(result_hits) != len(expected_hits) || (len(result_hits) > 0 && !reflect.DeepEqual(result_hits, expected_hits)) {
		t.Fatal("result_hits != expected_hits:", result_hits, expected_hits)
	}
}

func TestUseSelection(t *testing.T) {
	client, purge, err := initElasticSelectionTest()
	defer purge()
	if err != nil {
		t.Fatal()
	}
	element := `{
		"searchable": 	"foo",
		"selectable": 	"a",
		"number":     	42,
		"user": 		"userid",
		"list": [
			{"a": "a"},
			{"b": "b"},
			{"a": "c", "b": "d"}
		]
	}`
	_, err = client.Index().Index("test").Type(ElasticResourceType).BodyString(element).Do(context.Background())
	time.Sleep(time.Second * 2)
	if err != nil {
		t.Fatal()
	}

	queriesStr := `{
    	"test":{
      		"as_is": {
        		"selection": {"all": true},
        		"projection": ["*"]
      		},
			"eq_val_a":{
				"selection": {"condition": {"feature": "selectable", "operation": "==", "value": "a"}},
        		"projection": ["*"]
			},
			"eq_val_b":{
				"selection": {"condition": {"feature": "selectable", "operation": "==", "value": "b"}},
        		"projection": ["*"]
			},
			"uneq_val_b":{
				"selection": {"condition": {"feature": "selectable", "operation": "!=", "value": "b"}},
        		"projection": ["*"]
			},
			"uneq_val_a":{
				"selection": {"condition": {"feature": "selectable", "operation": "!=", "value": "a"}},
        		"projection": ["*"]
			},
			"list_val_ab":{
				"selection": {"condition": {"feature": "selectable", "operation": "any_value_in_feature", "value": ["a", "b"]}},
        		"projection": ["*"]
			},
			"list_val_xb":{
				"selection": {"condition": {"feature": "selectable", "operation": "any_value_in_feature", "value": ["x", "b"]}},
        		"projection": ["*"]
			},
			"list_list_val_xb":{
				"selection": {"condition": {"feature": "list.a", "operation": "any_value_in_feature", "value": ["x", "b"]}},
        		"projection": ["*"]
			},
			"list_list_val_ab":{
				"selection": {"condition": {"feature": "list.a", "operation": "any_value_in_feature", "value": ["a", "b"]}},
        		"projection": ["*"]
			},
			"list_list_ref":{
				"selection": {"condition": {"feature": "list.a", "operation": "any_value_in_feature", "ref": "ref1"}},
        		"projection": ["*"]
			},
			"eq_ref":{
				"selection": {"condition": {"feature": "selectable", "operation": "==", "ref": "ref1"}},
        		"projection": ["*"]
			},
			"uneq_ref":{
				"selection": {"condition": {"feature": "selectable", "operation": "!=", "ref": "ref1"}},
        		"projection": ["*"]
			},
			"user_ref":{
				"selection": {"condition": {"feature": "user", "operation": "==", "ref": "jwt.user"}},
        		"projection": ["*"]
			},
			"or":{
				"selection": {
					"or": [
            			{"condition": {"feature": "selectable", "operation": "==", "ref": "ref1"}},
            			{"condition": {"feature": "selectable", "operation": "==", "ref": "ref2"}}
          			]
				},
        		"projection": ["*"]
			},
			"and":{
				"selection": {
					"and": [
            			{"condition": {"feature": "selectable", "operation": "==", "ref": "ref1"}},
            			{"condition": {"feature": "number", "operation": "==", "ref": "ref2"}}
          			]
				},
        		"projection": ["*"]
			},
			"and_not_existing":{
				"selection": {
					"and": [
            			{"condition": {"feature": "selectable", "operation": "==", "ref": "ref1"}},
            			{"condition": {"feature": "not_existing", "operation": "==", "ref": "ref2"}}
          			]
				},
        		"projection": ["*"]
			}
		}
  	}`

	queryConfig := QueriesConfig{}
	err = json.Unmarshal([]byte(queriesStr), &queryConfig)
	if err != nil {
		t.Fatal()
	}

	Config.Queries = queryConfig

	_, _, err = selectionTestResult("not_existing_endpoint", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	if err == nil {
		t.Fatal("expected error not thrown")
	}

	all, hits, err := selectionTestResult("as_is", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("eq_val_a", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("eq_val_b", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("uneq_val_b", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("uneq_val_a", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("list_val_ab", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("list_val_xb", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("list_list_val_ab", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("list_list_val_xb", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("list_list_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{"ref1": {"x"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("list_list_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{"ref1": {"a,x,y"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("eq_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{"ref1": {"a"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("eq_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{"ref1": {"b"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("uneq_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{"ref1": {"b"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("uneq_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{"ref1": {"a"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("user_ref", jwt_http_router.Jwt{UserId: "userid"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("user_ref", jwt_http_router.Jwt{UserId: "x"}, url.Values{})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("or", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("or", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}, "ref2": {"x"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("or", jwt_http_router.Jwt{}, url.Values{"ref1": {"y"}, "ref2": {"x"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("or", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("and", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}, "ref2": {"42"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("and", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}, "ref2": {"13"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("and", jwt_http_router.Jwt{}, url.Values{"ref1": {"y"}, "ref2": {"x"}})
	if err == nil {
		t.Fatal("expected error missing")
	}

	all, hits, err = selectionTestResult("and", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

	all, hits, err = selectionTestResult("and_not_existing", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{element})

	all, hits, err = selectionTestResult("and_not_existing", jwt_http_router.Jwt{}, url.Values{"ref1": {"a"}, "ref2": {"13"}})
	testSelectionCheck(t, all, hits, err, []string{element}, []string{})

}

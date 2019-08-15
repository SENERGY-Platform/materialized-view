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
	"log"
	"runtime/debug"

	"testing"

	"encoding/json"

	"reflect"

	"time"

	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
)

func initElasticWhereTest() (elasticClient *elastic.Client, purge func(), err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	config := ConfigStruct{}
	err = json.Unmarshal([]byte(testDefaultConfig), &config)
	if err != nil {
		log.Fatalf("Could not unmarshal config: %s", err)
	}
	Config = &config

	elasticCloser, _, elasticIp, err := Elasticsearch(pool)
	Config.ElasticUrl = "http://" + elasticIp + ":9200"
	elasticClient = createClient()
	client = elasticClient
	purge = func() {
		if conn != nil {
			conn.Close()
		}
		elasticCloser()
	}
	if err != nil {
		purge()
		debug.PrintStack()
		log.Fatal(err)
	}
	return
}

func stringToWhereResult(str string, features map[string]interface{}) (all []string, hits []string, err error) {
	where := WhereConditions{}
	err = json.Unmarshal([]byte(str), &where)
	if err != nil {
		return all, hits, err
	}
	filter, err := where.ToElasticFilter(features)
	if err != nil {
		return all, hits, err
	}
	query := elastic.NewBoolQuery().Filter(filter...)
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

func TestWhere(t *testing.T) {
	client, purge, err := initElasticWhereTest()
	defer purge()
	if err != nil {
		t.Fatal()
	}
	element := `{
		"searchable": "foo",
		"selectable": "a",
		"number":     42,
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

	// == and != on string

	features := map[string]interface{}{"a": "a", "b": 42, "c": "b", "d": 13, "e": map[string]interface{}{"a": "a"}, "f": []string{"a", "x"}, "g": []string{"y", "x"}}
	all, result, err := stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "event_feature":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "event_feature":"b"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "event_feature":"c"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "event_feature":"b"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "event_feature":"c"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "event_feature":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	//not existing fields

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "event_feature":"x"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "nope", "operation":"==", "event_feature":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "event_feature":"x"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	// == and != on number

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "event_feature":"b"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, _, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "event_feature":"c"}]`, features)
	if err == nil {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "event_feature":"d"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, _, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "event_feature":"c"}]`, features)
	if err == nil {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "event_feature":"d"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "event_feature":"b"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// == and != with value number

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "value":42}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, _, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "value":"not_a_number"}]`, features)
	if err == nil {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "value":13}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, _, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "value":"not_a_number"}]`, features)
	if err == nil {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "value":13}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "value":42}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// == and != with value string

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "value":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "value":13}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "value":"nope"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "value":13}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "value":"nope"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "value":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// == and != with value null

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"==", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"==", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "not_existing", "operation":"==", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"!=", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "number", "operation":"!=", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "not_existing", "operation":"!=", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// unknown operation

	all, _, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"something_unknown", "value":null}]`, features)
	if err == nil {
		t.Fatal()
	}

	// on list element

	all, result, err = stringToWhereResult(`[{"target_feature": "list.a", "operation":"==", "value":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "list.a", "operation":"!=", "value":"b"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "list.b", "operation":"==", "event_feature":"c"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "list.a", "operation":"!=", "value":"a"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "list.a", "operation":"==", "value":"b"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "list.b", "operation":"!=", "event_feature":"c"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// any_target_in_event

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_event", "event_feature":"f"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_event", "event_feature":"g"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_event", "event_feature":"not_existent"}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// any_target_in_value

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_value", "value":["a", "x"]}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_value", "value":["y", "x"]}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_value", "value":[]}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_value", "value":null}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "selectable", "operation":"any_target_in_value", "value":[]}]`, features)
	if err != nil {
		t.Fatal()
	}
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 0 {
		t.Fatal()
	}

	// child object

	all, result, err = stringToWhereResult(`[{"target_feature": "list", "operation":"==", "event_feature":"e"}]`, features)
	if err == nil {
		t.Fatal("shuld throw err")
	}

	all, result, err = stringToWhereResult(`[{"target_feature": "list.a", "operation":"==", "event_feature":"e.a"}]`, features)
	if len(all) != 1 || !reflect.DeepEqual(all[0], element) {
		t.Fatal()
	}
	if len(result) != 1 || !reflect.DeepEqual(result[0], element) {
		t.Fatal()
	}

}

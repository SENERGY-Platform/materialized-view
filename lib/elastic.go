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
	"runtime/debug"

	"sync"

	"net/http"
	"syscall"
	"time"

	"github.com/olivere/elastic"
)

var client *elastic.Client
var once sync.Once

func GetClient() *elastic.Client {
	once.Do(func() {
		client = createClient()
	})
	return client
}

func createClient() (result *elastic.Client) {
	ctx := context.Background()
	result, err := elastic.NewClient(elastic.SetURL(Config.ElasticUrl), elastic.SetRetrier(newRetrier()))
	if err != nil {
		debug.PrintStack()
		log.Fatal(err, Config.ElasticUrl)
	}
	for kind := range Config.ElasticMapping {
		err = createIndex(kind, result, ctx)
		if err != nil {
			debug.PrintStack()
			log.Fatal(err)
		}
	}
	return
}

func createIndex(kind string, client *elastic.Client, ctx context.Context) (err error) {
	exists, err := client.IndexExists(kind).Do(ctx)
	if err != nil {
		return err
	}
	mapping, err := createMapping(kind)
	if err != nil {
		return err
	}
	mappingJson, _ := json.Marshal(mapping)
	log.Println("expected index setting ", kind, string(mappingJson))
	if !exists {
		createIndex, err := client.CreateIndex(kind + "_v1").BodyJson(mapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("index not acknowledged")
		}
		_, err = client.Alias().Add(kind+"_v1", kind).Do(ctx)
	}
	return
}

type MyRetrier struct {
	backoff elastic.Backoff
}

func newRetrier() *MyRetrier {
	return &MyRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

func (r *MyRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	// Fail hard on a specific error
	if err == syscall.ECONNREFUSED {
		return 0, false, errors.New("Elasticsearch or network down")
	}

	// Stop after n retries
	if int64(retry) >= Config.ElasticRetry {
		return 0, false, nil
	}

	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}

const ElasticResourceType = "resource"

func createMapping(kind string) (result map[string]map[string]map[string]map[string]interface{}, err error) {
	mapping, ok := Config.ElasticMapping[kind]
	if !ok {
		mapping = map[string]interface{}{}
	}
	mapping["feature_search"] = map[string]string{"type": "text", "analyzer": "autocomplete", "search_analyzer": "standard"}
	result = map[string]map[string]map[string]map[string]interface{}{
		"mappings": {
			ElasticResourceType: {
				"properties": mapping,
			},
		},
		"settings": {
			"analysis": {
				"filter": {
					"autocomplete_filter": map[string]interface{}{
						"type":     "edge_ngram",
						"min_gram": 1,
						"max_gram": 20,
					},
				},
				"analyzer": {
					"autocomplete": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"autocomplete_filter",
						},
					},
					"case_insensitive_sort": map[string]interface{}{
						"tokenizer": "keyword",
						"filter":    []string{"lowercase"},
					},
				},
			},
		},
	}
	foo, err := json.Marshal(result)
	log.Println("DEBUG:", string(foo))
	return result, nil
}

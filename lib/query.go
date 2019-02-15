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
	"net/url"

	"errors"

	"log"

	"strconv"

	"github.com/olivere/elastic"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func SearchSorted(target string, searchtext string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, limit string, offset string, orderBy string, asc bool) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("feature_search", searchtext))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Sort(orderBy, asc).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func SearchLimit(target string, searchtext string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, limit string, offset string) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("feature_search", searchtext))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func Search(target string, searchtext string, endpoint string, query url.Values, jwt jwt_http_router.Jwt) (result []interface{}, err error) {
	result = []interface{}{}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("feature_search", searchtext))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}

func Get(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt) (result []interface{}, err error) {
	result = []interface{}{}
	ctx := context.Background()
	q := elastic.NewBoolQuery()
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func GetLimit(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, limit string, offset string) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery()
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func GetSorted(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, limit string, offset string, orderBy string, asc bool) (result []interface{}, total int64, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, 0, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, 0, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery()
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, 0, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Sort(orderBy, asc).Do(ctx)
	if err != nil {
		return result, 0, err
	}
	total = resp.TotalHits()
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, 0, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, 0, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, 0, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}

func SelectField(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, field string, value string) (result []interface{}, err error) {
	result = []interface{}{}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewTermQuery(field, value))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func SelectFieldLimit(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, field string, value string, limit string, offset string) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewTermQuery(field, value))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func SelectFieldSorted(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, field string, value string, limit string, offset string, orderBy string, asc bool) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewTermQuery(field, value))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Sort(orderBy, asc).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}

func SelectFieldValues(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, field string, value []interface{}) (result []interface{}, err error) {
	result = []interface{}{}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewTermsQuery(field, value...))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func SelectFieldValuesLimit(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, field string, value []interface{}, limit string, offset string) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewTermsQuery(field, value...))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}
func SelectFieldValuesSorted(target string, endpoint string, query url.Values, jwt jwt_http_router.Jwt, field string, value []interface{}, limit string, offset string, orderBy string, asc bool) (result []interface{}, err error) {
	result = []interface{}{}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return result, err
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	q := elastic.NewBoolQuery().Must(elastic.NewTermsQuery(field, value...))
	err = UseSelection(q, target, endpoint, jwt, query)
	if err != nil {
		return result, err
	}
	resp, err := GetClient().Search().Index(target).Type(ElasticResourceType).Query(q).Size(l).From(o).Sort(orderBy, asc).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		entry := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return result, err
		}
		targetQuery, ok := Config.Queries[target]
		if !ok {
			log.Println("unknown target: " + target)
			return result, errors.New("unknown target: " + target)
		}
		endpointQuery, ok := targetQuery[endpoint]
		if !ok {
			log.Println("unknown endpoint: " + endpoint)
			return result, errors.New("unknown endpoint: " + endpoint)
		}
		result = append(result, endpointQuery.Projection.Use(entry))
	}
	return
}

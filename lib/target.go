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

	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
)

type Target struct {
	Features map[string]interface{}
	Version  *int64
	Id       string
	Changed  bool
	Removed  bool
	New      bool
	Name     string
}

type Sorting struct {
	Asc   bool   `json:"asc"`
	By    string `json:"by"`
	Limit int    `json:"limit"`
}

func GetTargetsWhere(targetName string, where WhereConditions, features map[string]interface{}) (result []Target, err error) {
	ctx := context.Background()
	filter, err := where.ToElasticFilter(features)
	if err != nil {
		return result, err
	}
	query := elastic.NewBoolQuery().Filter(filter...)
	resp, err := GetClient().Search().Index(targetName).Type(ElasticResourceType).Version(true).Query(query).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		fields := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &fields)
		if err != nil {
			return result, err
		}
		target := Target{Id: hit.Id, Version: hit.Version, Features: fields, Name: targetName}
		result = append(result, target)
	}
	return result, err
}

func GetTargetsWhereSorted(targetName string, where WhereConditions, features map[string]interface{}, sorting Sorting) (result []Target, err error) {
	if sorting.By == "" {
		return GetTargetsWhere(targetName, where, features)
	}
	if sorting.Limit == 0 {
		sorting.Limit = 1000
	}
	ctx := context.Background()
	filter, err := where.ToElasticFilter(features)
	if err != nil {
		return result, err
	}
	query := elastic.NewBoolQuery().Filter(filter...)
	resp, err := GetClient().Search().Index(targetName).Type(ElasticResourceType).Version(true).Query(query).Size(sorting.Limit).Sort(sorting.By, sorting.Asc).Do(ctx)
	if err != nil {
		return result, err
	}
	for _, hit := range resp.Hits.Hits {
		fields := map[string]interface{}{}
		err = json.Unmarshal(*hit.Source, &fields)
		if err != nil {
			return result, err
		}
		target := Target{Id: hit.Id, Version: hit.Version, Features: fields, Name: targetName}
		result = append(result, target)
	}
	return result, err
}

func GetTargetById(targetName string, idFeature string, features map[string]interface{}) (result Target, validRequest bool, err error) {
	ctx := context.Background()
	var idStr string
	if idFeature != "" {
		id, ok := features[idFeature]
		if !ok {
			log.Println("ERROR: event has no id feature", idFeature, features)
			err = errors.New("no valid id feature in event")
			validRequest = false
			return
		}
		idStr, ok = id.(string)
		if !ok {
			log.Println("ERROR: event has no id feature", idFeature, features)
			err = errors.New("no valid id feature in event")
			validRequest = false
			return
		}
	} else {
		idStr = uuid.NewV4().String()
	}
	exists, err := GetClient().Exists().Index(targetName).Type(ElasticResourceType).Id(idStr).Do(ctx)
	if err != nil {
		log.Println("ERROR: exists: ", err)
		return result, true, err
	}
	if exists {
		//resp, err := elastic.NewGetService(GetClient()).Index(targetName).Type(ElasticResourceType).Id(idStr).Do(ctx)
		resp, err := GetClient().Get().Index(targetName).Type(ElasticResourceType).Id(idStr).Do(ctx)
		if err != nil {
			log.Println("ERROR: save: ", err)
			return result, true, err
		}
		fields := map[string]interface{}{}
		err = json.Unmarshal(*resp.Source, &fields)
		if err != nil {
			return result, true, err
		}
		result = Target{Id: resp.Id, Version: resp.Version, Features: fields, Name: targetName}
	} else {
		result = Target{Id: idStr, Features: map[string]interface{}{}, Name: targetName, New: true, Changed: true}
	}
	return result, true, err
}

func (target Target) Save() (err error) {
	ctx := context.Background()
	if target.New {
		_, err = GetClient().Index().Index(target.Name).Type(ElasticResourceType).Id(target.Id).OpType("create").BodyJson(target.Features).Do(ctx)
	} else {
		_, err = GetClient().Index().Index(target.Name).Type(ElasticResourceType).Id(target.Id).Version(*target.Version).BodyJson(target.Features).Do(ctx)
	}
	if err != nil {
		log.Println("ERROR: save: ", err, target)
	}
	return
}

func (target Target) Remove() (err error) {
	if !target.New {
		_, err = elastic.NewDeleteService(GetClient()).Index(target.Name).Type(ElasticResourceType).Id(target.Id).Version(*target.Version).Do(context.Background())
	}
	return
}

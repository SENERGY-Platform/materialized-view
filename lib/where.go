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
	"errors"
	"log"

	"github.com/olivere/elastic"
)

type WhereOperationType string

const (
	WhereEqualOperation   WhereOperationType = "=="
	WhereUnequalOperation WhereOperationType = "!="
	WhereAnyTargetInEvent WhereOperationType = "any_target_in_event"
	WhereAnyTargetInValue WhereOperationType = "any_target_in_value"
)

type WhereCondition struct {
	TargetFeature string             `json:"target_feature"`
	Operation     WhereOperationType `json:"operation"`
	EventFeature  string             `json:"event_feature"`
	Value         interface{}        `json:"value"`
}

type WhereConditions []WhereCondition

func (this WhereConditions) ToElasticFilter(features map[string]interface{}) (result []elastic.Query, err error) {
	if len(this) == 0 {
		return []elastic.Query{elastic.NewMatchAllQuery()}, nil
	}
	for _, condition := range this {
		filter, err := condition.ToElasticFilter(features)
		if err != nil {
			return result, err
		}
		result = append(result, filter)
	}
	return
}

func (this WhereCondition) ToElasticFilter(features Features) (result elastic.Query, err error) {
	switch this.Operation {
	case WhereEqualOperation:
		if this.EventFeature == "" {
			if this.Value == nil {
				return elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery(this.TargetFeature)), err
			} else {
				return elastic.NewTermQuery(this.TargetFeature, this.Value), err
			}
		} else {
			if val, ok := features.Get(this.EventFeature); ok {
				return elastic.NewTermQuery(this.TargetFeature, val), err
			} else {
				return elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery(this.TargetFeature)), err
			}
		}
	case WhereUnequalOperation:
		if this.EventFeature == "" {
			if this.Value == nil {
				return elastic.NewExistsQuery(this.TargetFeature), err
			} else {
				return elastic.NewBoolQuery().MustNot(elastic.NewTermQuery(this.TargetFeature, this.Value)), err
			}
		} else {
			if val, ok := features.Get(this.EventFeature); ok {
				return elastic.NewBoolQuery().MustNot(elastic.NewTermQuery(this.TargetFeature, val)), err
			} else {
				return elastic.NewExistsQuery(this.TargetFeature), err
			}
		}
	case WhereAnyTargetInEvent:
		or := []elastic.Query{}
		list, ok := features.Get(this.EventFeature)
		if !ok {
			return elastic.NewMatchNoneQuery(), err
		}
		valueList, err := InterfaceSlice(list)
		if err != nil {
			log.Println("ERROR: ", err, this, features)
			return result, err
		}
		if len(valueList) == 0 {
			return elastic.NewMatchNoneQuery(), err
		}
		for _, value := range valueList {
			or = append(or, elastic.NewTermQuery(this.TargetFeature, value))
		}
		return elastic.NewBoolQuery().Filter(elastic.NewBoolQuery().Should(or...)), err
	case WhereAnyTargetInValue:
		or := []elastic.Query{}
		if this.Value == nil {
			return elastic.NewMatchNoneQuery(), err
		}
		valueList, err := InterfaceSlice(this.Value)
		if err != nil {
			log.Println("ERROR: ", err, this, features)
			return result, err
		}
		if len(valueList) == 0 {
			return elastic.NewMatchNoneQuery(), err
		}
		for _, value := range valueList {
			or = append(or, elastic.NewTermQuery(this.TargetFeature, value))
		}
		return elastic.NewBoolQuery().Filter(elastic.NewBoolQuery().Should(or...)), err
	}
	err = errors.New("unknown WhereOperator " + string(this.Operation))
	return
}

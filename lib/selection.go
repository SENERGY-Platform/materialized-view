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
	"net/url"

	"errors"

	"github.com/olivere/elastic"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"reflect"
	"strings"
)

type QueryOperationType string

const (
	QueryEqualOperation             QueryOperationType = "=="
	QueryUnequalOperation           QueryOperationType = "!="
	QueryAnyValueInFeatureOperation QueryOperationType = "any_value_in_feature"
)

type ConditionConfig struct {
	Feature   string             `json:"feature"`
	Operation QueryOperationType `json:"operation"`
	Value     interface{}        `json:"value"`
	Ref       string             `json:"ref"`
}

type SelectionConfig struct {
	And       []SelectionConfig `json:"and"`
	Or        []SelectionConfig `json:"or"`
	Condition ConditionConfig   `json:"condition"`
	All       bool              `json:"all"`
}

func UseSelection(query *elastic.BoolQuery, target string, endpoint string, jwt jwt_http_router.Jwt, params url.Values) (err error) {
	endpoints, ok := Config.Queries[target]
	if !ok {
		err = errors.New("unknown target: " + target)
		return
	}
	config, ok := endpoints[endpoint]
	if !ok {
		err = errors.New("unknown target endpoint: " + target + " " + endpoint)
		return
	}
	if !config.Selection.All {
		filter, err := config.Selection.GetFilter(jwt, params)
		if err != nil {
			return err
		}
		query.Filter(filter)
	}
	return err
}

func (this SelectionConfig) GetFilter(jwt jwt_http_router.Jwt, values url.Values) (result elastic.Query, err error) {
	if len(this.And) > 0 {
		and := []elastic.Query{}
		for _, sub := range this.And {
			andElement, err := sub.GetFilter(jwt, values)
			if err != nil {
				return result, err
			}
			and = append(and, andElement)
		}
		result = elastic.NewBoolQuery().Filter(and...)
		return
	}
	if len(this.Or) > 0 {
		or := []elastic.Query{}
		for _, sub := range this.Or {
			orElement, err := sub.GetFilter(jwt, values)
			if err != nil {
				return result, err
			}
			or = append(or, orElement)
		}
		result = elastic.NewBoolQuery().Should(or...)
		return
	}
	return this.Condition.GetFilter(jwt, values)
}

func (this ConditionConfig) GetFilter(jwt jwt_http_router.Jwt, values url.Values) (elastic.Query, error) {
	val := this.Value
	if val == nil || val == "" {
		switch this.Ref {
		case "jwt.user":
			val = jwt.UserId
		case "jwt.groups":
			val = jwt.RealmAccess.Roles
		default:
			val = values.Get(this.Ref)
		}
	}
	switch this.Operation {
	case QueryEqualOperation:
		if val == nil || val == "" {
			return elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery(this.Feature)), nil
		} else {
			return elastic.NewTermQuery(this.Feature, val), nil
		}
	case QueryUnequalOperation:
		if val == nil || val == "" {
			return elastic.NewExistsQuery(this.Feature), nil
		} else {
			return elastic.NewBoolQuery().MustNot(elastic.NewTermQuery(this.Feature, val)), nil
		}
	case QueryAnyValueInFeatureOperation:
		if reflect.TypeOf(val).Kind() == reflect.String {
			val = strings.Split(val.(string), ",")
		}
		arr, err := InterfaceSlice(val)
		if err != nil {
			return nil, err
		}
		return elastic.NewTermsQuery(this.Feature, arr...), nil
	}
	return nil, errors.New("unknown query opperation type " + string(this.Operation))
}

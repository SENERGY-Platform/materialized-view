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
	"log"
	"reflect"
	"strings"
)

type IfOperationType string

const (
	IfEqualOperation                      IfOperationType = "=="
	IfUnequalOperation                    IfOperationType = "!="
	IfFeatureStringContainsValueOperation IfOperationType = "feature_str_contains_value"
)

type IfCondition struct {
	Feature   string          `json:"feature"`
	Operation IfOperationType `json:"operation"`
	Value     interface{}     `json:"value"`
}

type IfConditions []IfCondition

func (this IfConditions) CheckFeatures(features Features) bool {
	for _, condition := range this {
		if !condition.CheckFeatures(features) {
			return false
		}
	}
	return true
}

func (this IfCondition) CheckFeatures(features Features) bool {
	val, ok := features.Get(this.Feature)
	if !ok {
		return this.Operation.Do(nil, this.Value)
	}
	return this.Operation.Do(val, this.Value)
}

func (this IfOperationType) Do(feature interface{}, value interface{}) bool {
	switch this {
	case IfEqualOperation:
		return reflect.DeepEqual(feature, value)
	case IfUnequalOperation:
		return !reflect.DeepEqual(feature, value)
	case IfFeatureStringContainsValueOperation:
		fstr, ok := feature.(string)
		if !ok {
			return false
		}
		vstr, ok := value.(string)
		if !ok {
			return false
		}
		return strings.Contains(fstr, vstr)
	}
	log.Println("WARNING: unknown if operation type", this)
	return false
}

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

type InitDefault []InitDefaultElement

type InitDefaultElement map[string]InitDefaultValue

type InitDefaultValue struct {
	Feature string      `json:"feature"`
	Value   interface{} `json:"value"`
}

func (this InitDefault) Use(features map[string]interface{}) (result []map[string]interface{}) {
	for _, value := range this {
		result = append(result, value.Use(features))
	}
	return
}

func (this InitDefaultElement) Use(features map[string]interface{}) (result map[string]interface{}) {
	result = map[string]interface{}{}
	for key, value := range this {
		result[key] = value.Use(features)
	}
	return
}

func (this InitDefaultValue) Use(features Features) interface{} {
	if this.Feature == "" {
		return this.Value
	} else {
		val, _ := features.Get(this.Feature)
		return val
	}
}

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
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

type Features map[string]interface{}

type Feature struct {
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Temp       bool        `json:"temp"`
	Omitempty  bool        `json:"omitempty"`
	Default    interface{} `json:"default"`
	DefaultRef string      `json:"default_ref"`
}

var featureRef = map[string]func() interface{}{
	"time.epoch_millis": func() interface{} {
		return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	},
	"time.epoch_second": func() interface{} {
		return time.Now().Unix()
	},
}

func MsgToFeatures(features []Feature, msg []byte) (temp Features, permanent Features, err error) {
	temp = map[string]interface{}{}
	permanent = map[string]interface{}{}
	for _, feature := range features {
		var pathResult interface{}
		if feature.Path == "" {
			pathResult = UseDefault(feature)
		} else {
			pathResult, err = UseJsonPath(msg, feature.Path)
			if err != nil {
				return temp, permanent, err
			}
			if pathResult == nil {
				pathResult = UseDefault(feature)
			}
		}
		if !(feature.Omitempty && isEmpty(pathResult)) {
			temp[feature.Name] = pathResult
			if !feature.Temp {
				permanent[feature.Name] = pathResult
			}
		}
	}
	return
}

func MsgStructToFeatures(features []Feature, msgStruct map[string]interface{}) (temp Features, permanent Features, err error) {
	msg, err := json.Marshal(msgStruct)
	if err != nil {
		return temp, permanent, err
	}
	return MsgToFeatures(features, msg)
}

func UseDefault(feature Feature) interface{} {
	if feature.Default != nil {
		return feature.Default
	}
	if ref, ok := featureRef[feature.DefaultRef]; ok {
		return ref()
	}
	return nil
}

func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	kind := reflect.TypeOf(value).Kind()
	switch kind {
	case reflect.Slice:
		return reflect.ValueOf(value).Len() == 0
	case reflect.String:
		v, ok := value.(string)
		return !ok || v == ""
	}
	return false
}

func (this Features) Get(location string) (result interface{}, ok bool) {
	path := strings.Split(location, ".")
	result = this
	for _, key := range path {
		ok = true
		if reflect.TypeOf(result).Kind() == reflect.Map && reflect.TypeOf(result).Key().Kind() == reflect.String {
			sub := reflect.ValueOf(result).MapIndex(reflect.ValueOf(key))
			if !sub.IsValid() {
				return nil, false
			}
			result = sub.Interface()
		} else {
			return nil, false
		}
	}
	return
}

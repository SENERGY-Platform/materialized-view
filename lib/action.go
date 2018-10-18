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
	"reflect"
)

type ActionType string

const (
	RemoveTargetAction ActionType = "remove_target"
	RemoveAction       ActionType = "remove"
	InsertAction       ActionType = "insert"
)

type ScaleType string

const (
	ScaleOne  ScaleType = "one"
	ScaleMany ScaleType = "many"
)

type Actions []Action

type Action struct {
	Type   ActionType   `json:"type"`
	If     IfConditions `json:"if"`
	Fields []string     `json:"fields"`
	Scale  ScaleType    `json:"scale"`
}

func (this Actions) Do(target Target, temp map[string]interface{}, perm map[string]interface{}) (result Target, err error) {
	result = target
	for _, action := range this {
		if action.If.CheckFeatures(temp) {
			result.Changed = true
			result, err = action.Do(result, temp, perm)
			if err != nil {
				log.Println("ERROR: while doing action", err)
				return
			}
		}
	}
	return
}

func (this Action) Do(target Target, temp map[string]interface{}, perm map[string]interface{}) (Target, error) {
	var err error
	switch this.Type {
	case InsertAction:
		for _, field := range this.Fields {
			switch this.Scale {
			case ScaleMany:
				err = target.addToList(field, perm)
			case ScaleOne:
				err = target.setField(field, perm)
			default:
				err = errors.New("unknown scale type " + string(this.Scale))
			}
			if err != nil {
				return target, err
			}
		}
	case RemoveAction:
		for _, field := range this.Fields {
			switch this.Scale {
			case ScaleMany:
				err = target.removeFromList(field, perm)
			case ScaleOne:
				err = target.removeField(field, perm)
			default:
				err = errors.New("unknown scale type " + string(this.Scale))
			}
			if err != nil {
				return target, err
			}
		}
	case RemoveTargetAction:
		target.Removed = true
		target.Changed = true
	default:
		err = errors.New("unknown action type " + string(this.Type))
	}
	return target, err
}

func (this *Target) addToList(field string, value interface{}) (err error) {
	list := this.Features[field]
	if reflect.TypeOf(list).Kind() == reflect.Slice {
		this.Features[field], err = InterfaceSliceAppend(list, value)
	} else {
		this.Features[field] = []interface{}{value}
	}
	return
}

func (this *Target) setField(field string, value map[string]interface{}) (err error) {
	if field == "" {
		for key, elem := range value {
			this.Features[key] = elem
		}
	} else {
		this.Features[field] = value
	}
	return
}

func (this *Target) removeFromList(field string, value map[string]interface{}) (err error) {
	list := this.Features[field]
	if reflect.TypeOf(list) != nil && reflect.TypeOf(list).Kind() == reflect.Slice {
		this.Features[field], err = InterfaceSliceRemove(list, value)
	} else {
		this.Features[field] = []interface{}{}
	}
	return
}

func (this *Target) removeField(field string, value map[string]interface{}) (err error) {
	if field != "" {
		delete(this.Features, field)
	} else {
		for key := range value {
			delete(this.Features, key)
		}
	}
	return
}

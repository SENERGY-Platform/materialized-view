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
)

type GroupType string

const (
	RootGroupType  GroupType = "root"
	ChildGroupType GroupType = "child"
)

type EventActionGroup struct {
	Type      GroupType         `json:"type"`
	Target    string            `json:"target"`
	IdFeature string            `json:"id_feature"`
	Where     WhereConditions   `json:"where"`
	If        IfConditions      `json:"if"`
	Features  []Feature         `json:"features"`
	Actions   Actions           `json:"actions"`
	Init      []InitActionGroup `json:"init"`
}

type InitActionGroup struct {
	Target    string          `json:"target"`
	Where     WhereConditions `json:"where"`
	Sorting   Sorting         `json:"sorting"`
	Default   InitDefault     `json:"default"`
	Transform []Feature       `json:"transform"`
	Actions   Actions         `json:"actions"`
}

func CreateGroupHandler(group EventActionGroup) (handler func(delivery []byte) error, err error) {
	return func(delivery []byte) error {
		temp, perma, err := MsgToFeatures(group.Features, delivery)
		if err != nil {
			return err
		}
		if !group.If.CheckFeatures(temp) {
			return nil
		}
		switch group.Type {
		case RootGroupType:
			return handleRoot(group, temp, perma)
		case ChildGroupType:
			return handleChild(group, temp, perma)
		default:
			log.Println("WARNING: unknown group type, will not be processed", group.Type)
		}
		return nil
	}, nil
}

func handleChild(group EventActionGroup, temp map[string]interface{}, perm map[string]interface{}) error {
	targets, err := GetTargetsWhere(group.Target, group.Where, temp)
	if err != nil {
		return err
	}
	for _, target := range targets {
		result, err := group.Actions.Do(target, temp, perm)
		if err != nil {
			return err
		}
		if result.Removed {
			err = result.Remove()
		} else if result.Changed {
			err = result.Save()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func handleRoot(group EventActionGroup, temp map[string]interface{}, perm map[string]interface{}) error {
	target, validRequest, err := GetTargetById(group.Target, group.IdFeature, temp)
	if err != nil {
		if validRequest {
			return err
		} else {
			log.Println("INVALID REQUEST: ", err)
			return nil
		}
	}
	result, err := group.Actions.Do(target, temp, perm)
	if err != nil {
		return err
	}
	if result.New {
		result, err = handleInit(result, group.Init, temp)
		if err != nil {
			return err
		}
	}
	if result.Removed {
		err = result.Remove()
	} else if result.Changed {
		err = result.Save()
	}
	return err
}

func handleInit(target Target, groups []InitActionGroup, temp map[string]interface{}) (Target, error) {
	for _, group := range groups {
		children, err := GetTargetsWhereSorted(group.Target, group.Where, temp, group.Sorting)
		if err != nil {
			return target, err
		}
		if len(children) == 0 {
			for _, defaultVal := range group.Default.Use(temp) {
				children = append(children, Target{Features: defaultVal})
			}
		}
		for _, child := range children {
			temp, perm, err := MsgStructToFeatures(group.Transform, child.Features)
			if err != nil {
				return target, err
			}
			target, err = group.Actions.Do(target, temp, perm)
			if err != nil {
				return target, err
			}
		}
	}
	return target, nil
}

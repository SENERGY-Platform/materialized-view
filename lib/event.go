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
	"github.com/SENERGY-Platform/materialized-view/lib/kafka"
	"log"
)

var conn kafka.Interface

func InitEventHandling() (err error) {
	conn, err = kafka.Init(Config.ZookeeperUrl, Config.ConsumerGroup, Config.Debug)
	if err != nil {
		log.Fatal("ERROR: while initializing amqp connection ", err)
		return
	}
	for topic, groupes := range Config.Events {
		handler, err := createHandler(groupes)
		if err != nil {
			log.Fatal("ERROR: while creating topic handler", topic, err)
			return err
		}
		err = conn.Consume(topic, handler)
		if err != nil {
			log.Fatal("ERROR: while initializing consumer", topic, err)
			return err
		}
	}
	return
}

func createHandler(groupes []EventActionGroup) (handler func(delivery []byte) error, err error) {
	groupHandlers := []func(delivery []byte) error{}
	for _, group := range groupes {
		groupHandler, err := CreateGroupHandler(group)
		if err != nil {
			return handler, err
		}
		groupHandlers = append(groupHandlers, groupHandler)
	}
	return func(delivery []byte) error {
		log.Println("DEBUG: consume", string(delivery))
		for _, handler := range groupHandlers {
			if err := handler(delivery); err != nil {
				return err
			}
		}
		return nil
	}, err
}

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
	"log"

	"github.com/SmartEnergyPlatform/amqp-wrapper-lib"
)

var conn *amqp_wrapper_lib.Connection

func InitEventHandling() (err error) {
	conn, err = amqp_wrapper_lib.Init(Config.AmqpUrl, Config.Events.GetTopicList(), Config.AmqpReconnectTimeout)
	if err != nil {
		log.Fatal("ERROR: while initializing amqp connection ", err, "CONFIG: ", Config.AmqpUrl, Config.Events.GetTopicList(), Config.AmqpReconnectTimeout)
		return
	}
	conn.SetMessageLogging(Config.AmqpLogging == "true")
	for topic, groupes := range Config.Events {
		handler, err := createHandler(groupes)
		if err != nil {
			log.Fatal("ERROR: while creating topic handler", topic, err)
			return err
		}
		err = conn.Consume(Config.AmqpConsumerName+"_"+topic, topic, handler)
		if err != nil {
			log.Fatal("ERROR: while initializing consumer", topic, err)
			return err
		}
	}
	return
}

func createHandler(groupes []EventActionGroup) (handler amqp_wrapper_lib.ConsumerFunc, err error) {
	groupHandlers := []amqp_wrapper_lib.ConsumerFunc{}
	for _, group := range groupes {
		groupHandler, err := CreateGroupHandler(group)
		if err != nil {
			return handler, err
		}
		groupHandlers = append(groupHandlers, groupHandler)
	}
	return func(delivery []byte) error {
		for _, handler := range groupHandlers {
			if err := handler(delivery); err != nil {
				return err
			}
		}
		return nil
	}, err
}

func sendEvent(topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Println("ERROR: event marshaling:", err)
		return err
	}
	log.Println("DEBUG: send amqp event: ", topic, string(payload))
	return conn.Publish(topic, payload)
}

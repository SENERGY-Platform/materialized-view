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
	"net/http"

	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"

	"encoding/json"

	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func StartApi() {
	log.Println("start server on port: ", Config.ServerPort)
	httpHandler := getRoutes()
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, Config.LogLevel)
	log.Println(http.ListenAndServe(":"+Config.ServerPort, logger))
}

func getRoutes() (router *jwt_http_router.Router) {
	router = jwt_http_router.New(jwt_http_router.JwtConfig{
		ForceUser: Config.ForceUser == "true",
		ForceAuth: Config.ForceAuth == "true",
		PubRsa:    Config.JwtPubRsa,
	})

	router.GET("/search/:target/:searchtext/:endpoint", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		searchtext := ps.ByName("searchtext")
		endpoint := ps.ByName("endpoint")
		result, err := Search(target, searchtext, endpoint, r.URL.Query(), jwt)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/:target/:searchtext/:endpoint/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		searchtext := ps.ByName("searchtext")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		result, err := SearchLimit(target, searchtext, endpoint, r.URL.Query(), jwt, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/:target/:searchtext/:endpoint/:limit/:offset/:order_by/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		searchtext := ps.ByName("searchtext")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, _, err := SearchSorted(target, searchtext, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/:target/:searchtext/:endpoint/:limit/:offset/:order_by/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		searchtext := ps.ByName("searchtext")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, _, err := SearchSorted(target, searchtext, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/search/:target/:searchtext/:endpoint/:limit/:offset/:order_by/asc/withtotal", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		searchtext := ps.ByName("searchtext")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, total, err := SearchSorted(target, searchtext, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(map[string]interface{}{"total": total, "result": result})
	})

	router.GET("/search/:target/:searchtext/:endpoint/:limit/:offset/:order_by/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		searchtext := ps.ByName("searchtext")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, total, err := SearchSorted(target, searchtext, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(map[string]interface{}{"total": total, "result": result})
	})

	router.GET("/get/:target/:endpoint", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		result, err := Get(target, endpoint, r.URL.Query(), jwt)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/get/:target/:endpoint/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		result, err := GetLimit(target, endpoint, r.URL.Query(), jwt, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/get/:target/:endpoint/:limit/:offset/:order_by/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, _, err := GetSorted(target, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/get/:target/:endpoint/:limit/:offset/:order_by/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, _, err := GetSorted(target, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/get/:target/:endpoint/:limit/:offset/:order_by/desc/withtotal", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, total, err := GetSorted(target, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(map[string]interface{}{"total": total, "result": result})
	})

	router.GET("/get/:target/:endpoint/:limit/:offset/:order_by/asc/withtotal", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		result, total, err := GetSorted(target, endpoint, r.URL.Query(), jwt, limit, offset, orderBy, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(map[string]interface{}{"total": total, "result": result})
	})

	router.GET("/select/field/:target/:endpoint/:field/:value", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		field := ps.ByName("field")
		value := ps.ByName("value")
		result, err := SelectField(target, endpoint, r.URL.Query(), jwt, field, value)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/field/:target/:endpoint/:field/:value/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		field := ps.ByName("field")
		value := ps.ByName("value")
		result, err := SelectFieldLimit(target, endpoint, r.URL.Query(), jwt, field, value, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/field/:target/:endpoint/:field/:value/:limit/:offset/:order_by/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		field := ps.ByName("field")
		value := ps.ByName("value")
		result, err := SelectFieldSorted(target, endpoint, r.URL.Query(), jwt, field, value, limit, offset, orderBy, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/select/field/:target/:endpoint/:field/:value/:limit/:offset/:order_by/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		field := ps.ByName("field")
		value := ps.ByName("value")
		result, err := SelectFieldSorted(target, endpoint, r.URL.Query(), jwt, field, value, limit, offset, orderBy, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.POST("/select/field/:target/:endpoint/:field", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		field := ps.ByName("field")
		value := []interface{}{}
		err := json.NewDecoder(r.Body).Decode(&value)
		if err != nil {
			log.Println("WARNING: error in user send data", err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := SelectFieldValues(target, endpoint, r.URL.Query(), jwt, field, value)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.POST("/select/field/:target/:endpoint/:field/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		field := ps.ByName("field")
		value := []interface{}{}
		err := json.NewDecoder(r.Body).Decode(&value)
		result, err := SelectFieldValuesLimit(target, endpoint, r.URL.Query(), jwt, field, value, limit, offset)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.POST("/select/field/:target/:endpoint/:field/:limit/:offset/:order_by/asc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		field := ps.ByName("field")
		value := []interface{}{}
		err := json.NewDecoder(r.Body).Decode(&value)
		result, err := SelectFieldValuesSorted(target, endpoint, r.URL.Query(), jwt, field, value, limit, offset, orderBy, true)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.POST("/select/field/:target/:endpoint/:field/:limit/:offset/:order_by/desc", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		target := ps.ByName("target")
		endpoint := ps.ByName("endpoint")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		orderBy := ps.ByName("order_by")
		field := ps.ByName("field")
		value := []interface{}{}
		err := json.NewDecoder(r.Body).Decode(&value)
		result, err := SelectFieldValuesSorted(target, endpoint, r.URL.Query(), jwt, field, value, limit, offset, orderBy, false)
		if err != nil {
			log.Println("ERROR: ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	return
}

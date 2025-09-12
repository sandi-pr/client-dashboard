/*
 * CR3 Client Dashboard
 * Package: gitlab.hexacode.org/cr3/client-dashboard
 * Maintainer: Dean Ramhan <dean@hexacode.org>
 *
 * Copyright (C) 2024 Hexacode Teknologi Indonesia
 * All Rights Reserved
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateCollectionApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	upload_by, ok := me.GetFloat64("id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	author, ok := me.GetString("username")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	name, ok := ctx.Params.GetString("name")
	if !ok {
		return ctx.Error(400, "Invalid parameter: name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		parent_id = 0
	}

	dataSend := U.Dict{
		"name": name,
		"type": "collection",
		"meta": U.Dict{
			"upload_by": upload_by,
			"author":    author,
		},
	}

	if description != "" {
		dataSend["description"] = description
	}

	if parent_id != 0 {
		dataSend["parent_id"] = parent_id
	}

	status, response, err := Strapi.Create("api-library-table", dataSend)
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func UpdateCollectionApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "Invalid parameter: id")
	}

	name, ok := ctx.Params.GetString("name")
	if !ok {
		name = ""
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		parent_id = 0
	}

	dataSend := U.Dict{}

	if name != "" {
		dataSend["name"] = name
	}

	if description != "" {
		dataSend["description"] = description
	}

	if parent_id != 0 {
		dataSend["parent_id"] = parent_id
	}

	status, response, err := Strapi.Update("api-library-table", int64(id), dataSend)
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func DeleteCollectionApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "Invalid parameter: id")
	}

	page := 1
	apiIds := U.List{}

	for {
		status, find, err := Strapi.Find("api-library-table", U.Dict{
			"filters": U.Dict{
				"parent_id": U.Dict{
					"$eq": id,
				},
			},
			"pagination": U.Dict{
				"page": page,
			},
		})
		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := find.GetList("data")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) <= 0 {
			break
		}

		for i := 0; i < len(list); i++ {
			dict, ok := list.GetDict(i)

			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}

			id, ok := dict.GetFloat64("id")

			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}

			apiIds = append(apiIds, id)
		}
	}

	for i := 0; i < len(apiIds); i++ {
		apiId, ok := apiIds.GetFloat64(i)

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		status, delete, err := Strapi.Delete("api-library-table", int64(apiId))

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			fmt.Println(delete)

			if err2, ok := delete.GetDict("error"); ok {
				if message, ok := err2.GetString("message"); ok {
					return ctx.Error(status, string(message))
				}
			}
		}
	}

	status, response, err := Strapi.Delete("api-library-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func GetCollectionApis(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		page_size = 10
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		parent_id = 0
	}

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

	dataSend := U.Dict{
		"filters": U.Dict{
			"type": "collection",
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	}

	if parent_id != 0 {
		dataSend["filters"].(U.Dict)["parent_id"] = U.Dict{
			"$eq": parent_id,
		}
	}

	if search != "" {
		dataSend["filters"].(U.Dict)["name"] = U.Dict{
			"$containsi": search,
		}
	}

	status, response, err := Strapi.Find("api-library-table", dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func GetCollectionApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "Invalid parameter: id")
	}

	status, response, err := Strapi.FindOne("api-library-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func CreateApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	upload_by, ok := me.GetFloat64("id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	author, ok := me.GetString("username")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	name, ok := ctx.Params.GetString("name")
	if !ok {
		return ctx.Error(400, "Invalid parameter: name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: parent_id")
	}

	url, ok := ctx.Params.GetString("url")
	if !ok {
		return ctx.Error(400, "Invalid parameter: url")
	}

	method, ok := ctx.Params.GetString("method")
	if !ok {
		return ctx.Error(400, "Invalid parameter: method")
	}

	auth_type, ok := ctx.Params.GetString("auth_type")
	if !ok {
		return ctx.Error(400, "Invalid parameter: auth_type")
	}

	username, ok := ctx.Params.GetString("username")
	if !ok {
		username = ""
	}

	password, ok := ctx.Params.GetString("password")
	if !ok {
		password = ""
	}

	token, ok := ctx.Params.GetString("token")
	if !ok {
		token = ""
	}

	params, ok := ctx.Params.GetString("params")
	if !ok {
		params = ""
	}

	body_type, ok := ctx.Params.GetString("body_type")
	if !ok {
		body_type = "none"
	}

	body, ok := ctx.Params.GetString("body")
	if !ok {
		body = ""
	}

	dataSend := U.Dict{
		"name":      name,
		"type":      "api",
		"parent_id": parent_id,
		"url":       url,
		"method":    method,
		"auth_type": auth_type,
		"body_type": body_type,
		"meta": U.Dict{
			"upload_by": upload_by,
			"author":    author,
		},
	}

	if description != "" {
		dataSend["description"] = description
	}

	if username != "" {
		dataSend["username"] = username
	}

	if password != "" {
		dataSend["password"] = password
	}

	if token != "" {
		dataSend["token"] = token
	}

	if params != "" {
		var paramsJson interface{}
		err := json.Unmarshal([]byte(params), &paramsJson)
		if err != nil {
			return ctx.Error(400, "Invalid parameter: params")
		}

		dataSend["params"] = paramsJson
	}

	if body != "" {
		var bodyJson interface{}
		err := json.Unmarshal([]byte(body), &bodyJson)
		if err != nil {
			return ctx.Error(400, "Invalid parameter: body")
		}

		dataSend["body"] = bodyJson
	}

	status, response, err := Strapi.Create("api-library-table", dataSend)
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func UpdateApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "Invalid parameter: id")
	}

	name, ok := ctx.Params.GetString("name")
	if !ok {
		name = ""
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	url, ok := ctx.Params.GetString("url")
	if !ok {
		url = ""
	}

	method, ok := ctx.Params.GetString("method")
	if !ok {
		method = ""
	}

	auth_type, ok := ctx.Params.GetString("auth_type")
	if !ok {
		auth_type = ""
	}

	username, ok := ctx.Params.GetString("username")
	if !ok {
		username = ""
	}

	password, ok := ctx.Params.GetString("password")
	if !ok {
		password = ""
	}

	token, ok := ctx.Params.GetString("token")
	if !ok {
		token = ""
	}

	params, ok := ctx.Params.GetString("params")
	if !ok {
		params = ""
	}

	body_type, ok := ctx.Params.GetString("body_type")
	if !ok {
		body_type = "none"
	}

	body, ok := ctx.Params.GetString("body")
	if !ok {
		body = ""
	}

	dataSend := U.Dict{}

	if name != "" {
		dataSend["name"] = name
	}

	if description != "" {
		dataSend["description"] = description
	}

	if url != "" {
		dataSend["url"] = url
	}

	if method != "" {
		dataSend["method"] = method
	}

	if auth_type != "" {
		dataSend["auth_type"] = auth_type
	}

	if username != "" {
		dataSend["username"] = username
	}

	if password != "" {
		dataSend["password"] = password
	}

	if token != "" {
		dataSend["token"] = token
	}

	if params != "" {
		var paramsJson map[string]interface{}
		err := json.Unmarshal([]byte(params), &paramsJson)
		if err != nil {
			return ctx.Error(400, "Invalid parameter: params")
		}
	}

	if body_type != "" {
		dataSend["body_type"] = body_type
	}

	if body != "" {
		var bodyJson map[string]interface{}
		err := json.Unmarshal([]byte(body), &bodyJson)
		if err != nil {
			return ctx.Error(400, "Invalid parameter: body")
		}

		dataSend["body"] = bodyJson
	}

	status, response, err := Strapi.Update("api-library-table", int64(id), dataSend)
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func DeleteApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "Invalid parameter: id")
	}

	status, response, err := Strapi.Delete("api-library-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func GetApis(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		page_size = 10
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: parent_id")
	}

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

	dataSend := U.Dict{
		"filters": U.Dict{
			"type":      U.Dict{"$eq": "api"},
			"parent_id": U.Dict{"$eq": parent_id},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	}

	if search != "" {
		dataSend["filters"].(U.Dict)["name"] = U.Dict{"$containsi": search}
	}

	status, response, err := Strapi.Find("api-library-table", dataSend)
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, response, nil
}

func GetApi(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "Invalid parameter: id")
	}

	status, response, err := Strapi.Find("api-library-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": id,
			},
		},
		"populate": "*",
	})
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := response.GetList("data")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, dict, nil
}

func IsCollectionApiNameAvailable(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	name, ok := ctx.Params.GetString("name")
	if !ok {
		return ctx.Error(400, "Invalid parameter: name")
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		parent_id = 0
	}

	collectionType, ok := ctx.Params.GetString("type")
	if !ok {
		return ctx.Error(400, "Invalid parameter: type")
	}

	dataSend := U.Dict{
		"filters": U.Dict{
			"name": U.Dict{"$eq": name},
			"type": U.Dict{"$eq": collectionType},
		},
	}

	if parent_id != 0 {
		dataSend["filters"].(U.Dict)["parent_id"] = U.Dict{"$eq": parent_id}
	}

	fmt.Println(dataSend)

	status, response, err := Strapi.Find("api-library-table", dataSend)
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(response)

		if err2, ok := response.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := response.GetList("data"); ok && len(list) > 0 {
		return 200, U.Dict{
			"available": false,
		}, nil
	}

	return 200, U.Dict{
		"available": true,
	}, nil
}

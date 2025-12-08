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
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

// {
// 	"name": "agent",
// 	"label": "Agent",
// 	"fields": {
// 		"name": {
// 			"type": "string",
// 			"required": true
// 		},
// 		"description": {
// 			"type": "text"
// 		},
// 		"system_prompt": {
// 			"type": "text"
// 		},
// 		"llm": {
// 			"type": "string",
// 			"required": true
// 		},
// 		"tools": {
// 			"type": "string"
// 		},
// 		"handover": {
// 			"type": "string",
// 			"required": false
// 		},
// 		"is_root_agent": {
// 			"type": "boolean",
// 			"default": false
// 		},
// 		"contract_rule_id": {
// 			"type": "integer"
// 		}
// 	}
// },

func CreateAgent(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "Invalid parameter: description")
	}

	system_prompt, ok := ctx.Params.GetString("system_prompt")
	if !ok {
		system_prompt = ""
	}

	llm, ok := ctx.Params.GetString("llm")
	if !ok {
		return ctx.Error(400, "Invalid parameter: llm")
	}

	tools, ok := ctx.Params.GetString("tools")
	if !ok {
		tools = ""
	}

	handover, ok := ctx.Params.GetString("handover")
	if !ok {
		handover = ""
	}

	is_root_agent, ok := ctx.Params.GetBool("is_root_agent")
	if !ok {
		is_root_agent = false
	}

	contract_rule_id, ok := ctx.Params.GetFloat64("contract_rule_id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: contract_rule_id")
	}

	dataSend := U.Dict{
		"name":             name,
		"description":      description,
		"llm":              llm,
		"is_root_agent":    is_root_agent,
		"contract_rule_id": contract_rule_id,
	}

	if system_prompt != "" {
		dataSend["system_prompt"] = system_prompt
	}

	if tools != "" {
		dataSend["tools"] = tools
	}

	if handover != "" {
		dataSend["handover"] = handover
	}

	status, response, err := Strapi.Create("agent-table", dataSend)
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

func UpdateAgent(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	system_prompt, ok := ctx.Params.GetString("system_prompt")
	if !ok {
		system_prompt = ""
	}

	llm, ok := ctx.Params.GetString("llm")
	if !ok {
		llm = ""
	}

	tools, ok := ctx.Params.GetString("tools")
	if !ok {
		tools = ""
	}

	handover, ok := ctx.Params.GetString("handover")
	if !ok {
		handover = ""
	}

	is_root_agent, ok := ctx.Params.GetBool("is_root_agent")
	if !ok {
		is_root_agent = false
	}

	dataSend := U.Dict{
		"name":          name,
		"description":   description,
		"llm":           llm,
		"is_root_agent": is_root_agent,
		"system_prompt": system_prompt,
		"tools":         tools,
		"handover":      handover,
	}

	status, response, err := Strapi.Update("agent-table", int64(id), dataSend)
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

func DeleteAgent(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, response, err := Strapi.Delete("agent-table", int64(id))
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

func GetAgents(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	contract_rule_id, ok := ctx.Params.GetFloat64("contract_rule_id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: contract_rule_id")
	}

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		page_size = 10
	}

	dataSend := U.Dict{
		"contract_rule_id": contract_rule_id,
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	}

	status, response, err := Strapi.Find("agent-table", dataSend)
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

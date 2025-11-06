/*
 * CR3 Client Dashboard
 * Package: gitlab.hexacode.org/cr3/client-dashboard
 * Maintainer: Azzis Arswendo <azzis@hexacode.org>
 * 						 Dean Ramhan <dean@hexacode.org>
 *
 * Copyright (C) 2024 Hexacode Teknologi Indonesia
 * All Rights Reserved
 */

package main

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.hexacode.org/go-libs/cr3call"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateContractRules(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	creator_id, ok := me.GetFloat64("id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	rule_name, ok := ctx.Params.GetString("rule_name")
	if !ok {
		return ctx.Error(400, "invalid parameters rule_name")
	}

	status, rule, err := Strapi.Create("contract-rules-table", U.Dict{
		"rule_name":  rule_name,
		"creator_id": creator_id,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, rule, nil
}

func UpdateContractRules(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	idString, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	rule_name, ok := ctx.Params.GetString("rule_name")
	if !ok {
		return ctx.Error(400, "invalid parameters rule_name")
	}

	status, rule, err := Strapi.Update("contract-rules-table", int64(id), U.Dict{
		"rule_name": rule_name,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, rule, nil
}

func GetContractRules(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	search, ok := ctx.Params.GetString("search")

	if !ok {
		search = ""
	}

	page, ok := ctx.Params.GetFloat64("page")

	if !ok {
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")

	if !ok {
		page_size = 10
	}

	if search == "" {
		status, rules, err := Strapi.Find("contract-rules-table", U.Dict{
			"pagination": U.Dict{
				"pageSize": page_size,
				"page":     page,
			},
			"populate": "*",
			"sort":     "id:desc",
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, rules, nil
	} else {
		status, rules, err := Strapi.Find("contract-rules-table", U.Dict{
			"filters": U.Dict{
				"rule_name": U.Dict{
					"$containsi": search,
				},
			},
			"pagination": U.Dict{
				"pageSize": page_size,
				"page":     page,
			},
			"populate": "*",
			"sort":     "id:desc",
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, rules, nil
	}
}

func GetContractRule(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	idString, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, rule, err := Strapi.Find("contract-rules-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": id,
			},
		},
		"populate": "*",
		"sort":     "id:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("rule", rule)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, rule, nil
}

func GetContractRuleSection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	contract_rule_idString, ok := ctx.Params.GetString("contract_rule_id")

	if !ok {
		return ctx.Error(400, "invalid parameters contract_rule_id")
	}

	contract_rule_id, err := strconv.Atoi(string(contract_rule_idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters contract_rule_id")
	}

	page, ok := ctx.Params.GetFloat64("page")

	if !ok {
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")

	if !ok {
		page_size = 10
	}

	status, section, err := Strapi.Find("contract-rule-sections-table", U.Dict{
		"filters": U.Dict{
			"contract_rule_id": U.Dict{
				"$eq": contract_rule_id,
			},
		},
		"pagination": U.Dict{
			"pageSize": page_size,
			"page":     page,
		},
		"populate": "*",
		"sort":     "id:asc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("section", section)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, section, nil
}

func DeleteRules(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	idString, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, _, err = Strapi.Delete("contract-rules-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	err = DeleteAllDocumentRules(strings.Replace(string(bearer), "Bearer ", "", 1), int(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	err = DeleteAllSectionRules(strings.Replace(string(bearer), "Bearer ", "", 1), int(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

func DeleteAllDocumentRules(token string, id int) error {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	page := 1
	for {
		status, documents, err := Strapi.Find("references-document-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"filters": U.Dict{
				"contract_rule_id": U.Dict{
					"$eq": int64(id),
				},
			},
		})

		if err != nil {
			fmt.Println(err)
			return err
		}

		if status != 200 {
			fmt.Println(status)
			return err
		}

		list, ok := documents.GetList("data")

		fmt.Println("list", list)

		if !ok {
			return err
		}

		if len(list) == 0 {
			break
		}

		for i := 0; i < len(list); i++ {
			dict, ok := list.GetDict(i)

			if !ok {
				return err
			}

			idDoc, ok := dict.GetFloat64("id")

			if !ok {
				return err
			}

			status, _, err := Strapi.Delete("references-document-table", int64(idDoc))

			if status != 200 {
				fmt.Println(status)
				return err
			}

			if err != nil {
				fmt.Println(err)
				return err
			}

			idDocStr := strconv.Itoa(int(idDoc))

			CR3 := cr3call.New(MODULE_MANAGER_URL)
			go CR3.Call("client-dashboard", "delete-document-collection-by-contract", cr3call.CR3CallData{
				Headers: U.Dict{
					"Authorization": token,
				},
				Params: U.Dict{
					"contract":        "REFERENCE",
					"contract_doc_id": idDocStr,
				},
			}, false)
		}
		page++
	}

	return nil
}

func DeleteAllSectionRules(token string, id int) error {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	page := 1
	for {
		status, documents, err := Strapi.Find("contract-rule-sections-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"filters": U.Dict{
				"contract_rule_id": U.Dict{
					"$eq": int64(id),
				},
			},
		})

		if err != nil {
			fmt.Println(err)
			return err
		}

		if status != 200 {
			fmt.Println(status)
			return err
		}

		list, ok := documents.GetList("data")

		fmt.Println("list", list)

		if !ok {
			return err
		}

		if len(list) == 0 {
			break
		}

		for i := 0; i < len(list); i++ {
			dict, ok := list.GetDict(i)

			if !ok {
				return err
			}

			idDoc, ok := dict.GetFloat64("id")

			if !ok {
				return err
			}

			status, _, err := Strapi.Delete("contract-rule-sections-table", int64(idDoc))

			if status != 200 {
				fmt.Println(status)
				return err
			}

			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		page++
	}

	return nil
}

func CreateContractRuleSection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	id, ok := me.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	contract_rule_idString, ok := ctx.Params.GetString("contract_rule_id")

	if !ok {
		return ctx.Error(400, "invalid parameters contract_rule_id")
	}

	contract_rule_id, err := strconv.Atoi(string(contract_rule_idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters contract_rule_id")
	}

	section_title, ok := ctx.Params.GetString("section_title")

	if !ok {
		return ctx.Error(400, "invalid parameters section_title")
	}

	section_prompt, ok := ctx.Params.GetString("section_prompt")

	if !ok {
		return ctx.Error(400, "invalid parameters section_prompt")
	}

	reference_prompt, ok := ctx.Params.GetString("reference_prompt")

	if !ok {
		return ctx.Error(400, "invalid parameters reference_prompt")
	}

	section_sample, ok := ctx.Params.GetString("section_sample")

	if !ok {
		return ctx.Error(400, "invalid parameters section_sample")
	}

	status, rule, err := Strapi.Create("contract-rule-sections-table", U.Dict{
		"contract_rule_id": contract_rule_id,
		"section_title":    section_title,
		"section_prompt":   section_prompt,
		"reference_prompt": reference_prompt,
		"section_sample":   section_sample,
		"creator_id":       id,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, rule, nil
}

func DeleteContractRuleSection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	idString, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, section, err := Strapi.Delete("contract-rule-sections-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, section, nil
}

func UpdateContractRuleSection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	idString, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idString))

	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	section_title, ok := ctx.Params.GetString("section_title")

	if !ok {
		return ctx.Error(400, "invalid parameters section_title")
	}

	section_prompt, ok := ctx.Params.GetString("section_prompt")

	if !ok {
		return ctx.Error(400, "invalid parameters section_prompt")
	}

	reference_prompt, ok := ctx.Params.GetString("reference_prompt")

	if !ok {
		return ctx.Error(400, "invalid parameters reference_prompt")
	}

	section_sample, ok := ctx.Params.GetString("section_sample")

	if !ok {
		return ctx.Error(400, "invalid parameters section_sample")
	}

	status, section, err := Strapi.Update("contract-rule-sections-table", int64(id), U.Dict{
		"section_title":    section_title,
		"section_prompt":   section_prompt,
		"reference_prompt": reference_prompt,
		"section_sample":   section_sample,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, section, nil
}

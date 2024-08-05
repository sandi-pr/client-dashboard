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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateReferencesDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	industry_id, ok := ctx.Params.GetFloat64("industry_id")
	if !ok {
		return ctx.Error(400, "invalid parameters industry_id")
	}

	referencesString, ok := ctx.Params.GetString("references")
	if !ok {
		return ctx.Error(400, "invalid parameters references")
	}

	var documents interface{}
	err = json.Unmarshal([]byte(referencesString), &documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

	datas := []U.Dict{}

	for _, document := range documents.([]interface{}) {

		documentBase64, ok := document.(map[string]interface{})["base64"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		documentFile, err := base64.StdEncoding.DecodeString(documentBase64)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(400, "invalid parameters documents")
		}

		documentName, ok := document.(map[string]interface{})["name"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		// file_size, ok := document.(map[string]interface{})["size"].(string)
		// if !ok {
		// 	return ctx.Error(400, "invalid parameters documents")
		// }

		role, ok := document.(map[string]interface{})["role"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		docType, ok := document.(map[string]interface{})["type"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		status, upload, err := Strapi.Upload(documentName, documentFile)

		fmt.Println("upload", upload)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		data, ok := upload.GetDict(0)

		if !ok {
			fmt.Println("failed to get dict")
			return ctx.Error(500, "Internal Server Error")
		}

		id, ok := data.GetFloat64("id")
		if !ok {
			fmt.Println("failed to get id")
			return ctx.Error(500, "Internal Server Error")
		}

		status, documents, err := Strapi.Create("references-document-table", U.Dict{
			"industry_id":        industry_id,
			"document_for":       role,
			"document_title":     docType,
			"document_name":      documentName,
			"document_file_name": documentName,
			"document_media_id":  id,
			"creator_id":         creator_id,
		})

		fmt.Println("document", documents)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			return ctx.Error(500, "Internal Server Error")
		}

		// CR3 := cr3call.New(MODULE_MANAGER_URL)
		// go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		// 	Headers: U.Dict{
		// 		"Authorization": bearer,
		// 	},
		// }, false)

		result, ok := documents.GetDict("data")
		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		datas = append(datas, result)
	}

	return 200, U.Dict{"data": datas}, nil
}

func GetReferencesDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	industry_id, ok := ctx.Params.GetFloat64("industry_id")
	if !ok {
		return ctx.Error(400, "invalid parameters industry_id")
	}

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		return ctx.Error(400, "invalid parameters page")
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		return ctx.Error(400, "invalid parameters page_size")
	}

	status, documents, err := Strapi.Find("references-document-table", U.Dict{
		"filters": U.Dict{
			"industry_id": U.Dict{
				"$eq": int64(industry_id),
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"populate": "*",
		"sort":     "id:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(documents)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := documents.GetList("data")
	if len(list) == 0 {
		return 200, U.Dict{"files": []U.Dict{}}, nil
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	var datas []U.Dict

	for i := 0; i < len(list); i++ {
		document_url := list[i].(U.Dict)["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		document_url = U.String(strings.TrimPrefix(string(document_url), "/"))

		resp, err := http.Get(STRAPI_URL + string(document_url))
		if err != nil {
			fmt.Printf("error saat mengunduh file: %v", err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error saat membaca isi file: %v", err)
		}

		base64Data := base64.StdEncoding.EncodeToString(data)

		datas = append(datas, U.Dict{
			"id":     list[i].(U.Dict)["id"].(U.Float64),
			"type":   list[i].(U.Dict)["attributes"].(U.Dict)["document_title"].(U.String),
			"role":   list[i].(U.Dict)["attributes"].(U.Dict)["document_for"].(U.String),
			"name":   list[i].(U.Dict)["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["name"].(U.String),
			"size":   list[i].(U.Dict)["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["size"].(U.Float64),
			"base64": base64Data,
		})
	}
	// dict, ok := list.GetDict(0)

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	return 200, U.Dict{"files": datas}, nil
}

func DeleteReferencesDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(idString))

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	status, _, err = Strapi.Delete("references-document-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

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

	page, ok := ctx.Params.GetFloat64("page")

	if !ok {
		return ctx.Error(400, "invalid parameters page")
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")

	if !ok {
		return ctx.Error(400, "invalid parameters page_size")
	}

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

	status, section, err := Strapi.Find("contract-rule-sections-table", U.Dict{
		"filters": U.Dict{
			"contract_rule_id": U.Dict{
				"$eq": contract_rule_id,
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
		fmt.Println("section", section)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, section, nil
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

	section_sample, ok := ctx.Params.GetString("section_sample")

	if !ok {
		return ctx.Error(400, "invalid parameters section_sample")
	}

	status, rule, err := Strapi.Create("contract-rule-sections-table", U.Dict{
		"contract_rule_id": contract_rule_id,
		"section_title":    section_title,
		"section_prompt":   section_prompt,
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

func CreateContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	project_name, ok := ctx.Params.GetString("project_name")
	if !ok {
		return ctx.Error(400, "invalid parameters project_name")
	}

	industry_type, ok := ctx.Params.GetString("industry_type")
	if !ok {
		return ctx.Error(400, "invalid parameters industry_type")
	}

	project_as, ok := ctx.Params.GetString("project_as")
	if !ok {
		return ctx.Error(400, "invalid parameters project_as")
	}

	project_start_date, ok := ctx.Params.GetString("project_start_date")
	if !ok {
		return ctx.Error(400, "invalid parameters project_start_date")
	}

	project_end_date, ok := ctx.Params.GetString("project_end_date")
	if !ok {
		return ctx.Error(400, "invalid parameters project_end_date")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	if description == "" {
		status, res, err := Strapi.Create("project-contract-table", U.Dict{
			"project_name":       project_name,
			"industry_type":      industry_type,
			"project_as":         project_as,
			"project_start_date": project_start_date,
			"project_end_date":   project_end_date,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, res, nil
	} else {
		status, res, err := Strapi.Create("project-contract-table", U.Dict{
			"project_name":       project_name,
			"industry_type":      industry_type,
			"project_as":         project_as,
			"project_start_date": project_start_date,
			"project_end_date":   project_end_date,
			"description":        description,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, res, nil
	}
}

func GetContractProjects(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, projects, err := Strapi.Find("project-contract-table", U.Dict{
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort":     "id:desc",
		"populate": "*",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, projects, nil
}

func GetContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	id, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, project, err := Strapi.Find("project-contract-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, project, nil
}

func GetDocumentContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	project_id, ok := ctx.Params.GetString("project_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, documents, err := Strapi.Find("document-project-table", U.Dict{
		"filters": U.Dict{
			"project_id": U.Dict{
				"$eq": project_id,
			},
		},
		"populate": "*",
	})

	fmt.Println(documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := documents.GetList("data")

	fmt.Println("list", list)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(500, "Internal Server Error")
	}

	var datas []U.Dict

	for i := 0; i < len(list); i++ {
		document_url := list[i].(U.Dict)["attributes"].(U.Dict)["document"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		document_url = U.String(strings.TrimPrefix(string(document_url), "/"))

		resp, err := http.Get(STRAPI_URL + string(document_url))
		if err != nil {
			fmt.Printf("error saat mengunduh file: %v", err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error saat membaca isi file: %v", err)
		}

		base64Data := base64.StdEncoding.EncodeToString(data)

		datas = append(datas, U.Dict{
			"type":   list[i].(U.Dict)["attributes"].(U.Dict)["document_type"].(U.String),
			"name":   list[i].(U.Dict)["attributes"].(U.Dict)["document"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["name"].(U.String),
			"size":   list[i].(U.Dict)["attributes"].(U.Dict)["file_size"].(U.String),
			"base64": base64Data,
			"id":     list[i].(U.Dict)["id"].(U.Float64),
		})
	}

	return 200, U.Dict{"files": datas}, nil
}

func CreateDocumentContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	project_id, ok := ctx.Params.GetString("project_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	documentString, ok := ctx.Params.GetString("document")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var document interface{}
	err = json.Unmarshal([]byte(documentString), &document)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	documentBase64, ok := document.(map[string]interface{})["base64"].(string)
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	documentFile, err := base64.StdEncoding.DecodeString(documentBase64)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

	documentName, ok := document.(map[string]interface{})["name"].(string)
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	document_type, ok := document.(map[string]interface{})["type"].(string)
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	file_size, ok := document.(map[string]interface{})["size"].(string)
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	status, upload, err := Strapi.Upload(documentName, documentFile)

	fmt.Println("upload", upload)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := upload.GetDict(0)

	if !ok {
		fmt.Println("failed to get dict")
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := data.GetFloat64("id")
	if !ok {
		fmt.Println("failed to get id")
		return ctx.Error(500, "Internal Server Error")
	}

	status, documentProject, err := Strapi.Create("document-project-table", U.Dict{
		"project_id":    project_id,
		"document_type": document_type,
		"document":      id,
		"file_size":     file_size,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, documentProject, nil
}

func DeleteDocumentContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(idString))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, _, err = Strapi.Delete("document-project-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

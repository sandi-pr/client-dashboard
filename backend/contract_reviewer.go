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

	id, ok := me.GetFloat64("id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	project_name, ok := ctx.Params.GetString("project_name")
	if !ok {
		return ctx.Error(400, "invalid parameters project_name")
	}

	industry_id, ok := ctx.Params.GetFloat64("industry_id")
	if !ok {
		return ctx.Error(400, "invalid parameters industry_id")
	}

	contract_rule_id, ok := ctx.Params.GetFloat64("contract_rule_id")
	if !ok {
		return ctx.Error(400, "invalid parameters contract_rule_id")
	}

	project_as, ok := ctx.Params.GetString("project_as")
	if !ok {
		return ctx.Error(400, "invalid parameters project_as")
	}

	duration_start, ok := ctx.Params.GetString("duration_start")
	if !ok {
		return ctx.Error(400, "invalid parameters duration_start")
	}

	duration_end, ok := ctx.Params.GetString("duration_end")
	if !ok {
		return ctx.Error(400, "invalid parameters duration_end")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	if description == "" {
		status, res, err := Strapi.Create("projects-table", U.Dict{
			"project_name":     project_name,
			"industry_id":      industry_id,
			"contract_rule_id": contract_rule_id,
			"project_as":       project_as,
			"duration_start":   duration_start,
			"duration_end":     duration_end,
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

		return 200, res, nil
	} else {
		status, res, err := Strapi.Create("projects-table", U.Dict{
			"project_name":     project_name,
			"industry_id":      industry_id,
			"contract_rule_id": contract_rule_id,
			"project_as":       project_as,
			"duration_start":   duration_start,
			"duration_end":     duration_end,
			"creator_id":       id,
			"description":      description,
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

	status, projects, err := Strapi.Find("projects-table", U.Dict{
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

	idString, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(idString))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, project, err := Strapi.Find("projects-table", U.Dict{
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

	project_idStr, ok := ctx.Params.GetString("project_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	project_id, err := strconv.Atoi(string(project_idStr))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, documents, err := Strapi.Find("project-documents-table", U.Dict{
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
		return 200, U.Dict{"files": []U.Dict{}}, nil
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

		ragID := ""
		if val, ok := list[i].(U.Dict)["attributes"].(U.Dict)["rag_id"].(U.String); ok {
			ragID = string(val)
		}

		status := ""
		if val, ok := list[i].(U.Dict)["attributes"].(U.Dict)["status"].(U.String); ok {
			status = string(val)
		}

		status_generated := ""
		if val, ok := list[i].(U.Dict)["attributes"].(U.Dict)["status_generated"].(U.String); ok {
			status_generated = string(val)
		}

		datas = append(datas, U.Dict{
			"type":             list[i].(U.Dict)["attributes"].(U.Dict)["document_title"].(U.String),
			"name":             list[i].(U.Dict)["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["name"].(U.String),
			"base64":           base64Data,
			"id":               list[i].(U.Dict)["id"].(U.Float64),
			"rag_id":           ragID,
			"status":           status,
			"status_generated": status_generated,
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

	id, ok := me.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	project_idStr, ok := ctx.Params.GetString("project_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	project_id, err := strconv.Atoi(string(project_idStr))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	content_date, ok := ctx.Params.GetString("content_date")

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

	docId, ok := data.GetFloat64("id")
	if !ok {
		fmt.Println("failed to get id")
		return ctx.Error(500, "Internal Server Error")
	}

	status, documentProject, err := Strapi.Create("project-documents-table", U.Dict{
		"project_id":         project_id,
		"content_date":       content_date,
		"document_title":     document_type,
		"document_name":      documentName,
		"document_file_name": documentName,
		"document_media_id":  docId,
		"creator_id":         id,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(documentProject)
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

	status, _, err = Strapi.Delete("project-documents-table", int64(id))
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

func UpdateDocumentContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	rag_id, ok := ctx.Params.GetString("rag_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, document, err := Strapi.Update("project-documents-table", int64(id), U.Dict{
		"rag_id": rag_id,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func UpdateDocumentContractProjectStatusGenerated(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status_generated, ok := ctx.Params.GetString("status_generated")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, document, err := Strapi.Update("project-documents-table", int64(id), U.Dict{
		"status_generated": status_generated,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func GetIdsProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, data, err := Strapi.Find("project-documents-table", U.Dict{
		"filters": U.Dict{
			"industry_id": U.Dict{
				"$eq": industry_id,
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
		fmt.Println(data)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := data.GetList("data")
	if !ok {
		return 200, U.Dict{"data": []U.Dict{}}, nil
	}

	var ids []string
	for i := 0; i < len(list); i++ {
		id, ok := list[i].(U.Dict).GetFloat64("id")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		ids = append(ids, strconv.Itoa(int(id)))
	}

	return 200, U.Dict{"ids_reference": ids}, nil
}

func CreateProjectResults(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	project_idStr, ok := ctx.Params.GetString("project_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	project_id, err := strconv.Atoi(string(project_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	project_document_idStr, ok := ctx.Params.GetString("project_document_id")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	project_document_id, err := strconv.Atoi(string(project_document_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	resultString, ok := ctx.Params.GetString("results")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var result interface{}
	err = json.Unmarshal([]byte(resultString), &result)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	contract_rule_idStr, ok := ctx.Params.GetString("contract_rule_id")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	contract_rule_id, err := strconv.Atoi(string(contract_rule_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	generated_at, ok := ctx.Params.GetString("generated_at")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	statusString, ok := ctx.Params.GetString("status")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, project_result, err := Strapi.Create("project-results-table", U.Dict{
		"project_id":          int64(project_id),
		"project_document_id": int64(project_document_id),
		"result":              result,
		"contract_rule_id":    int64(contract_rule_id),
		"generated_at":        generated_at,
		"status":              statusString,
		"creator_id":          int64(id),
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(project_result)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, project_result, nil
}

func GetContractProjectResultWithDocumentId(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	documentIdString, ok := ctx.Params.GetString("project_document_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(documentIdString))

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	status, documents, err := Strapi.Find("project-results-table", U.Dict{
		"filters": U.Dict{
			"project_document_id": U.Dict{
				"$eq": int64(id),
			},
		},
		"populate": "*",
		"sort":     "id:asc",
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

	return 200, documents, nil
}

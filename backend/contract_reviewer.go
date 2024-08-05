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

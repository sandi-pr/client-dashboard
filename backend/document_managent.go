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
	"fmt"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	document_title, ok := ctx.Params.GetString("document_title")
	if !ok {
		return ctx.Error(400, "invalid parameters document_title")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	document_type, ok := ctx.Params.GetString("document_type")
	if !ok {
		return ctx.Error(400, "invalid parameters document_type")
	}

	documentBase64, ok := ctx.Params.GetString("document")
	if !ok {
		return ctx.Error(400, "invalid parameters document")
	}

	document_name, ok := ctx.Params.GetString("document_name")
	if !ok {
		return ctx.Error(400, "invalid parameters document_name")
	}

	person, ok := ctx.Params.GetFloat64("person")
	if !ok {
		fmt.Println("not found person")
	}

	date_start, ok := ctx.Params.GetString("contract_start")
	if !ok {
		fmt.Println("invalid parameters date_start")
	}

	date_end, ok := ctx.Params.GetString("contract_end")
	if !ok {
		fmt.Println("invalid parameters date_end")
	}

	document, err := base64.StdEncoding.DecodeString(string(documentBase64))
	if err != nil {
		return ctx.Error(400, "failed to decode document")
	}

	status, upload, err := Strapi.Upload(string(document_name), document)

	fmt.Println("upload", upload)

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := upload.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := data.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if document_type == "GENERAL" {

		status, documents, err := Strapi.Create("document-management-table", U.Dict{
			"document_title": document_title,
			"description":    description,
			"document_type":  document_type,
			"document":       id,
		})

		fmt.Println("document", documents)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			return ctx.Error(500, "Internal Server Error")
		}

		return status, documents, nil
	} else {
		status, documents, err := Strapi.Create("document-management-table", U.Dict{
			"document_title": document_title,
			"description":    description,
			"document_type":  document_type,
			"document":       id,
			"person":         person,
			"contract_start": date_start,
			"contract_end":   date_end,
		})

		fmt.Println("document", documents)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			return ctx.Error(500, "Internal Server Error")
		}

		return status, documents, nil
	}
}

func GetDocuments(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		fmt.Println("invalid parameters")
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		page_size = 10
	}

	status, documents, err := Strapi.Find("document-management-table", U.Dict{
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

	return 200, documents, nil
}

func GetDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, document, err := Strapi.Find("document-management-table", U.Dict{
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
		fmt.Println(document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	// list, ok := document.GetList("data")
	// if len(list) == 0 {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// dict, ok := list.GetDict(0)

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// document_url, ok := dict.GetString("url")

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// fmt.Println(document_url)

	return 200, document, nil
}

func DeleteDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, document, err := Strapi.Delete("document-management-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

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
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func SetReferences(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	industry_type, ok := me.GetString("industry_type")
	if !ok {
		return ctx.Error(400, "invalid parameters industry_type")
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

		file_size, ok := document.(map[string]interface{})["size"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		role, ok := document.(map[string]interface{})["role"].(string)
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

		status, documents, err := Strapi.Create("reference-contract-table", U.Dict{
			"industry_type": industry_type,
			"document_as":   documentName,
			"document_type": role,
			"document":      id,
			"file_size":     file_size,
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

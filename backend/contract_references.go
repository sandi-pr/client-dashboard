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

	contract_rule_id, ok := ctx.Params.GetFloat64("contract_rule_id")
	if !ok {
		return ctx.Error(400, "invalid parameters contract_rule_id")
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
			"contract_rule_id":   contract_rule_id,
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

func UpdateReferencesDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	rag_id, ok := ctx.Params.GetString("rag_id")
	if !ok {
		return ctx.Error(400, "invalid parameters rag_id")
	}

	status, update, err := Strapi.Update("references-document-table", int64(id), U.Dict{
		"rag_id": rag_id,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(update)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{"data": update}, nil
}

func GetIdsReference(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, data, err := Strapi.Find("references-document-table", U.Dict{
		"filters": U.Dict{
			"industry_id": U.Dict{
				"$eq": industry_id,
			},
		},
		"populate": "*",
		"sort":     "id:desc",
	})

	fmt.Println("data", data)

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
		return 200, U.Dict{"ids_reference": []U.Dict{}}, nil
	}

	var ids []string
	for i := 0; i < len(list); i++ {
		attributes, ok := list[i].(U.Dict).GetDict("attributes")
		if !ok {
			fmt.Println("attributes", list[i])
			return 500, U.Dict{"ids_reference": []U.Dict{}}, nil
		}

		rag_id, ok := attributes.GetString("rag_id")
		if !ok {
			fmt.Println("rag_id", attributes)
			return 500, U.Dict{"ids_reference": []U.Dict{}}, nil
		}

		ids = append(ids, string(rag_id))
	}

	return 200, U.Dict{"ids_reference": ids}, nil
}

func CheckReferenceRag(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, data, err := Strapi.Find("references-document-table", U.Dict{
		"filters": U.Dict{
			"industry_id": U.Dict{
				"$eq": industry_id,
			},
		},
		"populate": "*",
		"sort":     "id:desc",
	})

	fmt.Println("data", data)

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
		return 200, U.Dict{"ready": false}, nil
	}

	if len(list) < 2 {
		return 200, U.Dict{"ready": false}, nil
	}

	if len(list) > 2 {
		return 200, U.Dict{"ready": false}, nil
	}

	for i := 0; i < len(list); i++ {
		attributes, ok := list[i].(U.Dict).GetDict("attributes")
		if !ok {
			fmt.Println("attributes", list[i])
			return 500, U.Dict{"ready": false}, nil
		}

		rag_id, ok := attributes.GetString("rag_id")
		if !ok {
			fmt.Println("rag_id", rag_id)
			return 500, U.Dict{"ready": false}, nil
		}
	}

	return 200, U.Dict{"ready": true}, nil
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

	contract_rule_id, ok := ctx.Params.GetFloat64("contract_rule_id")
	if !ok {
		return ctx.Error(400, "invalid parameters contract_rule_id")
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
			"contract_rule_id": U.Dict{
				"$eq": int64(contract_rule_id),
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
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

		status := ""
		if val, ok := list[i].(U.Dict)["attributes"].(U.Dict)["status"].(U.String); ok {
			status = string(val)
		}

		datas = append(datas, U.Dict{
			"id":     list[i].(U.Dict)["id"].(U.Float64),
			"type":   list[i].(U.Dict)["attributes"].(U.Dict)["document_title"].(U.String),
			"role":   list[i].(U.Dict)["attributes"].(U.Dict)["document_for"].(U.String),
			"name":   list[i].(U.Dict)["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["name"].(U.String),
			"size":   list[i].(U.Dict)["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["size"].(U.Float64),
			"base64": base64Data,
			"status": status,
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

func GetContractProjectResult(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, documents, err := Strapi.Find("project-results-table", U.Dict{
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
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, documents, nil
}

func UpdateContractProjectResult(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	resultStr, ok := ctx.Params.GetString("result")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var result interface{}
	err = json.Unmarshal([]byte(resultStr), &result)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	status, results, err := Strapi.Update("project-results-table", int64(id), U.Dict{
		"result": result,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, results, nil
}

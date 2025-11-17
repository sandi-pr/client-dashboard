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
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.hexacode.org/go-libs/cr3call"
	"gitlab.hexacode.org/go-libs/hctypes"
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

	is_compare, ok := ctx.Params.GetBool("is_compare")
	if !ok {
		is_compare = false
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
			"is_compare":       is_compare,
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
			"is_compare":       is_compare,
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

func UpdateContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	project_id_str, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "invalid parameters project_id")
	}

	project_id, err := strconv.Atoi(string(project_id_str))
	if err != nil {
		return ctx.Error(400, "invalid parameters project_id")
	}

	project_name, ok := ctx.Params.GetString("project_name")
	if !ok {
		return ctx.Error(400, "invalid parameters project_name")
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
		status, res, err := Strapi.Update("projects-table", int64(project_id), U.Dict{
			"project_name":   project_name,
			"duration_start": duration_start,
			"duration_end":   duration_end,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			fmt.Println(res)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, res, nil
	} else {
		status, res, err := Strapi.Update("projects-table", int64(project_id), U.Dict{
			"project_name":   project_name,
			"duration_start": duration_start,
			"duration_end":   duration_end,
			"description":    description,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			fmt.Println(res)
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

	is_compare, ok := ctx.Params.GetBool("is_compare")

	if !ok {
		is_compare = false
	}

	status, projects, err := Strapi.Find("projects-table", U.Dict{
		"filters": U.Dict{
			"is_compare": U.Dict{
				"$eq": is_compare,
			},
		},
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

func DeleteContractProject(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, _, err = Strapi.Delete("projects-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	err = DeleteAllDocumentProject(strings.Replace(string(bearer), "Bearer ", "", 1), int(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

func DeleteAllDocumentProject(token string, id int) error {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	page := 1
	for {
		status, documents, err := Strapi.Find("project-documents-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"filters": U.Dict{
				"project_id": U.Dict{
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

			err := DeleteProjectResult(int(idDoc), token)

			if err != nil {
				fmt.Println(err)
				return err
			}

			status, _, err := Strapi.Delete("project-documents-table", int64(idDoc))

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
					"contract":        "PROJECT",
					"contract_doc_id": idDocStr,
				},
			}, false)
		}
		page++
	}

	return nil
}

func DeleteProjectResult(id int, token string) error {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	page := 1
	for {

		status, result, err := Strapi.Find("project-results-table", U.Dict{
			"filters": U.Dict{
				"project_document_id": U.Dict{
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

		list, ok := result.GetList("data")

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

			idResult, ok := dict.GetFloat64("id")

			if !ok {
				return err
			}

			status, _, err := Strapi.Delete("project-results-table", int64(idResult))

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
			"name":             list[i].(U.Dict)["attributes"].(U.Dict)["document_file_name"].(U.String),
			"base64":           base64Data,
			"id":               list[i].(U.Dict)["id"].(U.Float64),
			"rag_id":           ragID,
			"status":           status,
			"status_generated": status_generated,
		})
	}

	return 200, U.Dict{"files": datas}, nil
}

func GetContractProjectWithDocumentId(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	id, ok := ctx.Params.GetFloat64("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, document, err := Strapi.Find("project-documents-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": int64(id),
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

	list, ok := document.GetList("data")
	if len(list) == 0 {
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	document_url := dict["attributes"].(U.Dict)["document_media_id"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
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

	fmt.Println(base64Data)

	return 200, U.Dict{"document": document, "base64": base64Data}, nil
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

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai-contract-review", "process-to-rag", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
		Params: U.Dict{
			"project_id": project_idStr,
		},
	}, false)

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

	return 200, documents, nil
}

func GetAllReportWithChart(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	bearer, ok := ctx.Headers.GetString("Authorization")
	if !ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(
		STRAPI_URL,
		strings.Replace(string(bearer), "Bearer ", "", 1),
	)

	status, me, err := Strapi.UsersMe()
	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println("ME:", me)

	// ---- GET ALL DATA FROM STRAPI USING PAGINATION ----
	allDocs, err := StrapiFindAll(Strapi, "project-results-table", 25)
	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	// ---- GROUP BY MONTH & PROJECT ----
	// output: map[month][][]tabular_result
	grouped := ExtractTabularResults(allDocs)

	// ---- BUILD 12 MONTH LIST ----
	last12 := GetLast12Months()

	chartData := []U.Dict{}
	var summaryProjects [][]hctypes.Dict

	// ---- GENERATE CHART DATA ----
	for _, month := range last12 {

		projects, exists := grouped[month]

		if exists {
			// projects = [][]hctypes.Dict (project → tabular_result[])
			total, valid, invalid := CalculatePeopleStatus_MultiFiles(projects)

			chartData = append(chartData, U.Dict{
				"month":        month,
				"total_people": total,
				"valid":        valid,
				"invalid":      invalid,
			})

			summaryProjects = append(summaryProjects, projects...)

		} else {
			chartData = append(chartData, U.Dict{
				"month":        month,
				"total_people": 0,
				"valid":        0,
				"invalid":      0,
			})
		}
	}

	// ---- SUMMARY: HITUNG TOTAL PER PROJECT (HANYA LAST 12 MONTH) ----
	allTotal, allValid, allInvalid := CalculatePeopleStatus_MultiFiles(summaryProjects)

	// ---- PERSENTASE: PERBANDINGAN BULAN INI VS BULAN LALU ----
	currentMonth := time.Now().Format("2006-01")
	previousMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

	currentTotal, currentValid, currentInvalid := CalculateMonthlyPeopleStatus(grouped, currentMonth)
	prevTotal, prevValid, prevInvalid := CalculateMonthlyPeopleStatus(grouped, previousMonth)

	return 200, U.Dict{
		"chart": chartData,
		"summary": U.Dict{
			"total_people":            allTotal,
			"valid":                   allValid,
			"invalid":                 allInvalid,
			"percentage_total_people": CalculatePercentageChange(prevTotal, currentTotal),
			"percentage_valid":        CalculatePercentageChange(prevValid, currentValid),
			"percentage_invalid":      CalculatePercentageChange(prevInvalid, currentInvalid),
		},
	}, nil
}

func ExtractTabularResults(strapiDocs hctypes.List) map[string][][]hctypes.Dict {
	result := map[string][][]hctypes.Dict{}

	for i := 0; i < len(strapiDocs); i++ {

		doc, ok := strapiDocs.GetDict(i)
		if !ok {
			continue
		}

		attr, ok := doc.GetDict("attributes")
		if !ok {
			continue
		}

		createdAtRaw, _ := attr.GetString("createdAt")
		month := GetMonthFromCreatedAt(string(createdAtRaw))
		if month == "" {
			continue
		}

		resultList, ok := attr.GetList("result")
		if !ok {
			continue
		}

		// satu project: list tabular_result
		oneProject := []hctypes.Dict{}

		for j := 0; j < len(resultList); j++ {
			item, ok := resultList.GetDict(j)
			if !ok {
				continue
			}

			gen, ok := item.GetDict("generated")
			if !ok {
				continue
			}

			tab, ok := gen.GetDict("tabular_result")
			if !ok {
				continue
			}

			oneProject = append(oneProject, tab)
		}

		result[month] = append(result[month], oneProject)
	}

	return result
}

func CalculatePeopleStatus_MultiFiles(allProjects [][]hctypes.Dict) (total int, totalValid int, totalInvalid int) {

	for _, tabs := range allProjects {

		statusByPerson := map[int][]string{}

		for _, tab := range tabs {

			logs, ok := tab.GetList("logs")
			if !ok {
				continue
			}

			for i := 0; i < len(logs); i++ {

				logEntry, ok := logs.GetDict(i)
				if !ok {
					continue
				}

				rowIdxFloat, ok := logEntry.GetFloat64("row_index")
				if !ok {
					continue
				}

				rowIndex := int(rowIdxFloat)

				statusRaw, _ := logEntry.GetString("status")
				status := strings.ToLower(strings.TrimSpace(string(statusRaw)))

				statusByPerson[rowIndex] = append(statusByPerson[rowIndex], status)
			}
		}

		for _, statuses := range statusByPerson {

			total++

			allValid := true

			for _, s := range statuses {
				if s == "invalid" || s == "unknown" {
					allValid = false
				}
			}

			if allValid {
				totalValid++
			} else {
				totalInvalid++
			}
		}
	}

	return
}

func CalculatePeopleStatus_PerProject(allProjects hctypes.List) (total int, totalValid int, totalInvalid int) {

	for i := 0; i < len(allProjects); i++ {

		project, ok := allProjects.GetDict(i)
		if !ok {
			continue
		}

		attr, ok := project.GetDict("attributes")
		if !ok {
			continue
		}

		resultList, ok := attr.GetList("result")
		if !ok {
			continue
		}

		statusByPerson := map[int][]string{}

		for j := 0; j < len(resultList); j++ {

			item, ok := resultList.GetDict(j)
			if !ok {
				continue
			}

			generated, ok := item.GetDict("generated")
			if !ok {
				continue
			}

			tab, ok := generated.GetDict("tabular_result")
			if !ok {
				continue
			}

			logs, ok := tab.GetList("logs")
			if !ok {
				continue
			}

			for k := 0; k < len(logs); k++ {

				logEntry, ok := logs.GetDict(k)
				if !ok {
					continue
				}

				rowIdxFloat, ok := logEntry.GetFloat64("row_index")
				if !ok {
					continue
				}

				rowIndex := int(rowIdxFloat)

				statusRaw, _ := logEntry.GetString("status")
				status := strings.ToLower(strings.TrimSpace(string(statusRaw)))

				statusByPerson[rowIndex] = append(statusByPerson[rowIndex], status)
			}
		}

		for _, statuses := range statusByPerson {

			total++

			allValid := true

			for _, s := range statuses {
				if s == "invalid" || s == "unknown" {
					allValid = false
				}
			}

			if allValid {
				totalValid++
			} else {
				totalInvalid++
			}
		}
	}

	return
}

func CalculateMonthlyPeopleStatus(grouped map[string][][]hctypes.Dict, month string) (total int, totalValid int, totalInvalid int) {
	projects, ok := grouped[month]
	if !ok {
		return
	}

	return CalculatePeopleStatus_MultiFiles(projects)
}

func CalculatePercentageChange(previous int, current int) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}

	change := (float64(current-previous) / float64(previous)) * 100
	return math.Round(change*100) / 100
}

func GetLast12Months() []string {
	now := time.Now()
	months := []string{}

	for i := 0; i < 12; i++ {
		t := now.AddDate(0, -i, 0)
		months = append(months, t.Format("2006-01"))
	}

	for i, j := 0, len(months)-1; i < j; i, j = i+1, j-1 {
		months[i], months[j] = months[j], months[i]
	}

	return months
}

func GetProjectsWithResults(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	bearer, ok := ctx.Headers.GetString("Authorization")
	if !ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(
		STRAPI_URL,
		strings.Replace(string(bearer), "Bearer ", "", 1),
	)

	status, _, err := Strapi.UsersMe()
	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	page := 1
	pageSize := 10

	pageStr, _ := ctx.Params.GetFloat64("page")
	pageSizeStr, _ := ctx.Params.GetFloat64("page_size")

	if pageStr > 0 {
		page = int(pageStr)
	}

	if pageSizeStr > 0 {
		if pageSizeStr > 25 {
			pageSize = 25
		} else {
			pageSize = int(pageSizeStr)
		}
	}

	datas, err := StrapiFindAll(Strapi, "project-results-table", 25)

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	projectIDSet := map[int]bool{}
	for _, item := range datas {
		attr, _ := item.(hctypes.Dict).GetDict("attributes")
		pid, ok := attr.GetFloat64("project_id")
		if ok {
			projectIDSet[int(pid)] = true
		}
	}

	if len(projectIDSet) == 0 {
		return 200, U.Dict{
			"meta": U.Dict{
				"page":      page,
				"pageSize":  pageSize,
				"total":     0,
				"pageCount": 0,
			},
			"data": U.List{},
		}, nil
	}

	projectIDs := []int64{}
	for id := range projectIDSet {
		projectIDs = append(projectIDs, int64(id))
	}

	status, projects, err := Strapi.Find("projects-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$in": projectIDs,
			},
		},
		"populate": "*",
		"pagination": U.Dict{
			"pageSize": pageSize,
			"page":     page,
		},
	})

	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}
	return 200, projects, nil
}

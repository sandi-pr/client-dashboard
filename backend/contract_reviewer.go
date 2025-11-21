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

func FilterByKeyValue(dicts []U.Dict, key string, value interface{}) []U.Dict {
	result := []U.Dict{}
	for _, d := range dicts {
		if v, ok := d[key]; ok && v == value {
			result = append(result, d)
		}
	}
	return result
}

func filterSummariesByStatus(dicts []U.Dict, target string) []U.Dict {
	result := []U.Dict{}
	normalizedTarget := strings.ToLower(strings.TrimSpace(target))
	if normalizedTarget == "" {
		return result
	}

	getString := func(value interface{}) string {
		if value == nil {
			return ""
		}
		switch v := value.(type) {
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprintf("%v", v)
		}
	}

	for _, d := range dicts {
		statusValue := strings.ToLower(strings.TrimSpace(getString(d["status"])))
		statusLabel := strings.ToLower(strings.TrimSpace(getString(d["status_label"])))
		if statusValue == normalizedTarget || statusLabel == normalizedTarget {
			result = append(result, d)
		}
	}

	return result
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

	logPrefix := "[GetContractProjectResult]"
	fmt.Println(me)

	idStringRaw, ok := ctx.URLParams.GetString("id")

	if !ok {
		idStringRaw, ok = ctx.Params.GetString("id")
	}

	if !ok {
		idStringRaw, ok = ctx.Params.GetString("ids")
	}

	if !ok {
		fmt.Println(logPrefix, "missing id parameter")
		return ctx.Error(400, "invalid parameters")
	}

	idParts := strings.Split(string(idStringRaw), ",")
	if len(idParts) == 0 {
		fmt.Println(logPrefix, "empty id parts", idStringRaw)
		return ctx.Error(400, "invalid parameters")
	}

	parsedIDs := []int64{}
	for _, part := range idParts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			fmt.Println(logPrefix, "empty id entry ignored")
			continue
		}

		idValue, err := strconv.Atoi(trimmed)
		if err != nil {
			fmt.Println(logPrefix, "failed to parse id", trimmed, err)
			return ctx.Error(500, "Internal Server Error")
		}

		parsedIDs = append(parsedIDs, int64(idValue))
	}

	if len(parsedIDs) == 0 {
		fmt.Println(logPrefix, "no valid id parsed from", idStringRaw)
		return ctx.Error(400, "invalid parameters")
	}

	idFilter := U.Dict{}
	if len(parsedIDs) == 1 {
		idFilter["$eq"] = parsedIDs[0]
	} else {
		idInterfaces := make([]interface{}, len(parsedIDs))
		for i, v := range parsedIDs {
			idInterfaces[i] = v
		}
		idFilter["$in"] = idInterfaces
	}

	status, documents, err := Strapi.Find("project-results-table", U.Dict{
		"filters": U.Dict{
			"id": idFilter,
		},
		"populate": "*",
		"sort":     "id:desc",
	})

	if err != nil {
		fmt.Println(logPrefix, "strapi find error", err)
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(logPrefix, "strapi status", status, "response", documents)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := documents.GetList("data")
	if !ok {
		fmt.Println(logPrefix, "failed to get data list from documents")
		return ctx.Error(500, "Internal Server Error")
	}

	meta, ok := documents.GetDict("meta")
	if !ok {
		fmt.Println(logPrefix, "failed to get meta from documents")
		return ctx.Error(500, "Internal Server Error")
	}

	aggregatedSummary := []U.Dict{}
	for i := 0; i < len(list); i++ {
		dict, ok := list.GetDict(i)
		if !ok {
			fmt.Println(logPrefix, "failed to get dict from data list at index", i)
			return ctx.Error(500, "Internal Server Error")
		}

		attributes, ok := dict.GetDict("attributes")
		if !ok {
			fmt.Println(logPrefix, "failed to get attributes from dict at index", i)
			return ctx.Error(500, "Internal Server Error")
		}

		tabulars := extractTabularResultsFromAttributes(attributes)
		entries := buildSummaryEntriesFromTabular(tabulars)
		aggregatedSummary = append(aggregatedSummary, entries...)
	}

	for idx := range aggregatedSummary {
		aggregatedSummary[idx]["no"] = idx + 1
	}

	validList := []U.Dict{}
	invalidList := []U.Dict{}
	for _, entry := range aggregatedSummary {
		status := strings.ToLower(strings.TrimSpace(getStringFromUDict(entry, "status")))
		switch status {
		case "valid":
			validList = append(validList, entry)
		case "invalid":
			invalidList = append(invalidList, entry)
		default:
			invalidList = append(invalidList, entry)
		}
	}

	meta["summary_of_findings"] = aggregatedSummary
	meta["summary_of_findings_valid"] = validList
	meta["summary_of_findings_invalid"] = invalidList

	documents["meta"] = meta
	fmt.Println(logPrefix, "final meta", meta)
	return 200, documents, nil
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

	// ---- SUMMARY: HITUNG PROJECT TERBARU ----
	latest, ok := GetLatestProject(allDocs)
	var recentTotal, recentValid, recentInvalid int

	if ok {
		// convert single project → multi-project format
		projectList := [][]hctypes.Dict{}

		// ekstrak tabular_result seperti ExtractTabularResults
		attr, _ := latest.GetDict("attributes")
		resultList, _ := attr.GetList("result")

		oneProject := []hctypes.Dict{}
		for i := 0; i < len(resultList); i++ {
			item, _ := resultList.GetDict(i)
			gen, _ := item.GetDict("generated")
			tab, _ := gen.GetDict("tabular_result")
			oneProject = append(oneProject, tab)
		}

		projectList = append(projectList, oneProject)
		projectList = projectList[len(projectList)-1:]

		recentTotal, recentValid, recentInvalid = CalculatePeopleStatus_MultiFiles(projectList)
	}

	// ---- PERSENTASE: PERBANDINGAN BULAN INI VS BULAN LALU ----
	currentMonth := time.Now().Format("2006-01")
	previousMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

	currentTotal, currentValid, currentInvalid := CalculateMonthlyPeopleStatus(grouped, currentMonth)
	prevTotal, prevValid, prevInvalid := CalculateMonthlyPeopleStatus(grouped, previousMonth)

	return 200, U.Dict{
		"chart": chartData,
		"summary": U.Dict{
			"total_people":            recentTotal,
			"valid":                   recentValid,
			"invalid":                 recentInvalid,
			"percentage_total_people": CalculatePercentageChange(prevTotal, currentTotal),
			"percentage_valid":        CalculatePercentageChange(prevValid, currentValid),
			"percentage_invalid":      CalculatePercentageChange(prevInvalid, currentInvalid),
		},
	}, nil
}

func GetLatestProject(allDocs hctypes.List) (hctypes.Dict, bool) {
	if len(allDocs) == 0 {
		return nil, false
	}

	latest := allDocs[0]
	latestTime := getCreatedAtTime(allDocs[0])

	for i := 1; i < len(allDocs); i++ {
		t := getCreatedAtTime(allDocs[i])
		if t.After(latestTime) {
			latest = allDocs[i]
			latestTime = t
		}
	}

	return latest.(hctypes.Dict), true
}

func getCreatedAtTime(doc hctypes.Item) time.Time {
	d, ok := doc.(hctypes.Dict)
	if !ok {
		return time.Time{}
	}

	attr, ok := d.GetDict("attributes")
	if !ok {
		return time.Time{}
	}

	raw, _ := attr.GetString("createdAt")
	t, _ := time.Parse(time.RFC3339, string(raw))
	return t
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

func extractTabularResultsFromAttributes(attr hctypes.Dict) []hctypes.Dict {
	results := []hctypes.Dict{}

	resultList, ok := attr.GetList("result")
	if !ok {
		return results
	}

	for i := 0; i < len(resultList); i++ {
		item, ok := resultList.GetDict(i)
		if !ok {
			continue
		}

		generated, ok := item.GetDict("generated")
		if !ok {
			continue
		}

		tabular, ok := generated.GetDict("tabular_result")
		if !ok {
			continue
		}

		results = append(results, tabular)
	}

	return results
}

type personSummaryBucket struct {
	rowIndex    int
	context     string
	sheet       string
	section     string
	statuses    []string
	suggestions []string
	statusIcon  string
	statusLabel string
}

func buildSummaryEntriesFromTabular(tabulars []hctypes.Dict) []U.Dict {
	buckets := map[int]*personSummaryBucket{}
	order := []int{}

	for _, tab := range tabulars {
		logs, ok := tab.GetList("logs")
		if !ok || len(logs) == 0 {
			continue
		}

		sheet := strings.TrimSpace(getStringFromDict(tab, "sheet"))
		sectionTitle := strings.TrimSpace(getStringFromDict(tab, "section_title"))

		for i := 0; i < len(logs); i++ {
			logEntry, ok := logs.GetDict(i)
			if !ok {
				continue
			}

			rowIndex := extractRowIndexFromLogEntry(logEntry)
			if rowIndex < 0 {
				continue
			}

			suggestion := extractSuggestionFromLog(logEntry, sectionTitle)
			if suggestion == "" {
				continue
			}

			bucket, exists := buckets[rowIndex]
			if !exists {
				bucket = &personSummaryBucket{rowIndex: rowIndex}
				order = append(order, rowIndex)
				buckets[rowIndex] = bucket
			}

			if bucket.context == "" {
				ctx := strings.TrimSpace(extractContextFromLog(logEntry, rowIndex))
				if ctx != "" {
					bucket.context = ctx
				}
			}

			if bucket.sheet == "" && sheet != "" {
				bucket.sheet = sheet
			}
			if bucket.section == "" && sectionTitle != "" {
				bucket.section = sectionTitle
			}

			statusValue := normalizeLogStatus(logEntry)
			bucket.statuses = append(bucket.statuses, statusValue)

			if bucket.statusLabel == "" {
				label := strings.TrimSpace(getStringFromDict(logEntry, "status_label"))
				if label != "" {
					bucket.statusLabel = label
				}
			}

			if bucket.statusIcon == "" {
				icon := strings.TrimSpace(getStringFromDict(logEntry, "status_icon"))
				if icon != "" {
					bucket.statusIcon = icon
				}
			}

			if !containsString(bucket.suggestions, suggestion) {
				bucket.suggestions = append(bucket.suggestions, suggestion)
			}
		}
	}

	entries := []U.Dict{}
	for _, rowIndex := range order {
		bucket := buckets[rowIndex]
		if bucket == nil || len(bucket.suggestions) == 0 {
			continue
		}

		entry := U.Dict{
			"row":     bucket.rowIndex,
			"context": bucket.context,
			"summary": strings.Join(bucket.suggestions, "\n"),
		}

		if entry["context"] == "" {
			entry["context"] = fmt.Sprintf("Row %d", bucket.rowIndex)
		}

		if bucket.sheet != "" {
			entry["sheet"] = bucket.sheet
		}
		if bucket.section != "" {
			entry["section_title"] = bucket.section
		}

		status := classifyStatuses(bucket.statuses)
		entry["status"] = status

		label := bucket.statusLabel
		if label == "" {
			label = prettifyLabel(status)
		}
		if label != "" {
			entry["status_label"] = label
		}

		entry["status_icon"] = getStatusIcon(status, bucket.statusIcon)

		entries = append(entries, entry)
	}

	return entries
}

func extractRowIndexFromLogEntry(logEntry hctypes.Dict) int {
	if rowIdx, ok := logEntry.GetFloat64("row_index"); ok {
		return int(rowIdx)
	}
	if rowNo, ok := logEntry.GetFloat64("no"); ok {
		return int(rowNo)
	}
	return -1
}

func normalizeLogStatus(logEntry hctypes.Dict) string {
	status := strings.ToLower(strings.TrimSpace(getStringFromDict(logEntry, "status")))
	if status != "" {
		return status
	}
	label := strings.ToLower(strings.TrimSpace(getStringFromDict(logEntry, "status_label")))
	if label != "" {
		return label
	}
	return "unknown"
}

func getStatusIcon(status string, fallback string) string {
	switch status {
	case "valid":
		return "✅"
	case "invalid":
		return "❌"
	case "unknown":
		return "⚠️"
	default:
		if fallback != "" {
			return fallback
		}
		return ""
	}
}

func classifyStatuses(statuses []string) string {
	if len(statuses) == 0 {
		return "unknown"
	}

	allValid := true
	hasInvalid := false

	for _, s := range statuses {
		switch s {
		case "valid":
			continue
		case "invalid":
			allValid = false
			hasInvalid = true
		default:
			allValid = false
		}
	}

	if allValid {
		return "valid"
	}
	if hasInvalid {
		return "invalid"
	}
	return "unknown"
}

var summaryContextPriorityKeys = []string{
	"nama lengkap",
	"nama_lengkap",
	"nama",
	"name",
	"context",
	"keterangan",
	"description",
	"deskripsi",
	"perusahaan",
	"company",
	"nik",
	"id",
}

func extractContextFromLog(logEntry hctypes.Dict, fallbackRow int) string {
	if rowData, ok := logEntry.GetDict("row_data"); ok {
		if ctx := pickContextFromRowData(rowData); ctx != "" {
			return ctx
		}
	}

	if section, ok := logEntry.GetDict("section"); ok {
		if value, ok := section.GetString("result"); ok {
			if trimmed := strings.TrimSpace(string(value)); trimmed != "" {
				return trimmed
			}
		}
		if title, ok := section.GetString("title"); ok {
			if trimmed := strings.TrimSpace(string(title)); trimmed != "" {
				return trimmed
			}
		}
	}

	if extraction, ok := logEntry.GetString("ai_extraction"); ok {
		if trimmed := strings.TrimSpace(string(extraction)); trimmed != "" {
			return trimmed
		}
	}

	if fallbackRow > 0 {
		return fmt.Sprintf("Row %d", fallbackRow)
	}

	return ""
}

func pickContextFromRowData(rowData hctypes.Dict) string {
	if len(rowData) == 0 {
		return ""
	}

	keyMapper := map[string]string{}
	firstColumnValue := ""
	firstColumnChecked := false

	for key := range rowData {
		keyMapper[strings.ToLower(key)] = key

		if !firstColumnChecked {
			firstColumnChecked = true
			firstColumnValue = extractRowDataValue(rowData, key)
		}
	}

	for _, wanted := range summaryContextPriorityKeys {
		if actual, ok := keyMapper[wanted]; ok {
			if val, ok := rowData.GetString(actual); ok {
				if trimmed := strings.TrimSpace(string(val)); trimmed != "" {
					return trimmed
				}
			}
		}
	}

	if strings.TrimSpace(firstColumnValue) != "" {
		return strings.TrimSpace(firstColumnValue)
	}

	for _, actual := range keyMapper {
		if val, ok := rowData.GetString(actual); ok {
			if trimmed := strings.TrimSpace(string(val)); trimmed != "" {
				return trimmed
			}
		}

		if num, ok := rowData.GetFloat64(actual); ok {
			return strconv.Itoa(int(num))
		}
	}

	return ""
}

func extractRowDataValue(rowData hctypes.Dict, key string) string {
	if val, ok := rowData.GetString(key); ok {
		return strings.TrimSpace(string(val))
	}

	if num, ok := rowData.GetFloat64(key); ok {
		return strings.TrimSpace(strconv.FormatFloat(float64(num), 'f', -1, 64))
	}

	if raw, exists := rowData[key]; exists && raw != nil {
		switch value := raw.(type) {
		case string:
			return strings.TrimSpace(value)
		case float64:
			return strings.TrimSpace(strconv.FormatFloat(value, 'f', -1, 64))
		case int:
			return strconv.Itoa(value)
		case int64:
			return strconv.FormatInt(value, 10)
		default:
			if str, ok := value.(fmt.Stringer); ok {
				return strings.TrimSpace(str.String())
			}
			return strings.TrimSpace(fmt.Sprintf("%v", value))
		}
	}

	return ""
}

func extractSuggestionFromLog(logEntry hctypes.Dict, sectionTitle string) string {
	if suggestion, ok := logEntry.GetString("ai_suggestion"); ok {
		if trimmed := strings.TrimSpace(string(suggestion)); trimmed != "" {
			return trimmed
		}
	}

	if statusLabel, ok := logEntry.GetString("status_label"); ok {
		label := prettifyLabel(string(statusLabel))
		if label != "" {
			if sectionTitle != "" {
				return fmt.Sprintf("%s - %s", sectionTitle, label)
			}
			return label
		}
	}

	if status, ok := logEntry.GetString("status"); ok {
		status = hctypes.String(strings.TrimSpace(string(status)))
		if status != "" {
			label := prettifyLabel(string(status))
			if sectionTitle != "" && label != "" {
				return fmt.Sprintf("%s - %s", sectionTitle, label)
			}
			return label
		}
	}

	return ""
}

func prettifyLabel(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	input = strings.ReplaceAll(input, "_", " ")
	input = strings.ReplaceAll(input, "-", " ")

	parts := strings.Fields(input)
	for i, part := range parts {
		lower := strings.ToLower(part)
		if len(lower) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(lower[:1]) + lower[1:]
	}

	return strings.Join(parts, " ")
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func getStringFromDict(dict hctypes.Dict, key string) string {
	if val, ok := dict.GetString(key); ok {
		return string(val)
	}

	if raw, exists := dict[key]; exists {
		switch value := raw.(type) {
		case string:
			return value
		case fmt.Stringer:
			return value.String()
		case float64:
			return strconv.Itoa(int(value))
		case int:
			return strconv.Itoa(value)
		default:
			if value == nil {
				return ""
			}
			return fmt.Sprintf("%v", value)
		}
	}

	return ""
}

func getStringFromUDict(dict U.Dict, key string) string {
	if raw, exists := dict[key]; exists {
		switch value := raw.(type) {
		case string:
			return value
		case fmt.Stringer:
			return value.String()
		case float64:
			return strconv.Itoa(int(value))
		case int:
			return strconv.Itoa(value)
		default:
			if value == nil {
				return ""
			}
			return fmt.Sprintf("%v", value)
		}
	}
	return ""
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
		if !ok {
			fmt.Println("[DEBUG] Skip: tidak ada project_id")
			continue
		}

		docID, ok := attr.GetFloat64("project_document_id")
		if !ok {
			fmt.Println("[DEBUG] Skip: project", pid, "tidak punya project_document_id")
			continue
		}

		statusDoc, docData, err := Strapi.Find("project-documents-table", U.Dict{
			"filters[id][$eq]": int(docID),
			"populate":         "*",
		})

		fmt.Printf("[DEBUG] Cek project_document_id %v → status %v err %v\n", docID, statusDoc, err)

		if err != nil || statusDoc != 200 {
			fmt.Println("[DEBUG] Skip: gagal fetch document", docID)
			continue
		}

		list, ok := docData.GetList("data")
		if !ok || len(list) == 0 {
			fmt.Println("[DEBUG] Skip: data dokumen kosong")
			continue
		}

		docAttr, _ := list[0].(hctypes.Dict).GetDict("attributes")

		mediaParent, ok := docAttr.GetDict("document_media_id")
		if !ok {
			fmt.Println("[DEBUG] Skip: document_media_id tidak ada")
			continue
		}

		mediaData, ok := mediaParent.GetDict("data")
		if !ok {
			fmt.Println("[DEBUG] Skip: document_media_id.data tidak ada")
			continue
		}

		mediaAttr, ok := mediaData.GetDict("attributes")
		if !ok {
			fmt.Println("[DEBUG] Skip: document_media_id.data.attributes tidak ada")
			continue
		}

		filename, _ := mediaAttr.GetString("name")
		ext, _ := mediaAttr.GetString("ext")

		fmt.Println("[DEBUG] Filename:", filename, "EXT:", ext)

		if strings.ToLower(string(ext)) != ".csv" {
			fmt.Println("[DEBUG] Skip: file bukan CSV →", filename)
			continue
		}

		projectIDSet[int(pid)] = true
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

	params := U.Dict{
		"populate": "*",
		"pagination": U.Dict{
			"page":     page,
			"pageSize": pageSize,
		},
	}

	for i, pid := range projectIDs {
		key := fmt.Sprintf("filters[id][$in][%d]", i)
		params[key] = pid
	}

	status, projects, err := Strapi.Find("projects-table", params)

	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}
	return 200, projects, nil
}

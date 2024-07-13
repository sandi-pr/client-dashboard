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

	"gitlab.hexacode.org/go-libs/cr3call"
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

	upload_by, ok := me.GetFloat64("id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	// document_title, ok := ctx.Params.GetString("document_title")
	// if !ok {
	// 	return ctx.Error(400, "invalid parameters document_title")
	// }

	// description, ok := ctx.Params.GetString("description")
	// if !ok {
	// 	return ctx.Error(400, "invalid parameters description")
	// }

	// document_type, ok := ctx.Params.GetString("document_type")
	// if !ok {
	// 	return ctx.Error(400, "invalid parameters document_type")
	// }

	// documentBase64, ok := ctx.Params.GetString("document")
	// if !ok {
	// 	return ctx.Error(400, "invalid parameters document")
	// }

	// document_name, ok := ctx.Params.GetString("document_name")
	// if !ok {
	// 	return ctx.Error(400, "invalid parameters document_name")
	// }

	// person, ok := ctx.Params.GetFloat64("person")
	// if !ok {
	// 	fmt.Println("not found person")
	// }

	// date_start, ok := ctx.Params.GetString("contract_start")
	// if !ok {
	// 	fmt.Println("invalid parameters date_start")
	// }

	// date_end, ok := ctx.Params.GetString("contract_end")
	// if !ok {
	// 	fmt.Println("invalid parameters date_end")
	// }

	documentsString, ok := ctx.Params.GetString("documents")
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	document_type, ok := ctx.Params.GetString("document_type")
	if !ok {
		return ctx.Error(400, "invalid parameters document_type")
	}

	collection_idStr, ok := ctx.Params.GetString("collection_id")
	if !ok {
		fmt.Println("invalid parameters collection_id")
	}

	collection_id, err := strconv.Atoi(string(collection_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters collection_id")
	}

	var documents interface{}
	err = json.Unmarshal([]byte(documentsString), &documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

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

		status, document, err := Strapi.Create("document-management-table", U.Dict{
			"document_title": documentName,
			"description":    "-",
			"document_type":  document_type,
			"document":       id,
			"collection_id":  collection_id,
			"upload_by":      upload_by,
			"file_size":      file_size,
		})

		fmt.Println("document", document)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			return ctx.Error(500, "Internal Server Error")
		}

		CR3 := cr3call.New(MODULE_MANAGER_URL)
		go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
			Headers: U.Dict{
				"Authorization": bearer,
			},
		}, false)
	}

	return 200, U.Dict{}, nil

	// document, err := base64.StdEncoding.DecodeString(string(documentBase64))
	// if err != nil {
	// 	return ctx.Error(400, "failed to decode document")
	// }

	// status, upload, err := Strapi.Upload(string(document_name), document)

	// fmt.Println("upload", upload)

	// if err != nil {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// if status != 200 {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// data, ok := upload.GetDict(0)

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// id, ok := data.GetFloat64("id")

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// if document_type == "GENERAL" {

	// 	status, documents, err := Strapi.Create("document-management-table", U.Dict{
	// 		"document_title": document_title,
	// 		"description":    description,
	// 		"document_type":  document_type,
	// 		"document":       id,
	// 		"collection_id":  collection_id,
	// 	})

	// 	fmt.Println("document", documents)

	// 	if err != nil {
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	if status != 200 {
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	// CR3 := cr3call.New(MODULE_MANAGER_URL)
	// 	// go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
	// 	// 	Headers: U.Dict{
	// 	// 		"Authorization": bearer,
	// 	// 	},
	// 	// }, false)

	// 	return status, documents, nil
	// } else {
	// 	status, documents, err := Strapi.Create("document-management-table", U.Dict{
	// 		"document_title": document_title,
	// 		"description":    description,
	// 		"document_type":  document_type,
	// 		"document":       id,
	// 		"person":         person,
	// 		"contract_start": date_start,
	// 		"contract_end":   date_end,
	// 		"collection_id":  collection_id,
	// 	})

	// 	fmt.Println("document", documents)

	// 	if err != nil {
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	if status != 200 {
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	return status, documents, nil
	// }
}

func GetCollectionById(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, collection, err := Strapi.FindOne("document-collection-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, collection, nil
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

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
	}, false)

	return 200, documents, nil
}

func GetDocumentsByCollectionId(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	collection_id, ok := ctx.Params.GetString("collection_id")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_id")
	}

	id, err := strconv.Atoi(string(collection_id))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, documents, err := Strapi.Find("document-management-table", U.Dict{
		"filters": U.Dict{
			"collection_id": U.Dict{
				"$eq": id,
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
		fmt.Println(documents)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
	}, false)

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

	document_url := dict["attributes"].(U.Dict)["document"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
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

func UpdateDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	detected_text, ok := ctx.Params.GetString("detected_text")
	if !ok {
		return ctx.Error(400, "invalid parameters detected_text")
	}

	status, document, err := Strapi.Update("document-management-table", int64(id), U.Dict{
		"detected_text": detected_text,
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

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "refresh-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
	}, false)

	return 200, U.Dict{}, nil
}

func GetCountDocumentsPerMonth(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	earlyDate, ok := ctx.Params.GetString("earlyDate")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	lastDate, ok := ctx.Params.GetString("lastDate")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, document, err := Strapi.Find("document-management-table", U.Dict{
		"filters": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	fmt.Println(document)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, general, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("general :", general)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, contract, err := Strapi.Find("document-collection-table", U.Dict{
		"filters[$and][0]": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
		"filters[$and][1]": U.Dict{
			"status": U.Dict{
				"$eq": "Ready",
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("contract :", contract)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{"document": document, "general": general, "contract": contract}, nil
}

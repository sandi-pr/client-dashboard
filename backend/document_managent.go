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

	category, ok := ctx.Params.GetString("category")
	if !ok {
		category = ""
	}

	var documents interface{}
	err = json.Unmarshal([]byte(documentsString), &documents)

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

		document_title, ok := document.(map[string]interface{})["document_title"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		description, ok := document.(map[string]interface{})["description"].(string)
		if !ok {
			description = ""
		}

		contract, ok := document.(map[string]interface{})["contract"].(string)
		if !ok {
			contract = ""
		}

		id_contract, ok := document.(map[string]interface{})["id"].(float64)
		if !ok {
			id_contract = 0
		}

		var reference_doc_id float64 = 0

		var project_doc_id float64 = 0

		if contract == "REFERENCE" {

			reference_doc_id = id_contract

		}

		if contract == "PROJECT" {

			project_doc_id = id_contract

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

		dataSend := U.Dict{
			"document_title":   document_title,
			"document_type":    document_type,
			"document":         id,
			"collection_id":    collection_id,
			"upload_by":        upload_by,
			"file_size":        file_size,
			"edited":           false,
			"inserted":         false,
			"reference_doc_id": reference_doc_id,
			"project_doc_id":   project_doc_id,
		}

		if category != "" {
			dataSend["category"] = category
		}

		if description != "" {
			dataSend["description"] = description
		}

		status, documents, err := Strapi.Create("document-management-table", dataSend)

		fmt.Println("document", documents)

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

		result, ok := documents.GetDict("data")
		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		datas = append(datas, result)
	}

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai-rag", "add-document-by-collection-id", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
		Params: U.Dict{
			"collection_id": collection_idStr,
		},
	}, false)

	return 200, U.Dict{"data": datas}, nil

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

func GetCollectionByName(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	collection_name, ok := ctx.Params.GetString("collection_name")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_name")
	}

	status, collection, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"collection_name": U.Dict{
				"$eq": collection_name,
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

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
	}, false)

	return 200, collection, nil
}

func GetCollectionByNamePublic(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(STRAPI_TOKEN), "Bearer ", "", 1))

	// status, me, err := Strapi.UsersMe()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// if status != 200 {
	// 	fmt.Println(status)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// fmt.Println(me)

	collection_name, ok := ctx.Params.GetString("collection_name")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_name")
	}

	status, collection, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"collection_name": U.Dict{
				"$eq": collection_name,
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

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": STRAPI_TOKEN,
		},
	}, false)

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

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

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
			"document_title": U.Dict{
				"$contains": search,
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

	list, ok := documents.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return 200, documents, nil
	}

	page2 := 1

	var allDoc U.List

	for {
		status, documents2, err := Strapi.Find("document-management-table", U.Dict{
			"filters": U.Dict{
				"collection_id": U.Dict{
					"$eq": id,
				},
			},
			"pagination": U.Dict{
				"page":     page2,
				"pageSize": 25,
			},
			"sort": "id:desc",
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := documents2.GetList("data")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) == 0 {
			break
		}

		allDoc = append(allDoc, list...)

		page2++
	}

	status_rag := "Ready"

	for i := 0; i < len(allDoc); i++ {
		attributes, ok := allDoc[i].(U.Dict).GetDict("attributes")
		if ok {
			status, ok := attributes.GetString("status")
			if ok {
				if status == "Reading" {
					status_rag = "Reading"
				}
			} else {
				status_rag = ""
				break
			}
		}
	}

	if status_rag != "" {
		documents["status_rag"] = status_rag
	}

	return 200, documents, nil
}

func GetDocumentsByCollectionIdPublic(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(STRAPI_TOKEN), "Bearer ", "", 1))

	// status, me, err := Strapi.UsersMe()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// if status != 200 {
	// 	fmt.Println(status)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// fmt.Println(me)

	// page, ok := ctx.Params.GetFloat64("page")
	// if !ok {
	// 	fmt.Println("invalid parameters")
	// 	page = 1
	// }

	// page_size, ok := ctx.Params.GetFloat64("page_size")
	// if !ok {
	// 	page_size = 10
	// }

	collection_id, ok := ctx.Params.GetString("collection_id")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_id")
	}

	id, err := strconv.Atoi(string(collection_id))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	page := 1

	allDocuments := U.List{}

	for {
		status, documents, err := Strapi.Find("document-management-table", U.Dict{
			"filters": U.Dict{
				"collection_id": U.Dict{
					"$eq": id,
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"sort": "id:desc",
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

		if list, ok := documents.GetList("data"); ok {
			if len(list) == 0 {
				break
			} else {
				allDocuments = append(allDocuments, list...)
			}
		}
		page++
	}

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": STRAPI_TOKEN,
		},
	}, false)

	statusRAG := true

	for i := 0; i < len(allDocuments); i++ {
		if data, ok := allDocuments.GetDict(i); ok {
			if attributes, ok := data.GetDict("attributes"); ok {
				status, ok := attributes.GetString("status")

				if !ok {
					statusRAG = false
					break
				}

				if status != "Ready" {
					statusRAG = false
					break
				}
			}
		}
	}

	return 200, U.Dict{"status": statusRAG}, nil
}

func GetDocumentsByCollectionIdWithFile(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	list, ok := documents.GetList("data")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	listDoc := U.List{}

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

		fmt.Println(base64Data)

		listDoc = append(listDoc, U.Dict{
			"id":     list[i].(U.Dict)["id"].(U.Float64),
			"name":   list[i].(U.Dict)["attributes"].(U.Dict)["document_title"].(U.String),
			"type":   list[i].(U.Dict)["attributes"].(U.Dict)["document_type"].(U.String),
			"size":   list[i].(U.Dict)["attributes"].(U.Dict)["file_size"].(U.String),
			"base64": base64Data,
			"ext":    list[i].(U.Dict)["attributes"].(U.Dict)["document"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["ext"].(U.String),
		})
	}

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
	}, false)

	return 200, U.Dict{"data": listDoc}, nil
}

func GetDocumentByCategories(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	collection_idStr, ok := ctx.Params.GetString("collection_id")

	if !ok {
		return ctx.Error(400, "invalid parameters collection")
	}

	collection_id, err := strconv.Atoi(string(collection_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters collection")
	}

	status, collection, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": collection_id,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := collection.GetList("data")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(data) == 0 {
		fmt.Println("collection not found")
		return 500, U.Dict{}, nil
	}

	dict, ok := data.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	allCategories, ok := attributes.GetList("categories")
	if !ok {
		return 200, U.Dict{"data": []U.Dict{}}, nil
	}

	if len(allCategories) == 0 {
		return 200, U.Dict{"data": []U.Dict{}}, nil
	}

	allDocumentByCategories := U.List{}

	for i := 0; i < len(allCategories); i++ {
		status, documents, err := Strapi.Find("document-management-table", U.Dict{
			"filters": U.Dict{
				"category": U.Dict{
					"$eq": allCategories[i],
				},
				"collection_id": U.Dict{
					"$eq": collection_id,
				},
			},
			"pagination": U.Dict{
				"page":     1,
				"pageSize": 10,
			},
			"sort":     "id:desc",
			"populate": "*",
		})

		fmt.Println(allCategories[i], documents)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(documents)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		meta, ok := documents.GetDict("meta")

		if !ok {
			fmt.Println("meta not found")
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := documents.GetList("data")
		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		allDocumentByCategories = append(allDocumentByCategories, U.Dict{
			"category": allCategories[i],
			"data":     list,
			"meta":     meta,
		})
	}

	return 200, U.Dict{"data": allDocumentByCategories}, nil
}

func GetDocumentByCategory(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	collection_idStr, ok := ctx.Params.GetString("collection_id")

	if !ok {
		return ctx.Error(400, "invalid parameters collection")
	}

	collection_id, err := strconv.Atoi(string(collection_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters collection")
	}

	category, ok := ctx.Params.GetString("category")

	if !ok {
		return ctx.Error(400, "invalid parameters category")
	}

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		fmt.Println("invalid parameters")
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		fmt.Println("invalid parameters")
		page_size = 10
	}

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

	status, documents, err := Strapi.Find("document-management-table", U.Dict{
		"filters": U.Dict{
			"category": U.Dict{
				"$eq": category,
			},
			"collection_id": U.Dict{
				"$eq": int64(collection_id),
			},
			"document_title": U.Dict{
				"$contains": search,
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

func UpdateDocumentHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	result_trocr, ok := ctx.Params.GetString("result_trocr")
	if !ok {
		return ctx.Error(400, "invalid parameters result_trocr")
	}

	edited, ok := ctx.Params.GetBool("edited")
	if !ok {
		return ctx.Error(400, "invalid parameters edited")
	}

	var result_trocr_json interface{}
	err = json.Unmarshal([]byte(result_trocr), &result_trocr_json)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	document_title, ok := ctx.Params.GetString("document_title")
	if !ok {
		return ctx.Error(400, "invalid parameters document_title")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	category, ok := ctx.Params.GetString("category")
	if !ok {
		category = ""
	}

	dataSend := U.Dict{
		"detected_text":  detected_text,
		"result_trocr":   result_trocr_json,
		"edited":         edited,
		"document_title": document_title,
	}

	if category != "" {
		dataSend["category"] = category
	}

	if description != "" {
		dataSend["description"] = description
	}

	status, document, err := Strapi.Update("document-management-table", int64(id), dataSend)
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

	document_title, ok := ctx.Params.GetString("document_title")
	if !ok {
		return ctx.Error(400, "invalid parameters document_title")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	category, ok := ctx.Params.GetString("category")
	if !ok {
		category = ""
	}

	dataSend := U.Dict{
		"document_title": document_title,
	}

	if category != "" {
		dataSend["category"] = category
	}

	if description != "" {
		dataSend["description"] = description
	}

	status, document, err := Strapi.Update("document-management-table", int64(id), dataSend)
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

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	resp, err := CR3.Call("hexa-ai-rag", "delete-document-by-id", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
		Params: U.Dict{
			"document_id": idStr,
		},
	}, false)

	fmt.Println(resp)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
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

	data, _ := document.GetList("data")

	collectionStatusMap := make(map[int]bool)

	for i := 0; i < len(data); i++ {
		attributes, ok := data[i].(U.Dict).GetDict("attributes")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		collectionID, ok := attributes.GetFloat64("collection_id")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		status, _ := attributes.GetString("status")

		fmt.Println(collectionID)

		if _, exists := collectionStatusMap[int(collectionID)]; !exists {
			collectionStatusMap[int(collectionID)] = true
		}

		// Set the map entry to false if any document is not "Ready"
		if status != "Ready" {
			collectionStatusMap[int(collectionID)] = false
		}
	}
	// Collect all collection IDs that are fully "Ready"
	readyCollectionIDs := []int{}
	for collectionID, allReady := range collectionStatusMap {
		if allReady {
			readyCollectionIDs = append(readyCollectionIDs, collectionID)
		}
	}

	fmt.Println("Collection IDs with all documents 'Ready':", readyCollectionIDs)

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
		"filters": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
			"id": U.Dict{
				"$in": readyCollectionIDs,
			},
		},
	})

	fmt.Println("contract :", contract)

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

func DeleteDocumentCollectionByContract(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	contract, ok := ctx.Params.GetString("contract")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	fmt.Println(contract)

	contract_doc_idStr, ok := ctx.Params.GetString("contract_doc_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	contract_doc_id, err := strconv.Atoi(string(contract_doc_idStr))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters")
	}

	var document U.Dict

	if contract == "REFERENCE" {
		fmt.Println("ayam")
		status, document, err = Strapi.Find("document-management-table", U.Dict{
			"filters": U.Dict{
				"reference_doc_id": U.Dict{
					"$eq": int64(contract_doc_id),
				},
			},
		})
	} else {
		fmt.Println("bebek")
		status, document, err = Strapi.Find("document-management-table", U.Dict{
			"filters": U.Dict{
				"project_doc_id": U.Dict{
					"$eq": int64(contract_doc_id),
				},
			},
		})

	}

	if err != nil {
		fmt.Println(document)
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := document.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(data) == 0 {
		return 200, document, nil
	}

	dict, ok := data.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	status, _, err = Strapi.Delete("document-management-table", int64(id))

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func DeleteDocumentByCategories(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	category, ok := ctx.Params.GetString("category")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	collection_idStr, ok := ctx.Params.GetString("collection_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	collection_id, err := strconv.Atoi(string(collection_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	page := 1
	var ids []int64

	for {
		status, document, err := Strapi.Find("document-management-table", U.Dict{
			"filters": U.Dict{
				"category": U.Dict{
					"$eq": category,
				},
				"collection_id": U.Dict{
					"$eq": int64(collection_id),
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		data, ok := document.GetList("data")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(data) == 0 {
			break
		}

		for _, dict := range data {
			id, ok := dict.(U.Dict).GetFloat64("id")
			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}
			ids = append(ids, int64(id))
		}
		page++
	}

	for _, id := range ids {
		status, _, err := Strapi.Delete("document-management-table", int64(id))
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}
		if status != 200 {
			return ctx.Error(500, "Internal Server Error")
		}
	}

	status, collection, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": collection_id,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := collection.GetList("data")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(data) == 0 {
		fmt.Println("collection not found")
		return 500, U.Dict{}, nil
	}

	dict, ok := data.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	categories, ok := attributes.GetList("categories")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(categories) == 0 {
		return ctx.Error(500, "Internal Server Error")
	}

	var newCategories U.List

	for i := 0; i < len(categories); i++ {
		if categories[i] != category {
			newCategories = append(newCategories, categories[i])
		}
	}

	status, updated, err := Strapi.Update("document-collection-table", int64(collection_id), U.Dict{
		"categories": newCategories,
	})

	fmt.Println("updated", updated)

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

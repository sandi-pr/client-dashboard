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

func CreateDocumentCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	// documentsString, ok := ctx.Params.GetString("documents")
	// if !ok {
	// 	return ctx.Error(400, "invalid parameters documents")
	// }

	status, documentCollection, err := Strapi.Create("document-collection-table", U.Dict{
		"collection_name": collection_name,
		"description":     description,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	// dataCollection, ok := documentCollection.GetDict("data")
	// if !ok {
	// 	fmt.Println("failed to get data")
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// idCollection, ok := dataCollection.GetFloat64("id")
	// if !ok {
	// 	fmt.Println("failed to get id")
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// if documentsString == "-" {
	// 	return 200, documentCollection, nil
	// }

	// var documents interface{}
	// err = json.Unmarshal([]byte(documentsString), &documents)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(400, "invalid parameters documents")
	// }

	// for _, document := range documents.([]interface{}) {

	// 	documentBase64, ok := document.(map[string]interface{})["base64"].(string)
	// 	if !ok {
	// 		return ctx.Error(400, "invalid parameters documents")
	// 	}

	// 	documentFile, err := base64.StdEncoding.DecodeString(documentBase64)

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return ctx.Error(400, "invalid parameters documents")
	// 	}

	// 	documentName, ok := document.(map[string]interface{})["name"].(string)
	// 	if !ok {
	// 		return ctx.Error(400, "invalid parameters documents")
	// 	}

	// 	status, upload, err := Strapi.Upload(documentName, documentFile)

	// 	fmt.Println("upload", upload)

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	if status != 200 {
	// 		fmt.Println(status)
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	data, ok := upload.GetDict(0)

	// 	if !ok {
	// 		fmt.Println("failed to get dict")
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	id, ok := data.GetFloat64("id")
	// 	if !ok {
	// 		fmt.Println("failed to get id")
	// 		return ctx.Error(500, "Internal Server Error")
	// 	}

	// 	status, document, err := Strapi.Create("document-management-table", U.Dict{
	// 		"document_title": documentName,
	// 		"description":    "-",
	// 		"document_type":  "GENERAL",
	// 		"document":       id,
	// 		"collection_id":  idCollection,
	// 	})

	// 	fmt.Println("document", document)

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
	// }

	return 200, documentCollection, nil
}

func GetDocumentCollections(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, documentCollections, err := Strapi.Find("document-collection-table", U.Dict{
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

	return status, documentCollections, nil
}

func GetDocumentCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": id,
			},
		},
		"populate": "*",
	})

	fmt.Println("document", document)

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
		fmt.Print("ayam 1")
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		fmt.Print("ayam 2")
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		fmt.Print("ayam 3")
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")
	if !ok {
		fmt.Print("ayam 4")
		return ctx.Error(500, "Internal Server Error")
	}

	documents, ok := attributes.GetDict("documents")
	if !ok {
		fmt.Print("ayam 5")
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := documents.GetList("data")
	if !ok {
		fmt.Print("ayam 6")
		return ctx.Error(500, "Internal Server Error")
	}

	var documentsBase64 []string

	for i := 0; i < len(data); i++ {
		document_url, ok := data[i].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		if !ok {
			fmt.Print("ayam 7")
			return ctx.Error(500, "Internal Server Error")
		}

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

		documentsBase64 = append(documentsBase64, base64Data)
	}

	return 200, U.Dict{"document": document, "base64": documentsBase64}, nil
}

func UpdateDocumentCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	collection_name, ok := ctx.Params.GetString("collection_name")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	status, document, err := Strapi.Update("document-collection-table", int64(id), U.Dict{
		"collection_name": collection_name,
		"description":     description,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func DeleteDocumentCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	var page = 1

	var idsDoc U.List

	for {
		status, documents, err := Strapi.Find("document-management-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
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

		list, ok := documents.GetList("data")

		if !ok {
			fmt.Println(ok)
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) < 1 {
			break
		}

		for i := 0; i < len(list); i++ {
			if dict, ok := list.GetDict(i); ok {
				if idDoc, ok := dict.GetFloat64("id"); ok {
					idsDoc = append(idsDoc, idDoc)
				}
			}
		}

		page++
	}

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	resp, err := CR3.Call("hexa-ai-rag", "delete-document-by-collection-id", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
		Params: U.Dict{
			"collection_id": idStr,
		},
	}, false)

	fmt.Println(resp)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	for i := 0; i < len(idsDoc); i++ {
		if idDoc, ok := idsDoc.GetFloat64(i); ok {
			status, document, err := Strapi.Delete("document-management-table", int64(idDoc))
			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println(document)
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}
		}
	}

	var page2 = 1
	var idsBot U.List

	for {
		status, data, err := Strapi.Find("hexabot-table", U.Dict{
			"filters": U.Dict{
				"collection_id": U.Dict{
					"$eq": id,
				},
			},
			"pagination": U.Dict{
				"page":     page2,
				"pageSize": 25,
			},
			"populate": "*",
			"sort":     "updatedAt:desc",
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := data.GetList("data")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) < 1 {
			break
		}

		for i := 0; i < len(list); i++ {
			if dict, ok := list.GetDict(i); ok {
				if idBot, ok := dict.GetFloat64("id"); ok {
					idsBot = append(idsBot, idBot)
				}
			}
		}

		page2++
	}

	for i := 0; i < len(idsBot); i++ {
		if idDoc, ok := idsBot.GetFloat64(i); ok {
			status, document, err := Strapi.Delete("hexabot-table", int64(idDoc))
			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println(document)
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}
		}
	}

	status, document, err := Strapi.Delete("document-collection-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func UpdateCollectionWebchat(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, allDocument, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"webchat_tmp": U.Dict{
				"$eq": true,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("ayam 1")
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println("ayam 2")
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println("allDocument", allDocument)

	list, ok := allDocument.GetList("data")

	if !ok {
		fmt.Println("ayam 3")
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) > 0 {
		if dict, ok := list.GetDict(0); ok {
			if id, ok := dict.GetFloat64("id"); ok {
				status, document, err := Strapi.Update("document-collection-table", int64(id), U.Dict{
					"webchat_tmp": false,
				})

				if err != nil {
					fmt.Println(err)
					fmt.Println("ayam 4")
					return ctx.Error(500, "Internal Server Error")
				}

				if status != 200 {
					fmt.Println(status)
					fmt.Println("ayam 5")
					return ctx.Error(500, "Internal Server Error")
				}

				fmt.Println(document)
			}
		}

	}

	collection_id, ok := ctx.Params.GetFloat64("collection_id")
	if !ok {
		fmt.Println("ayam 6")
		return ctx.Error(400, "invalid parameters collection_id")
	}

	status, document, err := Strapi.Update("document-collection-table", int64(collection_id), U.Dict{
		"webchat_tmp": true,
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("ayam 7")
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println("ayam 8")
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func GetCollectionWebchatPublic(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	status, allDocument, err := Strapi.Find("document-collection-table", U.Dict{
		"filters": U.Dict{
			"webchat_tmp": U.Dict{
				"$eq": true,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("ayam 1")
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println("ayam 2")
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println("allDocument", allDocument)

	list, ok := allDocument.GetList("data")

	if !ok {
		fmt.Println("ayam 3")
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) > 0 {
		if dict, ok := list.GetDict(0); ok {
			if id, ok := dict.GetFloat64("id"); ok {
				return 200, U.Dict{"id": id}, nil
			}
		}

	}

	return 404, U.Dict{}, nil
}

func GetCategoriesByCollectionId(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters")
	}

	collection_id, err := strconv.Atoi(string(collection_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
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
		return 200, U.Dict{"categories": []U.Dict{}}, nil
	}

	if len(categories) == 0 {
		return 200, U.Dict{"categories": []U.Dict{}}, nil
	}

	return 200, U.Dict{"categories": categories}, nil
}

func UpdateCategories(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters")
	}

	collection_id, err := strconv.Atoi(string(collection_idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	categoriesStr, ok := ctx.Params.GetString("categories")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var categories interface{}
	err = json.Unmarshal([]byte(categoriesStr), &categories)

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, _, err = Strapi.Update("document-collection-table", int64(collection_id), U.Dict{
		"categories": categories,
	})

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

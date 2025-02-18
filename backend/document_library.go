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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab.hexacode.org/go-libs/cr3call"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateLibrary(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	author, ok := me.GetString("username")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	library_name, ok := ctx.Params.GetString("library_name")
	if !ok {
		return ctx.Error(400, "invalid parameters library_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	document_path := DOCUMENT_PATH + "/" + string(library_name)

	dir := filepath.Dir(document_path)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Gagal membuat direktori:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	dataSend := U.Dict{
		"type":             "library",
		"name":             library_name,
		"path":             document_path,
		"library_industry": false,
		"system_item":      false,
		"meta": U.Dict{
			"upload_by": upload_by,
			"author":    author,
		},
	}

	if description != "" {
		dataSend["meta"].(U.Dict)["description"] = description
	}

	status, documentLibrary, err := Strapi.Create("document-library-table", dataSend)

	fmt.Println("result", documentLibrary)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)

		if err2, ok := documentLibrary.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, documentLibrary, nil
}

func CreateLibraryCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	author, ok := me.GetString("username")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	name, ok := ctx.Params.GetString("collection_name")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "invalid parameters parent_id")
	}

	status, library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": parent_id,
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

	var document_path string

	if list, ok := library.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			if attributes, ok := dict.GetDict("attributes"); ok {
				if path, ok := attributes.GetString("path"); ok {
					document_path = string(path) + "/" + string(name)
				}
			}
		}
	}

	dir := filepath.Dir(document_path)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Gagal membuat direktori:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	dataSend := U.Dict{
		"type":      "collection",
		"name":      name,
		"path":      document_path,
		"parent_id": parent_id,
		"meta": U.Dict{
			"upload_by": upload_by,
			"author":    author,
		},
	}

	if description != "" {
		dataSend["meta"].(U.Dict)["description"] = description
	}

	status, documentLibrary, err := Strapi.Create("document-library-table", dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)

		if err2, ok := documentLibrary.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	if bot_id, ok := ctx.Params.GetString("bot_id"); ok {
		fmt.Println("Received bot_id:", bot_id)

		dataCollection, ok := documentLibrary.GetDict("data")
		if !ok {
			fmt.Println("Error: Failed to get data collection")
			return ctx.Error(500, "Internal Server Error")
		}

		collection_id, ok := dataCollection.GetFloat64("id")
		if !ok {
			fmt.Println("Error: Failed to get collection_id from data collection")
			return ctx.Error(500, "Internal Server Error")
		}

		collection_idStr := strconv.Itoa(int(collection_id))

		botId, _ := strconv.Atoi(string(bot_id))
		fmt.Println("Parsed bot_id:", botId)

		status, bot, err := Strapi.Find("hexabot-table", U.Dict{
			"filters": U.Dict{
				"id": U.Dict{
					"$eq": botId,
				},
			},
		})

		if err != nil {
			fmt.Println("Error finding bot:", err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Printf("Unexpected status code: %d while finding bot\n", status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := bot.GetList("data")
		if !ok {
			fmt.Println("Error: Failed to get 'data' list from bot")
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) == 0 {
			fmt.Println("Error: No data found in the list")
			return ctx.Error(400, "invalid parameters")
		}

		dict, ok := list.GetDict(0)
		if !ok {
			fmt.Println("Error: Failed to get dictionary from list")
			return ctx.Error(500, "Internal Server Error")
		}

		attributes, ok := dict.GetDict("attributes")
		if !ok {
			fmt.Println("Error: Failed to get 'attributes' from dictionary")
			return ctx.Error(500, "Internal Server Error")
		}

		collection_ids, ok := attributes.GetString("collection_id")
		if !ok {
			fmt.Println("Error: Failed to get 'collection_id' from attributes")
			return ctx.Error(500, "Internal Server Error")
		}

		fmt.Println("Current collection_ids:", collection_ids)

		newCollectionIds := string(collection_ids) + "," + collection_idStr
		fmt.Println("Updated collection_ids:", newCollectionIds)

		status, _, err = Strapi.Update("hexabot-table", int64(botId), U.Dict{
			"collection_id": newCollectionIds,
		})

		if err != nil {
			fmt.Println("Error updating bot:", err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Printf("Unexpected status code: %d while updating bot\n", status)
			return ctx.Error(500, "Internal Server Error")
		}

		fmt.Println("Bot updated successfully")
	}

	return 200, documentLibrary, nil
}

func CreateLibraryDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	author, ok := me.GetString("username")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	name, ok := ctx.Params.GetString("document_name")
	if !ok {
		return ctx.Error(400, "invalid parameters document_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "invalid parameters parent_id")
	}

	document_str, ok := ctx.Params.GetString("documents")
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	document_type, ok := ctx.Params.GetString("document_type")
	if !ok {
		return ctx.Error(400, "invalid parameters document_type")
	}

	category, ok := ctx.Params.GetString("category")
	if !ok {
		category = ""
	}

	var documents interface{}
	err = json.Unmarshal([]byte(document_str), &documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

	status, library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": parent_id,
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

	if len(library) == 0 {
		return ctx.Error(400, "invalid parameters library_id")
	}

	if list, ok := library.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			idCollection, ok := dict.GetFloat64("id")

			if !ok {
				return ctx.Error(400, "invalid parameters collection")
			}

			if attributes, ok := dict.GetDict("attributes"); ok {
				if path, ok := attributes.GetString("path"); ok {
					for _, document := range documents.([]interface{}) {
						documentName, ok := document.(map[string]interface{})["name"].(string)
						if !ok {
							return ctx.Error(400, "invalid parameters documents")
						}

						document_path := string(path) + "/" + string(documentName)

						dir := filepath.Dir(document_path)

						err = os.MkdirAll(dir, os.ModePerm)
						if err != nil {
							fmt.Println("Gagal membuat direktori:", err)
							return ctx.Error(500, "Internal Server Error")
						}

						documentBase64, ok := document.(map[string]interface{})["base64"].(string)
						if !ok {
							return ctx.Error(400, "invalid parameters documents")
						}

						ext := filepath.Ext(documentName)

						ext = strings.TrimPrefix(ext, ".")

						if ext == "docx" || ext == "doc" {
							CR3 := cr3call.New(MODULE_MANAGER_URL)
							resp, err := CR3.Call("hexa-ai-rag", "convert-docx-to-pdf", cr3call.CR3CallData{
								Headers: U.Dict{
									"Authorization": bearer,
								},
								Params: U.Dict{
									"base64":        documentBase64,
									"document_name": documentName,
								},
							}, false)

							fmt.Printf("Debug: CR3.Call response: %v\n", resp)
							if err != nil {
								return ctx.Error(500, "Failed convert docx to pdf")
							}

							if result, ok := resp.GetDict("result"); ok {
								if base64, ok := result.GetString("base64"); ok {
									documentBase64 = string(base64)
								} else {
									return ctx.Error(500, "Failed convert docx to pdf")
								}

								if document_name, ok := result.GetString("document_name"); ok {
									documentName = string(document_name)
									document_path = string(path) + "/" + string(documentName)
								} else {
									return ctx.Error(500, "Failed convert docx to pdf")
								}
							} else {
								return ctx.Error(500, "Failed convert docx to pdf")
							}
						}

						documentFile, err := base64.StdEncoding.DecodeString(documentBase64)

						if _, err := os.Stat(document_path); err == nil {
							extension := filepath.Ext(documentName)
							name := documentName[0 : len(documentName)-len(extension)]
							timestamp := time.Now().Unix()
							documentName = fmt.Sprintf("%s_%d%s", name, timestamp, extension)
							document_path = string(path) + "/" + documentName
						}

						file, err := os.Create(document_path)
						if err != nil {
							fmt.Println("Gagal membuat file:", err)
							return ctx.Error(500, "Internal Server Error")
						}
						defer file.Close()

						_, err = file.Write(documentFile)
						if err != nil {
							fmt.Println("Gagal menulis ke file:", err)
							return ctx.Error(500, "Internal Server Error")
						}

						file_size, ok := document.(map[string]interface{})["size"].(string)
						if !ok {
							return ctx.Error(400, "invalid parameters documents")
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

						dataSend := U.Dict{
							"type":      "document",
							"name":      name,
							"path":      document_path,
							"parent_id": parent_id,
							"meta": U.Dict{
								"description":      description,
								"document_type":    document_type,
								"collection_id":    idCollection,
								"upload_by":        upload_by,
								"author":           author,
								"file_size":        file_size,
								"reference_doc_id": reference_doc_id,
								"project_doc_id":   project_doc_id,
							},
							"edited":   false,
							"inserted": false,
						}

						if category != "" {
							dataSend["meta"].(U.Dict)["category"] = category

							CR3 := cr3call.New(MODULE_MANAGER_URL)
							resp, err := CR3.Call("client-dashboard", "find-create-label", cr3call.CR3CallData{
								Headers: U.Dict{
									"Authorization": bearer,
								},
								Params: U.Dict{
									"names": category,
								},
							}, false)

							fmt.Printf("Debug: CR3.Call response: %v\n", resp)
							if err != nil {
								return ctx.Error(500, "Internal Server Error")
							}
						}

						if document_type == "Handwriting" || document_type == "Document Images" || document_type == "Document Audio" {
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

							dataSend["meta"].(U.Dict)["detected_text"] = detected_text
							dataSend["meta"].(U.Dict)["result_trocr"] = result_trocr_json
							dataSend["meta"].(U.Dict)["edited"] = edited
						}

						status, documentLibrary, err := Strapi.Create("document-library-table", dataSend)

						if err != nil {
							fmt.Printf("Error during Strapi.Create: %v\n", err)
							return ctx.Error(500, "Internal Server Error")
						}

						if status != 200 {

							if err2, ok := documentLibrary.GetDict("error"); ok {
								if message, ok := err2.GetString("message"); ok {
									fmt.Printf("Error message from Strapi response: %s\n", message)
									return ctx.Error(status, string(message))
								}
							}

							return ctx.Error(500, "Internal Server Error")
						}

						bot_ids, err := findBotIdByCollectionId(int64(parent_id), Strapi)

						if err != nil {
							fmt.Printf("Error in findBotIdByCollectionId: %v\n", err)
							return ctx.Error(500, "Internal Server Error")
						}

						if len(bot_ids) > 0 {
							for i := 0; i < len(bot_ids); i++ {
								if bot_id, ok := bot_ids.GetFloat64(i); ok {
									bot_idStr := strconv.Itoa(int(bot_id))

									CR3 := cr3call.New(MODULE_MANAGER_URL)
									resp, err := CR3.Call("hexa-ai-rag", "add-document-by-bot-id-params", cr3call.CR3CallData{
										Headers: U.Dict{
											"Authorization": bearer,
										},
										Params: U.Dict{
											"id": bot_idStr,
										},
									}, false)

									fmt.Printf("Debug: CR3.Call response: %v\n", resp)
									if err != nil {
										fmt.Printf("Error during CR3.Call for bot ID %s: %v\n", bot_idStr, err)
										return ctx.Error(500, "Internal Server Error")
									}

									go getStatusBotNew(int64(bot_id), Strapi, "0")
								} else {
									fmt.Printf("Error: Failed to get bot_id at index %d\n", i)
								}
							}
						} else {
							fmt.Println("Debug: No bot IDs found for collection")
						}

						return 200, documentLibrary, nil

					}
				}
			}
		}
	}

	return ctx.Error(500, "Internal Server Error")
}

func GetLibrarys(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	system_item, ok := ctx.Params.GetBool("system_item")
	if !ok {
		system_item = false
	}

	status, documentLibrarys, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
			},
			"system_item": U.Dict{
				"$eq": system_item,
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

	return status, documentLibrarys, nil
}

func GetLibrarysWithoutUsecase(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, documentLibrarys, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
			},
			"library_industry": U.Dict{
				"$eq": false,
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

	return status, documentLibrarys, nil
}

func GetLibraryCollections(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "invalid parameters parent_id")
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

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

	data := U.Dict{
		"filters": U.Dict{
			"parent_id": U.Dict{
				"$eq": int64(parent_id),
			},
			"type": U.Dict{
				"$eq": "collection",
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort":     "id:desc",
		"populate": "*",
	}

	if search != "" {
		data["filters"].(U.Dict)["name"] = U.Dict{
			"$containsi": search,
		}
	}

	fmt.Println("data", data)

	status, documentLibrarys, err := Strapi.Find("document-library-table", data)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, documentLibrarys, nil
}

func GetLibraryDocuments(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "invalid parameters parent_id")
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

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

	data := U.Dict{
		"filters": U.Dict{
			"parent_id": U.Dict{
				"$eq": int64(parent_id),
			},
			"type": U.Dict{
				"$eq": "document",
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort":     "id:desc",
		"populate": "*",
	}

	if search != "" {
		data["filters"].(U.Dict)["name"] = U.Dict{
			"$containsi": search,
		}
	}

	status, documentLibrarys, err := Strapi.Find("document-library-table", data)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, documentLibrarys, nil
}

func GetLibrary(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("document-library-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, dict, nil
}

func GetLibraryCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("document-library-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, dict, nil
}

func GetLibraryDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("document-library-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if attributes, ok := dict.GetDict("attributes"); ok {
		var document_path U.String

		document_path, ok = attributes.GetString("path")
		if !ok {
			document_path = ""
		}

		base64Data, err := getDocumentInBase64(string(document_path))

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		ext := filepath.Ext(string(document_path))

		return 200, U.Dict{
			"id":     dict["id"].(U.Float64),
			"data":   dict["attributes"].(U.Dict),
			"base64": base64Data,
			"ext":    ext,
		}, nil
	}

	return 200, dict, nil
}

func GetLibraryDocumentPub(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	bot_id, ok := ctx.URLParams.GetString("bot_id")
	if !ok {
		return ctx.Error(400, "invalid parameters bot_id")
	}

	status, document, err := Strapi.Find("document-library-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if attributes, ok := dict.GetDict("attributes"); ok {
		allIdCollection, err := getCollectionHexabotByBotId(string(bot_id), Strapi)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		parent_id, ok := attributes.GetFloat64("parent_id")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		found := false

		for _, v := range allIdCollection {
			if v == int(parent_id) {
				found = true
			}
		}

		if !found {
			return ctx.Error(400, "Unauthorized")
		}

		var document_path U.String

		document_path, ok = attributes.GetString("path")
		if !ok {
			document_path = ""
		}

		base64Data, err := getDocumentInBase64(string(document_path))

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		ext := filepath.Ext(string(document_path))

		return 200, U.Dict{
			"id":     dict["id"].(U.Float64),
			"data":   dict["attributes"].(U.Dict),
			"base64": base64Data,
			"ext":    ext,
		}, nil
	}

	return 200, dict, nil
}

func UpdateLibrary(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	library_name, ok := ctx.Params.GetString("library_name")
	if !ok {
		library_name = ""
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	dataSend := U.Dict{}

	if library_name != "" {
		dataSend["name"] = library_name
	}

	if description != "" {
		dataSend["meta"] = U.Dict{}
		dataSend["meta"].(U.Dict)["description"] = description
	}

	status, document, err := Strapi.Update("document-library-table", int64(id), dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)

		if err2, ok := document.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document, nil
}

func UpdateLibraryCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		collection_name = ""
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	parent_id, ok := ctx.Params.GetString("parent_id")
	if !ok {
		parent_id = ""
	}

	status, document, err := Strapi.Find("document-library-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dataSend := U.Dict{}

	if collection_name != "" {
		dataSend["name"] = collection_name
	}

	if description != "" {
		dataSend["meta"] = U.Dict{}
		dataSend["meta"].(U.Dict)["description"] = description
	}

	if parent_id != "" {
		dataSend["parent_id"] = parent_id
	}

	status, document2, err := Strapi.Update("document-library-table", int64(id), dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)

		if err2, ok := document.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document2, nil
}

func UpdateLibraryDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	document_name, ok := ctx.Params.GetString("document_name")
	if !ok {
		return ctx.Error(400, "invalid parameters document_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	detected_text, ok := ctx.Params.GetString("detected_text")
	if !ok {
		detected_text = ""
	}

	category, ok := ctx.Params.GetString("category")
	if !ok {
		category = ""
	}

	status, document, err := Strapi.Find("document-library-table", U.Dict{
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
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	parent_id, ok := attributes.GetFloat64("parent_id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	meta, ok := attributes.GetDict("meta")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	meta = U.Dict{
		"document_type":    meta["document_type"],
		"collection_id":    meta["collection_id"],
		"upload_by":        meta["upload_by"],
		"file_size":        meta["file_size"],
		"reference_doc_id": meta["reference_doc_id"],
		"project_doc_id":   meta["project_doc_id"],
	}

	if description != "" {
		meta["description"] = description
	}

	if detected_text != "" {
		meta["detected_text"] = detected_text
	}

	if category != "" {
		meta["category"] = category

		CR3 := cr3call.New(MODULE_MANAGER_URL)
		resp, err := CR3.Call("client-dashboard", "find-create-label", cr3call.CR3CallData{
			Headers: U.Dict{
				"Authorization": bearer,
			},
			Params: U.Dict{
				"names": category,
			},
		}, false)

		fmt.Printf("Debug: CR3.Call response: %v\n", resp)
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}
	}

	dataSend := U.Dict{
		"name":   document_name,
		"meta":   meta,
		"edited": true,
	}

	status, document2, err := Strapi.Update("document-library-table", int64(id), dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	bot_ids, err := findBotIdByCollectionId(int64(parent_id), Strapi)

	if err != nil {
		fmt.Printf("Error in findBotIdByCollectionId: %v\n", err)
		return ctx.Error(500, "Internal Server Error")
	}

	if len(bot_ids) > 0 {
		for i := 0; i < len(bot_ids); i++ {
			if bot_id, ok := bot_ids.GetFloat64(i); ok {
				bot_idStr := strconv.Itoa(int(bot_id))

				CR3 := cr3call.New(MODULE_MANAGER_URL)
				resp, err := CR3.Call("hexa-ai-rag", "add-document-by-bot-id-params", cr3call.CR3CallData{
					Headers: U.Dict{
						"Authorization": bearer,
					},
					Params: U.Dict{
						"id": bot_idStr,
					},
				}, false)

				fmt.Printf("Debug: CR3.Call response: %v\n", resp)
				if err != nil {
					fmt.Printf("Error during CR3.Call for bot ID %s: %v\n", bot_idStr, err)
					return ctx.Error(500, "Internal Server Error")
				}

				go getStatusBotNew(int64(bot_id), Strapi, "0")
			} else {
				fmt.Printf("Error: Failed to get bot_id at index %d\n", i)
			}
		}
	} else {
		fmt.Println("Debug: No bot IDs found for collection")
	}

	return 200, document2, nil
}

func UpdateLibraryDocumentParent(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	document_id, ok := ctx.Params.GetFloat64("document_id")
	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	var tmpParentPath U.String

	parent, err := getLibraryDict(int(parent_id), Strapi)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if attributes, ok := parent.GetDict("attributes"); ok {
		if path, ok := attributes.GetString("path"); ok {
			tmpParentPath = path
		} else {
			return ctx.Error(500, "Internal Server Error")
		}
	} else {
		return ctx.Error(500, "Internal Server Error")
	}

	var tmpLibraryType U.String
	var tmpLibraryPath U.String

	library, err := getLibraryDict(int(document_id), Strapi)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if attributes, ok := library.GetDict("attributes"); ok {
		if libraryType, ok := attributes.GetString("type"); ok {
			tmpLibraryType = libraryType
		} else {
			return ctx.Error(500, "Internal Server Error")
		}
		if libraryPath, ok := attributes.GetString("path"); ok {
			tmpLibraryPath = libraryPath
		} else {
			return ctx.Error(500, "Internal Server Error")
		}
	} else {
		return ctx.Error(500, "Internal Server Error")
	}

	if tmpLibraryType == "collection" {
		err := moveFolder(string(tmpLibraryPath), string(tmpParentPath))
		if err != nil {
			fmt.Println("Gagal memindahkan folder:", err)
		} else {
			fmt.Println("Folder berhasil dipindahkan")
		}
	} else if tmpLibraryType == "document" {
		err := moveFile(string(tmpLibraryPath), string(tmpParentPath))
		if err != nil {
			fmt.Println("Gagal memindahkan file:", err)
			return ctx.Error(500, "Failed to move document")
		} else {
			fmt.Println("File berhasil dipindahkan")
		}
	}

	filename := filepath.Base(string(tmpLibraryPath))

	newPath := tmpParentPath + "/" + U.String(filename)

	if tmpLibraryType == "collection" {
		err := updatePathContent(string(newPath), int64(document_id), Strapi)

		if err != nil {
			return ctx.Error(500, "Failed update path content")
		}
	}

	status, document2, err := Strapi.Update("document-library-table", int64(document_id), U.Dict{
		"parent_id": parent_id,
		"path":      newPath,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, document2, nil
}

func DeleteLibrary(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	for {
		status, collections, err := Strapi.Find("document-library-table", U.Dict{
			"filters": U.Dict{
				"type": U.Dict{
					"$eq": "collection",
				},
				"parent_id": U.Dict{
					"$eq": id,
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
			fmt.Println("collections", collections)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := collections.GetList("data")

		if !ok {
			fmt.Println(ok)
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) < 1 {
			break
		}

		for i := 0; i < len(list); i++ {
			if dict, ok := list.GetDict(i); ok {
				if idCollection, ok := dict.GetFloat64("id"); ok {

					var page4 = 1

					var idsCol U.List

					for {
						status, documents, err := Strapi.Find("document-library-table", U.Dict{
							"filters": U.Dict{
								"type": U.Dict{
									"$eq": "collection",
								},
								"parent_id": U.Dict{
									"$eq": idCollection,
								},
							},
							"pagination": U.Dict{
								"page":     page4,
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
							fmt.Println("ayam 1", documents)
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
								if idCol, ok := dict.GetFloat64("id"); ok {
									idsCol = append(idsCol, idCol)
								}
							}
						}

						page4++
					}

					idsCol = append(idsCol, idCollection)

					var idsDoc U.List

					for c := 0; c < len(idsCol); c++ {
						var page2 = 1

						for {
							status, documents, err := Strapi.Find("document-library-table", U.Dict{
								"filters": U.Dict{
									"type": U.Dict{
										"$eq": "document",
									},
									"parent_id": U.Dict{
										"$eq": idsCol[c],
									},
								},
								"pagination": U.Dict{
									"page":     page2,
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
								fmt.Println("documents", documents)
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

							page2++
						}
					}

					knowledge_ids, err := findKnowledgeIdBotByCollectionId(int64(idCollection), Strapi)

					if err != nil {
						fmt.Println(err)
						return ctx.Error(500, "Internal Server Error")
					}

					for i := 0; i < len(idsDoc); i++ {
						if idDoc, ok := idsDoc.GetFloat64(i); ok {
							for i := 0; i < len(knowledge_ids); i++ {
								if knowledge_id, ok := knowledge_ids.GetString(i); ok {
									idDocStr := strconv.Itoa(int(idDoc))

									// dataSend := U.Dict{
									// 	"knowledge_id": knowledge_id,
									// 	"document_id":  idDocStr,
									// }

									// http := requests.NewHttp(U.Dict{
									// 	"Content-type": "application/json",
									// })

									// _, err = http.Post(fmt.Sprintf("%s/api/%s", HEXA_RAG_URL, "remove-document"), nil, dataSend.ToJson())

									CR3 := cr3call.New(MODULE_MANAGER_URL)
									resp, err := CR3.Call("hexa-ai-rag", "remove-document", cr3call.CR3CallData{
										Headers: U.Dict{
											"Authorization": bearer,
										},
										Params: U.Dict{
											"knowledge_id": knowledge_id,
											"document_id":  idDocStr,
										},
									}, false)

									fmt.Printf("Debug: CR3.Call response: %v\n", resp)
									if err != nil {
										return ctx.Error(500, "Internal Server Error")
									}

									// var deleteResponse DeleteResponse
									// err = json.Unmarshal(response.Body, &deleteResponse)
									// if err != nil {
									// 	fmt.Println(err)
									// 	return ctx.Error(500, "Internal Server Error")
									// }

									// fmt.Println("delete response ", knowledge_id, deleteResponse)
								} else {
									fmt.Println("failed to get knowledge_id")
								}
							}

							status, document, err := Strapi.Delete("document-library-table", int64(idDoc))
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

					for c := 0; c < len(idsCol); c++ {
						var idDelete int

						idDelete = int(idsCol[c].(U.Float64))

						fmt.Println("Try delete: ", idDelete)

						status, delete, err := Strapi.Delete("document-library-table", int64(idDelete))
						if err != nil {
							fmt.Println(err)
							return ctx.Error(500, "Internal Server Error")
						}

						if status != 200 {
							fmt.Println("Failed delete: ", idDelete)
							fmt.Println(delete)
							fmt.Println(status)
							return ctx.Error(500, "Internal Server Error")
						}
					}

					page++
				}
			}
		}
	}

	var page3 = 1
	var idsBot U.List

	for {
		status, data, err := Strapi.Find("hexabot-table", U.Dict{
			"filters": U.Dict{
				"library_id": U.Dict{
					"$eq": id,
				},
			},
			"pagination": U.Dict{
				"page":     page3,
				"pageSize": 25,
			},
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			fmt.Println("data", data)
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

		page3++
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

	status, document, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
			},
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
		fmt.Println("document 2", document)
		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := document.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			if attributes, ok := dict.GetDict("attributes"); ok {
				if path, ok := attributes.GetString("path"); ok {
					find := true

					if _, err := os.Stat(string(path)); os.IsNotExist(err) {
						fmt.Println("Direktori atau file tidak ditemukan:", path)
						find = false
					}

					if find {
						err = os.RemoveAll(string(path))
						if err != nil {
							fmt.Println("Gagal menghapus direktori:", err)
							return ctx.Error(500, "Internal Server Error")
						}

						fmt.Println("Direktori berhasil dihapus:", path)
					}
				}
			}
		}
	}

	status, deleted, err := Strapi.Delete("document-library-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println("deleted", deleted)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

func DeleteLibraryCollection(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	var idsCol U.List

	for {
		status, documents, err := Strapi.Find("document-library-table", U.Dict{
			"filters": U.Dict{
				"type": U.Dict{
					"$eq": "collection",
				},
				"parent_id": U.Dict{
					"$eq": id,
				},
			},
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
			fmt.Println("ayam 1", documents)
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
				if idCol, ok := dict.GetFloat64("id"); ok {
					idsCol = append(idsCol, idCol)
				}
			}
		}

		page++
	}

	idsCol = append(idsCol, id)

	var idsDoc U.List

	for c := 0; c < len(idsCol); c++ {
		var page2 = 1

		for {
			status, documents, err := Strapi.Find("document-library-table", U.Dict{
				"filters": U.Dict{
					"type": U.Dict{
						"$eq": "document",
					},
					"parent_id": U.Dict{
						"$eq": idsCol[c],
					},
				},
				"pagination": U.Dict{
					"page":     page2,
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
				fmt.Println("ayam 1", documents)
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

			page2++
		}

	}

	knowledge_ids, err := findKnowledgeIdBotByCollectionId(int64(id), Strapi)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	for i := 0; i < len(idsDoc); i++ {
		if idDoc, ok := idsDoc.GetFloat64(i); ok {
			for i := 0; i < len(knowledge_ids); i++ {
				if knowledge_id, ok := knowledge_ids.GetString(i); ok {
					idDocStr := strconv.Itoa(int(idDoc))

					// dataSend := U.Dict{
					// 	"knowledge_id": knowledge_id,
					// 	"document_id":  idDocStr,
					// }

					// http := requests.NewHttp(U.Dict{
					// 	"Content-type": "application/json",
					// })

					// _, err = http.Post(fmt.Sprintf("%s/api/%s", HEXA_RAG_URL, "remove-document"), nil, dataSend.ToJson())

					CR3 := cr3call.New(MODULE_MANAGER_URL)
					resp, err := CR3.Call("hexa-ai-rag", "remove-document", cr3call.CR3CallData{
						Headers: U.Dict{
							"Authorization": bearer,
						},
						Params: U.Dict{
							"knowledge_id": knowledge_id,
							"document_id":  idDocStr,
						},
					}, false)

					fmt.Printf("Debug: CR3.Call response: %v\n", resp)
					if err != nil {
						return ctx.Error(500, "Internal Server Error")
					}

					// var deleteResponse DeleteResponse
					// err = json.Unmarshal(response.Body, &deleteResponse)
					// if err != nil {
					// 	fmt.Println(err)
					// 	return ctx.Error(500, "Internal Server Error")
					// }

					// fmt.Println("delete response ", knowledge_id, deleteResponse)
				} else {
					fmt.Println("failed to get knowledge_id")
				}
			}

			status, document, err := Strapi.Delete("document-library-table", int64(idDoc))
			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println("ayam 2", document)
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}
		}
	}

	status, document, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "collection",
			},
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
		fmt.Println("ayam 3", document)
		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := document.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			if attributes, ok := dict.GetDict("attributes"); ok {
				if path, ok := attributes.GetString("path"); ok {
					find := true

					if _, err := os.Stat(string(path)); os.IsNotExist(err) {
						fmt.Println("Direktori atau file tidak ditemukan:", path)
						find = false
					}

					if find {
						err = os.RemoveAll(string(path))
						if err != nil {
							fmt.Println("Gagal menghapus direktori:", err)
							return ctx.Error(500, "Internal Server Error")
						}

						fmt.Println("Direktori berhasil dihapus:", path)
					}
				}
			}
		}
	}

	for c := 0; c < len(idsCol); c++ {
		var idDelete int

		if idsCol[c] == id {
			idDelete = idsCol[c].(int)
		} else {
			idDelete = int(idsCol[c].(U.Float64))
		}

		status, _, err = Strapi.Delete("document-library-table", int64(idDelete))
		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println("ayam 4", document)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}
	}

	return 200, U.Dict{}, nil
}

func DeleteLibraryDocument(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "document",
			},
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

	if list, ok := document.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			if attributes, ok := dict.GetDict("attributes"); ok {
				if parent_id, ok := attributes.GetFloat64("parent_id"); ok {
					fmt.Println("found parent_id")
					knowledge_ids, err := findKnowledgeIdBotByCollectionId(int64(parent_id), Strapi)

					if err != nil {
						fmt.Println(err)
						return ctx.Error(500, "Internal Server Error")
					}

					for i := 0; i < len(knowledge_ids); i++ {
						if knowledge_id, ok := knowledge_ids.GetString(i); ok {
							idDocStr := strconv.Itoa(int(id))

							// dataSend := U.Dict{
							// 	"knowledge_id": knowledge_id,
							// 	"document_id":  idDocStr,
							// }

							// http := requests.NewHttp(U.Dict{
							// 	"Content-type": "application/json",
							// })

							// _, err = http.Post(fmt.Sprintf("%s/api/%s", HEXA_RAG_URL, "remove-document"), nil, dataSend.ToJson())

							CR3 := cr3call.New(MODULE_MANAGER_URL)
							resp, err := CR3.Call("hexa-ai-rag", "remove-document", cr3call.CR3CallData{
								Headers: U.Dict{
									"Authorization": bearer,
								},
								Params: U.Dict{
									"knowledge_id": knowledge_id,
									"document_id":  idDocStr,
								},
							}, false)

							fmt.Printf("Debug: CR3.Call response: %v\n", resp)
							if err != nil {
								return ctx.Error(500, "Internal Server Error")
							}

							// var deleteResponse DeleteResponse
							// err = json.Unmarshal(response.Body, &deleteResponse)
							// if err != nil {
							// 	fmt.Println(err)
							// 	return ctx.Error(500, "Internal Server Error")
							// }

							// fmt.Println("delete response ", knowledge_id, deleteResponse)
						} else {
							fmt.Println("failed to get knowledge_id")
						}
					}
				}

				parent_id2, ok := attributes.GetFloat64("parent_id")

				if !ok {
					fmt.Println("not found parent_id", parent_id2)
					return ctx.Error(500, "Internal Server Error")
				}

				if path, ok := attributes.GetString("path"); ok {
					find := true

					if _, err := os.Stat(string(path)); os.IsNotExist(err) {
						fmt.Println("Direktori atau file tidak ditemukan:", path)
						find = false
					}

					if find {
						err = os.RemoveAll(string(path))
						if err != nil {
							fmt.Println("Gagal menghapus direktori:", err)
							return ctx.Error(500, "Internal Server Error")
						}

						fmt.Println("Direktori berhasil dihapus:", path)
					}
				}
			}
		}
	}

	status, _, err = Strapi.Delete("document-library-table", int64(id))
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

func getDocumentInBase64(document_path string) (string, error) {
	file, err := os.Open(document_path)
	if err != nil {
		fmt.Println("Gagal membuka file:", err)
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Gagal mendapatkan informasi file:", err)
		return "", err
	}

	content := make([]byte, fileInfo.Size())
	_, err = file.Read(content)
	if err != nil {
		fmt.Println("Gagal membaca konten file:", err)
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

func GetStatusLibrary(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	page1 := 1
	colectionsIds := []int{}
	documentsIds := []int{}

	findNotReady := false
	statusRAG := "Ready"

	for {
		status, collections, err := Strapi.Find("document-library-table", U.Dict{
			"filters": U.Dict{
				"type": U.Dict{
					"$eq": "collection",
				},
				"parent_id": U.Dict{
					"$eq": id,
				},
			},
			"populate": "*",
			"pagination": U.Dict{
				"page":     page1,
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

		if list, ok := collections.GetList("data"); ok {
			if len(list) == 0 {
				break
			}

			for i := 0; i < len(list); i++ {
				if dict, ok := list.GetDict(i); ok {
					if id, ok := dict.GetFloat64("id"); ok {
						colectionsIds = append(colectionsIds, int(id))
					}
				}
			}
		}

		page1++
	}

	for i := 0; i < len(colectionsIds); i++ {
		page2 := 1
		for {
			status, collection, err := Strapi.Find("document-library-table", U.Dict{
				"filters": U.Dict{
					"type": U.Dict{
						"$eq": "document",
					},
					"parent_id": U.Dict{
						"$eq": colectionsIds[i],
					},
				},
				"populate": "*",
				"pagination": U.Dict{
					"page":     page2,
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

			if list, ok := collection.GetList("data"); ok {
				if len(list) == 0 {
					break
				}

				for i := 0; i < len(list); i++ {
					if dict, ok := list.GetDict(i); ok {
						if id, ok := dict.GetFloat64("id"); ok {
							documentsIds = append(documentsIds, int(id))
							if attributes, ok := dict.GetDict("attributes"); ok {
								if meta, ok := attributes.GetDict("meta"); ok {
									if status, ok := meta.GetString("status"); ok {
										if status != "Ready" {
											statusRAG = "Reading"
											findNotReady = true
											break
										}
									} else {
										statusRAG = "Reading"
										findNotReady = true
										break
									}
								}
							}
						}
					}
				}

				if findNotReady {
					break
				}

				page2++
			}
		}

		if findNotReady {
			break
		}
	}

	if len(documentsIds) == 0 {
		statusRAG = "Empty"
	}

	return 200, U.Dict{
		"status": statusRAG,
	}, nil
}

func GetCountLibrariesPerMonth(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, libraries, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
			},
			"system_item": U.Dict{
				"$eq": false,
			},
			"createdAt": U.Dict{
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
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, documents, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "document",
			},
			"createdAt": U.Dict{
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
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{
		"library":  libraries,
		"document": documents,
	}, nil

}

func FindCreateLibraryByName(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	name, ok := ctx.Params.GetString("library_name")

	if !ok {
		return ctx.Error(400, "invalid parameters library_name")
	}

	description, ok := ctx.Params.GetString("description")

	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	status, library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
			},
			"name": U.Dict{
				"$eq": name,
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

	if list, ok := library.GetList("data"); ok {
		if len(list) > 0 {
			if dict, ok := list.GetDict(0); ok {
				if id, ok := dict.GetFloat64("id"); ok {
					return 200, U.Dict{
						"id": id,
					}, nil
				}
			}
		} else {
			document_path := DOCUMENT_PATH + "/" + string(name)

			dir := filepath.Dir(document_path)

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				fmt.Println("Gagal membuat direktori:", err)
				return ctx.Error(500, "Internal Server Error")
			}

			status, _, err = Strapi.Create("document-library-table", U.Dict{
				"type": "library",
				"name": name,
				"path": document_path,
				"meta": U.Dict{
					"description": description,
				},
				"library_industry": true,
			})

			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}

			if data, ok := library.GetDict("data"); ok {
				if id, ok := data.GetFloat64("id"); ok {
					return 200, U.Dict{
						"id": id,
					}, nil
				}
			}
		}
	}

	return ctx.Error(500, "Internal Server Error")
}

func FindCreateLibraryCollectionByName(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")
	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, _, err := Strapi.UsersMe()
	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	name, ok := ctx.Params.GetString("collection_name")
	if !ok {
		return ctx.Error(400, "invalid parameters collection_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "invalid parameters parent_id")
	}

	status, library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type":      U.Dict{"$eq": "collection"},
			"name":      U.Dict{"$eq": name},
			"parent_id": U.Dict{"$eq": parent_id},
		},
	})

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := library.GetList("data"); ok && len(list) > 0 {
		if dict, ok := list.GetDict(0); ok {
			if id, ok := dict.GetFloat64("id"); ok {
				return 200, U.Dict{"id": id}, nil
			}
		}
	}

	var document_path string
	status, parent_library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{"$eq": parent_id},
		},
	})

	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := parent_library.GetList("data"); ok && len(list) > 0 {
		if dict, ok := list.GetDict(0); ok {
			if attributes, ok := dict.GetDict("attributes"); ok {
				if path, ok := attributes.GetString("path"); ok {
					document_path = string(path) + "/" + string(name)
				}
			}
		}
	} else {
		return ctx.Error(400, "Parent library not found")
	}

	dir := filepath.Dir(document_path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return ctx.Error(500, "Failed to create directory")
	}

	dataSend := U.Dict{
		"type":      "collection",
		"name":      name,
		"path":      document_path,
		"parent_id": parent_id,
		"meta": U.Dict{
			"description": description,
		},
	}

	status, documentLibrary, err := Strapi.Create("document-library-table", dataSend)
	if err != nil || status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	if data, ok := documentLibrary.GetDict("data"); ok {
		if id, ok := data.GetFloat64("id"); ok {
			return 200, U.Dict{"id": id}, nil
		}
	}

	return ctx.Error(500, "Internal Server Error")
}

func findKnowledgeIdBotByCollectionId(collection_id int64, Strapi *strapi.StrapiClassic) (U.List, error) {
	page := 1

	var knowledgeIds U.List

	for {
		status, data, err := Strapi.Find("hexabot-table", U.Dict{
			"filters": U.Dict{
				"collection_id": U.Dict{
					"$containsi": collection_id,
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"populate": "*",
			"sort":     "updatedAt:desc",
		})

		if err != nil {
			fmt.Printf("Error during Strapi.Find on page %d: %v\n", page, err)
			return U.List{}, err
		}

		if status != 200 {
			fmt.Printf("Error: Non-200 status from Strapi.Find: %d\n", status)
			fmt.Printf("Data received: %v\n", data)
			return U.List{}, fmt.Errorf("unexpected status: %d", status)
		}

		list, ok := data.GetList("data")
		if !ok {
			fmt.Printf("Error: Unable to retrieve 'data' list from response on page %d\n", page)
			return U.List{}, fmt.Errorf("failed to parse response data")
		}

		if len(list) == 0 {
			fmt.Println("Debug: No more data to process, exiting loop")
			break
		}

		for i := 0; i < len(list); i++ {
			if dict, ok := list.GetDict(i); ok {
				if attributes, ok := dict.GetDict("attributes"); ok {
					knowledge_id, ok := attributes.GetString("bot_id")

					if !ok {
						fmt.Printf("Error: Missing 'knowledge_id' field in item %d\n", i+1)
						return U.List{}, fmt.Errorf("missing 'knowledge_id' in data")
					}

					if collection_ids, ok := attributes.GetString("collection_id"); ok {
						collectionStrings := strings.Split(string(collection_ids), ",")

						for _, collectionStr := range collectionStrings {
							id, err := strconv.Atoi(strings.TrimSpace(collectionStr))
							if err != nil {
								fmt.Printf("Error converting '%s' to int: %v\n", collectionStr, err)
								continue
							}

							if id == int(collection_id) {
								knowledgeIds = append(knowledgeIds, knowledge_id)
							}
						}
					}
				}
			}
		}

		page++
	}

	return knowledgeIds, nil
}

func findBotIdByCollectionId(collection_id int64, Strapi *strapi.StrapiClassic) (U.List, error) {
	page := 1
	var botIds U.List

	for {
		status, data, err := Strapi.Find("hexabot-table", U.Dict{
			"filters": U.Dict{
				"collection_id": U.Dict{
					"$containsi": collection_id,
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"populate": "*",
			"sort":     "updatedAt:desc",
		})

		if err != nil {
			fmt.Printf("Error during Strapi.Find on page %d: %v\n", page, err)
			return U.List{}, err
		}

		if status != 200 {
			fmt.Printf("Error: Non-200 status from Strapi.Find: %d\n", status)
			fmt.Printf("Data received: %v\n", data)
			return U.List{}, fmt.Errorf("unexpected status: %d", status)
		}

		list, ok := data.GetList("data")
		if !ok {
			fmt.Printf("Error: Unable to retrieve 'data' list from response on page %d\n", page)
			return U.List{}, fmt.Errorf("failed to parse response data")
		}

		if len(list) == 0 {
			fmt.Println("Debug: No more data to process, exiting loop")
			break
		}

		for i := 0; i < len(list); i++ {
			if dict, ok := list.GetDict(i); ok {
				botId, ok := dict.GetFloat64("id")
				if !ok {
					fmt.Printf("Error: Missing 'id' field in item %d\n", i+1)
					return U.List{}, fmt.Errorf("missing 'id' in data")
				}

				if attributes, ok := dict.GetDict("attributes"); ok {
					if collection_ids, ok := attributes.GetString("collection_id"); ok {
						collectionStrings := strings.Split(string(collection_ids), ",")

						for _, collectionStr := range collectionStrings {
							id, err := strconv.Atoi(strings.TrimSpace(collectionStr))
							if err != nil {
								fmt.Printf("Error converting '%s' to int: %v\n", collectionStr, err)
								continue
							}

							if id == int(collection_id) {
								botIds = append(botIds, botId)
							}
						}
					} else {
						fmt.Println("Debug: 'collection_id' not found in attributes")
					}
				} else {
					fmt.Println("Debug: 'attributes' not found in item")
				}
			} else {
				fmt.Printf("Error: Item %d is not a dictionary\n", i+1)
			}
		}

		page++
	}

	return botIds, nil
}

func IsLibraryNameAvailable(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearerOk := ctx.Headers.GetString("Authorization")
	if !bearerOk {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	name, ok := ctx.Params.GetString("name")
	if !ok {
		return ctx.Error(400, "Invalid parameter: name")
	}

	status, library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"name": U.Dict{
				"$eq": name,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("Error response from Strapi:", status)
		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := library.GetList("data"); ok && len(list) > 0 {
		return 200, U.Dict{
			"available": false,
		}, nil
	}

	return 200, U.Dict{
		"available": true,
	}, nil
}

func FindLibraryByParentId(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearerOk := ctx.Headers.GetString("Authorization")
	if !bearerOk {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	parent_id, ok := ctx.Params.GetFloat64("parent_id")
	if !ok {
		return ctx.Error(400, "Invalid parameter: parent_id")
	}

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

	dataSend := U.Dict{
		"filters": U.Dict{
			"parent_id": U.Dict{
				"$eq": int64(parent_id),
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	}

	if search != "" {
		dataSend["filters"].(U.Dict)["name"] = U.Dict{
			"$containsi": search,
		}
	} else {
		dataSend["sort"] = "type:asc"
	}

	status, library, err := Strapi.Find("document-library-table", dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("Error response from Strapi:", status)
		return ctx.Error(500, "Internal Server Error")
	}
	return 200, library, nil
}

func FindLibrary(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearerOk := ctx.Headers.GetString("Authorization")
	if !bearerOk {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

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

	dataSend := U.Dict{
		"filters": U.Dict{},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	}

	if search != "" {
		dataSend["filters"].(U.Dict)["name"] = U.Dict{
			"$containsi": search,
		}
	} else {
		dataSend["sort"] = "type:asc"
	}

	status, library, err := Strapi.Find("document-library-table", dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("Error response from Strapi:", status)
		return ctx.Error(500, "Internal Server Error")
	}
	return 200, library, nil
}

func FindLibraryForBakrie(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearerOk := ctx.Headers.GetString("Authorization")
	if !bearerOk {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

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

	dataSend := U.Dict{
		"filters": U.Dict{
			"$or": U.List{
				U.Dict{
					"type": "collection",
				},
				U.Dict{
					"type": "document",
				},
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	}

	if search != "" {
		dataSend["filters"].(U.Dict)["name"] = U.Dict{
			"$containsi": search,
		}
	} else {
		dataSend["sort"] = "type:asc"
	}

	status, library, err := Strapi.Find("document-library-table", dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("Error response from Strapi:", status)
		return ctx.Error(500, "Internal Server Error")
	}
	return 200, library, nil
}

type DeleteResponse struct {
	DocumentId string `json:"document_id"`
}

func getLibraryDict(id int, strapi *strapi.StrapiClassic) (U.Dict, error) {
	status, document, err := strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": id,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return U.Dict{}, err
	}

	if status != 200 {
		fmt.Println("document", document)
		fmt.Println(status)
		return U.Dict{}, err
	}

	list, ok := document.GetList("data")

	if !ok {
		fmt.Println(ok)
		return U.Dict{}, err
	}

	if len(list) < 1 {
		return U.Dict{}, fmt.Errorf("not found")
	}

	if dict, ok := list.GetDict(0); ok {
		return dict, nil
	} else {
		return U.Dict{}, err
	}
}

func moveFile(sourcePath, destinationDir string) error {
	err := os.MkdirAll(destinationDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("gagal membuat folder tujuan: %w", err)
	}

	fileName := filepath.Base(sourcePath)

	destinationPath := filepath.Join(destinationDir, fileName)

	err = os.Rename(sourcePath, destinationPath)
	if err != nil {
		return fmt.Errorf("gagal memindahkan file: %w", err)
	}

	fmt.Println("File berhasil dipindahkan ke:", destinationPath)
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		} else {
			return copyFile(path, destPath)
		}
	})
}

func moveFolder(sourcePath, destinationDir string) error {
	folderName := filepath.Base(sourcePath)
	destinationPath := filepath.Join(destinationDir, folderName)

	err := copyDir(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	err = os.RemoveAll(sourcePath)
	if err != nil {
		return err
	}

	fmt.Println("Folder berhasil dipindahkan ke:", destinationPath)
	return nil
}

func updatePathContent(pathLibrary string, idLibrary int64, strapi *strapi.StrapiClassic) error {
	page := 1

	var collectionIds U.List

	for {
		status, documentLibrarys, err := strapi.Find("document-library-table", U.Dict{
			"filters": U.Dict{
				"parent_id": U.Dict{
					"$eq": int64(idLibrary),
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"sort":     "id:desc",
			"populate": "*",
		})

		if err != nil {
			fmt.Println(err)
			return err
		}

		if status != 200 {
			fmt.Println(status)
			fmt.Println(documentLibrarys)
			return fmt.Errorf("status not 200")
		}

		if list, ok := documentLibrarys.GetList("data"); ok {
			if len(list) < 1 {
				break
			}

			for i := 0; i < len(list); i++ {
				if dict, ok := list.GetDict(i); ok {
					id, ok := dict.GetFloat64("id")

					if !ok {
						return fmt.Errorf("failed to get id content")
					}

					if attributes, ok := dict.GetDict("attributes"); ok {
						if typeLibrary, ok := attributes.GetString("type"); ok {
							if pathContent, ok := attributes.GetString("path"); ok {
								filename := filepath.Base(string(pathContent))

								newPath := pathLibrary + "/" + filename

								status, documentLibrarys2, err := strapi.Update("document-library-table", int64(id), U.Dict{
									"path": newPath,
								})

								if err != nil {
									fmt.Println(err)
									return err
								}

								if status != 200 {
									fmt.Println(status)
									fmt.Println(documentLibrarys2)
									return fmt.Errorf("failed to update path content")
								}

								if typeLibrary == "collection" {
									collectionIds = append(collectionIds, U.Dict{
										"id":   id,
										"path": newPath,
									})
								}
							} else {
								return fmt.Errorf("failed to get path content")
							}
						} else {
							return fmt.Errorf("failed to get type content")
						}
					} else {
						return fmt.Errorf("failed to get attributes content")
					}
				} else {
					return fmt.Errorf("failed to get dict content")
				}
			}

			page = page + 1
		} else {
			return fmt.Errorf("failed to get data content")
		}
	}

	if len(collectionIds) == 0 {
		return nil
	} else {
		for j := 0; j < len(collectionIds); j++ {
			if dictContent, ok := collectionIds.GetDict(j); ok {
				id, id_ok := dictContent.GetFloat64("id")

				pathContent, path_ok := dictContent.GetString("path")

				if !id_ok || !path_ok {
					return fmt.Errorf("failed to get collection content, id: %v, path: %v", id, pathContent)
				}

				err := updatePathContent(string(pathContent), int64(id), strapi)

				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func FindCreateLibraryLlm(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, library, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
			},
			"system_item": U.Dict{
				"$eq": true,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("ayam")
		fmt.Println(status)
		fmt.Println(library)
		return ctx.Error(500, "Internal Server Error")
	}

	if list, ok := library.GetList("data"); ok {
		if len(list) > 0 {
			if dict, ok := list.GetDict(0); ok {
				if id, ok := dict.GetFloat64("id"); ok {
					idCollection, err := findCreateLibraryCollection(int64(id), Strapi)

					if err != nil {
						fmt.Println(err)
						return ctx.Error(500, "Internal Server Error")
					}

					return 200, U.Dict{
						"collection_id": idCollection,
						"library_id":    id,
					}, nil
				}
			}
		} else {
			name := "Library LLM"
			description := "This is a library for try LLM"

			document_path := DOCUMENT_PATH + "/" + string(name)

			dir := filepath.Dir(document_path)

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				fmt.Println("Gagal membuat direktori:", err)
				return ctx.Error(500, "Internal Server Error")
			}

			status, library, err = Strapi.Create("document-library-table", U.Dict{
				"type": "library",
				"name": name,
				"path": document_path,
				"meta": U.Dict{
					"description": description,
				},
				"system_item": true,
			})

			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println("bebek")
				fmt.Println(status)
				fmt.Println(library)
				return ctx.Error(500, "Internal Server Error")
			}

			if data, ok := library.GetDict("data"); ok {
				if id, ok := data.GetFloat64("id"); ok {
					idCollection, err := findCreateLibraryCollection(int64(id), Strapi)

					if err != nil {
						fmt.Println(err)
						return ctx.Error(500, "Internal Server Error")
					}

					return 200, U.Dict{
						"collection_id": idCollection,
						"library_id":    id,
					}, nil
				}
			}
		}
	}

	return ctx.Error(500, "Internal Server Error")
}

func findCreateLibraryCollection(id int64, strapi *strapi.StrapiClassic) (int, error) {
	data := U.Dict{
		"filters": U.Dict{
			"parent_id": U.Dict{
				"$eq": int64(id),
			},
			"type": U.Dict{
				"$eq": "collection",
			},
		},
		"sort":     "id:desc",
		"populate": "*",
	}

	fmt.Println("data", data)

	status, documentLibrarys, err := strapi.Find("document-library-table", data)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if status != 200 {
		fmt.Println("cicak")
		fmt.Println(status)
		fmt.Println(documentLibrarys)
		return 0, fmt.Errorf("status not 200")
	}

	if list, ok := documentLibrarys.GetList("data"); ok {
		if len(list) > 0 {
			if dict, ok := list.GetDict(0); ok {
				if id, ok := dict.GetFloat64("id"); ok {
					return int(id), nil
				}
			}
		} else {
			upload_by := 0

			author := "system"

			name := "Collection LLM"

			description := "this is collection for try LLM"

			status, library, err := strapi.Find("document-library-table", U.Dict{
				"filters": U.Dict{
					"id": U.Dict{
						"$eq": id,
					},
				},
			})

			if err != nil {
				fmt.Println(err)
				return 0, err
			}

			if status != 200 {
				fmt.Println(status)
				fmt.Println(library)
				return 0, fmt.Errorf("status not 200")
			}

			var document_path string

			if list, ok := library.GetList("data"); ok {
				if dict, ok := list.GetDict(0); ok {
					if attributes, ok := dict.GetDict("attributes"); ok {
						if path, ok := attributes.GetString("path"); ok {
							document_path = string(path) + "/" + string(name)
						}
					}
				}
			}

			dir := filepath.Dir(document_path)

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				fmt.Println("Gagal membuat direktori:", err)
				return 0, err
			}

			dataSend := U.Dict{
				"type":      "collection",
				"name":      name,
				"path":      document_path,
				"parent_id": id,
				"meta": U.Dict{
					"upload_by": upload_by,
					"author":    author,
				},
			}

			if description != "" {
				dataSend["meta"].(U.Dict)["description"] = description
			}

			status, documentLibrary, err := strapi.Create("document-library-table", dataSend)

			if err != nil {
				fmt.Println(err)
				return 0, err
			}

			if status != 200 {
				fmt.Println(status)

				if err2, ok := documentLibrary.GetDict("error"); ok {
					if message, ok := err2.GetString("message"); ok {
						return 0, fmt.Errorf(string(message))
					}
				}

				return 0, fmt.Errorf("status not 200")
			}

			if data, ok := documentLibrary.GetDict("data"); ok {
				if id, ok := data.GetFloat64("id"); ok {
					return int(id), nil
				}
			}

			return 0, fmt.Errorf("not found")
		}
	}

	return 0, fmt.Errorf("not found")
}

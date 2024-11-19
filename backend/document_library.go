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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	library_name, ok := ctx.Params.GetString("library_name")
	if !ok {
		return ctx.Error(400, "invalid parameters library_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	document_path := DOCUMENT_PATH + "/" + string(library_name)

	dir := filepath.Dir(document_path)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Gagal membuat direktori:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	dataSend := U.Dict{
		"type": "library",
		"name": library_name,
		"path": document_path,
		"meta": U.Dict{
			"description": description,
		},
		"library_industry": false,
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

	// 	// CR3 := cr3call.New(MODULE_MANAGER_URL)
	// 	// go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
	// 	// 	Headers: U.Dict{
	// 	// 		"Authorization": bearer,
	// 	// 	},
	// 	// }, false)

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
			"description": description,
		},
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

	// 	// CR3 := cr3call.New(MODULE_MANAGER_URL)
	// 	// go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
	// 	// 	Headers: U.Dict{
	// 	// 		"Authorization": bearer,
	// 	// 	},
	// 	// }, false)

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

	name, ok := ctx.Params.GetString("document_name")
	if !ok {
		return ctx.Error(400, "invalid parameters document_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
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
								"file_size":        file_size,
								"reference_doc_id": reference_doc_id,
								"project_doc_id":   project_doc_id,
							},
							"edited":   false,
							"inserted": false,
						}

						if category != "" {
							dataSend["meta"].(U.Dict)["category"] = category
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

						// 	// CR3 := cr3call.New(MODULE_MANAGER_URL)
						// 	// go CR3.Call("hexa-ai", "getting-new-document", cr3call.CR3CallData{
						// 	// 	Headers: U.Dict{
						// 	// 		"Authorization": bearer,
						// 	// 	},
						// 	// }, false)

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

	status, documentLibrarys, err := Strapi.Find("document-library-table", U.Dict{
		"filters": U.Dict{
			"type": U.Dict{
				"$eq": "library",
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
		return ctx.Error(400, "invalid parameters library_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	document_path := DOCUMENT_PATH + "/" + string(library_name)

	dir := filepath.Dir(document_path)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Gagal membuat direktori:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	dataSend := U.Dict{
		"name": library_name,
		"path": document_path,
		"meta": U.Dict{
			"description": description,
		},
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
		return ctx.Error(400, "invalid parameters collection_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		return ctx.Error(400, "invalid parameters description")
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

	document_path_old, ok := attributes.GetString("path")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dir_old := filepath.Dir(string(document_path_old))

	// folder_name_old := filepath.Base(string(document_path_old))

	document_path_new := filepath.Join(dir_old, string(collection_name))

	fmt.Println("Old path:", document_path_old)
	fmt.Println("New path:", document_path_new)

	dir := filepath.Dir(document_path_new)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Gagal membuat direktori:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	dataSend := U.Dict{
		"name": collection_name,
		"path": document_path_new,
		"meta": U.Dict{
			"description": description,
		},
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
					var page2 = 1

					var idsDoc U.List

					for {
						status, documents, err := Strapi.Find("document-library-table", U.Dict{
							"filters": U.Dict{
								"type": U.Dict{
									"$eq": "document",
								},
								"parent_id": U.Dict{
									"$eq": idCollection,
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

					// CR3 := cr3call.New(MODULE_MANAGER_URL)
					// resp, err := CR3.Call("hexa-ai-rag", "delete-document-by-collection-id", cr3call.CR3CallData{
					// 	Headers: U.Dict{
					// 		"Authorization": bearer,
					// 	},
					// 	Params: U.Dict{
					// 		"collection_id": idStr,
					// 	},
					// }, false)

					// fmt.Println(resp)

					// if err != nil {
					// 	fmt.Println(err)
					// 	return ctx.Error(500, "Internal Server Error")
					// }

					for i := 0; i < len(idsDoc); i++ {
						if idDoc, ok := idsDoc.GetFloat64(i); ok {
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

					status, document, err := Strapi.Delete("document-library-table", int64(idCollection))
					if err != nil {
						fmt.Println(err)
						return ctx.Error(500, "Internal Server Error")
					}

					if status != 200 {
						fmt.Println(document)
						fmt.Println(status)
						return ctx.Error(500, "Internal Server Error")
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

					// Menghapus direktori beserta isinya
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

	var page2 = 1

	var idsDoc U.List

	for {
		status, documents, err := Strapi.Find("document-library-table", U.Dict{
			"filters": U.Dict{
				"type": U.Dict{
					"$eq": "document",
				},
				"parent_id": U.Dict{
					"$eq": id,
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

	// CR3 := cr3call.New(MODULE_MANAGER_URL)
	// resp, err := CR3.Call("hexa-ai-rag", "delete-document-by-collection-id", cr3call.CR3CallData{
	// 	Headers: U.Dict{
	// 		"Authorization": bearer,
	// 	},
	// 	Params: U.Dict{
	// 		"collection_id": idStr,
	// 	},
	// }, false)

	// fmt.Println(resp)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	for i := 0; i < len(idsDoc); i++ {
		if idDoc, ok := idsDoc.GetFloat64(i); ok {
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

					// Menghapus direktori beserta isinya
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
		fmt.Println("ayam 4", document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
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

	// CR3 := cr3call.New(MODULE_MANAGER_URL)
	// resp, err := CR3.Call("hexa-ai-rag", "delete-document-by-collection-id", cr3call.CR3CallData{
	// 	Headers: U.Dict{
	// 		"Authorization": bearer,
	// 	},
	// 	Params: U.Dict{
	// 		"collection_id": idStr,
	// 	},
	// }, false)

	// fmt.Println(resp)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

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
				if path, ok := attributes.GetString("path"); ok {
					find := true

					if _, err := os.Stat(string(path)); os.IsNotExist(err) {
						fmt.Println("Direktori atau file tidak ditemukan:", path)
						find = false
					}

					// Menghapus direktori beserta isinya
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
	// Membuka file
	file, err := os.Open(document_path)
	if err != nil {
		fmt.Println("Gagal membuka file:", err)
		return "", err
	}
	defer file.Close()

	// Mendapatkan informasi file
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Gagal mendapatkan informasi file:", err)
		return "", err
	}

	// Mengalokasikan buffer untuk menyimpan konten
	content := make([]byte, fileInfo.Size())
	_, err = file.Read(content)
	if err != nil {
		fmt.Println("Gagal membaca konten file:", err)
		return "", err
	}

	// Mengonversi konten ke base64
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
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	fmt.Println(libraries)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, libraries, nil

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

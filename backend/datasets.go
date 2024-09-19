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
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateConversation(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	title, ok := ctx.Params.GetString("title")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	conversationId := uuid.New().String()

	status, conversation, err := Strapi.Create("conversation-fine-tuned-table", U.Dict{
		"title":           title,
		"user":            id,
		"conversation_id": conversationId,
	})

	fmt.Println(conversation)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{
		"conversation_id": conversationId,
	}, nil
}

func CreateMessage(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	messageString, ok := ctx.Params.GetString("message")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var message []interface{}
	err = json.Unmarshal([]byte(messageString), &message)
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, update, err := Strapi.Update("conversation-fine-tuned-table", int64(id), U.Dict{
		"message": message,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func AddMessage(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")

	if !ok {
		fmt.Println("Error getting attributes")
		return ctx.Error(500, "Internal Server Error")
	}

	already_messages, ok := attributes.GetList("message")

	if !ok {
		fmt.Println("Error getting already_messages")
		return ctx.Error(500, "Internal Server Error")
	}

	messageString, ok := ctx.Params.GetString("message")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var message U.List
	err = json.Unmarshal([]byte(messageString), &message)
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	already_messages = append(already_messages, message...)

	status, update, err := Strapi.Update("conversation-fine-tuned-table", int64(id), U.Dict{
		"message": already_messages,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func UpdateMessage(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")

	if !ok {
		fmt.Println("Error getting attributes")
		return ctx.Error(500, "Internal Server Error")
	}

	already_messages, ok := attributes.GetList("message")

	if !ok {
		fmt.Println("Error getting already_messages")
		return ctx.Error(500, "Internal Server Error")
	}

	messageString, ok := ctx.Params.GetString("message")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var message []interface{}
	err = json.Unmarshal([]byte(messageString), &message)
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	for i := 0; i < len(already_messages); i++ {
		for j := 0; j < len(message); j++ {
			alreadyID := already_messages[i].(U.Dict)["id"].(U.Float64)
			messageID := message[j].(map[string]interface{})["id"].(float64)

			if int64(alreadyID) == int64(messageID) {
				already_messages[i].(U.Dict)["message"] = message[j].(map[string]interface{})["message"]
				break
			}
		}
	}

	status, update, err := Strapi.Update("conversation-fine-tuned-table", int64(id), U.Dict{
		"message": already_messages,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func GetMessages(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	limit, ok := ctx.Params.GetFloat64("limit")
	if !ok || limit <= 0 {
		limit = 100
	}

	page, ok := ctx.Params.GetFloat64("page")
	if !ok || page <= 0 {
		page = 1
	}

	offset := int((page - 1) * limit)

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	message, ok := attributes.GetList("message")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	totalMessages := len(message)
	pageTotal := int(math.Ceil(float64(totalMessages) / float64(limit)))

	start := totalMessages - int(limit) - int(offset)
	if start < 0 {
		start = 0
	}
	end := totalMessages - int(offset)

	if start >= totalMessages {
		message = U.List{}
	} else if end > totalMessages {
		message = message[start:]
	} else {
		message = message[start:end]
	}

	attributes["message"] = message
	dict["attributes"] = attributes
	dict["page_total"] = pageTotal

	return 200, dict, nil
}

func UpdateTitle(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	title, ok := ctx.Params.GetString("title")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, update, err := Strapi.Update("conversation-fine-tuned-table", int64(id), U.Dict{
		"title": title,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func GetConversations(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		page_size = 25
	}

	search, ok := ctx.Params.GetString("search")

	if !ok {
		search = ""
	}

	if search != "" {
		status, conversations, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
			"filters": U.Dict{
				"title": U.Dict{
					"$contains": search,
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": page_size,
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

		return 200, conversations, nil
	} else {
		status, conversations, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": page_size,
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

		return 200, conversations, nil
	}
}

func DeleteConversation(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	status, _, err = Strapi.Delete("conversation-fine-tuned-table", int64(id))

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

func RequestFineTuning(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	page := 1

	var allConversations U.Dict

	allConversations = U.Dict{
		"data": []interface{}{},
	}

	for {
		status, conversations, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
			"populate": "*",
			"sort":     "updatedAt:asc",
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := conversations.GetList("data")
		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) == 0 {
			break
		}

		existingList, ok := allConversations.GetList("data")
		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		allConversations["data"] = append(existingList, list...)

		page++
	}

	transformedData := transformData(allConversations)
	transformedJSON, err := json.MarshalIndent(transformedData, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	status, upload, err := Strapi.Upload("request-fine-tuning.txt", transformedJSON)

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

	idFile, ok := data.GetFloat64("id")
	if !ok {
		fmt.Println("failed to get id")
		return ctx.Error(500, "Internal Server Error")
	}

	status, checkRequest, err := Strapi.Find("request-fine-tuning-table", U.Dict{
		"filters": U.Dict{
			"status": U.Dict{
				"$eq": "REQUEST",
			},
		},
		"sort": "updatedAt:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := checkRequest.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) > 0 {
		data, ok := list.GetDict(0)

		if !ok {
			fmt.Println("failed to get dict")
			return ctx.Error(500, "Internal Server Error")
		}

		idRequest, ok := data.GetFloat64("id")

		if !ok {
			fmt.Println("failed to get id")
			return ctx.Error(500, "Internal Server Error")
		}

		status, fine_tuning, err := Strapi.Update("request-fine-tuning-table", int64(idRequest), U.Dict{
			"user": int64(id),
			"file": int64(idFile),
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, fine_tuning, nil
	} else {
		status, fine_tuning, err := Strapi.Create("request-fine-tuning-table", U.Dict{
			"user":   int64(id),
			"file":   int64(idFile),
			"status": "REQUEST",
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, fine_tuning, nil
	}
}

func GetStatusRequestFineTuning(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, fine_tuning, err := Strapi.Find("request-fine-tuning-table", U.Dict{
		"sort": "updatedAt:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, fine_tuning, nil
}

func ImportConversationWithCsv(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	conversation_id, ok := ctx.Params.GetString("conversation_id")
	if !ok {
		return ctx.Error(400, "invalid parameters conversation_id")
	}

	documentsString, ok := ctx.Params.GetString("documents")
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	var documents interface{}
	err = json.Unmarshal([]byte(documentsString), &documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

	datas := U.List{}

	for _, document := range documents.([]interface{}) {

		documentBase64, ok := document.(map[string]interface{})["base64"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		chat := transformConversationFromCsv(documentBase64)

		fmt.Println("chat")

		datas = append(datas, chat...)
	}

	fmt.Println(datas)

	status, conversation, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"conversation_id": U.Dict{
				"$eq": conversation_id,
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

	list, ok := conversation.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "conversation not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	id, ok := dict.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := dict.GetDict("attributes")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	message, ok := attributes.GetList("message")

	if !ok {
		message = U.List{}
	}

	datas = append(message, datas...)

	status, update, err := Strapi.Update("conversation-fine-tuned-table", int64(id), U.Dict{
		"message": datas,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func ImportDatasetsWithCsv(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	idUser, ok := me.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	documentsString, ok := ctx.Params.GetString("documents")

	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	var documents interface{}
	err = json.Unmarshal([]byte(documentsString), &documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

	datas := U.List{}

	for _, document := range documents.([]interface{}) {

		documentBase64, ok := document.(map[string]interface{})["base64"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		conversation := transformDatasetFromCsv(documentBase64, Strapi, int64(idUser))

		datas = append(datas, conversation...)
	}

	return 200, U.Dict{"data": datas}, nil

}

func transformConversationFromCsv(base64String string) U.List {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println("Error decoding CSV:", err)
		return U.List{}
	}

	reader := csv.NewReader(bytes.NewReader(decoded))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return U.List{}
	}

	var result = U.List{}
	var id int
	var currentPrompts []string
	var lastResponse string

	for _, row := range records {
		if len(row) < 2 {
			fmt.Println("Format baris tidak valid:", row)
			continue
		}

		role := strings.TrimSpace(row[0])
		text := strings.TrimSpace(row[1])

		if role == "Sebagai user" || role == "user" {
			currentPrompts = append(currentPrompts, text)
		} else if role == "Sebagai bot" || role == "bot" {
			lastResponse = text
			for _, prompt := range currentPrompts {
				result = append(result, U.Dict{
					"id":      id,
					"role":    "user",
					"message": prompt,
				})
				id++
				result = append(result, U.Dict{
					"id":      id,
					"role":    "assistant",
					"message": lastResponse,
				})
				id++
			}
			currentPrompts = nil
		}
	}

	if len(currentPrompts) > 0 && lastResponse != "" {
		for _, prompt := range currentPrompts {
			result = append(result, U.Dict{
				"id":      id,
				"role":    "user",
				"message": prompt,
			})
			id++
			result = append(result, U.Dict{
				"id":      id,
				"role":    "assistant",
				"message": lastResponse,
			})
			id++
		}
	}

	return result
}

func transformDatasetFromCsv(base64String string, strapi *strapi.StrapiClassic, idUser int64) U.List {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println("Error decoding CSV:", err)
		return U.List{}
	}

	reader := csv.NewReader(bytes.NewReader(decoded))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return U.List{}
	}

	datas := U.List{}

	for index, record := range records {
		if len(record) < 3 {
			fmt.Println("Format baris tidak valid:", record)
			continue
		}

		if index == 0 {
			continue
		}

		var conversations U.List

		datasetName := strings.TrimSpace(record[0])
		userMessage := strings.TrimSpace(record[1])
		botResponse := strings.TrimSpace(record[2])

		conversations = append(conversations, U.Dict{
			"id":      1,
			"role":    "user",
			"message": userMessage,
		})
		conversations = append(conversations, U.Dict{
			"id":      2,
			"role":    "assistant",
			"message": botResponse,
		})

		conversationId := uuid.New().String()

		status, conversation, err := strapi.Create("conversation-fine-tuned-table", U.Dict{
			"title":           datasetName,
			"user":            idUser,
			"conversation_id": conversationId,
		})

		if err != nil {
			fmt.Println(err)
			return U.List{}
		}

		if status != 200 {
			fmt.Println(status)
			return U.List{}
		}

		dict, ok := conversation.GetDict("data")

		if !ok {
			fmt.Println("Error getting data")
			return U.List{}
		}

		id, ok := dict.GetFloat64("id")

		if !ok {
			fmt.Println("Error getting id")
			return U.List{}
		}

		status, update, err := strapi.Update("conversation-fine-tuned-table", int64(id), U.Dict{
			"message": conversations,
		})

		if err != nil {
			fmt.Println(err)
			return U.List{}
		}

		if status != 200 {
			fmt.Println(status)
			return U.List{}
		}

		datas = append(datas, update)
	}

	return datas
}

func transformData(allConversations U.Dict) U.List {
	var result U.List

	list, ok := allConversations.GetList("data")

	if !ok {
		return result
	}

	for i := 0; i < len(list); i++ {
		var prompt, response string

		data, ok := list.GetDict(i)

		if !ok {
			fmt.Println("Error getting data")
			return result
		}

		attributes, ok := data.GetDict("attributes")

		if !ok {
			fmt.Println("Error getting attributes")
			return result
		}

		message, ok := attributes.GetList("message")

		if !ok {
			fmt.Println("Error getting message")
			return result
		}

		for j := 0; j < len(message); j++ {
			msgMap, ok := message.GetDict(j)

			if !ok {
				fmt.Println("Error getting message map")
				return result
			}

			role, ok := msgMap.GetString("role")

			if !ok {
				fmt.Println("Error getting role")
				return result
			}

			message, ok := msgMap.GetString("message")

			if !ok {
				fmt.Println("Error getting message")
				return result
			}

			if role == "user" {
				prompt = string(message)
			} else if role == "assistant" {
				response = string(message)
				result = append(result, U.Dict{
					"prompt":   prompt,
					"response": response,
				})
				prompt = ""
				response = ""
			}
		}
	}

	return result
}

// report

func SendReport(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	messageStr, ok := ctx.Params.GetString("message")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	var message U.List

	err := json.Unmarshal([]byte(messageStr), &message)

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	collection_id, ok := ctx.Params.GetFloat64("collection_id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	model, ok := ctx.Params.GetString("model")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	status, report, err := Strapi.Create("report-conversation-table", U.Dict{
		"message":       message,
		"collection_id": int64(collection_id),
		"model":         model,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, report, nil
}

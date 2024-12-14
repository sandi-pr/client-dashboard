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
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func GetCountDatasetPerMount(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, dataset, err := Strapi.Find("conversation-fine-tuned-table", U.Dict{
		"filters": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	fmt.Println(dataset)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, dataset, nil
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

	firstLine := strings.Split(string(decoded), "\n")[0]
	delimiter := detectDelimiter(firstLine)

	reader := csv.NewReader(bytes.NewReader(decoded))
	reader.Comma = delimiter

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

	firstLine := strings.Split(string(decoded), "\n")[0]
	delimiter := detectDelimiter(firstLine)

	reader := csv.NewReader(bytes.NewReader(decoded))
	reader.Comma = delimiter

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

func detectDelimiter(line string) rune {
	delimiters := []rune{',', ';', '\t', '|'} // Daftar delimiter yang umum
	for _, delim := range delimiters {
		if strings.Count(line, string(delim)) > 0 {
			return delim // Jika ditemukan, return delimiter
		}
	}
	return ',' // Default ke koma jika tidak ada delimiter yang dikenali
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

func CreateDatasetHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(401, "Unauthorized")
	}

	name, ok := ctx.Params.GetString("name")
	if !ok {
		return ctx.Error(400, "invalid parameters name")
	}

	document_str, ok := ctx.Params.GetString("documents")
	if !ok {
		return ctx.Error(400, "invalid parameters documents")
	}

	var documents interface{}
	err = json.Unmarshal([]byte(document_str), &documents)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters documents")
	}

	for _, document := range documents.([]interface{}) {
		documentName, ok := document.(map[string]interface{})["name"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}
		document_path := DATASETS_PATH + "/Handwriting/" + string(documentName)

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
			document_path = DATASETS_PATH + "/Handwriting/" + documentName
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

		detected_text, ok := ctx.Params.GetString("detected_text")
		if !ok {
			return ctx.Error(400, "invalid parameters detected_text")
		}

		result_trocr, ok := ctx.Params.GetString("result_trocr")
		if !ok {
			return ctx.Error(400, "invalid parameters result_trocr")
		}

		var result_trocr_json interface{}
		err = json.Unmarshal([]byte(result_trocr), &result_trocr_json)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		dataSend := U.Dict{
			"name":          name,
			"path":          document_path,
			"user":          upload_by,
			"file_size":     file_size,
			"detected_text": detected_text,
			"result_trocr":  result_trocr_json,
		}

		status, documentLibrary, err := Strapi.Create("dataset-handwriting-table", dataSend)

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

		return 200, documentLibrary, nil

	}

	return ctx.Error(500, "Internal Server Error")
}

func GetDatasetsHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	search, ok := ctx.Params.GetString("search")
	if !ok {
		search = ""
	}

	data := U.Dict{
		"filters": U.Dict{},
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

	status, documentLibrarys, err := Strapi.Find("dataset-handwriting-table", data)

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

func GetDatasetHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("dataset-handwriting-table", U.Dict{
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

	return ctx.Error(500, "Internal Server Error")
}

func UpdateDatasetHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	name, ok := ctx.Params.GetString("name")
	if !ok {
		name = ""
	}

	detected_text, ok := ctx.Params.GetString("detected_text")
	if !ok {
		detected_text = ""
	}

	if name == "" && detected_text == "" {
		return ctx.Error(400, "invalid parameters")
	}

	dataSend := U.Dict{}

	if name != "" {
		dataSend["name"] = name
	}

	if detected_text != "" {
		dataSend["detected_text"] = detected_text
	}

	status, document2, err := Strapi.Update("dataset-handwriting-table", int64(id), dataSend)

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

func DeleteDatasetHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, document, err := Strapi.Find("dataset-handwriting-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": id,
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

	status, _, err = Strapi.Delete("dataset-handwriting-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{
		"status": "success",
	}, nil
}

func RequestFineTuningHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	user_id, ok := me.GetFloat64("id")
	if !ok {
		return ctx.Error(401, "Unauthorized")
	}

	dataSend := U.Dict{
		"user":   user_id,
		"status": "REQUEST",
	}

	status, _, err = Strapi.Create("request-fine-tuning-handwriting-table", dataSend)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{
		"status": "success",
	}, nil
}

func GetStatusRequestFineTuningHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, fine_tuning, err := Strapi.Find("request-fine-tuning-handwriting-table", U.Dict{
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

	list, ok := fine_tuning.GetList("data")

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

	statusFinetuning, ok := attributes.GetDict("status")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{
		"status": statusFinetuning,
	}, nil
}

func DownloadAllDatasetHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	page := 1

	var AllDatasets U.List

	for {
		status, document, err := Strapi.Find("dataset-handwriting-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
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

		if len(list) < 1 {
			break
		}

		for i := 0; i < len(list); i++ {

			dict, ok := list.GetDict(i)

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

				name, ok := attributes.GetString("name")
				if !ok {
					return ctx.Error(500, "Internal Server Error")
				}

				detected_text, ok := attributes.GetString("detected_text")
				if !ok {
					return ctx.Error(500, "Internal Server Error")
				}

				ext := filepath.Ext(string(document_path))

				AllDatasets = append(AllDatasets, U.Dict{
					"base64":        base64Data,
					"name":          string(name) + ext,
					"detected_text": detected_text,
				})
			}
		}

		page++
	}

	zipBase64, err := createZipAndConvertToBase64(AllDatasets)
	if err != nil {
		fmt.Println("Error zip:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{
		"base64": zipBase64,
	}, nil
}

func createZipAndConvertToBase64(datasets U.List) (string, error) {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	for _, dataset := range datasets {
		name, _ := dataset.(U.Dict).GetString("name")
		base64Content, _ := dataset.(U.Dict).GetString("base64")
		detectedText, _ := dataset.(U.Dict).GetString("detected_text")

		content, err := base64.StdEncoding.DecodeString(string(base64Content))
		if err != nil {
			return "", fmt.Errorf("failed to decode base64: %v", err)
		}

		zipFile, err := zipWriter.Create(string(name))
		if err != nil {
			return "", fmt.Errorf("failed to create file in ZIP: %v", err)
		}
		_, err = zipFile.Write(content)
		if err != nil {
			return "", fmt.Errorf("failed to write file to ZIP: %v", err)
		}

		textFileName := filepath.Base(string(name)) + ".txt"
		textFile, err := zipWriter.Create(textFileName)
		if err != nil {
			return "", fmt.Errorf("failed to create text file in ZIP: %v", err)
		}
		_, err = textFile.Write([]byte(detectedText))
		if err != nil {
			return "", fmt.Errorf("failed to write text file to ZIP: %v", err)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close ZIP writer: %v", err)
	}

	zipBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return zipBase64, nil
}

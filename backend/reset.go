/*
 * CR3 Client Dashboard
 * Package: gitlab.hexacode.org/cr3/client-dashboard
 * Maintainer: Azzis Arswendo <azzis@hexacode.org>
 *
 * Copyright (C) 2024 Hexacode Teknologi Indonesia
 * All Rights Reserved
 */

package main

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.hexacode.org/go-libs/cr3call"
	"gitlab.hexacode.org/go-libs/hctypes"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func getPageCount(table string, token string, page_size int) (int, error) {
	tb := fmt.Sprintf("%s-table", table)

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	page := 1

	status, result, err := Strapi.Find(tb, hctypes.Dict{
		"sort": hctypes.List{
			"id:desc",
		},
		"pagination": hctypes.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	})

	if err != nil {
		return -1, err
	}

	if status != 200 {
		return -1, fmt.Errorf("error getting page count: %d", status)
	}

	if meta, ok := result.GetDict("meta"); ok {
		if pagination, ok := meta.GetDict("pagination"); ok {
			if count, ok := pagination.GetFloat64("pageCount"); ok {
				return int(count), nil
			}
		}
	}

	return -1, fmt.Errorf("error getting page count")
}

func getIDsPerPage(table string, token string, page_size int, page int) ([]int64, error) {
	tb := fmt.Sprintf("%s-table", table)

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	status, result, err := Strapi.Find(tb, hctypes.Dict{
		"sort": hctypes.List{
			"id:desc",
		},
		"pagination": hctypes.Dict{
			"page":     page,
			"pageSize": page_size,
		},
	})

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("error getting ids: %d", status)
	}

	if data, ok := result.GetList("data"); ok {
		ids := []int64{}

		for i := 0; i < len(data); i++ {
			if row, ok := data.GetDict(i); ok {
				if id, ok := row.GetFloat64("id"); ok {
					ids = append(ids, int64(id))
				}
			}
		}

		return ids, nil
	}

	return nil, fmt.Errorf("error getting ids")
}

func getIDs(table string, token string) ([]int64, error) {
	page_size := 25
	page := 1

	page_count, err := getPageCount(table, token, page_size)
	if err != nil {
		return nil, err
	}

	ids := []int64{}

	for page <= page_count {
		ids_page, err := getIDsPerPage(table, token, page_size, page)
		if err != nil {
			return nil, err
		}

		ids = append(ids, ids_page...)

		page++
	}

	return ids, nil
}

func getIDsUsers(token string) ([]int64, error) {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	status, result, err := Strapi.UsersFind()
	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("error getting ids: %d", status)
	}

	ids := []int64{}

	for i := 0; i < len(result); i++ {
		if row, ok := result.GetDict(i); ok {
			if id, ok := row.GetFloat64("id"); ok {
				ids = append(ids, int64(id))
			}
		}
	}

	return ids, nil
}

func reset(table string, token string) error {
	tb := fmt.Sprintf("%s-table", table)
	ids, err := getIDs(table, token)
	if err != nil {
		return err
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	for i := 0; i < len(ids); i++ {
		status, _, err := Strapi.Delete(tb, ids[i])
		if err != nil {
			return err
		}

		if status != 200 {
			return fmt.Errorf("error deleting id %d: %d", ids[i], status)
		}

		if table == "document-management" {
			CR3 := cr3call.New(MODULE_MANAGER_URL)
			go CR3.Call("hexa-ai-rag", "delete-document-by-id", cr3call.CR3CallData{
				Headers: U.Dict{
					"Authorization": "Bearier " + token,
				},
				Params: U.Dict{
					"document_id": strconv.Itoa(int(ids[i])),
				},
			}, false)
		}
	}

	return nil
}

func resetUser(token string) error {
	ids, err := getIDsUsers(token)
	if err != nil {
		return err
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)

	for i := 0; i < len(ids); i++ {
		status, _, err := Strapi.UsersDelete(ids[i])
		if err != nil {
			return err
		}

		if status != 200 {
			return fmt.Errorf("error deleting id %d: %d", ids[i], status)
		}
	}

	return nil
}

func Reset(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	token := strings.Replace(string(bearer), "Bearer ", "", 1)
	if token == "" {
		return ctx.Error(401, "Unauthorized")
	}

	tables := []string{
		"calendar",
		"category-product",
		"person",
		"invoice",
		"order-booking",
		"product",
		"resource",
		"resource-type",
		"company-profile",
		"booking",
		"document-management",
		"document-library",
		"references-document",
		"industries",
		"contract-rules",
		"contract-rule-sections",
		"projects",
		"project-documents",
		"project-results",

		"conversation",
		"conversation-message",
		"ai-history",

		"conversation-rag",
		"conversation-message-rag",
		"ai-history-rag",
	}

	for i := 0; i < len(tables); i++ {
		err := reset(tables[i], token)
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}
	}

	err := resetUser(STRAPI_TOKEN)
	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{}, nil
}

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
	"fmt"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func GetLabels(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
	var allLabel U.List

	for {
		status, labels, err := Strapi.Find("label-table", U.Dict{
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
		})

		if err != nil {
			fmt.Println("Error:", err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println("Status code:", status)
			return ctx.Error(500, "Internal Server Error")
		}

		list, ok := labels.GetList("data")
		if !ok || len(list) < 1 {
			break
		}

		for i := 0; i < len(list); i++ {
			dict, ok := list.GetDict(i)
			if !ok {
				continue
			}

			attributes, ok := dict.GetDict("attributes")
			if !ok {
				continue
			}

			labelName, ok := attributes.GetString("name")
			if !ok {
				continue
			}

			allLabel = append(allLabel, labelName)
		}

		page++
	}

	return 200, U.Dict{"labels": allLabel}, nil
}

func FindCreateLabel(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	names, ok := ctx.Params.GetString("names")

	if !ok {
		return ctx.Error(400, "invalid parameters names")
	}

	nameArray := strings.Split(string(names), ",")

	for i := range nameArray {
		nameArray[i] = strings.TrimSpace(nameArray[i])
	}

	for i := 0; i < len(nameArray); i++ {
		status, find, err := Strapi.Find("label-table", U.Dict{
			"filters": U.Dict{
				"name": U.Dict{
					"$eq": nameArray[i],
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

		list, ok := find.GetList("data")
		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) < 1 {
			status, _, err := Strapi.Create("label-table", U.Dict{
				"name": nameArray[i],
			})

			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}

		} else {
			continue
		}
	}

	return 200, U.Dict{}, nil
}

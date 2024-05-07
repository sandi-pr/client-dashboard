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

func GetSchedule(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	start_date, ok := ctx.Params.GetString("start_date")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	end_date, ok := ctx.Params.GetString("end_date")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, calendar, err := Strapi.Find("calendar-table", U.Dict{
		"filtees": U.Dict{
			"date": U.Dict{
				"gte": start_date,
				"lte": end_date,
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

	return 200, calendar, nil
}

func CreateSchedule(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	title, ok := ctx.Params.GetString("title")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	date, ok := ctx.Params.GetString("date")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	start_time, ok := ctx.Params.GetString("start_time")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	end_time, ok := ctx.Params.GetString("end_time")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	// availaible_slot, ok := ctx.Params.GetString("availaible_slot")

	// if !ok {
	// 	return ctx.Error(400, "invalid parameters")
	// }

	// note, ok := ctx.Params.GetString("note")

	// if !ok {
	// 	return ctx.Error(400, "invalid parameters")
	// }

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, calendar, err := Strapi.Create("calendar-table", U.Dict{
		"title":      title,
		"date":       date,
		"start_time": start_time,
		"end_time":   end_time,
		// "availaible_slot": availaible_slot,
		// "note":      note,
		"is_active": is_active,
	})

	fmt.Println("calendar", calendar)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, calendar, nil
}

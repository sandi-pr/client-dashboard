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
	"encoding/json"
	"fmt"
	"strconv"
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

	status, booking, err := Strapi.Find("booking-table", U.Dict{
		"filters": U.Dict{
			"date": U.Dict{
				"$gte": start_date,
				"$lte": end_date,
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

	status, availability, err := Strapi.Find("calendar-table", U.Dict{})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	availabilityArray, _ := availability.GetList("data")

	calendar, _ := availabilityArray.GetDict(0)

	if calendar == nil {
		calendar = nil
	}

	return 200, U.Dict{"booking": booking, "calendar": calendar}, nil
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

	description, ok := ctx.Params.GetString("description")

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

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	personString, ok := ctx.Params.GetString("person")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var person []interface{}
	err = json.Unmarshal([]byte(personString), &person)
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	var booking U.Dict

	if len(person) > 0 {
		status, booking, err = Strapi.Create("booking-table", U.Dict{
			"title":       title,
			"description": description,
			"date":        date,
			"start_time":  start_time,
			"end_time":    end_time,
			"is_active":   is_active,
			"person":      person,
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
		status, booking, err = Strapi.Create("booking-table", U.Dict{
			"title":       title,
			"description": description,
			"date":        date,
			"start_time":  start_time,
			"end_time":    end_time,
			"is_active":   is_active,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}
	}

	fmt.Println("booking", booking)

	return 200, booking, nil
}

func DeleteSchedule(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, booking, err := Strapi.Find("booking-table", U.Dict{
		"filters": U.Dict{
			"date": U.Dict{
				"$eq": date,
			},
			"start_time": U.Dict{
				"$eq": start_time,
			},
			"end_time": U.Dict{
				"$eq": end_time,
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

	bookingArray, _ := booking.GetList("data")

	if len(bookingArray) == 0 {
		return 200, U.Dict{}, nil
	}

	dict := bookingArray[0].(U.Dict)

	id, _ := dict.GetFloat64("id")

	status, booking, err = Strapi.Delete("booking-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(booking)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, booking, nil
}

func CreateAvailability(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	availabilityString, ok := ctx.Params.GetString("availability")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var availability interface{}
	err = json.Unmarshal([]byte(availabilityString), &availability)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	individualMode, ok := ctx.Params.GetString("individualMode")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, calendar, err := Strapi.Create("calendar-table", U.Dict{
		"availability":   availability,
		"individualMode": individualMode,
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

func UpdateAvailability(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	availabilityId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(availabilityId))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	availabilityString, ok := ctx.Params.GetString("availability")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var availability interface{}
	err = json.Unmarshal([]byte(availabilityString), &availability)

	if err != nil {
		return ctx.Error(400, "invalid availability format")
	}

	individualMode, ok := ctx.Params.GetString("individualMode")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, calendar, err := Strapi.Update("calendar-table", int64(id), U.Dict{
		"individualMode": individualMode,
		"availability":   availability,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(calendar)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, calendar, nil
}

func GetCountSchedulePerMount(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, booking, err := Strapi.Find("booking-table", U.Dict{
		"filters": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	fmt.Println(booking)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, booking, nil
}

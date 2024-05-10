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

func GetOrders(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, orders, err := Strapi.Find("order-booking-table", U.Dict{
		"pagination": U.Dict{
			"page": page,
		},
		"populate": "*",
		"sort":     "id:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, orders, nil

}

func CreateOrder(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	order_date, ok := ctx.Params.GetString("order_date")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	order_statusString, ok := ctx.Params.GetString("order_status")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var order_status interface{}
	err = json.Unmarshal([]byte(order_statusString), &order_status)

	if err != nil {
		return ctx.Error(400, "invalid availability format")
	}

	person, ok := ctx.Params.GetFloat64("person")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	product, ok := ctx.Params.GetFloat64("product")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	resource, ok := ctx.Params.GetFloat64("resource")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	shipping_address, ok := ctx.Params.GetString("shipping_address")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	total_amount, ok := ctx.Params.GetFloat64("total_amount")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, order, err := Strapi.Create("order-booking-table", U.Dict{
		"order_date":       order_date,
		"order_status":     order_status,
		"person":           person,
		"product":          product,
		"resources":        resource,
		"shipping_address": shipping_address,
		"total_amount":     total_amount,
	})

	fmt.Println(order)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, order, nil
}

func GetOrder(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	userId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(userId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, order, err := Strapi.Find("order-booking-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": id,
			},
		},
		"populate": "*",
	})

	fmt.Println(order)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, order, nil
}

func UpdateOrder(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	orderId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(orderId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	order_date, ok := ctx.Params.GetString("order_date")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	booking_date, ok := ctx.Params.GetString("booking_date")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	total_amount, ok := ctx.Params.GetString("total_amount")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	note, ok := ctx.Params.GetString("note")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	order_status, ok := ctx.Params.GetString("order_status")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	booking_time, ok := ctx.Params.GetString("booking_time")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	shipping_address, ok := ctx.Params.GetString("shipping_address")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, order, err := Strapi.Update("order-booking-table", int64(id), U.Dict{
		"data": U.Dict{
			"order_date":       order_date,
			"booking_date":     booking_date,
			"total_amount":     total_amount,
			"note":             note,
			"order_status":     order_status,
			"booking_time":     booking_time,
			"shipping_address": shipping_address,
		},
	})

	fmt.Println(order)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, order, nil
}

func DeleteOrder(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	orderId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(orderId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, order, err := Strapi.Delete("order-booking-table", int64(id))

	fmt.Println(order)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, order, nil
}

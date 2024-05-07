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

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func GetPersons(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		fmt.Println("invalid parameters")
		page = 1
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

	status, persons, err := Strapi.Find("person-table", U.Dict{
		"pagination": U.Dict{
			"page": page,
		},
		"sort": "id:desc",
	})
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, persons, nil
}

func CreatePerson(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	fullname, ok := ctx.Params.GetString("fullname")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	email, ok := ctx.Params.GetString("email")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	phone, ok := ctx.Params.GetString("phone")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	address, ok := ctx.Params.GetString("address")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	city, ok := ctx.Params.GetString("city")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	state, ok := ctx.Params.GetString("state")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	zip_code, ok := ctx.Params.GetString("zip_code")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	country, ok := ctx.Params.GetString("country")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	date_of_birth, ok := ctx.Params.GetString("date_of_birth")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, persons, err := Strapi.Create("person-table", U.Dict{
		"fullname":      fullname,
		"email":         email,
		"phone":         phone,
		"address":       address,
		"city":          city,
		"state":         state,
		"zip_code":      zip_code,
		"country":       country,
		"date_of_birth": date_of_birth,
		"is_active":     is_active,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, persons, nil
}

func GetPerson(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	status, persons, err := Strapi.FindOne("person-table", int64(id))
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, persons, nil
}

func UpdatePerson(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	fullname, ok := ctx.Params.GetString("fullname")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	email, ok := ctx.Params.GetString("email")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	phone, ok := ctx.Params.GetString("phone")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	address, ok := ctx.Params.GetString("address")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	city, ok := ctx.Params.GetString("city")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	state, ok := ctx.Params.GetString("state")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	zip_code, ok := ctx.Params.GetString("zip_code")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	country, ok := ctx.Params.GetString("country")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	date_of_birth, ok := ctx.Params.GetString("date_of_birth")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, persons, err := Strapi.Update("person-table", int64(id), U.Dict{
		"fullname":      fullname,
		"email":         email,
		"phone":         phone,
		"address":       address,
		"city":          city,
		"state":         state,
		"zip_code":      zip_code,
		"country":       country,
		"date_of_birth": date_of_birth,
		"is_active":     is_active,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, persons, nil

}

func DeletePerson(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	status, _, err = Strapi.Delete("person-table", int64(id))

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

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
	"errors"
	"fmt"
	"os"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

var MODULE_MANAGER_URL = ""
var STRAPI_URL = ""

func init() {
	MODULE_MANAGER_URL = os.Getenv("MODULE_MANAGER_URL")
	if MODULE_MANAGER_URL == "" {
		panic(errors.New("Environment MODULE_MANAGER_URL is empty"))
	}

	STRAPI_URL = os.Getenv("STRAPI_URL")
	if STRAPI_URL == "" {
		panic(errors.New("Environment STRAPI_URL is empty"))
	}

	config.Schedule("", MODULE_MANAGER_URL)
}

func GetPersons(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, persons, err := Strapi.Find("person-table", U.Dict{})
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

	status, persons, err := Strapi.Create("person-table", U.Dict{
		"fullname":      ctx.Params["fullname"],
		"email":         ctx.Params["email"],
		"phone":         ctx.Params["phone"],
		"address":       ctx.Params["address"],
		"city":          ctx.Params["city"],
		"state":         ctx.Params["state"],
		"zip_code":      ctx.Params["zip_code"],
		"country":       ctx.Params["country"],
		"date_of_birth": ctx.Params["date_of_birth"],
		"is_active":     ctx.Params["is_active"],
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

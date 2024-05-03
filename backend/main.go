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
	"strconv"
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

func GetProducts(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, product, err := Strapi.Find("product-table", U.Dict{
		"pagination": U.Dict{
			"page": page,
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

	return 200, product, nil
}

func CreateProduct(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	product_name, ok := ctx.Params.GetString("product_name")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	description, ok := ctx.Params.GetString("description")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	price, ok := ctx.Params.GetFloat64("price")

	if !ok {
		return ctx.Error(400, "invalid parameters price")
	}

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	resp := U.Dict{}
	status, resp, err = Strapi.Create("product-table", U.Dict{
		"product_name": product_name,
		"description":  description,
		"price":        price,
		"is_active":    is_active,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, resp, nil
}

func GetProduct(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, product, err := Strapi.Find("product-table", U.Dict{
		"filters": U.Dict{
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

	return 200, product, nil
}

func UpdateProduct(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	pruduct_name, ok := ctx.Params.GetString("product_name")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	description, ok := ctx.Params.GetString("description")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	price, ok := ctx.Params.GetFloat64("price")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, _, err = Strapi.Update("product-table", int64(id), U.Dict{
		"product_name": pruduct_name,
		"description":  description,
		"price":        price,
		"is_active":    is_active,
	})

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

func DeleteProduct(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, _, err = Strapi.Delete("product-table", int64(id))
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

func CreateResource(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	resource_name, ok := ctx.Params.GetString("resource_name")

	if !ok {
		return ctx.Error(400, "invalid parameters resource_name")
	}

	description, ok := ctx.Params.GetString("description")

	if !ok {
		return ctx.Error(400, "invalid parameters description")
	}

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters is_active")
	}

	// resource_types, ok := ctx.Params.GetFloat64("resource_types")

	// if !ok {
	// 	return ctx.Error(400, "invalid parameters resource_types")
	// }

	image := ctx.Params["image"]

	if image == nil {
		image = nil
	}

	product, ok := ctx.Params.GetFloat64("product")

	if !ok {
		return ctx.Error(400, "invalid parameters product")
	}

	status, _, err = Strapi.Create("resource-table", U.Dict{
		"resource_name": resource_name,
		"description":   description,
		"is_active":     is_active,
		// "resource_types": resource_types,
		"image":   image,
		"product": product,
	})

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

func GetUsers(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, users, err := Strapi.Find("user-table", U.Dict{
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

	return 200, users, nil
}

func GetUser(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, user, err := Strapi.Find("user-table", U.Dict{
		"filter": U.Dict{
			"id": userId,
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

	return 200, user, nil
}

func UpdateUser(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	username, ok := ctx.Params.GetString("username")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	email, ok := ctx.Params.GetString("email")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	password, ok := ctx.Params.GetString("password")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, user, err := Strapi.Update("user-table", int64(id), U.Dict{
		"username": username,
		"email":    email,
		"password": password,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, user, nil
}

func DeleteUser(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, user, err := Strapi.Delete("user-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, user, nil
}

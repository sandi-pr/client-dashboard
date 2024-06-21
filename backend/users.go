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
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

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

	// page, ok := ctx.Params.GetFloat64("page")

	// if !ok {
	// 	fmt.Println("invalid parameters")
	// 	page = 1
	// }

	status, users, err := Strapi.UsersFind()

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	users_list := U.Dict{
		"data": users,
	}

	return 200, users_list, nil
}

func CreateUser(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	companyId, ok := me.GetFloat64("company")

	if !ok {
		return ctx.Error(400, "invalid parameters")
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

	fullname, ok := ctx.Params.GetString("fullname")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	users := U.Dict{}
	status, users, err = Strapi.UsersCreate(string(username), string(email), string(password), true, false, int64(1))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 201 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, user, err := Strapi.UsersUpdate(int64(users["id"].(U.Float64)), U.Dict{"fullname": fullname, "company": companyId})

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

	id, err := strconv.Atoi(string(userId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, user, err := Strapi.UsersFindOne(int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println("user", user)

	photo, ok := user.GetDict("photo")

	if !ok {
		return 200, U.Dict{"user": user, "photo": ""}, nil
	}

	photo_url := photo["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
	photo_url = U.String(strings.TrimPrefix(string(photo_url), "/"))

	resp, err := http.Get(STRAPI_URL + string(photo_url))
	if err != nil {
		fmt.Printf("error saat mengunduh file: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error saat membaca isi file: %v", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	fmt.Println(base64Data)

	return 200, U.Dict{"user": user, "photo": base64Data}, nil
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

	fullname, ok := ctx.Params.GetString("fullname")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	email, ok := ctx.Params.GetString("email")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	company, ok := ctx.Params.GetFloat64("company")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	updateData := U.Dict{
		"username":  username,
		"email":     email,
		"fullname":  fullname,
		"confirmed": true,
		"roles":     1,
		"company":   company,
	}

	password, ok := ctx.Params.GetString("password")

	if ok {
		updateData["password"] = password
	}

	status, user, err := Strapi.UsersUpdate(int64(id), updateData)

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

func UpdatePhotoUser(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	photoBase64, ok := ctx.Params.GetString("photo")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	photo, err := base64.StdEncoding.DecodeString(string(photoBase64))
	if err != nil {
		return ctx.Error(400, "failed to decode photo")
	}

	filename, ok := ctx.Params.GetString("filename")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(userId))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	status, upload, err := Strapi.Upload(string(filename), photo)

	fmt.Println("upload", upload)

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := upload.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	idPhoto, ok := data.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	status, user, err := Strapi.UsersUpdate(int64(id), U.Dict{
		"photo": idPhoto,
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

	status, user, err := Strapi.UsersDelete(int64(id))

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

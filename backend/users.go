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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	if checkUserId {
		foundMissingUserId := false

		for _, user := range users {
			createdAtStr := user.(U.Dict)["createdAt"].(U.String)

			idUser := user.(U.Dict)["id"].(U.Float64)

			_, ok := user.(U.Dict).GetString("user_id")

			if !ok {
				user_id, err := FillMissingUserIDsForUser(Strapi, string(createdAtStr), int64(idUser))
				status, _, err := Strapi.UsersUpdate(int64(idUser), U.Dict{"user_id": user_id})

				if err != nil {
					fmt.Println(err)
					return ctx.Error(500, "Internal Server Error")
				}

				if status != 200 {
					fmt.Println(status)
					return ctx.Error(500, "Internal Server Error")
				}

				foundMissingUserId = true
			}
		}

		if foundMissingUserId {
			status, users, err = Strapi.UsersFind()

			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}
		}

		checkUserId = false
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

	authority, ok := ctx.Params.GetString("authority")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	blocked, ok := ctx.Params.GetBool("blocked")

	if !ok {
		blocked = false
	}

	status, users, err := Strapi.UsersCreate(string(username), string(email), string(password), true, false, int64(1))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 201 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	user_id, err := GenerateID(Strapi)
	if err != nil {
		fmt.Println("Error:", err)
		return ctx.Error(429, err.Error())
	}

	fmt.Println("user_id", user_id)

	status, user, err := Strapi.UsersUpdate(int64(users["id"].(U.Float64)), U.Dict{"fullname": fullname, "company": companyId, "authority": authority, "blocked": blocked, "user_id": user_id})

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

	authority, ok := ctx.Params.GetString("authority")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	blocked, ok := ctx.Params.GetBool("blocked")

	if !ok {
		blocked = false
	}

	updateData := U.Dict{
		"username":  username,
		"email":     email,
		"fullname":  fullname,
		"confirmed": true,
		"roles":     1,
		"company":   company,
		"authority": authority,
		"blocked":   blocked,
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

func GetCompanyType(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	status, company, err := Strapi.Find("company-profile-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": companyId,
			},
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

	list, ok := company.GetList("data")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	if len(list) == 0 {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var base64Data string

	if dict["attributes"].(U.Dict)["company_logo"].(U.Dict)["data"] != nil {
		fmt.Println(dict["attributes"].(U.Dict)["company_logo"].(U.Dict)["data"])
		document_url := dict["attributes"].(U.Dict)["company_logo"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		document_url = U.String(strings.TrimPrefix(string(document_url), "/"))

		resp, err := http.Get(STRAPI_URL + string(document_url))
		if err != nil {
			fmt.Printf("error saat mengunduh file: %v", err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error saat membaca isi file: %v", err)
		}

		base64Data = base64.StdEncoding.EncodeToString(data)
	}

	return 200, U.Dict{"data": dict, "company_logo": base64Data}, nil
}

func UpdateCompanyType(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	company_size, ok := ctx.Params.GetString("company_size")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	company_industry, ok := ctx.Params.GetString("company_industry")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, company, err := Strapi.Update("company-profile-table", int64(companyId), U.Dict{
		"company_size":     company_size,
		"company_industry": company_industry,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, company, nil
}

func UpdateTheme(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	sidebar_title, ok := ctx.Params.GetString("sidebar_title")

	if !ok {
		sidebar_title = ""
	}

	theme_color, ok := ctx.Params.GetString("theme_color")

	if !ok {
		theme_color = ""
	}

	company_logo_str, ok := ctx.Params.GetString("company_logo")

	if !ok {
		company_logo_str = ""
	}

	var company_logo_id U.Float64

	company_logo_id = 0

	if company_logo_str != "" {
		var company_logo U.Dict

		err = json.Unmarshal([]byte(company_logo_str), &company_logo)
		if err != nil {
			return ctx.Error(400, "invalid parameters")
		}

		if company_logo_base64, ok := company_logo.GetString("base64Data"); ok {

			if !ok {
				return ctx.Error(400, "invalid parameters")
			}

			file_logo, err := base64.StdEncoding.DecodeString(string(company_logo_base64))
			if err != nil {
				return ctx.Error(400, "failed to decode photo")
			}

			file_name, ok := company_logo.GetString("name")

			if !ok {
				return ctx.Error(400, "invalid parameters")
			}

			status, upload, err := Strapi.Upload(string(file_name), file_logo)

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

			company_logo_id, ok = data.GetFloat64("id")

			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}
		} else {
			company_logo_str = ""
		}
	}

	updates := U.Dict{}
	if sidebar_title != "" {
		updates["sidebar_title"] = sidebar_title
	}
	if theme_color != "" {
		updates["theme_color"] = theme_color
	}
	if company_logo_id != 0 {
		updates["company_logo"] = int64(company_logo_id)
	}

	status, company, err := Strapi.Update("company-profile-table", int64(companyId), updates)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, company, nil
}

func GetId(token string) (int64, error) {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, token)
	status, me, err := Strapi.UsersMe()

	if err != nil {
		return -1, err
	}

	if status != 200 {
		return -1, errors.New("failed to get user")
	}

	if id, ok := me.GetFloat64("id"); ok {
		return int64(id), nil
	}

	return -1, errors.New("failed to get user id")
}

func GenerateID(strapi *strapi.StrapiClassic) (string, error) {
	_, users, err := strapi.UsersFind()
	if err != nil {
		return "", err
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	var latestOrder int
	for _, user := range users {
		createdAtStr := user.(U.Dict)["createdAt"].(U.String)

		createdAt, err := time.Parse(time.RFC3339, string(createdAtStr))
		if err != nil {
			return "", fmt.Errorf("failed to parse createdAt: %v", err)
		}

		if createdAt.Format("2006-01-02") == today {
			latestOrder++
		}
	}

	if latestOrder >= 99 {
		return "", fmt.Errorf("batas harian tercapai: hanya 99 ID yang diizinkan per hari")
	}

	return fmt.Sprintf("%s%02d", now.Format("02012006"), latestOrder+1), nil
}

func FillMissingUserIDsForUser(strapi *strapi.StrapiClassic, createdAtStr string, userId int64) (string, error) {
	_, users, err := strapi.UsersFind()
	if err != nil {
		return "", fmt.Errorf("failed to fetch users: %v", err)
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse createdAt: %v", err)
	}
	today := createdAt.Format("2006-01-02")

	var order int

	for _, user := range users {
		userCreatedAtStr := user.(U.Dict)["createdAt"].(U.String)
		userCreatedAt, err := time.Parse(time.RFC3339, string(userCreatedAtStr))
		if err != nil {
			return "", fmt.Errorf("failed to parse user createdAt: %v", err)
		}

		_, ok := user.(U.Dict).GetString("user_id")

		if !ok {
			continue
		}

		if userCreatedAt.Format("2006-01-02") == today {
			order++
		}
	}

	return fmt.Sprintf("%s%02d", createdAt.Format("02012006"), order+1), nil
}

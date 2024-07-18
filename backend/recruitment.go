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
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateRecruitment(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters fullname")
	}

	email, ok := ctx.Params.GetString("email")
	if !ok {
		return ctx.Error(400, "invalid parameters email")
	}

	phone, ok := ctx.Params.GetString("phone")
	if !ok {
		return ctx.Error(400, "invalid parameters phone")
	}

	address, ok := ctx.Params.GetString("address")
	if !ok {
		return ctx.Error(400, "invalid parameters address")
	}

	job_level, ok := ctx.Params.GetString("job_level")
	if !ok {
		return ctx.Error(400, "invalid parameters job_level")
	}

	job_position, ok := ctx.Params.GetString("job_position")
	if !ok {
		return ctx.Error(400, "invalid parameters job_position")
	}

	education_history, ok := ctx.Params.GetString("education_history")
	if !ok {
		return ctx.Error(400, "invalid parameters education_history")
	}

	json_skills, ok := ctx.Params.GetString("skills")
	if !ok {
		return ctx.Error(400, "invalid parameters skills")
	}

	var skills interface{}
	err = json.Unmarshal([]byte(json_skills), &skills)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	experience, ok := ctx.Params.GetFloat64("experience")
	if !ok {
		return ctx.Error(400, "invalid parameters experience")
	}

	resume_name, ok := ctx.Params.GetString("resume_name")
	if !ok {
		return ctx.Error(400, "invalid parameters resume_name")
	}

	resume, ok := ctx.Params.GetString("resume")
	if !ok {
		return ctx.Error(400, "invalid parameters resume")
	}

	document, err := base64.StdEncoding.DecodeString(string(resume))
	if err != nil {
		return ctx.Error(400, "failed to decode document")
	}

	status, upload, err := Strapi.Upload(string(resume_name), document)

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

	id, ok := data.GetFloat64("id")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	status, recruitment, err := Strapi.Create("recruitment-table", U.Dict{
		"fullname":          fullname,
		"email":             email,
		"phone":             phone,
		"address":           address,
		"job_level":         job_level,
		"job_position":      job_position,
		"education_history": education_history,
		"skills":            skills,
		"experience":        experience,
		"resume":            id,
	})

	fmt.Println("recruitment", recruitment)

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	status, checkPerson, err := Strapi.Find("person-table", U.Dict{
		"filters": U.Dict{
			"email": U.Dict{
				"$eq": email,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(checkPerson)
		return ctx.Error(500, "Internal Server Error")
	}

	data2, ok := checkPerson.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(data2) > 0 {
		person, ok := data2.GetDict(0)

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		idPerson, ok := person.GetFloat64("id")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		id, ok := recruitment["data"].(U.Dict).GetFloat64("id")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		status, updateRequitment, err := Strapi.Update("recruitment-table", int64(id), U.Dict{
			"person": idPerson,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(updateRequitment)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, recruitment, nil
	} else {
		fmt.Print("recruitment")
		fmt.Print(recruitment)
		id, ok := recruitment["data"].(U.Dict).GetFloat64("id")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		status, createPerson, err := Strapi.Create("person-table", U.Dict{
			"fullname":    fullname,
			"email":       email,
			"phone":       phone,
			"address":     address,
			"is_active":   "Yes",
			"recruitment": id,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(createPerson)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, recruitment, nil
	}
}

func GetRecruitments(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	m_authorization, ok := ctx.Headers.GetString("Authorization")
	if !ok {
		return ctx.Error(401, "Unauthorized")
	}

	authorization := strings.Replace(string(m_authorization), "Bearer ", "", 1)

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, authorization)

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

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		fmt.Println("invalid parameters")
		page_size = 10
	}

	status, recruitments, err := Strapi.Find("recruitment-table", U.Dict{
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
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

	return status, recruitments, nil
}

func GetRecruitment(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, document, err := Strapi.Find("recruitment-table", U.Dict{
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
		fmt.Println(document)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := document.GetList("data")
	if len(list) == 0 {
		return ctx.Error(500, "Internal Server Error")
	}

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	document_url := dict["attributes"].(U.Dict)["resume"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
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

	base64Data := base64.StdEncoding.EncodeToString(data)

	fmt.Println(base64Data)

	return 200, U.Dict{"recruitment": document, "base64": base64Data}, nil
}

func UpdateRecruitment(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	fullname, ok := ctx.Params.GetString("fullname")
	if !ok {
		return ctx.Error(400, "invalid parameters fullname")
	}

	email, ok := ctx.Params.GetString("email")
	if !ok {
		return ctx.Error(400, "invalid parameters email")
	}

	phone, ok := ctx.Params.GetString("phone")
	if !ok {
		return ctx.Error(400, "invalid parameters phone")
	}

	address, ok := ctx.Params.GetString("address")
	if !ok {
		return ctx.Error(400, "invalid parameters address")
	}

	job_level, ok := ctx.Params.GetString("job_level")
	if !ok {
		return ctx.Error(400, "invalid parameters job_level")
	}

	job_position, ok := ctx.Params.GetString("job_position")
	if !ok {
		return ctx.Error(400, "invalid parameters job_position")
	}

	education_history, ok := ctx.Params.GetString("education_history")
	if !ok {
		return ctx.Error(400, "invalid parameters education_history")
	}

	skills, ok := ctx.Params.GetString("skills")
	if !ok {
		return ctx.Error(400, "invalid parameters skills")
	}

	experience, ok := ctx.Params.GetFloat64("experience")
	if !ok {
		return ctx.Error(400, "invalid parameters experience")
	}

	status, recruitment, err := Strapi.Update("recruitment-table", int64(id), U.Dict{
		"fullname":          fullname,
		"email":             email,
		"phone":             phone,
		"address":           address,
		"job_level":         job_level,
		"job_position":      job_position,
		"education_history": education_history,
		"skills":            skills,
		"experience":        experience,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(recruitment)

	return 200, recruitment, nil
}

func DeleteRecruitment(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	idStr, ok := ctx.URLParams.GetString("id")
	if !ok {
		return ctx.Error(400, "invalid parameters id")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, recruitment, err := Strapi.Delete("recruitment-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(recruitment)

	return 200, recruitment, nil
}

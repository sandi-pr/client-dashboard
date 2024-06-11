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
		"job_position":      job_position,
		"education_history": education_history,
		"skills":            skills,
		"resume":            id,
	})

	fmt.Println("recruitment", recruitment)

	if err != nil {
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		return ctx.Error(500, "Internal Server Error")
	}

	return status, recruitment, nil
}

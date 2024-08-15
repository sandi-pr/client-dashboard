package main

import (
	"fmt"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func GetModelTrocr(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	status, company, err := Strapi.FindOne("company-profile-table", int64(companyId))

	if err != nil {
		fmt.Println(err)
		fmt.Println(company)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	// modelTrocr := company["data"].(U.Dict)["attributes"].(U.Dict)["model_trocr"].(U.List)

	data, ok := company.GetDict("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	attributes, ok := data.GetDict("attributes")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	modelTrocr, ok := attributes.GetString("model_trocr")

	if !ok {
		return 200, U.Dict{"model_trocr": nil}, nil
	}

	return 200, U.Dict{"model_trocr": modelTrocr}, nil
}

func UpdateModelTrocr(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	model_trocr, ok := ctx.Params.GetString("model_trocr")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	companyId, ok := me.GetFloat64("company")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, company, err := Strapi.Update("company-profile-table", int64(companyId), U.Dict{
		"model_trocr": model_trocr,
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println(company)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(company)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, company, nil
}

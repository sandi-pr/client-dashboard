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
	"os"

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
	if token, ok := ctx.Params.GetString("token"); ok {
		Strapi := strapi.NewStrapiWithToken(STRAPI_URL, string(token))

		status, data, err := Strapi.Find("/api/person-table", U.Dict{})
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		return status, data, nil
	}

	return ctx.Error(500, "Internal Server Error")
}

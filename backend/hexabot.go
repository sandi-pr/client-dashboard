package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateBot(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

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

	bot_name, ok := ctx.Params.GetString("bot_name")
	if !ok {
		return ctx.Error(400, "invalid parameters bot_name")
	}

	description, ok := ctx.Params.GetString("description")
	if !ok {
		description = ""
	}

	bot_profileStr, ok := ctx.Params.GetString("bot_profile")
	if !ok {
		bot_profileStr = ""
	}

	var bot_profile_id U.Float64

	bot_profile_id = 0

	if bot_profileStr != "" {
		var bot_profile U.Dict

		err = json.Unmarshal([]byte(bot_profileStr), &bot_profile)
		if err != nil {
			return ctx.Error(400, "invalid parameters")
		}

		if bot_profile_base64, ok := bot_profile.GetString("base64Data"); ok {

			if !ok {
				return ctx.Error(400, "invalid parameters")
			}

			file_logo, err := base64.StdEncoding.DecodeString(string(bot_profile_base64))
			if err != nil {
				return ctx.Error(400, "failed to decode photo")
			}

			file_name, ok := bot_profile.GetString("name")

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

			bot_profile_id, ok = data.GetFloat64("id")

			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}
		} else {
			bot_profileStr = ""
		}
	}

	ai_engine, ok := ctx.Params.GetString("ai_engine")

	if !ok {
		return ctx.Error(400, "invalid parameters ai_engine")
	}

	collection_id, ok := ctx.Params.GetFloat64("collection_id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	additional_prompt, ok := ctx.Params.GetString("additional_prompt")

	if !ok {
		additional_prompt = ""
	}

	var data U.Dict

	data = U.Dict{
		"bot_name":       bot_name,
		"bot_profile_id": int64(bot_profile_id),
		"ai_engine":      ai_engine,
		"collection_id":  int64(collection_id),
	}

	if description != "" {
		data["description"] = description
	}

	if bot_profileStr != "" {
		data["bot_profile"] = int64(bot_profile_id)
	}

	if additional_prompt != "" {
		data["additional_prompt"] = additional_prompt
	}

	status, bot, err := Strapi.Create("hexabot-table", data)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, bot, nil
}

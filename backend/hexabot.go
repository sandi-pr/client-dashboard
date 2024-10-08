package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

	bot_profileStr, ok := ctx.Params.GetString("bot_logo")
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

	need_auth, ok := ctx.Params.GetBool("need_auth")
	if !ok {
		need_auth = false
	}

	var data U.Dict

	data = U.Dict{
		"bot_id":           uuid.New().String(),
		"bot_name":         bot_name,
		"ai_engine":        ai_engine,
		"collection_id":    int64(collection_id),
		"background_align": "center",
		"color_theme":      "purple",
		"layout":           "Modern Minimalism",
		"overlay_color":    "purple",
		"username":         "User",
		"need_auth":        need_auth,
	}

	if description != "" {
		data["description"] = description
	}

	// if bot_profileStr != "" {
	// 	data["bot_logo"] = int64(bot_profile_id)
	// }

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

func GetBots(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")

	if !ok {
		page_size = 10
	}

	search, ok := ctx.Params.GetString("search")

	if !ok {
		search = ""
	}

	status, data, err := Strapi.Find("hexabot-table", U.Dict{
		"filters": U.Dict{
			"bot_name": U.Dict{
				"$contains": search,
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"populate": "*",
		"sort":     "updatedAt:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := data.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	meta, ok := data.GetDict("meta")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	var bots U.Dict

	bots = U.Dict{
		"data": U.List{},
		"meta": meta,
	}

	for i := 0; i < len(list); i++ {
		if list[i].(U.Dict)["attributes"].(U.Dict)["bot_profile"].(U.Dict)["data"] != nil {

			document_url := list[i].(U.Dict)["attributes"].(U.Dict)["bot_profile"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
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

			existingData, ok := bots.GetList("data")
			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}

			newData := U.Dict{
				"id":          list[i].(U.Dict)["id"],
				"attributes":  list[i].(U.Dict)["attributes"],
				"bot_profile": base64Data,
			}

			bots["data"] = append(existingData, newData)
		} else {
			existingData, ok := bots.GetList("data")
			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}

			bots["data"] = append(existingData, list[i])
		}
	}

	return 200, bots, nil
}

func GetBot(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(idStr))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, data, err := Strapi.Find("hexabot-table", U.Dict{
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

	list, ok := data.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if dict["attributes"].(U.Dict)["bot_logo"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["bot_logo"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["bot_logo"] = base64Data
	}

	if dict["attributes"].(U.Dict)["bot_profile"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["bot_profile"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["bot_profile"] = base64Data
	}

	if dict["attributes"].(U.Dict)["user_profile"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["user_profile"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["user_profile"] = base64Data
	}

	if dict["attributes"].(U.Dict)["background"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["background"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["background"] = base64Data
	}

	return 200, dict, nil
}

func GetBotPublic(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	id, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, data, err := Strapi.Find("hexabot-table", U.Dict{
		"filters": U.Dict{
			"bot_id": U.Dict{
				"$eq": id,
			},
		},
		"populate": "*",
	})

	if err != nil {
		fmt.Println(data)
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(data)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	list, ok := data.GetList("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if len(list) == 0 {
		return ctx.Error(404, "not found")
	}

	dict, ok := list.GetDict(0)

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if dict["attributes"].(U.Dict)["bot_logo"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["bot_logo"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["bot_logo"] = base64Data
	}

	if dict["attributes"].(U.Dict)["bot_profile"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["bot_profile"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["bot_profile"] = base64Data
	}

	if dict["attributes"].(U.Dict)["user_profile"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["user_profile"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["user_profile"] = base64Data
	}

	if dict["attributes"].(U.Dict)["background"].(U.Dict)["data"] != nil {
		document_url := dict["attributes"].(U.Dict)["background"].(U.Dict)["data"].(U.Dict)["attributes"].(U.Dict)["url"].(U.String)
		base64Data, err := getFileByUrl(document_url)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		dict["background"] = base64Data
	}

	return 200, dict, nil
}

func UpdateBot(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	additional_prompt, ok := ctx.Params.GetString("additional_prompt")

	if !ok {
		additional_prompt = ""
	}

	ai_engine, ok := ctx.Params.GetString("ai_engine")

	if !ok {
		return ctx.Error(400, "invalid parameters ai_engine")
	}

	backgroundStr, ok := ctx.Params.GetString("background")

	if !ok {
		backgroundStr = ""
	}

	background_align, ok := ctx.Params.GetString("background_align")

	if !ok {
		background_align = ""
	}

	bot_logoStr, ok := ctx.Params.GetString("bot_logo")

	if !ok {
		bot_logoStr = ""
	}

	bot_name, ok := ctx.Params.GetString("bot_name")

	if !ok {
		return ctx.Error(400, "invalid parameters bot_name")
	}

	bot_profileStr, ok := ctx.Params.GetString("bot_profile")

	if !ok {
		bot_profileStr = ""
	}

	color_theme, ok := ctx.Params.GetString("color_theme")

	if !ok {
		return ctx.Error(400, "invalid parameters color_theme")
	}

	description, ok := ctx.Params.GetString("description")

	if !ok {
		description = ""
	}

	greeting_text, ok := ctx.Params.GetString("greeting_text")

	if !ok {
		greeting_text = ""
	}

	layout, ok := ctx.Params.GetString("layout")

	if !ok {
		return ctx.Error(400, "invalid parameters layout")
	}

	logo_name, ok := ctx.Params.GetString("logo_name")

	if !ok {
		logo_name = ""
	}

	logo_only, ok := ctx.Params.GetBool("logo_only")

	if !ok {
		logo_only = false
	}

	overlay_color, ok := ctx.Params.GetString("overlay_color")

	if !ok {
		return ctx.Error(400, "invalid parameters overlay_color")
	}

	user_profileStr, ok := ctx.Params.GetString("user_profile")

	if !ok {
		user_profileStr = ""
	}

	deployed_platform, ok := ctx.Params.GetString("deployed_platform")

	if !ok {
		return ctx.Error(400, "invalid parameters deployed_platform")
	}

	username, ok := ctx.Params.GetString("username")

	if !ok {
		return ctx.Error(400, "invalid parameters username")
	}

	need_auth, ok := ctx.Params.GetBool("need_auth")

	if !ok {
		need_auth = false
	}

	var background_id int64

	if backgroundStr != "" {
		var background U.Dict

		err := json.Unmarshal([]byte(backgroundStr), &background)
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		background_id, err = uploadFile(background, Strapi)
		if err != nil {
			backgroundStr = ""
		}
	}

	var bot_logo_id int64

	if bot_logoStr != "" {
		var bot_logo U.Dict

		err := json.Unmarshal([]byte(bot_logoStr), &bot_logo)
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		bot_logo_id, err = uploadFile(bot_logo, Strapi)
		if err != nil {
			bot_logoStr = ""
		}
	}

	var bot_profile_id int64

	if bot_profileStr != "" {
		var bot_profile U.Dict

		err := json.Unmarshal([]byte(bot_profileStr), &bot_profile)
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		bot_profile_id, err = uploadFile(bot_profile, Strapi)
		if err != nil {
			bot_profileStr = ""
		}
	}

	var user_profile_id int64

	if user_profileStr != "" {
		var user_profile U.Dict

		err := json.Unmarshal([]byte(user_profileStr), &user_profile)
		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		user_profile_id, err = uploadFile(user_profile, Strapi)
		if err != nil {
			user_profileStr = ""
		}
	}

	var data U.Dict

	data = U.Dict{
		"ai_engine":         ai_engine,
		"bot_name":          bot_name,
		"color_theme":       color_theme,
		"layout":            layout,
		"logo_only":         logo_only,
		"overlay_color":     overlay_color,
		"deployed_platform": deployed_platform,
		"username":          username,
		"need_auth":         need_auth,
	}

	if backgroundStr != "" {
		data["background"] = background_id
	}

	if bot_logoStr != "" {
		data["bot_logo"] = bot_logo_id
	}

	if bot_profileStr != "" {
		data["bot_profile"] = bot_profile_id
	}

	if user_profileStr != "" {
		data["user_profile"] = user_profile_id
	}

	if additional_prompt != "" {
		data["additional_prompt"] = additional_prompt
	}

	if background_align != "" {
		data["background_align"] = background_align
	}

	if description != "" {
		data["description"] = description
	}

	if greeting_text != "" {
		data["greeting_text"] = greeting_text
	}

	if logo_name != "" {
		data["logo_name"] = logo_name
	}

	status, update, err := Strapi.Update("hexabot-table", int64(id), data)
	if err != nil {
		fmt.Println(err)
		fmt.Println(update)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func DeleteBot(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return ctx.Error(400, "invalid parameters id")
	}

	status, _, err = Strapi.Delete("hexabot-table", int64(id))
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

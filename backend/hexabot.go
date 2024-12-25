package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.hexacode.org/go-libs/cr3call"
	"gitlab.hexacode.org/go-libs/hctypes"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/requests"
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

	library_id, ok := ctx.Params.GetFloat64("library_id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	collection_id, ok := ctx.Params.GetString("collection_id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	if collection_id == "all" {
		collection_id, err = getAllIdCollection(int(library_id), Strapi)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		collection_id_array := strings.Split(string(collection_id), ",")
		collection_id_int := make([]int, len(collection_id_array))

		for i, str := range collection_id_array {
			num, err := strconv.Atoi(strings.TrimSpace(str))
			if err != nil {
				return ctx.Error(500, "Internal Server Error")
			}
			collection_id_int[i] = num
		}

		var allIdsSub U.String

		for i := 0; i < len(collection_id_int); i++ {
			collection_id_sub, err := getAllIdCollection(int(collection_id_int[i]), Strapi)
			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if collection_id_sub != "" {
				allIdsSub = allIdsSub + "," + collection_id_sub
			}
		}

		if allIdsSub != "" {
			collection_id = collection_id + "," + allIdsSub
		}
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

	bot_id := uuid.New().String()

	data = U.Dict{
		"bot_id":           bot_id,
		"bot_name":         bot_name,
		"ai_engine":        ai_engine,
		"library_id":       int64(library_id),
		"collection_id":    collection_id,
		"background_align": "center",
		"color_theme":      "purple",
		"layout":           "Modern Minimalism",
		"overlay_color":    "purple",
		"username":         "User",
		"need_auth":        need_auth,
		"status":           "Inactive",
		"manage_document":  false,
	}

	if description != "" {
		data["description"] = description
	}

	// if bot_profileStr != "" {
	// 	data["bot_logo"] = int64(bot_profile_id)
	// }

	if bot_profileStr != "" {
		data["bot_profile"] = bot_profile_id
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

		if err2, ok := bot.GetDict("error"); ok {
			if message, ok := err2.GetString("message"); ok {
				return ctx.Error(status, string(message))
			}
		}

		return ctx.Error(500, "Internal Server Error")
	}

	http := requests.NewHttp(hctypes.Dict{
		"Content-type": "application/json",
	})

	resp, err := http.Post(fmt.Sprintf("%s/api/%s", HEXA_RAG_URL, "create-knowledge"), nil, hctypes.Dict{
		"knowledge_id":  bot_id,
		"system_prompt": additional_prompt,
	}.ToJson())

	fmt.Println("create-knowledge", resp)

	if err != nil {
		fmt.Println(err)
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

			// botStatus, _, err := getStatusBot(int64(list[i].(U.Dict)["id"].(U.Float64)), Strapi)

			// fmt.Println("botStatus", botStatus, list[i].(U.Dict)["id"].(U.Float64))

			// if err != nil {
			// 	fmt.Println(err)
			// 	return ctx.Error(500, "Internal Server Error")
			// }

			// newData["attributes"].(U.Dict)["status"] = botStatus

			bots["data"] = append(existingData, newData)
		} else {
			existingData, ok := bots.GetList("data")
			if !ok {
				return ctx.Error(500, "Internal Server Error")
			}

			// botStatus, _, err := getStatusBot(int64(list[i].(U.Dict)["id"].(U.Float64)), Strapi)

			// fmt.Println("botStatus", botStatus, list[i].(U.Dict)["id"].(U.Float64))

			// if err != nil {
			// 	fmt.Println(err)
			// 	return ctx.Error(500, "Internal Server Error")
			// }

			// list[i].(U.Dict)["attributes"].(U.Dict)["status"] = botStatus

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

	// botStatus, err := getStatusBot(int64(id), Strapi)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// dict["attributes"].(U.Dict)["status"] = botStatus

	return 200, dict, nil
}

func GetBotByName(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	name, ok := ctx.Params.GetString("name")

	if !ok {
		return ctx.Error(400, "invalid parameters name")
	}

	status, data, err := Strapi.Find("hexabot-table", U.Dict{
		"filters": U.Dict{
			"bot_name": U.Dict{
				"$eq": name,
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

	// botStatus, err := getStatusBot(int64(id), Strapi)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// dict["attributes"].(U.Dict)["status"] = botStatus

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

	if attributes, ok := dict.GetDict("attributes"); ok {
		if collection_id, ok := attributes.GetFloat64("collection_id"); ok {
			status, collection, err := Strapi.FindOne("document-collection-table", int64(collection_id))

			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if status != 200 {
				fmt.Println(status)
				return ctx.Error(500, "Internal Server Error")
			}

			if data, ok := collection.GetDict("data"); ok {
				if attributes, ok := data.GetDict("attributes"); ok {
					if categories, ok := attributes.GetString("categories"); ok {
						dict["categories"] = categories
					}
				}
			}
		}
	}

	// botStatus, err := getStatusBot(int64(dict["id"].(float64)), Strapi)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return ctx.Error(500, "Internal Server Error")
	// }

	// dict["attributes"].(U.Dict)["status"] = botStatus

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

	library_id, ok := ctx.Params.GetFloat64("library_id")
	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	collection_id, ok := ctx.Params.GetString("collection_id")

	if !ok {
		return ctx.Error(400, "invalid parameters collection_id")
	}

	if collection_id == "all" {
		collection_id, err = getAllIdCollection(int(library_id), Strapi)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		collection_id_array := strings.Split(string(collection_id), ",")
		collection_id_int := make([]int, len(collection_id_array))

		for i, str := range collection_id_array {
			num, err := strconv.Atoi(strings.TrimSpace(str))
			if err != nil {
				return ctx.Error(500, "Internal Server Error")
			}
			collection_id_int[i] = num
		}

		var allIdsSub U.String

		for i := 0; i < len(collection_id_int); i++ {
			collection_id_sub, err := getAllIdCollection(int(collection_id_int[i]), Strapi)
			if err != nil {
				fmt.Println(err)
				return ctx.Error(500, "Internal Server Error")
			}

			if collection_id_sub != "" {
				allIdsSub = allIdsSub + "," + collection_id_sub
			}
		}

		if allIdsSub != "" {
			collection_id = collection_id + "," + allIdsSub
		}
	}

	manage_document, ok := ctx.URLParams.GetBool("manage_document")

	if !ok {
		manage_document = false
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
		"collection_id":     collection_id,
		"manage_document":   manage_document,
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
		fmt.Println(update)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, update, nil
}

func UpdateBotStatus(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	statusBot, ok := ctx.Params.GetString("status")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, update, err := Strapi.Update("hexabot-table", int64(id), U.Dict{
		"status": statusBot,
	})

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

func GetCountBotPerMount(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	earlyDate, ok := ctx.Params.GetString("earlyDate")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	lastDate, ok := ctx.Params.GetString("lastDate")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, bots, err := Strapi.Find("hexabot-table", U.Dict{
		"filters": U.Dict{
			"updatedAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	fmt.Println(bots)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, bots, nil
}

func getAllIdCollection(library_id int, strapi *strapi.StrapiClassic) (U.String, error) {
	var allIds U.List
	page := 1

	for {
		status, documentLibrarys, err := strapi.Find("document-library-table", U.Dict{
			"filters": U.Dict{
				"parent_id": U.Dict{
					"$eq": int64(library_id),
				},
				"type": U.Dict{
					"$eq": "collection",
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
			},
		})

		if err != nil {
			fmt.Println(err)
			return "", err
		}

		if status != 200 {
			fmt.Println(status)
			return "", err
		}

		if list, ok := documentLibrarys.GetList("data"); ok {
			if len(list) < 1 {
				break
			}

			for i := 0; i < len(list); i++ {
				if dict, ok := list.GetDict(i); ok {
					if id, ok := dict.GetFloat64("id"); ok {
						allIds = append(allIds, id)
						page = page + 1
					} else {
						return "", err
					}
				} else {
					return "", err
				}
			}
		} else {

			return "", err
		}
	}

	var strNums []string
	for i := 0; i < len(allIds); i++ {
		if num, ok := allIds.GetFloat64(i); ok {
			strNums = append(strNums, strconv.Itoa(int(num)))
		}
	}

	var result U.String
	result = U.String(strings.Join(strNums, ","))

	return result, nil
}

func getStatusBot(bot_id int64, strapi *strapi.StrapiClassic) (string, string, error) {
	status, documentLibrarys, err := strapi.Find("hexabot-table", U.Dict{
		"filters": U.Dict{
			"id": U.Dict{
				"$eq": bot_id,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	if status != 200 {
		fmt.Println(status)
		return "", "", err
	}

	if list, ok := documentLibrarys.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			if attributes, ok := dict.GetDict("attributes"); ok {
				bot_name, ok := attributes.GetString("bot_name")

				if !ok {
					fmt.Println("ayam 1: Gagal mendapatkan bot_name")
					return "", "", err
				}

				if collection, ok := attributes.GetString("collection_id"); ok {
					collection_id_array := strings.Split(string(collection), ",")
					collection_id_int := make([]int, len(collection_id_array))

					for i, str := range collection_id_array {
						num, err := strconv.Atoi(strings.TrimSpace(str))
						if err != nil {
							return "", "", err
						}
						collection_id_int[i] = num
					}

					for i := 0; i < len(collection_id_int); i++ {
						page := 1
						for {
							status, documentLibrarys, err := strapi.Find("document-library-table", U.Dict{
								"filters": U.Dict{
									"parent_id": U.Dict{
										"$eq": int64(collection_id_int[i]),
									},
									"type": U.Dict{
										"$eq": "document",
									},
								},
								"pagination": U.Dict{
									"page":     page,
									"pageSize": 25,
								},
							})

							if err != nil {
								fmt.Println(err)
								return "", "", err
							}

							if status != 200 {
								fmt.Println(status)
								return "", "", err
							}

							if list, ok := documentLibrarys.GetList("data"); ok {
								if len(list) < 1 {
									// return "Ready", string(bot_name), nil
									break
								}

								for j := 0; j < len(list); j++ {
									if dict, ok := list.GetDict(j); ok {
										// idDoc, ok := dict.GetFloat64("id")

										// if !ok {
										// 	fmt.Println("failed get id")
										// }

										if attributes, ok := dict.GetDict("attributes"); ok {
											if status, ok := attributes.GetString("status"); ok {
												if status != "Ready" && status != "Error" {
													// fmt.Println("not ready", collection_id_int[i], idDoc)
													return string(status), string(bot_name), nil
												} else {
													// fmt.Println("ready", collection_id_int[i], idDoc)
													page = page + 1
													continue
												}
											} else {
												return "Inactive", string(bot_name), nil
											}
										} else {
											return "", "", err
										}
									} else {
										return "", "", err
									}
								}
							} else {
								return "", "", err
							}
						}
					}
					return "Ready", string(bot_name), nil
				} else {
					return "", "", err
				}
			} else {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	} else {
		return "", "", err
	}

	return "", "", err
}

func AddDocumentByBotID(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	userId, err := GetId(strings.Replace(string(bearer), "Bearer ", "", 1))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	userIdStr := strconv.Itoa(int(userId))

	idStr, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	bot_id, err := strconv.Atoi(string(idStr))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	fmt.Println("idStr", idStr)

	CR3 := cr3call.New(MODULE_MANAGER_URL)
	go CR3.Call("hexa-ai-rag", "add-document-by-bot-id-params", cr3call.CR3CallData{
		Headers: U.Dict{
			"Authorization": bearer,
		},
		Params: U.Dict{
			"id": idStr,
		},
	}, false)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	go getStatusBotNew(int64(bot_id), Strapi, userIdStr)

	return 200, U.Dict{
		"success": true,
	}, nil
}

func getStatusBotNew(bot_id int64, strapi *strapi.StrapiClassic, userId string) {

	tmpBotStatus := ""

	for {
		botStatus, bot_name, err := getStatusBot(int64(bot_id), strapi)

		if err != nil {
			fmt.Println(err)
			break
		}

		if botStatus != tmpBotStatus {
			var msg NotificationMessage
			if botStatus == "Ready" {
				msg = NotificationMessage{
					UserID:  "0",
					Title:   string(bot_name) + " documents are ready",
					Content: "Hexacube has finished analyzing the uploaded documents. You can now proceed to explore the insights!",
					Time:    time.Now().String(),
					Type:    "hexabot",
					Meta: U.Dict{
						"bot_id": bot_id,
					},
				}

				notificationWS.SendNotification(msg)
				fmt.Println("send notification")

				status, update, err := strapi.Update("hexabot-table", int64(bot_id), U.Dict{
					"status": "Ready",
				})

				fmt.Println("update hexabot", update)

				if err != nil {
					fmt.Println(err)
					break
				}

				if status != 200 {
					fmt.Println("failed update hexabot-table")
					break
				}

				tmpBotStatus = ""
				break

			} else if botStatus == "Reading" {
				// msg = NotificationMessage{
				// 	UserID:  "0",
				// 	Title:   "Reading " + string(bot_name) + " documents",
				// 	Content: "Hexacube is now analyzing the documents uploaded by HexaBot. Sit tight while we process the information!",
				// 	Time:    time.Now().String(),
				// 	Type:    "hexabot",
				// 	Meta: U.Dict{
				// 		"bot_id": bot_id,
				// 	},
				// }

				// notificationWS.SendNotification(msg)
				// fmt.Println("send notification")

				status, _, err := strapi.Update("hexabot-table", int64(bot_id), U.Dict{
					"status": "Reading",
				})

				if err != nil {
					fmt.Println(err)
					break
				}

				if status != 200 {
					fmt.Println("failed update hexabot-table")
					break
				}
			}
		}

		tmpBotStatus = botStatus
	}
}

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
	"fmt"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func Notification(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	notificationWS.ServeNotifications(ctx, ctx.Writer, ctx.Request)
	return 0, nil, nil
}

func GetNotifications(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	// id, ok := me.GetFloat64("id")

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

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

	status, notifications, err := Strapi.Find("notification-table", U.Dict{
		// "filters": U.Dict{
		// 	"user_id": U.Dict{
		// 		"$eq": id,
		// 	},
		// },
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort": "updatedAt:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, notifications, nil
}

func GetNotificationsUnread(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	// id, ok := me.GetFloat64("id")

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

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

	status, notifications, err := Strapi.Find("notification-table", U.Dict{
		"filters": U.Dict{
			// "user_id": U.Dict{
			// 	"$eq": id,
			// },
			"read": U.Dict{
				"$eq": false,
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort": "updatedAt:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, notifications, nil
}

func GetNotificationsRead(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	// id, ok := me.GetFloat64("id")

	// if !ok {
	// 	return ctx.Error(500, "Internal Server Error")
	// }

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

	status, notifications, err := Strapi.Find("notification-table", U.Dict{
		"filters": U.Dict{
			// "user_id": U.Dict{
			// 	"$eq": id,
			// },
			"read": U.Dict{
				"$eq": true,
			},
		},
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort": "updatedAt:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, notifications, nil
}

func ReadNotification(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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
		return ctx.Error(500, "Internal Server Error")
	}

	status, notifications, err := Strapi.Update("notification-table", int64(id), U.Dict{
		"read": true,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, notifications, nil
}

func ReadAllNotification(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	// id, ok := me.GetFloat64("id")

	// if !ok {
	// 	return ctx.Error(400, "invalid parameters id")
	// }

	page := 1

	for {
		status, notifications, err := Strapi.Find("notification-table", U.Dict{
			"filters": U.Dict{
				// "user_id": U.Dict{
				// 	"$eq": id,
				// },
				"read": U.Dict{
					"$eq": false,
				},
			},
			"pagination": U.Dict{
				"page":     page,
				"pageSize": 25,
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

		list, ok := notifications.GetList("data")

		if !ok {
			return ctx.Error(500, "Internal Server Error")
		}

		if len(list) == 0 {
			break
		}

		for i := 0; i < len(list); i++ {
			if dict, ok := list.GetDict(i); ok {
				if idNotif, ok := dict.GetFloat64("id"); ok {
					status, _, err := Strapi.Update("notification-table", int64(idNotif), U.Dict{
						"read": true,
					})

					if err != nil {
						fmt.Println(err)
						return ctx.Error(500, "Internal Server Error")
					}

					if status != 200 {
						fmt.Println(status)
						return ctx.Error(500, "Internal Server Error")
					}

					page = page + 1
				}
			}
		}
	}

	return 200, U.Dict{}, nil
}

func UpdateContentProgressByName(title string, content string, meta U.Dict) error {
	strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	// Step 1: Fetch data from Strapi
	status, data, err := strapi.Find("notification-table", U.Dict{
		"filters": U.Dict{
			"title": U.Dict{
				"$eq": title,
			},
			"content": U.Dict{
				"$eq": content,
			},
		},
	})

	if err != nil {
		fmt.Println("Error during Strapi Find:", err)
		return err
	}

	if status != 200 {
		fmt.Printf("Error: Unexpected status code from Strapi Find: %d\n", status)
		return fmt.Errorf("unexpected status code: %d", status)
	}

	// Step 2: Extract notification ID
	var idNotif int64
	if list, ok := data.GetList("data"); ok {
		if dict, ok := list.GetDict(0); ok {
			if id, ok := dict.GetFloat64("id"); ok {
				idNotif = int64(id)
			} else {
				fmt.Println("Error: Failed to extract ID from dictionary")
				return fmt.Errorf("failed to extract ID")
			}
		} else {
			fmt.Println("Error: Data list is empty or not a dictionary")
			return fmt.Errorf("data list is empty or invalid")
		}
	} else {
		fmt.Println("Error: Failed to retrieve data list from response")
		return fmt.Errorf("failed to retrieve data list")
	}

	// Step 3: Update content in Strapi
	status2, _, err := strapi.Update("notification-table", idNotif, U.Dict{
		"meta": meta,
	})

	if err != nil {
		fmt.Println("Error during Strapi Update:", err)
		return err
	}

	if status2 != 200 {
		fmt.Printf("Error: Unexpected status code from Strapi Update: %d\n", status2)
		return fmt.Errorf("unexpected status code: %d", status2)
	}

	return nil
}

func CheckNotification(title string, content string) (bool, error) {
	strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	status, data, err := strapi.Find("notification-table", U.Dict{
		"filters": U.Dict{
			"title": U.Dict{
				"$eq": title,
			},
			"content": U.Dict{
				"$eq": content,
			},
		},
	})

	if err != nil {
		fmt.Println("Error during Strapi Find:", err)
		return false, err
	}

	if status != 200 {
		fmt.Printf("Error: Unexpected status code from Strapi Find: %d\n", status)
		return false, fmt.Errorf("unexpected status code: %d", status)
	}

	if list, ok := data.GetList("data"); ok {
		if len(list) > 0 {
			return true, nil
		}
	}

	return false, nil
}

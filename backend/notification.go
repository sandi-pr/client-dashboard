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
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
)

func Notification(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	notificationWS.ServeNotifications(ctx, ctx.Writer, ctx.Request)
	return 0, nil, nil
}

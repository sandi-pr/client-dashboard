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
	"gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
)

func init() {

}

func multiply(a int, b int) int {
	return a * b
}

func Multiply(ctx *config.ConfigContext) (int, hctypes.Dict, *config.ReturnHTTPError) {
	a, ok_a := ctx.Params.GetFloat64("a")
	b, ok_b := ctx.Params.GetFloat64("b")

	if !ok_a || !ok_b {
		return ctx.Error(400, "invalid parameters")
	}

	return 200, hctypes.Dict{
		"result": multiply(int(a), int(b)),
	}, nil
}

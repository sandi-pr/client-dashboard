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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func AddDataHandwriting(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearerOk := ctx.Headers.GetString("Authorization")
	if !bearerOk {
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

	imagesString, ok := ctx.Params.GetString("images")
	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	var images interface{}
	err = json.Unmarshal([]byte(imagesString), &images)
	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	for _, document := range images.([]interface{}) {
		base64Str, ok := document.(map[string]interface{})["base64"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters")
		}

		// Step 1: Decode the base64 string
		imageData, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			fmt.Println("Failed to decode base64 string:", err)
			return ctx.Error(400, "invalid parameters")
		}

		// Step 2: Decode the image data
		img, format, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			fmt.Println("Failed to decode image data:", err)
			return ctx.Error(400, "invalid parameters")
		}

		// Step 3: Convert the image to grayscale
		grayImg := image.NewGray(img.Bounds())
		draw.Draw(grayImg, img.Bounds(), img, image.Point{0, 0}, draw.Src)

		// Step 4: Encode the grayscale image back to base64
		var buf bytes.Buffer
		switch strings.ToLower(format) {
		case "jpeg", "jpg":
			err = jpeg.Encode(&buf, grayImg, nil)
		case "png":
			err = png.Encode(&buf, grayImg)
		default:
			fmt.Println("Unsupported image format:", format)
			return ctx.Error(400, "unsupported image format")
		}
		if err != nil {
			fmt.Println("Failed to encode grayscale image:", err)
			return ctx.Error(500, "Internal Server Error")
		}

		grayBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
		fmt.Println("Grayscale image in base64:", grayBase64Str)

		imageFile, err := base64.StdEncoding.DecodeString(grayBase64Str)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(400, "invalid parameters documents")
		}

		fileName, ok := document.(map[string]interface{})["name"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		source_text, ok := document.(map[string]interface{})["sourceText"].(string)
		if !ok {
			return ctx.Error(400, "invalid parameters documents")
		}

		status, upload, err := Strapi.Upload(fileName, imageFile)

		fmt.Println("upload", upload)

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		data, ok := upload.GetDict(0)

		if !ok {
			fmt.Println("failed to get dict")
			return ctx.Error(500, "Internal Server Error")
		}

		id, ok := data.GetFloat64("id")
		if !ok {
			fmt.Println("failed to get id")
			return ctx.Error(500, "Internal Server Error")
		}

		status, dataHandwriting, err := Strapi.Create("data-handwriting-table", U.Dict{
			"name":        fileName,
			"image":       id,
			"source_text": source_text,
		})

		fmt.Println("data handwriting", dataHandwriting)

		if err != nil {
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			return ctx.Error(500, "Internal Server Error")
		}

	}

	return 200, U.Dict{"status": "success"}, nil
}

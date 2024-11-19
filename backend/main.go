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
	"errors"
	"os"

	"gitlab.hexacode.org/go-libs/microservice/config"
)

var MODULE_MANAGER_URL = ""
var STRAPI_URL = ""
var STRAPI_TOKEN = ""
var DOCUMENT_PATH = ""
var HEXA_RAG_URL = ""

func init() {
	MODULE_MANAGER_URL = os.Getenv("MODULE_MANAGER_URL")
	if MODULE_MANAGER_URL == "" {
		panic(errors.New("Environment MODULE_MANAGER_URL is empty"))
	}

	STRAPI_URL = os.Getenv("STRAPI_URL")
	if STRAPI_URL == "" {
		panic(errors.New("Environment STRAPI_URL is empty"))
	}

	STRAPI_TOKEN = os.Getenv("STRAPI_TOKEN")
	if STRAPI_TOKEN == "" {
		panic(errors.New("Environment STRAPI_TOKEN is empty"))
	}

	DOCUMENT_PATH = os.Getenv("DOCUMENT_PATH")
	if DOCUMENT_PATH == "" {
		panic(errors.New("Environment DOCUMENT_PATH is empty"))
	}

	HEXA_RAG_URL = os.Getenv("HEXA_RAG_URL")
	if HEXA_RAG_URL == "" {
		panic(errors.New("Environment HEXA_RAG_URL is empty"))
	}

	config.Schedule("", MODULE_MANAGER_URL)
}

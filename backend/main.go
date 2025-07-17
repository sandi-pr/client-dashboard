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

	"gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
)

var MODULE_MANAGER_URL = ""
var STRAPI_URL = ""
var STRAPI_TOKEN = ""
var DOCUMENT_PATH = ""
var DATASETS_PATH = ""
var USAGE_MEMORY_PATH = ""
var USAGE_URL = ""
var HEXA_RAG_URL = ""
var SHUTDOWN_HOST = ""
var SHUTDOWN_TOKEN = ""
var botIdProccess = hctypes.List{}
var notificationWS *NotificationWebSocket
var usageWS *UsageWebSocket

var checkUserId = true
var metaDataCSV = hctypes.Dict{}

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

	DATASETS_PATH = os.Getenv("DATASETS_PATH")
	if DATASETS_PATH == "" {
		panic(errors.New("Environment DATASETS_PATH is empty"))
	}

	USAGE_MEMORY_PATH = os.Getenv("USAGE_MEMORY_PATH")
	if USAGE_MEMORY_PATH == "" {
		panic(errors.New("Environment USAGE_MEMORY_PATH is empty"))
	}

	USAGE_URL = os.Getenv("USAGE_URL")
	if USAGE_URL == "" {
		panic(errors.New("Environment USAGE_URL is empty"))
	}

	SHUTDOWN_HOST = os.Getenv("SHUTDOWN_HOST")
	if SHUTDOWN_HOST == "" {
		panic(errors.New("Environment SHUTDOWN_HOST is empty"))
	}

	SHUTDOWN_TOKEN = os.Getenv("SHUTDOWN_TOKEN")
	if SHUTDOWN_TOKEN == "" {
		panic(errors.New("Environment SHUTDOWN_TOKEN is empty"))
	}

	config.Schedule("", MODULE_MANAGER_URL)

	notificationWS = NewNotificationWebSocket()
	usageWS = NewUsageWebSocket()

	go notificationWS.HandleMessages()
}

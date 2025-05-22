package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/strapi"
)

func getFileByUrl(url U.String) (string, error) {
	if strings.HasPrefix(string(url), "/") {
		url = U.String(strings.TrimPrefix(string(url), "/"))
	}
	resp, err := http.Get(STRAPI_URL + string(url))
	if err != nil {
		fmt.Printf("error saat mengunduh file: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error saat membaca isi file: %v", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	return base64Data, nil
}

func uploadFile(file U.Dict, Strapi *strapi.StrapiClassic) (int64, error) {
	if base64Data, ok := file.GetString("base64Data"); ok {

		fmt.Println("base64Data", base64Data)

		file_logo, err := base64.StdEncoding.DecodeString(string(base64Data))
		if err != nil {
			fmt.Println("error saat mengunduh file", err)
			return 0, fmt.Errorf("invalid parameters")
		}

		file_name, ok := file.GetString("name")

		if !ok {
			fmt.Println("error saat mengambil nama", err)
			return 0, fmt.Errorf("invalid parameters")
		}

		status, upload, err := Strapi.Upload(string(file_name), file_logo)

		fmt.Println("upload", upload)

		if err != nil {
			fmt.Println("error saat mengupload file", err)
			return 0, fmt.Errorf("internal server error")
		}

		if status != 200 {
			fmt.Println("error status", err)
			return 0, fmt.Errorf("internal server error")
		}

		data, ok := upload.GetDict(0)

		if !ok {
			return 0, fmt.Errorf("internal server error")
		}

		id, ok := data.GetFloat64("id")

		if !ok {
			return 0, fmt.Errorf("internal server error")
		}

		return int64(id), nil
	} else {
		return 0, fmt.Errorf("internal server error")
	}
}

func formatFileSize(size int) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func convertBytesToMB(b string) string {
	val, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return "-"
	}
	return fmt.Sprintf("%.0f", val/1024/1024)
}

func detectDiskTypeFromName(name string) string {
	if strings.HasPrefix(name, "nvme") {
		return "SSD"
	} else if strings.HasPrefix(name, "sd") {
		return "HDD"
	} else if strings.HasPrefix(name, "vd") {
		return "Virtual"
	}
	return "Unknown"
}

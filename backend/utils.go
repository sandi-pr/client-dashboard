package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.hexacode.org/go-libs/hctypes"
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

func StrapiFindAll(Strapi *strapi.StrapiClassic, collection string, pageSize int) (hctypes.List, error) {
	allItems := hctypes.List{}

	page := 1

	for {
		status, resp, err := Strapi.Find(collection, U.Dict{
			"sort": "createdAt:desc",
			"pagination": U.Dict{
				"page":     page,
				"pageSize": pageSize,
			},
		})

		if err != nil || status != 200 {
			return allItems, fmt.Errorf("strapi error")
		}

		data, _ := resp.GetList("data")
		for i := 0; i < len(data); i++ {
			item, _ := data.GetDict(i)
			allItems = allItems.Append(item)
		}

		meta, _ := resp.GetDict("meta")
		pg, _ := meta.GetDict("pagination")

		pageCount, _ := pg.GetFloat64("pageCount")
		currentPage, _ := pg.GetFloat64("page")

		if int(currentPage) >= int(pageCount) {
			break
		}

		page++
	}

	return allItems, nil
}

func GetMonthFromCreatedAt(createdAt string) string {
	if len(createdAt) < 7 {
		return ""
	}
	return createdAt[:7]
}

func IsInLast12MonthsExact(month string) bool {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return false
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	start := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -11, 0)

	end := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)

	return !t.Before(start) && !t.After(end)
}

type NormalizedResult struct {
	Result     U.List
	Summary    U.List
	CSVSummary U.Dict
}

func normalizeResultGo(raw U.Dict) NormalizedResult {
	result := U.List{}
	rawSummary := U.List{}
	total, invalid := 0, 0

	for cif, formsVal := range raw {

		forms, ok := formsVal.(U.Dict)
		if !ok {
			continue
		}
		formsDict := U.Dict(forms)

		for form, contentVal := range formsDict {
			if !strings.HasPrefix(form, "F") {
				continue
			}

			contentMap, ok := contentVal.(U.Dict)
			if !ok {
				continue
			}

			content := U.Dict(contentMap)

			rowData, ok := content.GetList("row_data")
			if !ok || len(rowData) == 0 {
				continue
			}

			total++
			invalid++

			firstRow, _ := rowData.GetDict(0)

			logItem := U.Dict{
				"no":       1,
				"row_data": firstRow,
				"reference_document": U.Dict{
					"id": nil, "name": "No Reference", "prompt": "", "result": "",
				},
				"ai_extraction": "",
				"status_icon":   "⚠️",
				"status_label":  "Needs Review",
				"ai_suggestion": "Review required",
			}

			result = append(result, U.Dict{
				"title": fmt.Sprintf("CIF %s - %s", cif, form),
				"generated": U.Dict{
					"results_query": U.List{},
					"tabular_result": U.Dict{
						"logs": U.List{logItem},
					},
				},
			})

			if findVal, ok := formsDict["finding"]; ok {
				finding, _ := findVal.(U.Dict)

				if cmpMap, ok := finding["comparison"].(U.Dict); ok {
					if v, exists := cmpMap[form]; exists {
						rawSummary = append(rawSummary, U.Dict{
							"cif":     cif,
							"form":    form,
							"context": "Comparison",
							"summary": U.Dict(v.(U.Dict))["reason"],
						})
					}
				}

				if accMap, ok := finding["accuracy"].(U.Dict); ok {
					if accVal, exists := accMap[form]; exists {
						m := accVal.(U.Dict)
						var first string
						for k := range m {
							first = k
							break
						}
						rawSummary = append(rawSummary, U.Dict{
							"cif":     cif,
							"form":    form,
							"context": "Accuracy",
							"summary": m[first],
						})
					}
				}

				if compMap, ok := finding["completeness"].(U.Dict); ok {
					if v, exists := compMap[form]; exists {
						rawSummary = append(rawSummary, U.Dict{
							"cif":     cif,
							"form":    form,
							"context": "Completeness",
							"summary": v,
						})
					}
				}
			}
		}
	}

	grouped := U.List{}
	temp := map[string][]string{}

	for _, v := range rawSummary {
		item := v.(U.Dict)
		cif := item["cif"].(string)
		line := fmt.Sprintf("%s : %v", item["form"], item["summary"])
		temp[cif] = append(temp[cif], line)
	}

	i := 1
	for cif, entries := range temp {
		grouped = append(grouped, U.Dict{
			"id":      i,
			"cif":     cif,
			"summary": strings.Join(entries, " \n\n ------------------------------------------- \n\n "),
		})
		i++
	}

	return NormalizedResult{
		Result:  result,
		Summary: grouped,
		CSVSummary: U.Dict{
			"valid":   total - invalid,
			"invalid": invalid,
			"total":   total,
		},
	}
}

func buildLogsGo(result U.List) U.List {
	output := U.List{}
	id := 1

	for _, itemVal := range result {
		item := itemVal.(U.Dict)

		generated, _ := item.GetDict("generated")
		tabular, _ := generated.GetDict("tabular_result")
		logs, _ := tabular.GetList("logs")

		for _, logVal := range logs {
			log := logVal.(U.Dict)

			output = append(output, U.Dict{
				"id": id,
				"no": log["no"],
				"section": U.Dict{
					"name":   item["title"],
					"prompt": "",
					"result": log["row_data"],
				},
				"reference_documents": log["reference_document"],
				"ai_extraction":       log["ai_extraction"],
				"ai_status": U.Dict{
					"icon":       log["status_icon"],
					"label":      log["status_label"],
					"suggestion": log["ai_suggestion"],
				},
			})

			id++
		}
	}

	return output
}

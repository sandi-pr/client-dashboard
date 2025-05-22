package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/requests"
	"gitlab.hexacode.org/go-libs/strapi"
)

func Usage(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	tokenEncoded, ok := ctx.URLParams.GetString("token")
	if !ok {
		fmt.Println("Unauthorized: No token provided")
		return 0, nil, nil
	}

	token, err := url.QueryUnescape(string(tokenEncoded))
	if err != nil {
		fmt.Println("Error decoding token:", err)
		return 0, nil, nil
	}

	_, err = GetId(string(token))
	if err != nil {
		fmt.Println("Unauthorized")
		return 0, nil, nil
	}

	dataType, ok := ctx.URLParams.GetString("dataType")
	if !ok {
		return 0, nil, nil
	}

	usageWS.Serve(ctx.Writer, ctx.Request, string(dataType))
	return 0, nil, nil
}

func GetCPUGeneralInfo(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	info, _ := cpu.Info()
	coreCount, _ := cpu.Counts(false)
	logicalCount, _ := cpu.Counts(true)

	// fmt.Println("== CPU General Information ==")
	// fmt.Println("Vendor:", cpuid.CPU.VendorID)
	// fmt.Println("Model:", info[0].ModelName)
	// fmt.Println("Base Speed:", info[0].Mhz, "MHz")
	// fmt.Println("Physical Cores:", coreCount)
	// fmt.Println("Logical Cores:", logicalCount)
	// fmt.Println("Supports Hyper-Threading:", cpuid.CPU.ThreadsPerCore > 1)
	// fmt.Println("Supports Virtualization:", cpuid.CPU.Supports(cpuid.VMX) || cpuid.CPU.Supports(cpuid.SVM))
	// fmt.Println("L1 Data Cache:", cpuid.CPU.Cache.L1D/1024, "KB")
	// fmt.Println("L1 Instruction Cache:", cpuid.CPU.Cache.L1I/1024, "KB")
	// fmt.Println("L2 Cache:", cpuid.CPU.Cache.L2/1024, "KB")
	// fmt.Println("L3 Cache:", cpuid.CPU.Cache.L3/1024/1024, "MB")

	return 200, U.Dict{
		"vendor":                  cpuid.CPU.VendorID,
		"model":                   info[0].ModelName,
		"base_speed":              info[0].Mhz,
		"physical_cores":          coreCount,
		"logical_cores":           logicalCount,
		"supports_hyperthreading": cpuid.CPU.ThreadsPerCore > 1,
		"supports_virtualization": cpuid.CPU.Supports(cpuid.VMX) || cpuid.CPU.Supports(cpuid.SVM),
		"l1_data_cache":           cpuid.CPU.Cache.L1D / 1024,
		"l1_instruction_cache":    cpuid.CPU.Cache.L1I / 1024,
		"l2_cache":                cpuid.CPU.Cache.L2 / 1024,
		"l3_cache":                cpuid.CPU.Cache.L3 / 1024 / 1024,
	}, nil
}

func GetMemoryGeneralInfo(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
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

	out, err := os.ReadFile(USAGE_MEMORY_PATH)
	if err != nil {
		fmt.Println("DEBUG: Failed to read dmidecode output:", err)
		return ctx.Error(500, "Failed to read memory info")
	}

	lines := strings.Split(string(out), "\n")
	var slotsUsed int
	var speed string
	var formFactor string
	var model string
	var totalSize int

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Size:") && !strings.Contains(line, "No Module") {
			slotsUsed++
			sizeStr := strings.TrimSpace(strings.TrimPrefix(line, "Size:"))
			if strings.HasSuffix(sizeStr, "MB") {
				sizeStr = strings.TrimSuffix(sizeStr, "MB")
				if sizeMB, err := strconv.Atoi(strings.TrimSpace(sizeStr)); err == nil {
					totalSize += sizeMB / 1024
				}
			} else if strings.HasSuffix(sizeStr, "GB") {
				sizeStr = strings.TrimSuffix(sizeStr, "GB")
				if sizeGB, err := strconv.Atoi(strings.TrimSpace(sizeStr)); err == nil {
					totalSize += sizeGB
				}
			}
		}
		// if strings.HasPrefix(line, "Size:") && !strings.Contains(line, "No Module") {
		// 	slotsUsed++
		// }
		if strings.HasPrefix(line, "Type:") && model == "" {
			model = strings.TrimSpace(strings.TrimPrefix(line, "Type:"))
		}
		if strings.HasPrefix(line, "Speed:") && speed == "" {
			speed = strings.TrimSpace(strings.TrimPrefix(line, "Speed:"))
		}
		if strings.HasPrefix(line, "Form Factor:") && formFactor == "" {
			formFactor = strings.TrimSpace(strings.TrimPrefix(line, "Form Factor:"))
		}
	}

	return 200, U.Dict{
		"model":             model,
		"speed":             speed,
		"slots_used":        slotsUsed,
		"form_factor":       formFactor,
		"total_installed":   fmt.Sprintf("%d GB", totalSize),
		"hardware_reserved": nil,
	}, nil
}

func GetDiskGeneralInfo(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")
	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil || status != 200 {
		fmt.Println("Strapi error:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	cmd := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MOUNTPOINT")
	out, err := cmd.Output()
	if err != nil {
		return ctx.Error(500, "Failed to read disk info")
	}

	var parsed LsblkOutput
	if err := json.Unmarshal(out, &parsed); err != nil {
		return ctx.Error(500, "Failed to parse disk info")
	}

	var disks []U.Dict

	for _, d := range parsed.BlockDevices {
		if d.Type == "disk" && d.Size != "0B" {
			diskInfo := U.Dict{
				"model":              "Disk " + detectDiskTypeFromName(d.Name) + " " + d.Size,
				"capacity":           d.Size,
				"formatted_capacity": "-",
				"system_disk":        "No",
				"page_file":          "No",
				"type":               detectDiskTypeFromName(d.Name),
			}

			hasMount := false
			for _, part := range d.Children {
				if part.Mountpoint != "" {
					hasMount = true
					break
				}
			}

			if hasMount {
				diskInfo["system_disk"] = "Yes"
			}

			for _, part := range d.Children {
				if part.Mountpoint != "" {
					dfCmd := exec.Command("df", "-h", part.Mountpoint)
					if dfOut, err := dfCmd.Output(); err == nil {
						lines := strings.Split(string(dfOut), "\n")
						if len(lines) > 1 {
							fields := strings.Fields(lines[1])
							if len(fields) > 1 {
								diskInfo["formatted_capacity"] = fields[1]
								break
							}
						}
					}
				}
			}

			swapCmd := exec.Command("swapon", "--show", "--noheadings", "--bytes")
			if swapOut, err := swapCmd.Output(); err == nil && len(swapOut) > 0 {
				diskInfo["page_file"] = "Yes"
			}

			disks = append(disks, diskInfo)
		}
	}

	return 200, U.Dict{
		"list": disks,
	}, nil
}

func GetGpuGeneralInfo(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")
	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))
	status, me, err := Strapi.UsersMe()
	if err != nil || status != 200 {
		fmt.Println("Strapi error:", err)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	http := requests.NewHttp(U.Dict{
		"Content-type": "application/json",
	})
	resp, err := http.Post(fmt.Sprintf("%s/api/get-gpu-general", HEXA_RAG_URL), nil, U.Dict{}.ToJson())

	if err != nil {
		fmt.Println("Error getting GPU general:", err)
		return ctx.Error(500, "Internal Server Error")
	}
	if resp.Status != 200 {
		fmt.Println("Error getting GPU general:", resp.Status)
		return ctx.Error(500, "Internal Server Error")
	}

	dict := U.NewDictFromJson(resp.Body)

	model, _ := dict.GetString("model")
	driverVersion, _ := dict.GetString("driver_version")
	// driverDate := "unknown"                                        // Tidak tersedia di Linux, hanya tersedia di Windows
	// directXVersion := "Not Available on Linux"                     // DirectX hanya tersedia di Windows
	busID, _ := dict.GetString("bus_id")
	physicalSlot := fmt.Sprintf("Bus ID %s", busID)
	hardwareReserved, _ := dict.GetString("hardware_reserved_memory")

	return 200, U.Dict{
		"model":          model,
		"driver_version": driverVersion,
		// "driver_date":          driverDate,
		// "directx_version":      directXVersion,
		"physical_allocation": physicalSlot,
		"hardware_reserved":   hardwareReserved,
	}, nil
}

type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

type BlockDevice struct {
	Name       string        `json:"name"`
	Size       string        `json:"size"`
	Type       string        `json:"type"`
	Mountpoint string        `json:"mountpoint"`
	Children   []BlockDevice `json:"children,omitempty"`
}

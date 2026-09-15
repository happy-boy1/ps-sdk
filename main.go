// ps-sdk 命令行入口：把各光伏云平台的电站与设备同步到本地数据库。
//
// 用法：
//
//	go run . -platform solarman        # 只同步指定平台
//	go run . -platform all             # 同步所有已配置凭据的平台
//	go run . -list                     # 只读库，打印已同步的数据
//
// 平台凭据优先取环境变量（PS_<平台短码大写>_APPID 等），
// 其次取数据库 platform_auth 表。
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"ps-sdk/pkg/database"
	"ps-sdk/service"
)

func main() {
	var (
		platform = flag.String("platform", "all", "平台短码：solarman|sungrow|ginlong|huawei|all")
		list     = flag.Bool("list", false, "只读取并打印库内数据，不访问平台")
		seeds    = flag.Bool("seed", false, "只写入平台/设备类型基础数据")
		timeout  = flag.Duration("timeout", 30*time.Second, "单次请求超时")
		debug    = flag.Bool("debug", false, "打印请求日志")
		devMode  = flag.Bool("dev", false, "打印完整请求/响应报文")
	)
	flag.Parse()

	if err := run(*platform, *list, *seeds, *timeout, *debug, *devMode); err != nil {
		log.Fatalf("[ps-sdk] %v", err)
	}
}

func run(platform string, list, seeds bool, timeout time.Duration, debug, devMode bool) error {
	if err := database.Init(); err != nil {
		return err
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("[ps-sdk] 关闭数据库失败: %v", err)
		}
	}()

	if err := service.Init(database.DB); err != nil {
		return fmt.Errorf("写入基础数据失败: %w", err)
	}
	if seeds {
		log.Println("[ps-sdk] 平台与设备类型基础数据已写入")
		return nil
	}

	opts := service.Options{Timeout: timeout, Debug: debug, DevMode: devMode}
	if list {
		return printStored(platform)
	}
	return syncPlatforms(platform, opts)
}

// syncPlatforms 同步指定平台，all 表示遍历全部平台
func syncPlatforms(platform string, opts service.Options) error {
	if strings.EqualFold(platform, "all") {
		results, err := service.SyncAll(opts)
		for _, result := range results {
			fmt.Println(result)
			printErrors(result)
		}
		return err
	}

	client, err := service.PlatformClient(platform, opts)
	if err != nil {
		return err
	}
	result, err := client.Sync()
	fmt.Println(result)
	printErrors(result)
	return err
}

// printStored 打印库内数据，便于同步后核对
func printStored(platform string) error {
	ids := make([]int16, 0, len(service.Platforms))
	if strings.EqualFold(platform, "all") {
		for _, p := range service.Platforms {
			ids = append(ids, p.ID)
		}
	} else {
		p, ok := service.PlatformByCode(platform)
		if !ok {
			return fmt.Errorf("%w: %s", service.ErrPlatformUnknown, platform)
		}
		ids = append(ids, p.ID)
	}

	for _, id := range ids {
		p, _ := service.PlatformByID(id)
		stations, err := service.ListStations(id)
		if err != nil {
			return err
		}
		fmt.Printf("== %s 电站 %d 个 ==\n", p.NameZH, len(stations))
		for _, station := range stations {
			devices, err := service.ListDevices(id, station.StationID)
			if err != nil {
				return err
			}
			fmt.Printf("  [%d] %s(%s) %.2fkW 状态 %s 设备 %d 台\n",
				station.StationID, station.StationName, station.StationIDOrigin,
				valueOrZero(station.CapacityKwp), service.StationStatusName(station.Status), len(devices))
			for _, device := range devices {
				fmt.Printf("      - %-18s %-20s %s %s\n",
					device.DeviceSN, device.DeviceIDOrigin,
					service.DeviceStatusName(device.Status), device.DeviceName)
			}
		}
	}
	return nil
}

// printErrors 打印同步过程中的单项失败与告警
func printErrors(result *service.SyncResult) {
	for _, err := range result.Warnings {
		log.Printf("[ps-sdk][告警] %v", err)
	}
	for _, err := range result.Errors {
		log.Printf("[ps-sdk][失败] %v", err)
	}
}

// valueOrZero 解引用可空数值，空值按 0 展示
func valueOrZero(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

// ps-sdk 命令行入口：把各光伏云平台的电站与设备同步到本地数据库。
//
// 用法：
//
//	go run . -platform solarman        # 只同步指定平台
//	go run . -platform all             # 同步所有已配置凭据的平台
//	go run . -list                     # 只读库，打印已同步的数据
//	go run . -config ./config.toml     # 指定配置文件
//
// 配置优先级由低到高：config.toml → 数据库 platform_auth 表（app_key/app_secret）→ 环境变量。
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"ps-sdk/pkg/config"
	"ps-sdk/pkg/database"
	"ps-sdk/service"
)

func main() {
	var (
		platform   = flag.String("platform", "all", "平台短码：solarman|sungrow|ginlong|huawei|all")
		configPath = flag.String("config", "", "配置文件路径，缺省依次查找 PS_CONFIG、./config.toml、./pkg/config/config.toml")
		list       = flag.Bool("list", false, "只读取并打印库内数据，不访问平台")
		seeds      = flag.Bool("seed", false, "只写入平台/设备类型基础数据")
	)
	flag.Parse()

	if err := run(*platform, *configPath, *list, *seeds); err != nil {
		log.Fatalf("[ps-sdk] %v", err)
	}
}

func run(platform, configPath string, list, seeds bool) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	log.Printf("[ps-sdk] 配置文件: %s", cfg.Path)

	if err := database.Init(cfg.Database); err != nil {
		return err
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("[ps-sdk] 关闭数据库失败: %v", err)
		}
	}()

	if err := service.InitWith(database.DB, cfg); err != nil {
		return fmt.Errorf("写入基础数据失败: %w", err)
	}
	if seeds {
		log.Println("[ps-sdk] 平台与设备类型基础数据已写入")
		return nil
	}

	if list {
		return printStored(platform)
	}
	return syncPlatforms(platform, service.OptionsFrom(cfg))
}

// syncPlatforms 同步指定平台，all 表示遍历全部平台
func syncPlatforms(platform string, opts service.Options) error {
	if strings.EqualFold(platform, "all") {
		results, err := service.SyncAll(opts)
		for _, result := range results {
			fmt.Println(result)
			printIssues(result)
		}
		return err
	}

	client, err := service.PlatformClient(platform, opts)
	if err != nil {
		return err
	}
	result, err := client.Sync()
	fmt.Println(result)
	printIssues(result)
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

// printIssues 打印同步过程中的告警与失败明细
func printIssues(result *service.SyncResult) {
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

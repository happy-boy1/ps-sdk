package solarman

import "fmt"

// ---------------------------------------------------------------------------
// 4.17 查询电站实时天气
// ---------------------------------------------------------------------------

// StationWeatherRequest 4.17 查询电站实时天气请求
type StationWeatherRequest struct {
	StationID Int64 `json:"stationId"` // 电站 ID，必填
}

// StationWeatherResult 4.17 查询电站实时天气返回
type StationWeatherResult struct {
	Response           // 公共响应字段，成功时 code 为 null
	rawHolder          // 原始数据兜底
	WeatherTime Str    `json:"weatherTime"` // 天气时间 yyyyMMddHHmmss
	WeatherID   Int64  `json:"weatherId"`   // 天气 ID
	WeatherCode string `json:"weatherCode"` // 天气 code，如 Clouds
	Temp        Num    `json:"temp"`        // 温度，单位 °C
	Pressure    Num    `json:"pressure"`    // 气压
	Humidity    Num    `json:"humidity"`    // 湿度
	WindSpeed   Num    `json:"windSpeed"`   // 风速
	WindDeg     Num    `json:"windDeg"`     // 风向，角度
	Clouds      Num    `json:"clouds"`      // 云量
	Sunrise     Int64  `json:"sunrise"`     // 日出时间戳
	Sunset      Int64  `json:"sunset"`      // 日落时间戳
	TempMin     Num    `json:"tempMin"`     // 最小温度，单位 °C
	TempMax     Num    `json:"tempMax"`     // 最大温度，单位 °C
	WeatherPic  Str    `json:"weatherPic"`  // 天气图标名，如 cloudy
	DateTime    Int64  `json:"dateTime"`    // 天气时间戳
}

// StationWeather 4.17 查询电站所在地实时天气
func (sdk *SolarmanSDK) StationWeather(req StationWeatherRequest) (*StationWeatherResult, error) {
	if req.StationID <= 0 {
		return nil, fmt.Errorf("[Solarman]: stationId 是空值或非法")
	}
	var out StationWeatherResult
	if err := sdk.do(PathStationWeather, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.18 查询实时天气
// ---------------------------------------------------------------------------

// CurrentWeatherRequest 4.18 按区域查询实时天气请求
type CurrentWeatherRequest struct {
	RegionNationID Int64 `json:"regionNationId"` // 国家 ID，必填
	RegionLevel1   Int64 `json:"regionLevel1"`   // 省份 ID，必填
	RegionLevel2   Int64 `json:"regionLevel2"`   // 城市 ID，必填
}

// CurrentWeatherResult 4.18 按区域查询实时天气返回
type CurrentWeatherResult struct {
	Response           // 公共响应字段，成功时 code 为 null
	rawHolder          // 原始数据兜底
	WeatherID   Int64  `json:"weatherId"`   // 天气 ID
	WeatherCode string `json:"weatherCode"` // 天气 code，如 Clear
	Temp        Num    `json:"temp"`        // 温度，单位 °C
	Pressure    Num    `json:"pressure"`    // 气压
	Humidity    Num    `json:"humidity"`    // 湿度
	WindSpeed   Num    `json:"windSpeed"`   // 风速
	WindDeg     Num    `json:"windDeg"`     // 风向，角度
	Clouds      Num    `json:"clouds"`      // 云量
	Sunrise     Int64  `json:"sunrise"`     // 日出时间戳
	Sunset      Int64  `json:"sunset"`      // 日落时间戳
	TempMin     Num    `json:"tempMin"`     // 最小温度，单位 °C
	TempMax     Num    `json:"tempMax"`     // 最大温度，单位 °C
	WeatherPic  Str    `json:"weatherPic"`  // 天气图标名，如 sunny
	DateTime    Int64  `json:"dateTime"`    // 天气时间戳
}

// CurrentWeather 4.18 按经纬度所在地区域查询实时天气
func (sdk *SolarmanSDK) CurrentWeather(req CurrentWeatherRequest) (*CurrentWeatherResult, error) {
	if req.RegionNationID <= 0 {
		return nil, fmt.Errorf("[Solarman]: regionNationId 是空值或非法")
	}
	if req.RegionLevel1 <= 0 {
		return nil, fmt.Errorf("[Solarman]: regionLevel1 是空值或非法")
	}
	if req.RegionLevel2 <= 0 {
		return nil, fmt.Errorf("[Solarman]: regionLevel2 是空值或非法")
	}
	var out CurrentWeatherResult
	if err := sdk.do(PathWeatherCurrent, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------------------------------------------------------------------------
// 4.19 查询未来天气
// ---------------------------------------------------------------------------

// ForecastWeatherRequest 4.19 按区域查询未来天气请求
type ForecastWeatherRequest struct {
	RegionNationID Int64 `json:"regionNationId"` // 国家 ID，必填
	RegionLevel1   Int64 `json:"regionLevel1"`   // 省份 ID，必填
	RegionLevel2   Int64 `json:"regionLevel2"`   // 城市 ID，必填
}

// ForecastWeatherItem 4.19 未来天气的单个逐日预报项
type ForecastWeatherItem struct {
	rawHolder
	WeatherID   Int64  `json:"weatherId"`   // 天气 ID
	WeatherCode string `json:"weatherCode"` // 天气 code，如 Clouds
	Temp        Num    `json:"temp"`        // 温度，单位 °C
	Pressure    Num    `json:"pressure"`    // 气压
	Humidity    Num    `json:"humidity"`    // 湿度
	WindSpeed   Num    `json:"windSpeed"`   // 风速
	WindDeg     Num    `json:"windDeg"`     // 风向，角度
	Clouds      Num    `json:"clouds"`      // 云量
	TempMin     Num    `json:"tempMin"`     // 最小温度，单位 °C
	TempMax     Num    `json:"tempMax"`     // 最大温度，单位 °C
	WeatherPic  Str    `json:"weatherPic"`  // 天气图标名，如 cloudy
	DateTime    Int64  `json:"dateTime"`    // 天气时间戳
}

// ForecastWeatherResult 4.19 按区域查询未来天气返回
type ForecastWeatherResult struct {
	Response                           // 公共响应字段，成功时 code 为 null
	rawHolder                          // 原始数据兜底
	TWeatherList []ForecastWeatherItem `json:"tweatherList"` // 未来天气列表
}

// ForecastWeather 4.19 按经纬度所在地区域查询未来天气
func (sdk *SolarmanSDK) ForecastWeather(req ForecastWeatherRequest) (*ForecastWeatherResult, error) {
	if req.RegionNationID <= 0 {
		return nil, fmt.Errorf("[Solarman]: regionNationId 是空值或非法")
	}
	if req.RegionLevel1 <= 0 {
		return nil, fmt.Errorf("[Solarman]: regionLevel1 是空值或非法")
	}
	if req.RegionLevel2 <= 0 {
		return nil, fmt.Errorf("[Solarman]: regionLevel2 是空值或非法")
	}
	var out ForecastWeatherResult
	if err := sdk.do(PathWeatherForecast, req, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

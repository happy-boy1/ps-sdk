package huawei

// 电站实时数据 key
const (
	KpiRealHealthState = "real_health_state"  // 电站健康状态(1断连 2故障 3健康)
	KpiDayPower        = "day_power"          // 当日发电量 kWh
	KpiMonthPower      = "month_power"        // 当月发电量 kWh
	KpiTotalPower      = "total_power"        // 总发电量 kWh
	KpiDayIncome       = "day_income"         // 当日收益
	KpiTotalIncome     = "total_income"       // 总收益
	KpiDayOnGridEnergy = "day_on_grid_energy" // 当日上网电量 kWh
	KpiDayUseEnergy    = "day_use_energy"     // 当日用电量 kWh
)

// 电站报表数据 key
const (
	KpiRadiationIntensity = "radiation_intensity" // 总辐照量 kWh/m²
	KpiTheoryPower        = "theory_power"        // 理论发电量 kWh
	KpiInverterPower      = "inverter_power"      // 逆变器发电量 kWh(命名有歧义，推荐 PVYield)
	KpiOngridPower        = "ongrid_power"        // 上网电量 kWh
	KpiPowerProfit        = "power_profit"        // 发电收益
	KpiChargeCap          = "chargeCap"           // 充电电量 kWh
	KpiDischargeCap       = "dischargeCap"        // 放电电量 kWh
	KpiSelfProvide        = "selfProvide"         // 自给自足电量 kWh
	KpiPVYield            = "PVYield"             // PV 发电量 kWh
	KpiInverterYield      = "inverterYield"       // 逆变器发电量 kWh
)

// StationRealtime 电站实时数据（无效值按 0 处理）
type StationRealtime struct {
	RealHealthState int     // 健康状态 1断连 2故障 3健康
	DayPower        float64 // 当日发电量 kWh
	MonthPower      float64 // 当月发电量 kWh
	TotalPower      float64 // 总发电量 kWh
	DayIncome       float64 // 当日收益
	TotalIncome     float64 // 总收益
	DayOnGridEnergy float64 // 当日上网电量 kWh
	DayUseEnergy    float64 // 当日用电量 kWh
}

// ParseStationRealtime 解析电站实时数据项
func (m ItemMap) ParseStationRealtime() StationRealtime {
	state, _ := m.Int(KpiRealHealthState)
	var r StationRealtime
	r.RealHealthState = state
	r.DayPower, _ = m.Float64(KpiDayPower)
	r.MonthPower, _ = m.Float64(KpiMonthPower)
	r.TotalPower, _ = m.Float64(KpiTotalPower)
	r.DayIncome, _ = m.Float64(KpiDayIncome)
	r.TotalIncome, _ = m.Float64(KpiTotalIncome)
	r.DayOnGridEnergy, _ = m.Float64(KpiDayOnGridEnergy)
	r.DayUseEnergy, _ = m.Float64(KpiDayUseEnergy)
	return r
}

// KPI 解析该电站的实时指标
func (r StationRealKpi) KPI() StationRealtime { return r.DataItemMap.ParseStationRealtime() }

// StationReport 电站小时/日/月/年报表数据（无效值按 0 处理）
type StationReport struct {
	RadiationIntensity float64 // 总辐照量 kWh/m²
	TheoryPower        float64 // 理论发电量 kWh
	InverterPower      float64 // 逆变器发电量 kWh
	OngridPower        float64 // 上网电量 kWh
	PowerProfit        float64 // 发电收益
	ChargeCap          float64 // 充电电量 kWh
	DischargeCap       float64 // 放电电量 kWh
	SelfProvide        float64 // 自给自足电量 kWh
	PVYield            float64 // PV 发电量 kWh
	InverterYield      float64 // 逆变器发电量 kWh
}

// ParseStationReport 解析电站报表数据项
func (m ItemMap) ParseStationReport() StationReport {
	var r StationReport
	r.RadiationIntensity, _ = m.Float64(KpiRadiationIntensity)
	r.TheoryPower, _ = m.Float64(KpiTheoryPower)
	r.InverterPower, _ = m.Float64(KpiInverterPower)
	r.OngridPower, _ = m.Float64(KpiOngridPower)
	r.PowerProfit, _ = m.Float64(KpiPowerProfit)
	r.ChargeCap, _ = m.Float64(KpiChargeCap)
	r.DischargeCap, _ = m.Float64(KpiDischargeCap)
	r.SelfProvide, _ = m.Float64(KpiSelfProvide)
	r.PVYield, _ = m.Float64(KpiPVYield)
	r.InverterYield, _ = m.Float64(KpiInverterYield)
	return r
}

// KPI 解析该条记录的报表指标
func (r StationKpiItem) KPI() StationReport { return r.DataItemMap.ParseStationReport() }

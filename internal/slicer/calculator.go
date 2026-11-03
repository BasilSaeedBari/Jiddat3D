package slicer

import (
	"fmt"
	"math"
)

// SlicerConfig holds configuration parameters, typically loaded from PocketBase slicer_settings
type SlicerConfig struct {
	Name                 string  `json:"name"`
	LayerHeightMm        float64 `json:"layer_height_mm"`         // e.g. 0.20 mm
	WallPerimeters       int     `json:"wall_perimeters"`          // e.g. 3 walls
	WallThicknessMm      float64 `json:"wall_thickness_mm"`       // e.g. 1.2 mm
	InfillPercentage     float64 `json:"infill_percentage"`       // e.g. 20%
	TimeMultiplier       float64 `json:"time_multiplier"`         // e.g. 1.5x buffer timer
	WeightMultiplier     float64 `json:"weight_multiplier"`       // e.g. 1.5x weight buffer
	MachineHourlyRate    int64   `json:"machine_hourly_rate"`     // e.g. 250 Rs/hr
	LaborOverhead        int64   `json:"labor_overhead"`          // e.g. 100 Rs prep
	ElectricityHourly    int64   `json:"electricity_hourly"`      // e.g. 35 Rs/hr
	FailureRatePct       int64   `json:"failure_rate_pct"`        // e.g. 5%
	ProfitMultiplierPct  int64   `json:"profit_multiplier_pct"`   // e.g. 25%
	DefaultFilamentPrice int64   `json:"default_filament_price"`  // e.g. 8 Rs/gram for White Basic PLA
}

// FilamentProfile describes a selectable material option
type FilamentProfile struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`           // e.g. "White Basic PLA"
	Material     string  `json:"material"`       // e.g. "PLA"
	ColorName    string  `json:"color_name"`     // e.g. "Signal White"
	ColorHex     string  `json:"color_hex"`      // e.g. "#FFFFFF"
	DensityGCm3  float64 `json:"density_g_cm3"`  // e.g. 1.24 g/cm3
	PricePerGram int64   `json:"price_per_gram"` // in Rs per gram, e.g. 8
	IsDefault    bool    `json:"is_default"`
}

// DefaultSlicerConfig provides reliable fallbacks
func DefaultSlicerConfig() SlicerConfig {
	return SlicerConfig{
		Name:                 "Jiddat CoreXY / Mark-1 0.20mm Production Profile",
		LayerHeightMm:        0.20,
		WallPerimeters:       3,
		WallThicknessMm:      1.2,
		InfillPercentage:     20.0,
		TimeMultiplier:       1.5, // 1.5x user-specified buffer timer
		WeightMultiplier:     1.5, // 1.5x user-specified buffer weight
		MachineHourlyRate:    250, // Rs 250 per print hour
		LaborOverhead:        50,  // Rs 50 bench handling
		ElectricityHourly:    30,  // Rs 30/hr electricity
		FailureRatePct:       5,   // 5% allowance
		ProfitMultiplierPct:  25,  // 25% margin
		DefaultFilamentPrice: 8,   // Rs 8 / gram
	}
}

// DefaultFilaments provides the standard in-house materials
func DefaultFilaments() []FilamentProfile {
	return []FilamentProfile{
		{
			ID:           "white-basic-pla",
			Name:         "White Basic PLA",
			Material:     "PLA",
			ColorName:    "Signal White",
			ColorHex:     "#FAF7F2",
			DensityGCm3:  1.24,
			PricePerGram: 8,
			IsDefault:    true,
		},
		{
			ID:           "black-tough-pla",
			Name:         "Black Tough PLA+",
			Material:     "PLA+",
			ColorName:    "Matte Black",
			ColorHex:     "#1A1A1A",
			DensityGCm3:  1.25,
			PricePerGram: 10,
			IsDefault:    false,
		},
		{
			ID:           "high-temp-petg",
			Name:         "High-Temp PETG",
			Material:     "PETG",
			ColorName:    "Terracotta Amber",
			ColorHex:     "#D9530F",
			DensityGCm3:  1.27,
			PricePerGram: 12,
			IsDefault:    false,
		},
		{
			ID:           "flex-tpu-95a",
			Name:         "TPU 95A Flexible",
			Material:     "TPU",
			ColorName:    "Emerald Flex",
			ColorHex:     "#17C3B2",
			DensityGCm3:  1.21,
			PricePerGram: 18,
			IsDefault:    false,
		},
	}
}

// SlicedQuoteResult holds the complete calculated slice and commercial invoice estimate
type SlicedQuoteResult struct {
	Mesh            MeshStats       `json:"mesh"`
	Filament        FilamentProfile `json:"filament"`
	Config          SlicerConfig    `json:"config"`

	// Kinematics & Slicing Metrics
	LayerCount      int     `json:"layer_count"`
	RawVolumeCm3    float64 `json:"raw_volume_cm3"`
	ExtrudedVolCm3  float64 `json:"extruded_vol_cm3"`
	RawGrams        float64 `json:"raw_grams"`
	BufferedGrams   float64 `json:"buffered_grams"`   // Displayed grams with 1.5x buffer
	RawHours        float64 `json:"raw_hours"`
	BufferedHours   float64 `json:"buffered_hours"`   // Displayed hours with 1.5x buffer
	FormattedTime   string  `json:"formatted_time"`   // e.g. "3 hrs 45 mins"

	// LayerInvoice Commercial Cost Breakdown (in PKR / Rs)
	RawFilamentCost     int64 `json:"raw_filament_cost"`
	MachineTimeCost     int64 `json:"machine_time_cost"`
	LaborCost           int64 `json:"labor_cost"`
	ElectricityCost     int64 `json:"electricity_cost"`
	BaseProductionCost  int64 `json:"base_production_cost"`
	FailureAllowance    int64 `json:"failure_allowance"`
	TotalProductionCost int64 `json:"total_production_cost"`
	ProfitMarkup        int64 `json:"profit_markup"`
	FinalPrice          int64 `json:"final_price"` // Total quote in Rs
}

// SliceAndCalculate evaluates the 3D model geometry and applies LayerInvoice math
func SliceAndCalculate(mesh MeshStats, fil FilamentProfile, cfg SlicerConfig) SlicedQuoteResult {
	if cfg.LayerHeightMm <= 0 {
		cfg.LayerHeightMm = 0.20
	}
	if cfg.TimeMultiplier <= 0 {
		cfg.TimeMultiplier = 1.5
	}
	if cfg.WeightMultiplier <= 0 {
		cfg.WeightMultiplier = 1.5
	}
	if fil.DensityGCm3 <= 0 {
		fil.DensityGCm3 = 1.24
	}
	if fil.PricePerGram <= 0 {
		fil.PricePerGram = cfg.DefaultFilamentPrice
		if fil.PricePerGram <= 0 {
			fil.PricePerGram = 8
		}
	}

	layerCount := int(math.Ceil(mesh.SizeZ / cfg.LayerHeightMm))
	if layerCount < 1 {
		layerCount = 1
	}

	// Calculate extruded volume:
	// Surface shell: area (cm^2) * wall_thickness (cm)
	wallCm := cfg.WallThicknessMm / 10.0
	shellVolCm3 := mesh.SurfaceAreaCm2 * wallCm
	if shellVolCm3 > mesh.VolumeCm3 {
		shellVolCm3 = mesh.VolumeCm3
	}

	// Infill volume: remainder * infill_pct
	infillRatio := cfg.InfillPercentage / 100.0
	infillVolCm3 := (mesh.VolumeCm3 - shellVolCm3) * infillRatio
	if infillVolCm3 < 0 {
		infillVolCm3 = 0
	}

	extrudedVolCm3 := shellVolCm3 + infillVolCm3

	// Raw filament mass
	rawGrams := extrudedVolCm3 * fil.DensityGCm3
	// Apply buffer multiplier (e.g. 1.5x)
	bufferedGrams := rawGrams * cfg.WeightMultiplier
	if bufferedGrams < 5.0 && mesh.VolumeCm3 > 0 {
		bufferedGrams = 5.0 // minimum 5g threshold for functional print
	}

	// Calculate Raw Print Hours based on extrusion rate and kinematic speed:
	// Modern corexy volumetric speed ~12 mm^3/sec = 0.012 cm^3/sec = 43.2 cm^3/hour
	// Wall printing speed ~60mm/s, Infill speed ~140mm/s, plus layer change overhead ~2.5 sec/layer
	volumetricSpeedCm3Hr := 16.0 // avg 16 cm^3/hr effective deposition
	depositionHours := extrudedVolCm3 / volumetricSpeedCm3Hr
	layerOverheadHours := (float64(layerCount) * 2.5) / 3600.0 // 2.5s z-hop, retraction, acceleration per layer
	rawHours := depositionHours + layerOverheadHours
	if rawHours < 0.25 {
		rawHours = 0.25 // minimum 15 mins
	}

	// Apply user-specified time buffer multiplier (e.g. 1.5x)
	bufferedHours := rawHours * cfg.TimeMultiplier

	// Format display time
	totalMinutes := int(math.Round(bufferedHours * 60))
	hrs := totalMinutes / 60
	mins := totalMinutes % 60
	formattedTime := ""
	if hrs > 0 {
		formattedTime = fmt.Sprintf("%d hr", hrs)
		if hrs > 1 {
			formattedTime += "s"
		}
		if mins > 0 {
			formattedTime += fmt.Sprintf(" %d min", mins)
			if mins > 1 {
				formattedTime += "s"
			}
		}
	} else {
		formattedTime = fmt.Sprintf("%d mins", mins)
	}

	// --- LayerInvoice Cost Formulation ---
	// 1. Raw Filament Cost
	rawFilamentCost := int64(math.Round(bufferedGrams * float64(fil.PricePerGram)))

	// 2. Machine Time Cost
	machineTimeCost := int64(math.Round(bufferedHours * float64(cfg.MachineHourlyRate)))

	// 3. Labor Cost
	laborCost := cfg.LaborOverhead

	// 4. Electricity Cost
	electricityCost := int64(math.Round(bufferedHours * float64(cfg.ElectricityHourly)))

	// 5. Base Production Cost
	baseProductionCost := rawFilamentCost + machineTimeCost + laborCost + electricityCost

	// 6. Failure Allowance
	failureAllowance := int64(math.Round(float64(baseProductionCost) * float64(cfg.FailureRatePct) / 100.0))

	// 7. Total Production Cost
	totalProductionCost := baseProductionCost + failureAllowance

	// 8. Profit Markup
	profitMarkup := int64(math.Round(float64(totalProductionCost) * float64(cfg.ProfitMultiplierPct) / 100.0))

	// 9. Final Customer Price (rounded to nearest 10 Rs)
	unroundedPrice := totalProductionCost + profitMarkup
	finalPrice := ((unroundedPrice + 5) / 10) * 10
	if finalPrice < 150 {
		finalPrice = 150 // Minimum shop order threshold
	}

	return SlicedQuoteResult{
		Mesh:                mesh,
		Filament:            fil,
		Config:              cfg,
		LayerCount:          layerCount,
		RawVolumeCm3:        mesh.VolumeCm3,
		ExtrudedVolCm3:      extrudedVolCm3,
		RawGrams:            rawGrams,
		BufferedGrams:       bufferedGrams,
		RawHours:            rawHours,
		BufferedHours:       bufferedHours,
		FormattedTime:       formattedTime,
		RawFilamentCost:     rawFilamentCost,
		MachineTimeCost:     machineTimeCost,
		LaborCost:           laborCost,
		ElectricityCost:     electricityCost,
		BaseProductionCost:  baseProductionCost,
		FailureAllowance:    failureAllowance,
		TotalProductionCost: totalProductionCost,
		ProfitMarkup:        profitMarkup,
		FinalPrice:          finalPrice,
	}
}

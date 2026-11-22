package slicer

import (
	"fmt"
	"math"
)

// PrinterPreset encodes hardware limits and process speeds/accelerations/rules
// exported from OrcaSlicer (Cbot CN-1 0.4 nozzle + @Cbot CN-1 0.4 nozzle V11).
type PrinterPreset struct {
	Name         string `json:"name"`
	PrinterModel string `json:"printer_model"`
	GcodeFlavor  string `json:"gcode_flavor"`

	// Build envelope & nozzle geometry (mm)
	NozzleDiameterMm float64 `json:"nozzle_diameter_mm"`
	BuildXMm         float64 `json:"build_x_mm"`
	BuildYMm         float64 `json:"build_y_mm"`
	BuildZMm         float64 `json:"build_z_mm"`
	MinLayerHeightMm float64 `json:"min_layer_height_mm"`
	MaxLayerHeightMm float64 `json:"max_layer_height_mm"`

	// Hardware firmware limits (M201 / M203 / M205)
	MaxSpeedXMmS          float64 `json:"max_speed_x_mm_s"`
	MaxSpeedYMmS          float64 `json:"max_speed_y_mm_s"`
	MaxSpeedZMmS          float64 `json:"max_speed_z_mm_s"`
	MaxSpeedEMmS          float64 `json:"max_speed_e_mm_s"`
	MaxAccelXMmS2         float64 `json:"max_accel_x_mm_s2"`
	MaxAccelYMmS2         float64 `json:"max_accel_y_mm_s2"`
	MaxAccelZMmS2         float64 `json:"max_accel_z_mm_s2"`
	MaxAccelEMmS2         float64 `json:"max_accel_e_mm_s2"`
	MaxAccelExtrudingMmS2 float64 `json:"max_accel_extruding_mm_s2"`
	MaxAccelRetractingMmS2 float64 `json:"max_accel_retracting_mm_s2"`
	MaxAccelTravelMmS2    float64 `json:"max_accel_travel_mm_s2"`
	MaxJerkXMmS           float64 `json:"max_jerk_x_mm_s"`
	MaxJerkYMmS           float64 `json:"max_jerk_y_mm_s"`
	MaxJerkZMmS           float64 `json:"max_jerk_z_mm_s"`
	MaxJerkEMmS           float64 `json:"max_jerk_e_mm_s"`

	// Retraction, Z-hop, wipe & startup compensation
	RetractionLengthMm     float64 `json:"retraction_length_mm"`
	RetractionSpeedMmS     float64 `json:"retraction_speed_mm_s"`
	DeretractionSpeedMmS   float64 `json:"deretraction_speed_mm_s"`
	ZHopMm                 float64 `json:"z_hop_mm"`
	WipeDistanceMm         float64 `json:"wipe_distance_mm"`
	PrepareCompensationSec float64 `json:"prepare_compensation_sec"`

	// Shell, line width & adhesion rules (V11 + process template)
	InitialLayerHeightMm    float64 `json:"initial_layer_height_mm"`
	OuterWallLineWidthMm    float64 `json:"outer_wall_line_width_mm"`
	InnerWallLineWidthMm    float64 `json:"inner_wall_line_width_mm"`
	SparseInfillLineWidthMm float64 `json:"sparse_infill_line_width_mm"`
	SolidInfillLineWidthMm  float64 `json:"solid_infill_line_width_mm"`
	TopSurfaceLineWidthMm   float64 `json:"top_surface_line_width_mm"`
	InitialLayerLineWidthMm float64 `json:"initial_layer_line_width_mm"`
	TopShellLayers          int     `json:"top_shell_layers"`
	TopShellMinThicknessMm  float64 `json:"top_shell_min_thickness_mm"`
	BottomShellLayers       int     `json:"bottom_shell_layers"`
	InfillPattern           string  `json:"infill_pattern"`
	WallGenerator           string  `json:"wall_generator"`
	BrimWidthMm             float64 `json:"brim_width_mm"`
	SkirtLoops              int     `json:"skirt_loops"`
	SlowDownLayers          int     `json:"slow_down_layers"`
	SlowDownLayerTimeSec    float64 `json:"slow_down_layer_time_sec"`
	SlowDownMinSpeedMmS     float64 `json:"slow_down_min_speed_mm_s"`

	// Process feature speeds (mm/s)
	OuterWallSpeedMmS          float64 `json:"outer_wall_speed_mm_s"`
	InnerWallSpeedMmS          float64 `json:"inner_wall_speed_mm_s"`
	SparseInfillSpeedMmS       float64 `json:"sparse_infill_speed_mm_s"`
	InternalSolidInfillSpeedMmS float64 `json:"internal_solid_infill_speed_mm_s"`
	TopSurfaceSpeedMmS         float64 `json:"top_surface_speed_mm_s"`
	GapInfillSpeedMmS          float64 `json:"gap_infill_speed_mm_s"`
	BridgeSpeedMmS             float64 `json:"bridge_speed_mm_s"`
	InitialLayerSpeedMmS       float64 `json:"initial_layer_speed_mm_s"`
	InitialLayerInfillSpeedMmS float64 `json:"initial_layer_infill_speed_mm_s"`
	TravelSpeedMmS             float64 `json:"travel_speed_mm_s"`
	InitialLayerTravelSpeedMmS float64 `json:"initial_layer_travel_speed_mm_s"`

	// Process feature accelerations (mm/s^2)
	DefaultAccelMmS2             float64 `json:"default_accel_mm_s2"`
	OuterWallAccelMmS2           float64 `json:"outer_wall_accel_mm_s2"`
	InnerWallAccelMmS2           float64 `json:"inner_wall_accel_mm_s2"`
	SparseInfillAccelMmS2        float64 `json:"sparse_infill_accel_mm_s2"`
	InternalSolidInfillAccelMmS2 float64 `json:"internal_solid_infill_accel_mm_s2"`
	TopSurfaceAccelMmS2          float64 `json:"top_surface_accel_mm_s2"`
	BridgeAccelMmS2              float64 `json:"bridge_accel_mm_s2"`
	InitialLayerAccelMmS2        float64 `json:"initial_layer_accel_mm_s2"`
	TravelAccelMmS2              float64 `json:"travel_accel_mm_s2"`
	InitialLayerTravelAccelMmS2  float64 `json:"initial_layer_travel_accel_mm_s2"`

	// Process feature jerks (mm/s)
	DefaultJerkMmS            float64 `json:"default_jerk_mm_s"`
	OuterWallJerkMmS          float64 `json:"outer_wall_jerk_mm_s"`
	InnerWallJerkMmS          float64 `json:"inner_wall_jerk_mm_s"`
	InfillJerkMmS             float64 `json:"infill_jerk_mm_s"`
	TopSurfaceJerkMmS         float64 `json:"top_surface_jerk_mm_s"`
	InitialLayerJerkMmS       float64 `json:"initial_layer_jerk_mm_s"`
	TravelJerkMmS             float64 `json:"travel_jerk_mm_s"`
	InitialLayerTravelJerkMmS float64 `json:"initial_layer_travel_jerk_mm_s"`
}

// DefaultCbotCN1Preset returns the calibrated OrcaSlicer Cbot CN-1 0.4mm nozzle (V11) preset.
func DefaultCbotCN1Preset() PrinterPreset {
	return PrinterPreset{
		Name:         "Cbot CN-1 0.4 nozzle (V11)",
		PrinterModel: "Cbot CN-1",
		GcodeFlavor:  "marlin",

		// Build volume: 200x200x200 mm, 0.4 mm brass nozzle
		NozzleDiameterMm: 0.4,
		BuildXMm:         200,
		BuildYMm:         200,
		BuildZMm:         200,
		MinLayerHeightMm: 0.08,
		MaxLayerHeightMm: 0.32,

		// Firmware M201 / M203 / M205 hardware limits
		MaxSpeedXMmS:           200,
		MaxSpeedYMmS:           200,
		MaxSpeedZMmS:           10,
		MaxSpeedEMmS:           40,
		MaxAccelXMmS2:          2000,
		MaxAccelYMmS2:          2000,
		MaxAccelZMmS2:          100,
		MaxAccelEMmS2:          3000,
		MaxAccelExtrudingMmS2:  1500,
		MaxAccelRetractingMmS2: 3000,
		MaxAccelTravelMmS2:     2000,
		MaxJerkXMmS:            35,
		MaxJerkYMmS:            35,
		MaxJerkZMmS:            0.4,
		MaxJerkEMmS:            5,

		// Direct Drive retraction, Z-hop, wipe & startup heat/home/prime compensation
		RetractionLengthMm:     6.0,
		RetractionSpeedMmS:     35,
		DeretractionSpeedMmS:   35,
		ZHopMm:                 0.4,
		WipeDistanceMm:         0.4,
		PrepareCompensationSec: 260,

		// V11 + template shell, line-width & adhesion rules
		InitialLayerHeightMm:    0.20,
		OuterWallLineWidthMm:    0.42,
		InnerWallLineWidthMm:    0.45,
		SparseInfillLineWidthMm: 0.45,
		SolidInfillLineWidthMm:  0.45,
		TopSurfaceLineWidthMm:   0.42,
		InitialLayerLineWidthMm: 0.50,
		TopShellLayers:          5,
		TopShellMinThicknessMm:  0.60,
		BottomShellLayers:       3,
		InfillPattern:           "gyroid",
		WallGenerator:           "arachne",
		BrimWidthMm:             2.5,
		SkirtLoops:              1,
		SlowDownLayers:          1,
		SlowDownLayerTimeSec:    1.0,
		SlowDownMinSpeedMmS:     20,

		// V11 process speeds (mm/s)
		OuterWallSpeedMmS:           50,
		InnerWallSpeedMmS:           80,
		SparseInfillSpeedMmS:        150,
		InternalSolidInfillSpeedMmS: 120,
		TopSurfaceSpeedMmS:          40,
		GapInfillSpeedMmS:           40,
		BridgeSpeedMmS:              25,
		InitialLayerSpeedMmS:        25,
		InitialLayerInfillSpeedMmS:  25,
		TravelSpeedMmS:              150,
		InitialLayerTravelSpeedMmS:  100,

		// V11 process accelerations (mm/s^2)
		DefaultAccelMmS2:             1500,
		OuterWallAccelMmS2:           800,
		InnerWallAccelMmS2:           1200,
		SparseInfillAccelMmS2:        1500,
		InternalSolidInfillAccelMmS2: 1500,
		TopSurfaceAccelMmS2:          800,
		BridgeAccelMmS2:              800,
		InitialLayerAccelMmS2:        500,
		TravelAccelMmS2:              1500,
		InitialLayerTravelAccelMmS2:  1000,

		// V11 process jerks (mm/s)
		DefaultJerkMmS:            18,
		OuterWallJerkMmS:          15,
		InnerWallJerkMmS:          15,
		InfillJerkMmS:             20,
		TopSurfaceJerkMmS:         8,
		InitialLayerJerkMmS:       5,
		TravelJerkMmS:             30,
		InitialLayerTravelJerkMmS: 20,
	}
}

// SlicerConfig holds configuration parameters, typically loaded from PocketBase slicer_settings
type SlicerConfig struct {
	Name                 string        `json:"name"`
	LayerHeightMm        float64       `json:"layer_height_mm"`         // 0.20 mm
	WallPerimeters       int           `json:"wall_perimeters"`         // 2 walls (V11)
	WallThicknessMm      float64       `json:"wall_thickness_mm"`       // 0.87 mm (0.42 outer + 0.45 inner)
	InfillPercentage     float64       `json:"infill_percentage"`       // 10% (V11 gyroid)
	TimeMultiplier       float64       `json:"time_multiplier"`         // e.g. 1.5x buffer timer
	WeightMultiplier     float64       `json:"weight_multiplier"`       // e.g. 1.5x weight buffer
	MachineHourlyRate    int64         `json:"machine_hourly_rate"`     // e.g. 250 Rs/hr
	LaborOverhead        int64         `json:"labor_overhead"`          // e.g. 50 Rs prep
	ElectricityHourly    int64         `json:"electricity_hourly"`      // e.g. 30 Rs/hr
	FailureRatePct       int64         `json:"failure_rate_pct"`        // e.g. 5%
	ProfitMultiplierPct  int64         `json:"profit_multiplier_pct"`   // e.g. 25%
	DefaultFilamentPrice int64         `json:"default_filament_price"`  // e.g. 8 Rs/gram for White Basic PLA
	Preset               PrinterPreset `json:"preset"`                  // Hardware & process kinematics preset
}

// FilamentProfile describes a selectable material option
type FilamentProfile struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`                       // e.g. "White Basic PLA"
	Material               string  `json:"material"`                   // e.g. "PLA"
	ColorName              string  `json:"color_name"`                 // e.g. "Signal White"
	ColorHex               string  `json:"color_hex"`                  // e.g. "#FAF7F2"
	DensityGCm3            float64 `json:"density_g_cm3"`              // e.g. 1.24 g/cm3
	FlowRatio              float64 `json:"flow_ratio"`                 // e.g. 0.979 (Cbot CN-1 White PLA)
	MaxVolumetricSpeedMm3S float64 `json:"max_volumetric_speed_mm3_s"` // e.g. 12 mm^3/s (Cbot CN-1 White PLA)
	PricePerGram           int64   `json:"price_per_gram"`             // in Rs per gram, e.g. 8
	IsDefault              bool    `json:"is_default"`
}

// DefaultSlicerConfig provides reliable fallbacks aligned with Cbot CN-1 0.4mm V11.
func DefaultSlicerConfig() SlicerConfig {
	preset := DefaultCbotCN1Preset()
	wallPerimeters := 2
	wallThickness := preset.OuterWallLineWidthMm + float64(wallPerimeters-1)*preset.InnerWallLineWidthMm
	return SlicerConfig{
		Name:                 "Cbot CN-1 0.4mm Nozzle (V11 Production Profile)",
		LayerHeightMm:        0.20,
		WallPerimeters:       wallPerimeters,
		WallThicknessMm:      wallThickness,
		InfillPercentage:     10.0, // V11 sparse_infill_density: 10% gyroid
		TimeMultiplier:       1.5,  // 1.5x user-specified buffer timer
		WeightMultiplier:     1.5,  // 1.5x user-specified buffer weight
		MachineHourlyRate:    250,  // Rs 250 per print hour
		LaborOverhead:        50,   // Rs 50 bench handling
		ElectricityHourly:    30,   // Rs 30/hr electricity
		FailureRatePct:       5,    // 5% allowance
		ProfitMultiplierPct:  25,   // 25% margin
		DefaultFilamentPrice: 8,    // Rs 8 / gram
		Preset:               preset,
	}
}

// DefaultFilaments provides the standard in-house materials tuned for Cbot CN-1 0.4mm nozzle.
func DefaultFilaments() []FilamentProfile {
	return []FilamentProfile{
		{
			ID:                     "white-basic-pla",
			Name:                   "White Basic PLA",
			Material:               "PLA",
			ColorName:              "Signal White",
			ColorHex:               "#FAF7F2",
			DensityGCm3:            1.24,
			FlowRatio:              0.979,
			MaxVolumetricSpeedMm3S: 12.0,
			PricePerGram:           8,
			IsDefault:              true,
		},
		{
			ID:                     "black-tough-pla",
			Name:                   "Black Tough PLA+",
			Material:               "PLA+",
			ColorName:              "Matte Black",
			ColorHex:               "#1A1A1A",
			DensityGCm3:            1.25,
			FlowRatio:              0.981,
			MaxVolumetricSpeedMm3S: 12.0,
			PricePerGram:           10,
			IsDefault:              false,
		},
		{
			ID:                     "high-temp-petg",
			Name:                   "High-Temp PETG",
			Material:               "PETG",
			ColorName:              "Terracotta Amber",
			ColorHex:               "#D9530F",
			DensityGCm3:            1.27,
			FlowRatio:              0.98,
			MaxVolumetricSpeedMm3S: 10.0,
			PricePerGram:           12,
			IsDefault:              false,
		},
		{
			ID:                     "flex-tpu-95a",
			Name:                   "TPU 95A Flexible",
			Material:               "TPU",
			ColorName:              "Emerald Flex",
			ColorHex:               "#17C3B2",
			DensityGCm3:            1.21,
			FlowRatio:              1.0,
			MaxVolumetricSpeedMm3S: 4.0,
			PricePerGram:           18,
			IsDefault:              false,
		},
	}
}

// SlicedQuoteResult holds the complete calculated slice and commercial invoice estimate
type SlicedQuoteResult struct {
	Mesh     MeshStats       `json:"mesh"`
	Filament FilamentProfile `json:"filament"`
	Config   SlicerConfig    `json:"config"`

	// Kinematics & Slicing Metrics
	LayerCount         int     `json:"layer_count"`
	RawVolumeCm3       float64 `json:"raw_volume_cm3"`
	ExtrudedVolCm3     float64 `json:"extruded_vol_cm3"`
	RawGrams           float64 `json:"raw_grams"`
	BufferedGrams      float64 `json:"buffered_grams"`       // Displayed grams with 1.5x buffer
	RawHours           float64 `json:"raw_hours"`
	BufferedHours      float64 `json:"buffered_hours"`       // Displayed hours with 1.5x buffer
	FormattedTime      string  `json:"formatted_time"`       // e.g. "3 hrs 45 mins"
	ExceedsBuildVolume bool    `json:"exceeds_build_volume"` // Exceeds Cbot CN-1 200x200x200mm bed

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

// normalizePreset ensures a zero-value PrinterPreset falls back to the Cbot CN-1 V11 defaults.
func normalizePreset(p PrinterPreset) PrinterPreset {
	if p.NozzleDiameterMm <= 0 || p.OuterWallSpeedMmS <= 0 || p.DefaultAccelMmS2 <= 0 {
		return DefaultCbotCN1Preset()
	}
	return p
}

// minPositive returns the smaller of two positive floats (ignoring <= 0 values).
func minPositive(a, b float64) float64 {
	if a <= 0 {
		return b
	}
	if b <= 0 {
		return a
	}
	if a < b {
		return a
	}
	return b
}

// effectiveSegmentSpeedMmS computes the physically achievable average speed (mm/s)
// over a segment of length segLenMm under a trapezoidal/triangular velocity profile,
// enforcing machine axis speed limits, filament volumetric flow caps, acceleration, and Marlin jerk.
func effectiveSegmentSpeedMmS(
	nominalSpeed, maxAxisSpeed, accel, maxAccel, jerk, maxJerk,
	lineWidthMm, layerHeightMm, maxVolumetricMm3S, segLenMm float64,
) float64 {
	v := minPositive(nominalSpeed, maxAxisSpeed)
	if maxVolumetricMm3S > 0 && lineWidthMm > 0 && layerHeightMm > 0 {
		volCap := maxVolumetricMm3S / (lineWidthMm * layerHeightMm)
		if volCap < v {
			v = volCap
		}
	}
	if v <= 0 {
		return 1.0
	}

	a := minPositive(accel, maxAccel)
	if a <= 0 {
		return v
	}

	j := minPositive(jerk, maxJerk)
	// Marlin classic jerk allows partial corner retention; 0.5*jerk models average turn deceleration.
	v0 := math.Min(v, j*0.5)
	if v0 < 0 {
		v0 = 0
	}

	if segLenMm <= 0 {
		return v
	}

	dRamp := (v*v - v0*v0) / a
	var tSeg float64
	if segLenMm >= dRamp {
		// Trapezoidal profile: accelerates to cruise speed v, cruises, decelerates to v0.
		tSeg = 2.0*(v-v0)/a + (segLenMm-dRamp)/v
	} else {
		// Triangular profile: acceleration-limited, peaks below nominal cruise speed.
		vPeak := math.Sqrt(v0*v0 + a*segLenMm)
		tSeg = 2.0 * (vPeak - v0) / a
	}
	if tSeg <= 0 {
		return v
	}
	return segLenMm / tSeg
}

// travelMoveSeconds computes the time for a single non-extruding XY travel move.
func travelMoveSeconds(distMm, speed, maxSpeed, accel, maxAccel, jerk, maxJerk float64) float64 {
	if distMm <= 0 {
		return 0
	}
	vEff := effectiveSegmentSpeedMmS(speed, maxSpeed, accel, maxAccel, jerk, maxJerk, 0, 0, 0, distMm)
	return distMm / vEff
}

// zHopCycleSeconds computes up+down Z-hop duration under Z speed and Z acceleration limits.
func zHopCycleSeconds(zHopMm, maxZSpeed, maxZAccel, maxZJerk float64) float64 {
	if zHopMm <= 0 || maxZSpeed <= 0 {
		return 0
	}
	oneWaySpeed := effectiveSegmentSpeedMmS(maxZSpeed, maxZSpeed, maxZAccel, maxZAccel, maxZJerk, maxZJerk, 0, 0, 0, zHopMm)
	return 2.0 * (zHopMm / oneWaySpeed)
}

// retractionCycleSeconds computes retract + deretract duration under E speed and acceleration limits.
func retractionCycleSeconds(lengthMm, retractSpeed, deretractSpeed, maxESpeed, maxRetractAccel, maxEJerk float64) float64 {
	if lengthMm <= 0 {
		return 0
	}
	vRet := effectiveSegmentSpeedMmS(retractSpeed, maxESpeed, maxRetractAccel, maxRetractAccel, maxEJerk, maxEJerk, 0, 0, 0, lengthMm)
	vDeret := effectiveSegmentSpeedMmS(deretractSpeed, maxESpeed, maxRetractAccel, maxRetractAccel, maxEJerk, maxEJerk, 0, 0, 0, lengthMm)
	return (lengthMm / vRet) + (lengthMm / vDeret)
}

// SliceAndCalculate evaluates the 3D model geometry using the Cbot CN-1 V11 kinematics preset
// and applies LayerInvoice commercial pricing math.
func SliceAndCalculate(mesh MeshStats, fil FilamentProfile, cfg SlicerConfig) SlicedQuoteResult {
	preset := normalizePreset(cfg.Preset)
	cfg.Preset = preset

	if cfg.LayerHeightMm <= 0 {
		cfg.LayerHeightMm = 0.20
	}
	if preset.MinLayerHeightMm > 0 && cfg.LayerHeightMm < preset.MinLayerHeightMm {
		cfg.LayerHeightMm = preset.MinLayerHeightMm
	}
	if preset.MaxLayerHeightMm > 0 && cfg.LayerHeightMm > preset.MaxLayerHeightMm {
		cfg.LayerHeightMm = preset.MaxLayerHeightMm
	}
	if cfg.WallPerimeters <= 0 {
		cfg.WallPerimeters = 2
	}
	if cfg.WallThicknessMm <= 0 {
		cfg.WallThicknessMm = preset.OuterWallLineWidthMm + float64(cfg.WallPerimeters-1)*preset.InnerWallLineWidthMm
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
	if fil.FlowRatio <= 0 {
		fil.FlowRatio = 0.979
	}
	if fil.MaxVolumetricSpeedMm3S <= 0 {
		// Default volumetric cap by material family (White PLA = 12 mm^3/s from Cbot CN-1 bundle)
		switch fil.Material {
		case "TPU":
			fil.MaxVolumetricSpeedMm3S = 4.0
		case "PETG":
			fil.MaxVolumetricSpeedMm3S = 10.0
		default:
			fil.MaxVolumetricSpeedMm3S = 12.0
		}
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

	exceedsBuildVolume := (preset.BuildXMm > 0 && mesh.SizeX > preset.BuildXMm) ||
		(preset.BuildYMm > 0 && mesh.SizeY > preset.BuildYMm) ||
		(preset.BuildZMm > 0 && mesh.SizeZ > preset.BuildZMm)

	// --- 1. Feature-Split Extrusion Volume (Walls, Top/Bottom Shells, Infill, Brim/Skirt) ---
	// Estimate horizontal footprint (top/bottom area) vs vertical side-wall surface area.
	bboxXYCm2 := (mesh.SizeX * mesh.SizeY) / 100.0
	horizontalAreaCm2 := math.Min(bboxXYCm2, mesh.SurfaceAreaCm2*0.45)
	if horizontalAreaCm2 < 0 {
		horizontalAreaCm2 = 0
	}
	sideAreaCm2 := mesh.SurfaceAreaCm2 - 2.0*horizontalAreaCm2
	if sideAreaCm2 < 0 {
		sideAreaCm2 = 0
	}

	// Side walls: 1 outer wall loop (0.42mm) + remaining inner wall loops (0.45mm each)
	outerWallCm := preset.OuterWallLineWidthMm / 10.0
	innerWallCm := float64(cfg.WallPerimeters-1) * preset.InnerWallLineWidthMm / 10.0
	if innerWallCm < 0 {
		innerWallCm = 0
	}
	outerWallVolCm3 := sideAreaCm2 * outerWallCm
	innerWallVolCm3 := sideAreaCm2 * innerWallCm

	// Top & bottom solid shells (V11: 5 top layers >= 0.6mm, 3 bottom layers)
	topThicknessMm := math.Max(float64(preset.TopShellLayers)*cfg.LayerHeightMm, preset.TopShellMinThicknessMm)
	bottomThicknessMm := float64(preset.BottomShellLayers) * cfg.LayerHeightMm
	topSurfaceVolCm3 := horizontalAreaCm2 * (cfg.LayerHeightMm / 10.0) // Topmost visible layer at top_surface_speed
	internalSolidMm := math.Max(0, topThicknessMm-cfg.LayerHeightMm) + bottomThicknessMm
	internalSolidVolCm3 := horizontalAreaCm2 * (internalSolidMm / 10.0)

	// Smooth 3D offset-shell saturation:
	// Flat-plate area*thickness overcounts inward corners, ridges, and sculpted details
	// (which otherwise clamps interior volume to 0 on organic meshes).
	// Using V_shell = V * (1 - exp(-V_slab / V)) matches exact cube/sphere shells within ~1-2%
	// while preserving a realistic infill core on high-surface-area meshes.
	rawShellVolCm3 := outerWallVolCm3 + innerWallVolCm3 + topSurfaceVolCm3 + internalSolidVolCm3
	effectiveShellVolCm3 := 0.0
	if mesh.VolumeCm3 > 0 && rawShellVolCm3 > 0 {
		effectiveShellVolCm3 = mesh.VolumeCm3 * (1.0 - math.Exp(-rawShellVolCm3/mesh.VolumeCm3))
		scale := effectiveShellVolCm3 / rawShellVolCm3
		outerWallVolCm3 *= scale
		innerWallVolCm3 *= scale
		topSurfaceVolCm3 *= scale
		internalSolidVolCm3 *= scale
	}

	// Sparse infill volume: interior core * infill ratio
	infillRatio := cfg.InfillPercentage / 100.0
	sparseInfillVolCm3 := (mesh.VolumeCm3 - effectiveShellVolCm3) * infillRatio
	if sparseInfillVolCm3 < 0 {
		sparseInfillVolCm3 = 0
	}

	// Adhesion brim + skirt on layer 1 (V11: 2.5mm outer brim + 1 skirt loop, 0.50mm initial width)
	approxPerimeterMm := 2.0 * (mesh.SizeX + mesh.SizeY)
	brimLenMm := 0.0
	if preset.BrimWidthMm > 0 && preset.InitialLayerLineWidthMm > 0 && approxPerimeterMm > 0 {
		brimLoops := math.Ceil(preset.BrimWidthMm / preset.InitialLayerLineWidthMm)
		brimLenMm = brimLoops * approxPerimeterMm
	}
	skirtLenMm := float64(preset.SkirtLoops) * (approxPerimeterMm + 40.0)
	if approxPerimeterMm <= 0 {
		skirtLenMm = 0
	}
	adhesionVolCm3 := ((brimLenMm + skirtLenMm) * preset.InitialLayerLineWidthMm * preset.InitialLayerHeightMm) / 1000.0

	// Apply calibrated filament flow ratio (0.979 for Cbot CN-1 White PLA)
	extrudedVolCm3 := (effectiveShellVolCm3 + sparseInfillVolCm3 + adhesionVolCm3) * fil.FlowRatio

	// Raw filament mass & 1.5x shop buffer (display truthful weight down to 0.1g; shop floor enforced via Rs 150 min order)
	rawGrams := extrudedVolCm3 * fil.DensityGCm3
	bufferedGrams := rawGrams * cfg.WeightMultiplier
	if bufferedGrams < 0.1 && mesh.VolumeCm3 > 0 {
		bufferedGrams = 0.1
	}

	// --- 2. Kinematics-Driven Print Time (Speeds + Accelerations + Jerks + Volumetric Cap) ---
	maxXYSpeed := minPositive(preset.MaxSpeedXMmS, preset.MaxSpeedYMmS)
	maxXYAccel := minPositive(preset.MaxAccelXMmS2, preset.MaxAccelYMmS2)
	maxExtrudeAccel := minPositive(preset.MaxAccelExtrudingMmS2, maxXYAccel)
	maxTravelAccel := minPositive(preset.MaxAccelTravelMmS2, maxXYAccel)
	maxXYJerk := minPositive(preset.MaxJerkXMmS, preset.MaxJerkYMmS)

	avgXY := (mesh.SizeX + mesh.SizeY) / 2.0
	wallSegLenMm := math.Max(4.0, avgXY*0.65)
	solidSegLenMm := math.Max(4.0, avgXY*0.55)
	// Gyroid infill curves continuously, giving shorter effective straight-line acceleration spans.
	infillSegLenMm := math.Max(3.0, math.Min(14.0, avgXY*0.30))

	h := cfg.LayerHeightMm
	qMax := fil.MaxVolumetricSpeedMm3S

	// Convert each feature's cm^3 volume into mm of toolpath length: len_mm = (vol_cm3 * 1000) / (width * h)
	outerWallLenMm := (outerWallVolCm3 * 1000.0) / (preset.OuterWallLineWidthMm * h)
	innerWallLenMm := (innerWallVolCm3 * 1000.0) / (preset.InnerWallLineWidthMm * h)
	topSurfaceLenMm := (topSurfaceVolCm3 * 1000.0) / (preset.TopSurfaceLineWidthMm * h)
	internalSolidLenMm := (internalSolidVolCm3 * 1000.0) / (preset.SolidInfillLineWidthMm * h)
	sparseInfillLenMm := (sparseInfillVolCm3 * 1000.0) / (preset.SparseInfillLineWidthMm * h)

	vOuter := effectiveSegmentSpeedMmS(preset.OuterWallSpeedMmS, maxXYSpeed, preset.OuterWallAccelMmS2, maxExtrudeAccel, preset.OuterWallJerkMmS, maxXYJerk, preset.OuterWallLineWidthMm, h, qMax, wallSegLenMm)
	vInner := effectiveSegmentSpeedMmS(preset.InnerWallSpeedMmS, maxXYSpeed, preset.InnerWallAccelMmS2, maxExtrudeAccel, preset.InnerWallJerkMmS, maxXYJerk, preset.InnerWallLineWidthMm, h, qMax, wallSegLenMm)
	vTop := effectiveSegmentSpeedMmS(preset.TopSurfaceSpeedMmS, maxXYSpeed, preset.TopSurfaceAccelMmS2, maxExtrudeAccel, preset.TopSurfaceJerkMmS, maxXYJerk, preset.TopSurfaceLineWidthMm, h, qMax, solidSegLenMm)
	vSolid := effectiveSegmentSpeedMmS(preset.InternalSolidInfillSpeedMmS, maxXYSpeed, preset.InternalSolidInfillAccelMmS2, maxExtrudeAccel, preset.InfillJerkMmS, maxXYJerk, preset.SolidInfillLineWidthMm, h, qMax, solidSegLenMm)
	vInfill := effectiveSegmentSpeedMmS(preset.SparseInfillSpeedMmS, maxXYSpeed, preset.SparseInfillAccelMmS2, maxExtrudeAccel, preset.InfillJerkMmS, maxXYJerk, preset.SparseInfillLineWidthMm, h, qMax, infillSegLenMm)
	vInitial := effectiveSegmentSpeedMmS(preset.InitialLayerSpeedMmS, maxXYSpeed, preset.InitialLayerAccelMmS2, maxExtrudeAccel, preset.InitialLayerJerkMmS, maxXYJerk, preset.InitialLayerLineWidthMm, preset.InitialLayerHeightMm, qMax, wallSegLenMm)

	extrusionSec := (outerWallLenMm / vOuter) +
		(innerWallLenMm / vInner) +
		(topSurfaceLenMm / vTop) +
		(internalSolidLenMm / vSolid) +
		(sparseInfillLenMm / vInfill)

	// Slow-down layer cooling floor (Generic PLA template: min 1s/layer, min speed 20 mm/s)
	if preset.SlowDownLayerTimeSec > 0 {
		minCoolingSec := float64(layerCount) * preset.SlowDownLayerTimeSec
		if extrusionSec < minCoolingSec && mesh.VolumeCm3 > 0 {
			extrusionSec = minCoolingSec
		}
	}

	// Initial layer slowdown + brim/skirt printing at initial_layer_speed (25 mm/s, 500 mm/s^2)
	adhesionSec := (brimLenMm + skirtLenMm) / vInitial
	firstLayerExtraSec := 0.0
	if layerCount > 0 && preset.SlowDownLayers > 0 {
		firstLayerVolCm3 := (outerWallVolCm3 + innerWallVolCm3 + internalSolidVolCm3) / float64(layerCount)
		firstLayerLenMm := (firstLayerVolCm3 * 1000.0) / (preset.InitialLayerLineWidthMm * preset.InitialLayerHeightMm)
		firstLayerAtSlow := firstLayerLenMm / vInitial
		firstLayerAtNormal := extrusionSec / float64(layerCount)
		if firstLayerAtSlow > firstLayerAtNormal {
			firstLayerExtraSec = float64(preset.SlowDownLayers) * (firstLayerAtSlow - firstLayerAtNormal)
		}
	}

	// Per-layer travel, retraction (6mm @ 35mm/s), wipe (0.4mm), and Z-hop (0.4mm @ 10mm/s, 100mm/s^2)
	retractSec := retractionCycleSeconds(
		preset.RetractionLengthMm, preset.RetractionSpeedMmS, preset.DeretractionSpeedMmS,
		preset.MaxSpeedEMmS, preset.MaxAccelRetractingMmS2, preset.MaxJerkEMmS,
	)
	zHopSec := zHopCycleSeconds(preset.ZHopMm, preset.MaxSpeedZMmS, preset.MaxAccelZMmS2, preset.MaxJerkZMmS)
	wipeSec := 0.0
	if preset.WipeDistanceMm > 0 && vOuter > 0 {
		wipeSec = preset.WipeDistanceMm / (vOuter * 0.8)
	}
	avgTravelDistMm := math.Max(6.0, avgXY*0.5)
	singleTravelSec := travelMoveSeconds(
		avgTravelDistMm, preset.TravelSpeedMmS, maxXYSpeed,
		preset.TravelAccelMmS2, maxTravelAccel, preset.TravelJerkMmS, maxXYJerk,
	)
	// Typical layer executes ~2.2 retraction+hop+travel cycles (seam, inner-to-infill, island hops)
	perLayerMotionSec := 2.2*(retractSec+zHopSec+wipeSec+singleTravelSec) + (cfg.LayerHeightMm / math.Max(1.0, preset.MaxSpeedZMmS))
	motionOverheadSec := float64(layerCount) * perLayerMotionSec

	// Total raw hours including machine_prepare_compensation_time (260s heating + G28 + prime lines)
	totalRawSec := extrusionSec + adhesionSec + firstLayerExtraSec + motionOverheadSec + preset.PrepareCompensationSec
	rawHours := totalRawSec / 3600.0
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
		ExceedsBuildVolume:  exceedsBuildVolume,
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
